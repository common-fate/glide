package types

import (
	"context"
	"testing"
	"time"

	"github.com/common-fate/iso8601"
	"github.com/stretchr/testify/assert"
)

func TestValidateGrant(t *testing.T) {
	type testcase struct {
		name    string
		input   CreateGrant
		wantErr error
	}

	testcases := []testcase{
		{
			name:    "empty id",
			input:   CreateGrant{},
			wantErr: ErrInvalidGrantID,
		},
		{
			name: "start time after end",
			input: CreateGrant{
				Id:       "abcd",
				Provider: "test",
				Subject:  "test@acme.com",
				Start:    iso8601.New(time.Date(2022, 1, 1, 10, 10, 0, 0, time.UTC)),
				End:      iso8601.New(time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)),
			},
			wantErr: ErrInvalidGrantTime{"grant start time must be earlier than end time"},
		},
		{
			name: "start time equal to end",
			input: CreateGrant{
				Id:       "abcd",
				Provider: "test",
				Subject:  "test@acme.com",
				Start:    iso8601.New(time.Date(2022, 1, 1, 10, 10, 0, 0, time.UTC)),
				End:      iso8601.New(time.Date(2022, 1, 1, 10, 10, 0, 0, time.UTC)),
			},
			wantErr: ErrInvalidGrantTime{"grant start and end time cannot be equal"},
		},
		{
			name: "end time in the past",
			input: CreateGrant{
				Id:       "abcd",
				Provider: "test",
				Subject:  "test@acme.com",
				Start:    iso8601.New(time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC)),
				End:      iso8601.New(time.Date(2022, 1, 1, 9, 10, 0, 0, time.UTC)),
			},
			wantErr: ErrInvalidGrantTime{"grant finish time is in the past"},
		},
	}

	ctx := context.Background()

	// testing time is 1st Jan 2022, 10:00am UTC
	now := time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.Validate(ctx, now)
			if tc.wantErr == nil && err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}
		})
	}
}
