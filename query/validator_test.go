package query_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	. "github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TODO 不十分
func TestQueryValidator_Query(t *testing.T) {
	fc := RandomFactory()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc, RandomConfig())

	// GenesisCommit
	authorizer := NewAccountWithPri("authorizer@com")
	genesisCommit(t, rp, authorizer)

	qv := NewQueryValidator(rp, fc, RandomConfig())

	query := GetAccountQuery(t, authorizer, "target@com")
	err := qv.Validate(query)
	require.NoError(t, err)

	q2 := GetAccountQuery(t, authorizer, "target1@com")
	err = qv.Validate(q2)
	require.NoError(t, err)

	tmpub, tmpri := RandomKeyPairs()
	q3 := GetAccountQuery(t, &AccountWithPri{authorizer.AccountId, tmpub, tmpri}, "target@com")
	err = qv.Validate(q3)
	assert.EqualError(t, errors.Cause(err), core.ErrQueryProcessorNotSignedAuthorizer.Error())

	q4 := GetAccountQuery(t, &AccountWithPri{"authorizer1@com", tmpub, tmpri}, "target@com")
	err = qv.Validate(q4)
	assert.EqualError(t, errors.Cause(err), core.ErrQueryProcessorNotExistAuthoirizer.Error())
}
