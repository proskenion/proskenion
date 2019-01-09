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

const authorizerId = "authorizer@com"

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
		CreateAccount(authorizerId, authorizerId, []model.PublicKey{}, 0).
		CreateAccount(authorizerId, "account1@com", []model.PublicKey{}, 0).
		CreateAccount(authorizerId, "account2@com", []model.PublicKey{}, 0).
		CreateAccount(authorizerId, "account3@com", []model.PublicKey{}, 0).
		CreateAccount(authorizerId, "account4@com", []model.PublicKey{}, 0).
		CreateAccount(authorizerId, "account5@com", []model.PublicKey{}, 0).
		Build()
	for _, cmd := range tx.GetPayload().GetCommands() {
		require.NoError(t, cmd.Execute(wsv))
	}
}

func prePareAddBalance(t *testing.T, fc model.ModelFactory, wsv core.WSV) {
	tx := fc.NewTxBuilder().
		AddBalance(authorizerId, 1000).
		AddBalance("account1@com", 100).
		AddBalance("account2@com", 100).
		Build()
	executeCommands(t, tx, wsv)
}

func prePareCreateStorage(t *testing.T, fc model.ModelFactory, wsv model.ObjectFinder) {
	mtSt := RandomLandStorage("none", authorizerId, 0, make([]model.Object, 0))
	walletIds := []string{
		"account1@com/land",
		"account2@com/land",
		"account3@com/land",
		"account4@com/land",
		"account5@com/land",
	}
	builder := fc.NewTxBuilder().
		DefineStorage(authorizerId, "/land", mtSt)
	for _, id := range walletIds {
		builder = builder.CreateStorage(authorizerId, id)
	}
	executeCommands(t, builder.Build(), wsv)
}

func prePareAddPeer(t *testing.T, fc model.ModelFactory, wsv model.ObjectFinder) {
	tx := fc.NewTxBuilder().
		AddPeer(authorizerId, "peer1@com", "0.0.0.0:5050", RandomPublicKey()).
		AddPeer(authorizerId, "peer2@com", "0.0.0.1:5050", RandomPublicKey()).
		AddPeer(authorizerId, "peer3@com", "0.0.0.2:5050", RandomPublicKey()).
		AddPeer(authorizerId, "peer4@com", "0.0.0.3:5050", RandomPublicKey()).
		Build()
	executeCommands(t, tx, wsv)

}

