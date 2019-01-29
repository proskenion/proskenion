package command_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommandValidator_Tx(t *testing.T) {
	fc := RandomFactory()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc, RandomConfig())

	acs := []*AccountWithPri{
		NewAccountWithPri("authorizer@com"),
		NewAccountWithPri("target1@com"),
		NewAccountWithPri("target2@com"),
		NewAccountWithPri("target3@com"),
	}
	GenesisCommitFromAccounts(t, rp, acs)

	top, ok := rp.Top()
	require.True(t, ok)
	rx, err := rp.Begin()
	require.NoError(t, err)
	wsv, err := rx.WSV(top.GetPayload().GetWSVHash())
	require.NoError(t, err)
	txh, err := rx.TxHistory(top.GetPayload().GetTxHistoryHash())

	require.NoError(t, err)
	for _, c := range []struct {
		name string
		tx   model.Transaction
		err  error
	}{
		{
			"case 1 correct",
			TxSign(t,
				fc.NewTxBuilder().CreateAccount("authorizer@com", "a@b", []model.PublicKey{}, 0).Build(),
				[]model.PublicKey{acs[0].Pubkey},
				[]model.PrivateKey{acs[0].Prikey}),
			nil,
		},
		{
			"case 2 different key",
			TxSign(t,
				fc.NewTxBuilder().CreateAccount("authorizer@com", "a@b", []model.PublicKey{}, 0).Build(),
				[]model.PublicKey{acs[1].Pubkey},
				[]model.PrivateKey{acs[1].Prikey}),
			core.ErrTxValidateNotSignedAuthorizer,
		},
		{
			"case 3 different key",
			fc.NewTxBuilder().CreateAccount("authorizer@com", "a@b", []model.PublicKey{}, 0).Build(),
			core.ErrTxValidateNotSignedAuthorizer,
		},// TODO "case 4 n-m multi-sig"
	} {
		t.Run(c.name, func(t *testing.T) {
			err := c.tx.Validate(wsv, txh)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func prePareCommandValidator(t *testing.T) (model.ModelFactory, core.CommandValidator, core.Repository) {
	fc, _, val, _, rp, _, _ := NewTestFactories()
	return fc, val, rp
}

func TestCommandValidator_CreateAccount(t *testing.T) {
	fc, val, rp := prePareCommandValidator(t)
	prePareCreateAccounts(t, fc, rp)

	_, wsv := prePareGetDtxWSV(t, rp)

	for _, c := range []struct {
		name           string
		exAuthoirzerId string
		exTargetId     string
		exErr          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			"target1@com",
			nil,
		},
		{
			"case 2 : duplicate error",
			authorizerId,
			"account1@com",
			core.ErrCommandExecutorCreateAccountAlreadyExistAccount,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_CreateAccount{
						CreateAccount: &proskenion.CreateAccount{},
					},
					TargetId:     c.exTargetId,
					AuthorizerId: c.exAuthoirzerId,
				}}
			err := val.CreateAccount(wsv, cmd)
			if c.exErr != nil {
				assert.EqualError(t, errors.Cause(err), c.exErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
