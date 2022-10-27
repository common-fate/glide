package slacknotifier

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/notifiers"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func (n *SlackWebhookNotifier) HandleRequestEvent(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {
	var requestEvent gevent.RequestEventPayload
	err := json.Unmarshal(event.Detail, &requestEvent)
	if err != nil {
		return err
	}
	req := requestEvent.Request

	ruleQuery := storage.GetAccessRuleVersion{ID: req.Rule, VersionID: req.RuleVersion}
	_, err = n.DB.Query(ctx, &ruleQuery)
	if err != nil {
		return errors.Wrap(err, "getting access rule")
	}
	rule := *ruleQuery.Result

	userQuery := storage.GetUser{ID: req.RequestedBy}
	_, err = n.DB.Query(ctx, &userQuery)
	if err != nil {
		return errors.Wrap(err, "getting requestor")
	}

	switch event.DetailType {
	case gevent.RequestCreatedType:
		if ruleQuery.Result.Approval.IsRequired() {
			msg := fmt.Sprintf("Your request to access *%s* requires approval. We've notified the approvers and will let you know once your request has been reviewed.", ruleQuery.Result.Name)
			fallback := fmt.Sprintf("Your request to access %s requires approval.", ruleQuery.Result.Name)

			_, err = SendMessage(ctx, n.client, userQuery.Result.Email, msg, fallback)
			if err != nil {
				log.Errorw("Failed to send direct message", "email", userQuery.Result.Email, "msg", msg, "error", err)
			}

			// Notify approvers
			reviewURL, err := notifiers.ReviewURL(n.FrontendURL, req.ID)
			if err != nil {
				return errors.Wrap(err, "building review URL")
			}

			// get the requestor's Slack user ID if it exists to render it nicely in the message to approvers.
			var slackUserID string
			requestor, err := n.client.GetUserByEmailContext(ctx, userQuery.Result.Email)
			if err != nil {
				zap.S().Infow("couldn't get slack user from requestor - falling back to email address", "requestor.id", userQuery.Result.ID, zap.Error(err))
			}
			if requestor != nil {
				slackUserID = requestor.ID
			}

			var wg sync.WaitGroup

			reviewers := storage.ListRequestReviewers{RequestID: req.ID}
			_, err = n.DB.Query(ctx, &reviewers)

			if err != nil {
				return errors.Wrap(err, "getting reviewers")
			}

			log.Infow("messaging reviewers", "reviewers", reviewers)

			for _, usr := range reviewers.Result {
				if usr.ReviewerID == req.RequestedBy {
					log.Infow("skipping sending approval message to requestor", "user.id", usr)
					continue
				}

				wg.Add(1)
				go func(usr access.Reviewer) {
					defer wg.Done()
					approver := storage.GetUser{ID: usr.ReviewerID}
					_, err := n.DB.Query(ctx, &approver)
					if err != nil {
						log.Errorw("failed to fetch user by id while trying to send message in slack", "user.id", usr, zap.Error(err))
						return
					}
					requestArguments, err := n.RenderRequestArguments(ctx, log, req, rule)
					if err != nil {
						log.Errorw("failed to generate request arguments, skipping including them in the slack message", "error", err)
					}
					summary, msg := BuildRequestMessage(RequestMessageOpts{
						Request:          req,
						RequestArguments: requestArguments,
						Rule:             rule,
						RequestorSlackID: slackUserID,
						RequestorEmail:   userQuery.Result.Email,
						ReviewURLs:       reviewURL,
					})

					ts, err := SendMessageBlocks(ctx, n.client, approver.Result.Email, msg, summary)
					if err != nil {
						log.Errorw("failed to send request approval message", "user", usr, "msg", msg, zap.Error(err))
					}

					updatedUsr := usr
					updatedUsr.Notifications = access.Notifications{
						SlackMessageID: &ts,
					}
					log.Infow("updating reviewer with slack msg id", "updatedUsr.SlackMessageID", ts)

					err = n.DB.Put(ctx, &updatedUsr)

					if err != nil {
						log.Errorw("failed to update reviewer", "user", usr, zap.Error(err))
					}
				}(usr)
			}

			wg.Wait()
		} else {
			//Review not required
			msg := fmt.Sprintf(":white_check_mark: Your request to access *%s* has been automatically approved. Hang tight - we're provisioning the role now and will let you know when it's ready.", ruleQuery.Result.Name)
			fallback := fmt.Sprintf("Your request to access %s has been automatically approved.", ruleQuery.Result.Name)
			_ = n.SendDMWithLogOnError(ctx, log, req.RequestedBy, msg, fallback)
		}
	case gevent.RequestApprovedType:
		msg := fmt.Sprintf("Your request to access *%s* has been approved. Hang tight - we're provisioning the access now and will let you know when it's ready.", ruleQuery.Result.Name)
		fallback := fmt.Sprintf("Your request to access %s has been approved.", ruleQuery.Result.Name)
		_ = n.SendDMWithLogOnError(ctx, log, req.RequestedBy, msg, fallback)

		// Loop over the request reviewers
		reviewers := storage.ListRequestReviewers{RequestID: req.ID}
		_, err = n.DB.Query(ctx, &reviewers)
		if err != nil {
			return errors.Wrap(err, "getting reviewers")
		}

		log.Infow("messaging reviewers", "reviewers", reviewers.Result)

		for _, rev := range reviewers.Result {
			err := n.UpdateSlackMessage(ctx, log, UpdateSlackMessageOpts{
				Review:            rev,
				Request:           req,
				RequestReviewerId: requestEvent.ReviewerID,
				Rule:              rule,
				DbRequestor:       userQuery.Result,
			})
			if err != nil {
				log.Errorw("failed to update slack message", "user", rev, zap.Error(err))
			}
		}
	case gevent.RequestCancelledType:
		// Loop over the request reviewers
		reviewers := storage.ListRequestReviewers{RequestID: req.ID}
		_, err = n.DB.Query(ctx, &reviewers)
		if err != nil {
			return errors.Wrap(err, "getting reviewers")
		}
		log.Infow("messaging reviewers", "reviewers", reviewers.Result)

		for _, usr := range reviewers.Result {
			err := n.UpdateSlackMessage(ctx, log,
				UpdateSlackMessageOpts{
					Review:            usr,
					Request:           req,
					RequestReviewerId: req.RequestedBy, // requestor ~= reviewer (they cancelled their own)
					Rule:              rule,
					DbRequestor:       userQuery.Result,
				})
			if err != nil {
				log.Errorw("failed to update slack message", "user", usr, "req", req, zap.Error(err))
			}
		}
	case gevent.RequestDeclinedType:
		msg := fmt.Sprintf("Your request to access *%s* has been declined.", ruleQuery.Result.Name)
		fallback := fmt.Sprintf("Your request to access %s has been declined.", ruleQuery.Result.Name)
		_ = n.SendDMWithLogOnError(ctx, log, req.RequestedBy, msg, fallback)

		// Loop over the request reviewers
		reviewers := storage.ListRequestReviewers{RequestID: req.ID}
		_, err = n.DB.Query(ctx, &reviewers)
		if err != nil {
			return errors.Wrap(err, "getting reviewers")
		}

		log.Infow("messaging reviewers", "reviewers", reviewers.Result)

		for _, usr := range reviewers.Result {
			err := n.UpdateSlackMessage(ctx, log,
				UpdateSlackMessageOpts{
					Review:            usr,
					Request:           req,
					RequestReviewerId: requestEvent.ReviewerID,
					Rule:              rule,
					DbRequestor:       userQuery.Result,
				})
			if err != nil {
				log.Errorw("failed to update slack message", "user", usr, zap.Error(err))
			}
		}
	}
	return nil
}