func executeCommands(t *testing.T, tx model.Transaction, wsv model.ObjectFinder) {
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
			authorizerId,
			authorizerId,
			nil,
		},
		{
			"case 2 : duplicate error",
			authorizerId,
			authorizerId,
			core.ErrCommandExecutorCreateAccountAlreadyExistAccount,
		},
		{
			"case 3 : no error",
			authorizerId,
			"account1@com",
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
				err = wsv.Query(model.MustAddress(model.MustAddress(c.exTargetId).AccountId()), ac)
				require.NoError(t, err)
				assert.Equal(t, c.exTargetId, ac.GetAccountId())
				assert.Equal(t, int64(0), ac.GetBalance())
				assert.Equal(t, make([]model.PublicKey, 0), ac.GetPublicKeys())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_AddBalance(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)

	for _, c := range []struct {
		name           string
		exAuthoirzerId string
		exTargetId     string
		addBalance     int64
		exBalance      int64
		exErr          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			authorizerId,
			10,
			10,
			nil,
		},
		{
			"case 2 : over plus no error",
			authorizerId,
			authorizerId,
			10,
			20,
			nil,
		},
		{
			"case 3 : no error",
			"account1@com",
			"account1@com",
			10,
			10,
			nil,
		},
		{
			"case 4 : no account add asset error",
			authorizerId,
			"unk@unk",
			10,
			10,
			core.ErrCommandExecutorAddBalanceNotExistAccount,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_AddBalance{
						AddBalance: &proskenion.AddBalance{
							Balance: c.addBalance,
						},
					},
					TargetId:     c.exTargetId,
					AuthorizerId: c.exAuthoirzerId,
				}}
			err := ex.AddBalance(wsv, cmd)
			if c.exErr != nil {
				assert.EqualError(t, errors.Cause(err), c.exErr.Error())
			} else {
				assert.NoError(t, err)

				ac := fc.NewEmptyAccount()
				err = wsv.Query(model.MustAddress(model.MustAddress(c.exTargetId).AccountId()), ac)
				require.NoError(t, err)
				assert.Equal(t, c.exTargetId, ac.GetAccountId())
				assert.Equal(t, c.exBalance, ac.GetBalance())
				assert.Equal(t, make([]model.PublicKey, 0), ac.GetPublicKeys())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_TransferBalance(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)
	prePareAddBalance(t, fc, wsv)

	for _, c := range []struct {
		name            string
		exAuthoirzerId  string
		exTargetId      string
		exDestAccountId string
		transBalance    int64
		exSrcBalance    int64
		exDestBalance   int64
		exErr           error
	}{
		{
			"case 1 : no error",
			authorizerId,
			authorizerId,
			"account1@com",
			100,
			900,
			200,
			nil,
		},
		{
			"case 1 : no error",
			authorizerId,
			"account1@com",
			"account3@com",
			100,
			100,
			100,
			nil,
		},
		{
			"case 2 : no src account",
			authorizerId,
			"unk@unk",
			"account3@com",
			100,
			100,
			100,
			core.ErrCommandExecutorTransferBalanceNotFoundSrcAccountId,
		},
		{
			"case 3 : no dest account",
			authorizerId,
			authorizerId,
			"unk@unk",
			100,
			100,
			100,
			core.ErrCommandExecutorTransferBalanceNotFoundDestAccountId,
		},
		{
			"case 4 : not enough src account ammount",
			authorizerId,
			"account4@com",
			"account3@com",
			100,
			100,
			100,
			core.ErrCommandExecutorTransferBalanceNotEnoughSrcAccountBalance,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_TransferBalance{
						TransferBalance: &proskenion.TransferBalance{
							DestAccountId: c.exDestAccountId,
							Balance:       c.transBalance,
						},
					},
					TargetId:     c.exTargetId,
					AuthorizerId: c.exAuthoirzerId,
				}}
			err := ex.TransferBalance(wsv, cmd)
			if c.exErr != nil {
				assert.EqualError(t, errors.Cause(err), c.exErr.Error())
			} else {
				assert.NoError(t, err)

				srcAc := fc.NewEmptyAccount()
				err = wsv.Query(model.MustAddress(model.MustAddress(c.exTargetId).AccountId()), srcAc)
				require.NoError(t, err)
				assert.Equal(t, c.exTargetId, srcAc.GetAccountId())
				assert.Equal(t, c.exSrcBalance, srcAc.GetBalance())
				assert.Equal(t, make([]model.PublicKey, 0), srcAc.GetPublicKeys())

				destAc := fc.NewEmptyAccount()
				err = wsv.Query(model.MustAddress(model.MustAddress(c.exDestAccountId).AccountId()), destAc)
				require.NoError(t, err)
				assert.Equal(t, c.exDestAccountId, destAc.GetAccountId())
				assert.Equal(t, c.exDestBalance, destAc.GetBalance())
				assert.Equal(t, make([]model.PublicKey, 0), destAc.GetPublicKeys())

			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_AddPublicKey(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)

	keys := []model.PublicKey{
		RandomPublicKey(),
		RandomPublicKey(),
		RandomPublicKey(),
	}

	for _, c := range []struct {
		name         string
		authorizerId string
		targetId     string
		key          model.PublicKey
		exKeys       []model.PublicKey
		err          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			authorizerId,
			keys[0],
			[]model.PublicKey{keys[0]},
			nil,
		},
		{
			"case 2 : no error",
			authorizerId,
			authorizerId,
			keys[1],
			[]model.PublicKey{keys[0], keys[1]},
			nil,
		},
		{
			"case 3 : no error",
			authorizerId,
			"account1@com",
			keys[2],
			[]model.PublicKey{keys[2]},
			nil,
		},
		{
			"case 3 : no error",
			authorizerId,
			"account1@com",
			keys[0],
			[]model.PublicKey{keys[2], keys[0]},
			nil,
		},
		{
			"case 4 : no target account",
			authorizerId,
			"unk@unk",
			keys[2],
			nil,
			core.ErrCommandExecutorAddPublicKeyNotExistAccount,
		},
		{
			"case 5 : duplicate pubkey",
			authorizerId,
			authorizerId,
			keys[1],
			nil,
			core.ErrCommandExecutorAddPublicKeyDuplicatePubkey,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_AddPublicKeys{
						AddPublicKeys: &proskenion.AddPublicKeys{
							PublicKeys: [][]byte{c.key},
						},
					},
					TargetId:     c.targetId,
					AuthorizerId: c.authorizerId,
				}}
			err := ex.AddPublicKeys(wsv, cmd)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)

				ac := fc.NewEmptyAccount()
				err = wsv.Query(model.MustAddress(model.MustAddress(c.targetId).AccountId()), ac)
				require.NoError(t, err)
				assert.Equal(t, c.targetId, ac.GetAccountId())
				assert.ElementsMatch(t, c.exKeys, ac.GetPublicKeys())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_DefineStorage(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)

	for _, c := range []struct {
		name         string
		authorizerId string
		storageId    string
		storage      model.Storage
		err          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			"/wallet",
			fc.NewStorageBuilder().
				Str("name", "default").
				Address("id", "authorizer@com").
				Build(),
			nil,
		},
		{
			"case 2 : no error",
			authorizerId,
			"/plane",
			fc.NewStorageBuilder().
				Str("address", "0-index").
				Address("owner", "none").
				Int64("value", 0).
				Build(),
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_DefineStorage{
						DefineStorage: &proskenion.DefineStorage{
							Storage: c.storage.(*convertor.Storage).Storage,
						},
					},
					TargetId:     c.storageId,
					AuthorizerId: c.authorizerId,
				}}
			err := ex.DefineStorage(wsv, cmd)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)
				st := fc.NewEmptyStorage()
				err = wsv.Query(model.MustAddress(c.storageId), st)
				require.NoError(t, err)
				assert.Equal(t, c.storage.Hash(), st.Hash())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_CreateStorage(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)

	for _, c := range []struct {
		name         string
		authorizerId string
		storageId    string
		storage      model.Storage
		err          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			"account1@com/wallet",
			fc.NewStorageBuilder().
				Str("name", "account1").
				Address("id", "account1@com").
				Build(),
			nil,
		},
		{
			"case 2 : no error",
			authorizerId,
			"account1@com/plane",
			fc.NewStorageBuilder().
				Str("address", "1-1-1").
				Address("owner", "account1@com").
				Int64("value", 11111).
				Build(),
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			id := model.MustAddress(c.storageId)
			tx := fc.NewTxBuilder().DefineStorage(c.authorizerId, "/"+id.Storage(), c.storage).Build()
			tx.GetPayload().GetCommands()[0].Execute(wsv)

			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_CreateStorage{
						CreateStorage: &proskenion.CreateStorage{},
					},
					TargetId:     c.storageId,
					AuthorizerId: c.authorizerId,
				}}
			err := ex.CreateStorage(wsv, cmd)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)

				st := fc.NewEmptyStorage()
				err = wsv.Query(model.MustAddress(c.storageId), st)
				require.NoError(t, err)
				assert.Equal(t, c.storage.Hash(), st.Hash())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_UpdateObjectStorage(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)
	prePareCreateStorage(t, fc, wsv)

	for _, c := range []struct {
		name         string
		authorizerId string
		storageId    string
		key          string
		object       model.Object
		err          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			"account1@com/land",
			"address",
			fc.NewObjectBuilder().Str("1-1-1").Build(),
			nil,
		},
		{
			"case 2 : no error",
			authorizerId,
			"account2@com/land",
			"value",
			fc.NewObjectBuilder().Int64(1000).Build(),
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_UpdateObject{
						UpdateObject: &proskenion.UpdateObject{
							Key:    c.key,
							Object: c.object.(*convertor.Object).Object,
						},
					},
					TargetId:     c.storageId,
					AuthorizerId: c.authorizerId,
				}}
			err := ex.UpdateObject(wsv, cmd)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)

				st := fc.NewEmptyStorage()
				err = wsv.Query(model.MustAddress(c.storageId), st)
				require.NoError(t, err)
				assert.Equal(t, c.object.Hash(), st.GetObject()[c.key].Hash())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_AddObjectStorage(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)
	prePareCreateStorage(t, fc, wsv)

	for _, c := range []struct {
		name         string
		authorizerId string
		storageId    string
		key          string
		object       model.Object
		objects      []model.Object
		err          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			"account1@com/land",
			"list",
			fc.NewObjectBuilder().Str("a-a").Build(),
			[]model.Object{fc.NewObjectBuilder().Str("a-a").Build()},
			nil,
		},
		{
			"case 2 : no error",
			authorizerId,
			"account1@com/land",
			"list",
			fc.NewObjectBuilder().Str("b-b").Build(),
			[]model.Object{
				fc.NewObjectBuilder().Str("a-a").Build(),
				fc.NewObjectBuilder().Str("b-b").Build(),
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := &convertor.Command{
				Command: &proskenion.Command{
					Command: &proskenion.Command_AddObject{
						AddObject: &proskenion.AddObject{
							Key:    c.key,
							Object: c.object.(*convertor.Object).Object,
						},
					},
					TargetId:     c.storageId,
					AuthorizerId: c.authorizerId,
				}}
			err := ex.AddObject(wsv, cmd)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)

				st := fc.NewEmptyStorage()
				err = wsv.Query(model.MustAddress(c.storageId), st)
				require.NoError(t, err)
				assert.Equal(t, fc.NewObjectBuilder().List(c.objects).Build().Hash(),
					st.GetObject()[c.key].Hash())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_TransferObjectStorage(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)
	prePareCreateStorage(t, fc, wsv)

	executeCommands(t, fc.NewTxBuilder().
		AddObject(authorizerId, "account1@com/land",
			"list", fc.NewObjectBuilder().Str("aa").Build()).
		AddObject(authorizerId, "account1@com/land",
			"list", fc.NewObjectBuilder().Str("bb").Build()).
		AddObject(authorizerId, "account1@com/land",
			"list", fc.NewObjectBuilder().Str("cc").Build()).
		Build(), wsv)

	for _, c := range []struct {
		name         string
		authorizerId string
		storageId    string
		destId       string
		key          string
		object       model.Object
		srcObs       []model.Object
		dstObs       []model.Object
		err          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			"account1@com/land",
			"account2@com/land",
			"list",
			fc.NewObjectBuilder().Str("bb").Build(),
			[]model.Object{
				fc.NewObjectBuilder().Str("aa").Build(),
				fc.NewObjectBuilder().Str("cc").Build(),
			},
			[]model.Object{fc.NewObjectBuilder().Str("bb").Build()},
			nil,
		},
		{
			"case 2 : no error",
			authorizerId,
			"account1@com/land",
			"account2@com/land",
			"list",
			fc.NewObjectBuilder().Str("aa").Build(),
			[]model.Object{fc.NewObjectBuilder().Str("cc").Build()},
			[]model.Object{
				fc.NewObjectBuilder().Str("bb").Build(),
				fc.NewObjectBuilder().Str("aa").Build(),
			},
			nil,
		},
		{
			"case 3 : no error",
			authorizerId,
			"account1@com/land",
			"account2@com/land",
			"list",
			fc.NewObjectBuilder().Str("cc").Build(),
			make([]model.Object, 0),
			[]model.Object{
				fc.NewObjectBuilder().Str("bb").Build(),
				fc.NewObjectBuilder().Str("aa").Build(),
				fc.NewObjectBuilder().Str("cc").Build(),
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := fc.NewTxBuilder().
				TransferObject(c.authorizerId, c.storageId, c.destId, c.key, c.object).
				Build().
				GetPayload().
				GetCommands()[0]
			err := ex.TransferObject(wsv, cmd)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)

				srcSt := fc.NewEmptyStorage()
				err = wsv.Query(model.MustAddress(c.storageId), srcSt)
				require.NoError(t, err)
				assert.Equal(t, fc.NewObjectBuilder().List(c.srcObs).Build().Hash(),
					srcSt.GetObject()[c.key].Hash())

				dstSt := fc.NewEmptyStorage()
				err = wsv.Query(model.MustAddress(c.destId), dstSt)
				require.NoError(t, err)
				assert.Equal(t, fc.NewObjectBuilder().List(c.dstObs).Build().Hash(),
					dstSt.GetObject()[c.key].Hash())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_AddPeer(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)

	for _, c := range []struct {
		name         string
		authorizerId string
		peerId       string
		address      string
		publicKey    model.PublicKey
		err          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			"peer1@com",
			"127.0.0.1",
			RandomPublicKey(),
			nil,
		},
		{
			"case 1 : no error",
			authorizerId,
			"peer2@com",
			"127.0.0.2",
			RandomPublicKey(),
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := fc.NewTxBuilder().
				AddPeer(c.authorizerId, c.peerId, c.address, c.publicKey).
				Build().
				GetPayload().
				GetCommands()[0]
			err := ex.AddPeer(wsv, cmd)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)

				pr := fc.NewEmptyPeer()
				err = wsv.Query(model.MustAddress(model.MustAddress(c.peerId).PeerId()), pr)
				require.NoError(t, err)
				assert.Equal(t, c.peerId, pr.GetPeerId())
				assert.Equal(t, c.address, pr.GetAddress())
				assert.Equal(t, c.publicKey, pr.GetPublicKey())
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

func TestCommandExecutor_Consign(t *testing.T) {
	fc, ex, dtx, wsv := prePareCommandExecutor(t)
	prePareCreateAccounts(t, fc, wsv)
	prePareAddPeer(t, fc, wsv)

	for _, c := range []struct {
		name         string
		authorizerId string
		accountId    string
		peerId       string
		err          error
	}{
		{
			"case 1 : no error",
			authorizerId,
			"account1@com",
			"peer1@com",
			nil,
		},
		{
			"case 1 : no error",
			authorizerId,
			"account1@com",
			"peer1@com",
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			cmd := fc.NewTxBuilder().
				Consign(c.authorizerId, c.accountId, c.peerId).
				Build().
				GetPayload().
				GetCommands()[0]
			err := ex.Consign(wsv, cmd)
			if c.err != nil {
				assert.EqualError(t, errors.Cause(err), c.err.Error())
			} else {
				assert.NoError(t, err)

				ac := fc.NewEmptyAccount()
				err = wsv.Query(model.MustAddress(model.MustAddress(c.accountId).AccountId()), ac)
				require.NoError(t, err)
				assert.Equal(t, ac.GetDelegatePeerId(), c.peerId)
			}
		})
	}
	require.NoError(t, dtx.Commit())
}

/*
CreateStorage(ObjectFinder, Command) error
UpdateObject(ObjectFinder, Command) error
AddObject(ObjectFinder, Command) error
TransferObject(ObjectFinder, Command) error
AddPeer(ObjectFinder, Command) error
Consign(ObjectFinder, Command) error
*/
