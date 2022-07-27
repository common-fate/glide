package ad

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-multierror"
)

type ADErr struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Date            string `json:"date"`
			RequestID       string `json:"request-id"`
			ClientRequestID string `json:"client-request-id"`
		} `json:"innerError"`
	} `json:"error"`
}

// Validate the access against AzureAD without actually granting it.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// keep a running track of validation errors.
	var result error

	// The user should exist in azure.
	_, err = p.client.GetUser(ctx, subject)
	if err != nil {
		var adError ADErr
		err = json.Unmarshal([]byte(err.Error()), &adError)
		if err != nil {
			return err
		}
		if adError.Error.Code == "Request_ResourceNotFound" {
			err = &UserNotFoundError{User: subject}

		}

		result = multierror.Append(result, err)
	}

	// The group we are trying to grant access to should exist in AzureAD.
	_, err = p.client.GetGroup(ctx, a.GroupID)
	if err != nil {
		var adError ADErr
		err = json.Unmarshal([]byte(err.Error()), &adError)
		if err != nil {
			return err
		}
		if adError.Error.Code == "Request_BadRequest" {
			err = &GroupNotFoundError{Group: a.GroupID}

		}

		// add the error to our list.
		result = multierror.Append(result, err)
	}

	return result
}
