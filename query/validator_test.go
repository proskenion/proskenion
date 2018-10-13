package query_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/query"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO test 不十分
func TestNewQueryValidator(t *testing.T) {
	for _, c := range []struct {
		name         string
		targetId     string
		authorizerId string
		objectCode   model.ObjectCode
		err          error
	}{
		{
			"case 1",
			"target@com",
			"authorizer@com",
			model.AccountObjectCode,
			nil,
		},
		{
			"case 2 authorizer not account",
			"target@com",
			"authoirizer",
			model.AccountObjectCode,
			query.ErrQueryValidateAuthorizerIdNotAccountId,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			query := NewTestFactory().NewQueryBuilder().
				TargetId(c.targetId).
				AuthorizerId(c.authorizerId).
				CreatedTime(RandomNow()).
				Build()
			err := query.Validate()
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}