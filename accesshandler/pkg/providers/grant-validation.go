package providers

import (
	"context"
	"errors"
	"sync"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/hashicorp/go-multierror"
)

type GrantValidationResult struct {
	Name string
	Logs diagnostics.Logs
}
type GrantValidationSteps map[string]GrantValidationStep
type GrantValidationResults map[string]GrantValidationResult
type GrantValidationStep struct {
	Name string
	Run  func(ctx context.Context, subject string, args []byte) diagnostics.Logs
}

func (s GrantValidationSteps) Run(ctx context.Context, subject string, args []byte) GrantValidationResults {
	validationResults := make(GrantValidationResults)
	var mu sync.Mutex
	var wg sync.WaitGroup
	for key, val := range s {
		k := key
		v := val
		wg.Add(1)
		go func() {
			logs := v.Run(ctx, subject, args)
			func(key string, value GrantValidationStep, logs diagnostics.Logs) {
				mu.Lock()
				defer mu.Unlock()
				validationResults[key] = GrantValidationResult{Logs: logs, Name: v.Name}
			}(k, v, logs)
			wg.Done()
		}()
	}
	wg.Wait()
	return validationResults
}

func (r GrantValidationResults) Failed() bool {
	for _, v := range r {
		if !v.Logs.HasSucceeded() {
			return true
		}
	}
	return false
}

// FailureMessage returns an error string containing the names of the failed validation steps, else an empty string
func (r GrantValidationResults) FailureMessage() string {
	if !r.Failed() {
		return ""
	}
	var message error
	for _, v := range r {
		if !v.Logs.HasSucceeded() {
			message = multierror.Append(message, errors.New(v.Name))
		}
	}
	return message.Error()
}