type UpdateSlackMessageOpts struct {
	Review            access.Reviewer
	Request           access.Request
	RequestReviewerId string
	Rule              rule.AccessRule
	DbRequestor       *identity.User
}

func (n *SlackWebhookNotifier) UpdateSlackMessage(ctx context.Context, log *zap.SugaredLogger, opts UpdateSlackMessageOpts) error {

	// Skip if requestor == reviewer

	// Get the reviewers email from db
	reviewerQuery := storage.GetUser{ID: opts.Review.ReviewerID}
	_, err := n.DB.Query(ctx, &reviewerQuery)
	if err != nil {
		return errors.Wrap(err, "getting reviewer")
	}
	// do the same but for the request reveiwer
	reqReviewer := storage.GetUser{ID: opts.RequestReviewerId}
	_, err = n.DB.Query(ctx, &reqReviewer)
	if err != nil && opts.Request.Status != access.CANCELLED {
		return errors.Wrap(err, "getting reviewer 2")
	}

	// get the requestor's Slack user ID if it exists to render it nicely in the message to approvers.
	var slackUserID string
	requestor, err := n.client.GetUserByEmailContext(ctx, opts.DbRequestor.Email)
	if err != nil {
		// log this instead of returning
		log.Errorw("failed to get slack user id, defaulting to email", "user", opts.DbRequestor.Email, zap.Error(err))
	}
	if requestor != nil {
		slackUserID = requestor.ID
	}
	reviewURL, err := notifiers.ReviewURL(n.FrontendURL, opts.Request.ID)
	if err != nil {
		return errors.Wrap(err, "building review URL")
	}

	requestArguments, err := n.RenderRequestArguments(ctx, log, opts.Request, opts.Rule)
	if err != nil {
		log.Errorw("failed to generate request arguments, skipping including them in the slack message", "error", err)
	}
	// Here we want to update the original approvers slack messages
	_, msg := BuildRequestMessage(RequestMessageOpts{
		Request:          opts.Request,
		RequestArguments: requestArguments,
		Rule:             opts.Rule,
		RequestorSlackID: slackUserID,
		RequestorEmail:   opts.DbRequestor.Email,
		ReviewURLs:       reviewURL,
		Reviewer:         reviewerQuery.Result,
		RequestReviewer:  reqReviewer.Result,
	})
	msg.Timestamp = *opts.Review.Notifications.SlackMessageID

	err = UpdateMessageBlocks(ctx, n.client, reviewerQuery.Result.Email, msg)
	if err != nil {
		return errors.Wrap(err, "failed to send updated request approval message")
	}
	return nil
}

