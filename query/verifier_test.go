package query_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/query"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO test 不十分
func TestNewQueryVerifier(t *testing.T) {
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
			ErrQueryVerifyAuthorizerIdNotAccountId,
		},
		{
			"case 3 peer",
			"peer1:50051",
			"authorizer@com",
			model.PeerObjectCode,
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			query := NewTestFactory().NewQueryBuilder().
				FromId(c.targetId).
				AuthorizerId(c.authorizerId).
				CreatedTime(RandomNow()).
				RequestCode(c.objectCode).
				Build()
			err := NewQueryVerifier().Verify(query)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
