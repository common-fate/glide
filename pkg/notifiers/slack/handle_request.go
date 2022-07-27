package slacknotifier

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/notifiers"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func (n *Notifier) HandleRequestEvent(ctx context.Context, log *zap.SugaredLogger, slackClient *slack.Client, event events.CloudWatchEvent) error {
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

			_, err = SendMessage(ctx, slackClient, userQuery.Result.Email, msg, fallback)
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
			requestor, err := slackClient.GetUserByEmailContext(ctx, userQuery.Result.Email)
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
				go func() {
					defer wg.Done()
					approver := storage.GetUser{ID: usr.ReviewerID}
					_, err := n.DB.Query(ctx, &approver)
					if err != nil {
						log.Errorw("failed to fetch user by id while trying to send message in slack", "user.id", usr, zap.Error(err))
						return
					}

					summary, msg := BuildRequestMessage(RequestMessageOpts{
						Request:          req,
						Rule:             rule,
						RequestorSlackID: slackUserID,
						RequestorEmail:   userQuery.Result.Email,
						ReviewURLs:       reviewURL,
					})

					ts, err := SendMessageBlocks(ctx, slackClient, approver.Result.Email, msg, summary)
					if err != nil {
						log.Errorw("failed to send request approval message", "user", usr, zap.Error(err))
					}

					usr.SlackMessageID = ts
					log.Infow("updating reviewer with slack msg id", "usr.SlackMessageID", ts)

					err = n.DB.Put(ctx, &usr)

					if err != nil {
						log.Errorw("failed to update reviewer", "user", usr, zap.Error(err))
					}
				}()
			}

			wg.Wait()
		} else {
			//Review not required
			msg := fmt.Sprintf(":white_check_mark: Your request to access *%s* has been automatically approved. Hang tight - we're provisioning the role now and will let you know when it's ready.", ruleQuery.Result.Name)
			fallback := fmt.Sprintf("Your request to access %s has been automatically approved.", ruleQuery.Result.Name)
			n.SendDMWithLogOnError(ctx, slackClient, log, userQuery.Result.Email, msg, fallback)
		}
	case gevent.RequestApprovedType:
		msg := fmt.Sprintf("Your request to access *%s* has been approved. Hang tight - we're provisioning the access now and will let you know when it's ready.", ruleQuery.Result.Name)
		fallback := fmt.Sprintf("Your request to access %s has been approved.", ruleQuery.Result.Name)
		n.SendDMWithLogOnError(ctx, slackClient, log, userQuery.Result.Email, msg, fallback)

		// Loop over the request approvers
		reviewers := storage.ListRequestReviewers{RequestID: req.ID}
		// reviewers, err := rulesvc.GetApprovers(ctx, n.DB, rule)
		_, err = n.DB.Query(ctx, &reviewers)

		if err != nil {
			return errors.Wrap(err, "getting reviewers")
		}

		for _, usr := range reviewers.Result {
			if usr.ReviewerID == req.RequestedBy {
				// Skip the requestor
				continue
			}
			// get the requestor's Slack user ID if it exists to render it nicely in the message to approvers.
			var slackUserID string
			requestor, err := slackClient.GetUserByEmailContext(ctx, userQuery.Result.Email)
			if err != nil {
				zap.S().Infow("couldn't get slack user from requestor - falling back to email address", "requestor.id", userQuery.Result.ID, zap.Error(err))
			}
			if requestor != nil {
				slackUserID = requestor.ID
			}
			reviewURL, err := notifiers.ReviewURL(n.FrontendURL, req.ID)
			if err != nil {
				return errors.Wrap(err, "building review URL")
			}
			// Here we want to update the original approvers slack messages
			_, msg := BuildRequestMessage(RequestMessageOpts{
				Request:          req,
				Rule:             rule,
				RequestorSlackID: slackUserID,
				RequestorEmail:   userQuery.Result.Email,
				ReviewURLs:       reviewURL,
			})
			approver := storage.GetUser{ID: usr.ReviewerID}
			_, err = n.DB.Query(ctx, &approver)
			if err != nil {
				log.Errorw("failed to fetch user by id while trying to send message in slack", "user.id", usr, zap.Error(err))
				continue
			}
			// @TODO ENSURE THAT THIS UPDATES RATHER THAN SENDS NEW
			msg.Timestamp = usr.SlackMessageID
			err = UpdateMessageBlocks(ctx, slackClient, approver.Result.Email, msg)
			if err != nil {
				log.Errorw("failed to send updated request approval message", "user", usr, zap.Error(err))
			}
		}

	case gevent.RequestCancelledType:
		// TODO: handle update
	case gevent.RequestDeclinedType:
		// TODO: handle update
	}
	return nil
}

type RequestMessageOpts struct {
	Request          access.Request
	Rule             rule.AccessRule
	ReviewURLs       notifiers.ReviewURLs
	RequestorSlackID string
	RequestorEmail   string
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
			Text: fmt.Sprintf("*Status:*\n%s", o.Request.Status),
		},
	}

	if o.Request.Data.Reason != nil {
		requestDetails = append(requestDetails, &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Reason:*\n%s", *o.Request.Data.Reason),
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
