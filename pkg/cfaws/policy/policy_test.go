package policy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddExpiryCondition(t *testing.T) {
	p := Policy{
		Statements: []Statement{
			{
				Effect:   "Allow",
				Action:   []string{"*"},
				Resource: []string{"*"},
			},
		},
	}

	exp := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)

	AddExpiryCondition(&p, exp)
	cond := p.Statements[0].Condition.DateLessThan

	assert.Equal(t, exp, cond.Time)
}