type RequestMessageOpts struct {
	Request          access.Request
	RequestArguments []types.With
	Rule             rule.AccessRule
	ReviewURLs       notifiers.ReviewURLs
	RequestorSlackID string
	RequestorEmail   string
	Reviewer         *identity.User
	RequestReviewer  *identity.User
}

func BuildRequestMessage(o RequestMessageOpts) (summary string, msg slack.Message) {
	requestor := o.RequestorEmail
	if o.RequestorSlackID != "" {
		requestor = fmt.Sprintf("<@%s>", o.RequestorSlackID)
	}

	summary = fmt.Sprintf("New request for %s from %s", o.Rule.Name, o.RequestorEmail)

	when := "ASAP"
	if o.Request.RequestedTiming.StartTime != nil {
		t := o.Request.RequestedTiming.StartTime
		when = fmt.Sprintf("<!date^%d^{date_short_pretty} at {time}|%s>", t.Unix(), t.String())
	}

	status := strings.ToLower(string(o.Request.Status))
	status = strings.ToUpper(string(status[0])) + status[1:]

	requestDetails := []*slack.TextBlockObject{
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*When:*\n%s", when),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Duration:*\n%s", o.Request.RequestedTiming.Duration),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Status:*\n%s", status),
		},
	}

	for _, v := range o.RequestArguments {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%s:*\n%s", v.Title, v.Label),
		})
	}

	// Only show the Request reason if it is not empty
	if o.Request.Data.Reason != nil && len(*o.Request.Data.Reason) > 0 {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Request Reason:*\n%s", *o.Request.Data.Reason),
		})
	}

	msg = slack.NewBlockMessage(
		slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: fmt.Sprintf("*<%s|New request for %s> from %s*", o.ReviewURLs.Review, o.Rule.Name, requestor),
			},
		},
		slack.SectionBlock{
			Type:   slack.MBTSection,
			Fields: requestDetails,
		},
	)

	if o.Reviewer != nil || o.Request.Status == access.CANCELLED {
		t := time.Now()
		when = fmt.Sprintf("<!date^%d^{date_short_pretty} at {time}|%s>", t.Unix(), t.String())

		text := fmt.Sprintf("*Reviewed by* %s at %s", o.RequestReviewer.Email, when)

		if o.Request.Status == access.CANCELLED {
			text = fmt.Sprintf("*Cancelled by* %s at %s", o.RequestorEmail, when)
		}

		reviewContextBlock := slack.NewContextBlock("", slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: text,
		})

		msg.Blocks.BlockSet = append(msg.Blocks.BlockSet, reviewContextBlock)
	}

	// If the request has just been sent (PENDING), then append Action Blocks
	if o.Request.Status == access.PENDING {
		msg.Blocks.BlockSet = append(msg.Blocks.BlockSet, slack.NewActionBlock("review_actions",
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Approve"},
				Style:    slack.StylePrimary,
				ActionID: "approve",
				Value:    "approve",
				URL:      o.ReviewURLs.Approve,
			},
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Close Request"},
				Style:    slack.StyleDanger,
				ActionID: "deny",
				Value:    "deny",
				URL:      o.ReviewURLs.Deny,
			},
		))

	}

	return summary, msg
}

