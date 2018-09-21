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

func prePareCommandExecutor(t *testing.T) (model.ModelFactory, core.CommandExecutor, core.DBATx, core.WSV) {
	fc := NewTestFactory()
	cryptor := RandomCryptor()
	ex := RandomCommandExecutor()

	dtx := RandomDBATx()
	wsv, err := repository.NewWSV(dtx, cryptor, nil)
	require.NoError(t, err)
	return fc, ex, dtx, wsv
}

func prePareCreateAccounts(t *testing.T, fc model.ModelFactory, wsv core.WSV) {
	tx := fc.NewTxBuilder().
		CreateAccount("authorizer", "authorizer").
		CreateAccount("authorizer", "account1").
		CreateAccount("authorizer", "account2").
		CreateAccount("authorizer", "account3").
		CreateAccount("authorizer", "account4").
		CreateAccount("authorizer", "account5").
		Build()
	for _, cmd := range tx.GetPayload().GetCommands() {
		require.NoError(t, cmd.Execute(wsv))
	}
}

func prePareAddAsset(t *testing.T, fc model.ModelFactory, wsv core.WSV) {
	tx := fc.NewTxBuilder().
		AddAsset("authorizer", 1000).
		AddAsset("account1", 100).
		AddAsset("account2", 100).
		Build()
	for _, cmd := range tx.GetPayload().GetCommands() {
		require.NoError(t, cmd.Execute(wsv))
	}
}

func TestCommandExecutor_CreateAccount(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)

	for _, c := range []struct {
		name           string
		exAuthoirzerId string
		exTargetId     string
		exErr          error
	}{
		{
			"case 1 : no error",
			"authorizer",
			"authorizer",
			nil,
		},
		{
			"case 2 : duplicate error",
			"authorizer",
			"authorizer",
			core.ErrCommandExecutorCreateAccountAlreadyExistAccount,
		},
		{
			"case 3 : no error",
			"authorizer",
			"account1",
			nil,
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
			err := ex.CreateAccount(wsv, cmd)
			if c.exErr != nil {
				assert.EqualError(t, errors.Cause(err), c.exErr.Error())
			} else {
				assert.NoError(t, err)

				ac := fc.NewEmptyAccount()
				err = wsv.Query(c.exTargetId, ac)
				require.NoError(t, err)
				assert.Equal(t, c.exTargetId, ac.GetAccountId())
				assert.Equal(t, int64(0), ac.GetAmount())
				assert.Equal(t, make([]model.PublicKey, 0), ac.GetPublicKeys())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_AddAsset(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)

	for _, c := range []struct {
		name           string
		exAuthoirzerId string
		exTargetId     string
		addAmount      int64
		exAmount       int64
		exErr          error
	}{
		{
			"case 1 : no error",
			"authorizer",
			"authorizer",
			10,
			10,
			nil,
		},
		{
			"case 2 : over plus no error",
			"authorizer",
			"authorizer",
			10,
			20,
			nil,
		},
		{
			"case 3 : no error",
			"account1",
			"account1",
			10,
			10,
			nil,
		},
		{
			"case 4 : no account add asset error",
			"authorizer",
			"unk",
			10,
			10,
			core.ErrCommandExecutorAddAssetNotExistAccount,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_AddAsset{
						AddAsset: &proskenion.AddAsset{
							Amount: c.addAmount,
						},
					},
					TargetId:     c.exTargetId,
					AuthorizerId: c.exAuthoirzerId,
				}}
			err := ex.AddAsset(wsv, cmd)
			if c.exErr != nil {
				assert.EqualError(t, errors.Cause(err), c.exErr.Error())
			} else {
				assert.NoError(t, err)

				ac := fc.NewEmptyAccount()
				err = wsv.Query(c.exTargetId, ac)
				require.NoError(t, err)
				assert.Equal(t, c.exTargetId, ac.GetAccountId())
				assert.Equal(t, c.exAmount, ac.GetAmount())
				assert.Equal(t, make([]model.PublicKey, 0), ac.GetPublicKeys())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_Transfer(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)
	prePareAddAsset(t, fc, wsv)

	for _, c := range []struct {
		name            string
		exAuthoirzerId  string
		exTargetId      string
		exDestAccountId string
		transAmount     int64
		exSrcAmount     int64
		exDestAmount    int64
		exErr           error
	}{
		{
			"case 1 : no error",
			"authorizer",
			"authorizer",
			"account1",
			100,
			900,
			200,
			nil,
		},
		{
			"case 1 : no error",
			"authorizer",
			"account1",
			"account3",
			100,
			100,
			100,
			nil,
		},
		{
			"case 2 : no src account",
			"authorizer",
			"unk",
			"account3",
			100,
			100,
			100,
			core.ErrCommandExecutorTransferNotFoundSrcAccountId,
		},
		{
			"case 3 : no dest account",
			"authorizer",
			"authorizer",
			"unk",
			100,
			100,
			100,
			core.ErrCommandExecutorTransferNotFoundDestAccountId,
		},
		{
			"case 4 : not enough src account ammount",
			"authorizer",
			"account4",
			"account3",
			100,
			100,
			100,
			core.ErrCommandExecutorTransferNotEnoughSrcAccountAmount,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_Transfer{
						Transfer: &proskenion.Transfer{
							DestAccountId: c.exDestAccountId,
							Amount:        c.transAmount,
						},
					},
					TargetId:     c.exTargetId,
					AuthorizerId: c.exAuthoirzerId,
				}}
			err := ex.Transfer(wsv, cmd)
			if c.exErr != nil {
				assert.EqualError(t, errors.Cause(err), c.exErr.Error())
			} else {
				assert.NoError(t, err)

				srcAc := fc.NewEmptyAccount()
				err = wsv.Query(c.exTargetId, srcAc)
				require.NoError(t, err)
				assert.Equal(t, c.exTargetId, srcAc.GetAccountId())
				assert.Equal(t, c.exSrcAmount, srcAc.GetAmount())
				assert.Equal(t, make([]model.PublicKey, 0), srcAc.GetPublicKeys())

				destAc := fc.NewEmptyAccount()
				err = wsv.Query(c.exDestAccountId, destAc)
				require.NoError(t, err)
				assert.Equal(t, c.exDestAccountId, destAc.GetAccountId())
				assert.Equal(t, c.exDestAmount, destAc.GetAmount())
				assert.Equal(t, make([]model.PublicKey, 0), destAc.GetPublicKeys())

			}
		})
	}
	require.NoError(t, dtx.Commit())
}