// @TODO this method maps request arguments in a deprecated way.
// it shoudl be replaced eventually with a cache lookup for the options available for the access rule
func (n *SlackWebhookNotifier) RenderRequestArguments(ctx context.Context, log *zap.SugaredLogger, request access.Request, rule rule.AccessRule) ([]types.With, error) {
	// Consider adding a fallback if the cache lookup fails
	pq := storage.ListCachedProviderOptions{
		ProviderID: rule.Target.ProviderID,
	}
	_, err := n.DB.Query(ctx, &pq)
	if err != nil && err != ddb.ErrNoItems {
		log.Errorw("failed to fetch provider options while trying to send message in slack", "provider.id", rule.Target.ProviderID, zap.Error(err))
	}
	var labelArr []types.With
	// Lookup the provider, ignore errors
	// if provider is not found, fallback to using the argument key as the title
	_, provider, _ := providerregistry.Registry().GetLatestByShortType(rule.Target.ProviderType)
	for k, v := range request.SelectedWith {
		with := types.With{
			Label: v.Label,
			Value: v.Value,
			Title: k,
		}
		// attempt to get the title for the argument from the provider arg schema
		if provider != nil {
			if s, ok := provider.Provider.(providers.ArgSchemarer); ok {
				t, ok := s.ArgSchema()[k]
				if ok {
					with.Title = t.Title
				}
			}
		}
		labelArr = append(labelArr, with)
	}

	for k, v := range rule.Target.With {
		// only include the with values if it does not have any groups selected,
		// if it does have groups selected, it means that it was a selectable field
		// so this check avoids duplicate/inaccurate values in the slack message
		if _, ok := rule.Target.WithArgumentGroupOptions[k]; !ok {
			with := types.With{
				Value: v,
				Title: k,
				Label: v,
			}
			// attempt to get the title for the argument from the provider arg schema
			if provider != nil {
				if s, ok := provider.Provider.(providers.ArgSchemarer); ok {
					t, ok := s.ArgSchema()[k]
					if ok {
						with.Title = t.Title
					}
				}
			}
			for _, ao := range pq.Result {
				// if a value is found, set it to true with a label
				if ao.Arg == k && ao.Value == v {
					with.Label = ao.Label
					break
				}
			}
			labelArr = append(labelArr, with)
		}
	}

	// now sort labelArr by Title
	sort.Slice(labelArr, func(i, j int) bool {
		return labelArr[i].Title < labelArr[j].Title
	})
	return labelArr, nil
}
