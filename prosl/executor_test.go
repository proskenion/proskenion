package prosl_test

import (
	"encoding/hex"
	"fmt"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/prosl"
	"github.com/proskenion/proskenion/proto"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

const genesisRootId = "root@com"

func Initalize() (core.Repository, model.ModelFactory, *config.Config) {
	dba := RandomDBA()
	cryptor := RandomCryptor()
	fc := RandomFactory()
	conf := RandomConfig()
	rp := repository.NewRepository(dba, cryptor, fc, conf)
	return rp, fc, conf
}

func NewQuerycutor(qp core.QueryProcessor, qv core.QueryValidator, qc core.QueryVerifier) core.Querycutor {
	return struct {
		core.QueryProcessor
		core.QueryValidator
		core.QueryVerifier
	}{qp, qv, qc}
}

var (
	authorizer AccountWithPri
	peer       PeerWithPri
	acs        []AccountWithPri
)

func FileToBianry(t *testing.T, filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	ret, err := hex.DecodeString(string(data))
	require.NoError(t, err)
	return ret
}

type KeyPair struct {
	model.PublicKey
	model.PrivateKey
}

func InitializeObjects(t *testing.T) {
	keypairs := make([]*KeyPair, 0)
	for _, name := range []string{
		"key_auth", "key_peer", "key1", "key2", "key3",
	} {
		pub := FileToBianry(t, "./test_yaml/"+name+".pub")
		pri := FileToBianry(t, "./test_yaml/"+name+".pri")
		keypairs = append(keypairs, &KeyPair{pub, pri})
	}
	authorizer = AccountWithPri{
		"authorizer@com",
		keypairs[0].PublicKey,
		keypairs[0].PrivateKey,
	}
	peer = PeerWithPri{
		RandomFactory().NewPeer("root@peer", "127.0.0.1:50055", keypairs[1].PublicKey),
		keypairs[1].PrivateKey,
	}
	acs = []AccountWithPri{
		{
			"account1@com",
			keypairs[2].PublicKey,
			keypairs[2].PrivateKey,
		},
		{
			"account2@com",
			keypairs[3].PublicKey,
			keypairs[3].PrivateKey,
		},
		{
			"account3@com",
			keypairs[4].PublicKey,
			keypairs[4].PrivateKey,
		},
	}
}

func testConvertProsl(t *testing.T, filename string) *proskenion.Prosl {
	buf, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	prosl, err := ConvertYamlToProtobuf(buf)
	require.NoError(t, err)
	return prosl
}

func testGenesisExecuteProsl(t *testing.T, filename string, fc model.ModelFactory, rp core.Repository, conf *config.Config) {
	top, _ := rp.Top()
	wsv, err := rp.TopWSV()
	require.NoError(t, err)

	value := InitProslStateValue(fc, wsv, top, RandomCryptor(), conf)

	prosl := testConvertProsl(t, filename)
	state := ExecuteProsl(prosl, value)
	require.NoError(t, state.Err)
	require.NotNil(t, state.ReturnObject)

	expB := RandomFactory().NewTxBuilder().
		AddPeer(genesisRootId,
			peer.GetPeerId(),
			peer.GetAddress(),
			peer.GetPublicKey()).
		CreateAccount(genesisRootId, authorizer.AccountId, []model.PublicKey{authorizer.Pubkey}, 1)
	for _, ac := range acs {
		expB = expB.CreateAccount(genesisRootId, ac.AccountId, []model.PublicKey{ac.Pubkey}, 1)
	}
	for i, ac := range acs {
		expB = expB.AddBalance(genesisRootId, ac.AccountId, int64(10000*(i+1)))
	}
	expTx := expB.DefineStorage(genesisRootId, "/degraders",
		RandomFactory().NewStorageBuilder().List("acs", make([]model.Object, 0)).Build()).
		CreateStorage(genesisRootId, "root@com/degraders").
		Build()

	actualTx := state.ReturnObject.GetTransaction()
	assert.Equal(t, expTx.Hash(), actualTx.Hash())

	txList := EmptyTxList()
	require.NoError(t, txList.Push(actualTx))
	CommitTxWrapBlock(t, rp, state.Fc, actualTx)
}

func testGetAccountsExecuteProsl(t *testing.T, filename string, fc model.ModelFactory, rp core.Repository, conf *config.Config) {
	top, _ := rp.Top()
	wsv, err := rp.TopWSV()
	require.NoError(t, err)

	value := InitProslStateValue(fc, wsv, top, RandomCryptor(), conf)

	prosl := testConvertProsl(t, filename)
	state := ExecuteProsl(prosl, value)
	require.NoError(t, state.Err)
	require.NotNil(t, state.ReturnObject)

	exIds := make([]string, 3)
	for i, ac := range acs {
		exIds[2-i] = ac.AccountId
	}
	actualList := state.ReturnObject.GetList()
	for i, id := range exIds {
		assert.Equal(t, id, actualList[i].GetAccount().GetAccountId())
	}
}

func accountsToObjectList(accounts []model.Account) model.Object {
	obs := make([]model.Object, 0)
	for _, ac := range accounts {
		obs = append(obs, RandomFactory().NewObjectBuilder().Account(ac))
	}
	return RandomFactory().NewObjectBuilder().List(obs)
}

func testIncentiveExecuteProsl(t *testing.T, filename string, fc model.ModelFactory, rp core.Repository, conf *config.Config, expTx model.Transaction) {
	top, _ := rp.Top()
	wsv, err := rp.TopWSV()
	require.NoError(t, err)

	value := InitProslStateValue(fc, wsv, top, RandomCryptor(), conf)

	prosl := testConvertProsl(t, filename)
	state := ExecuteProsl(prosl, value)
	require.NoError(t, state.Err)

	actualTx := state.ReturnObject.GetTransaction()
	require.NotNil(t, actualTx)

	assert.Equal(t, expTx.Hash(), actualTx.Hash())
	fmt.Println(expTx)
	fmt.Println(actualTx)
	CommitTxWrapBlock(t, rp, state.Fc, actualTx)
}

func TestExecuteProsl(t *testing.T) {
	rp, fc, conf := Initalize()
	InitializeObjects(t)

	testGenesisExecuteProsl(t, "./test_yaml/genesis.yaml", fc, rp, conf)

	testGetAccountsExecuteProsl(t, "./test_yaml/test_1.yaml", fc, rp, conf)

	expTx := fc.NewTxBuilder().
		UpdateObject(genesisRootId, "root@com/degraders", "acs",
			accountsToObjectList(
				[]model.Account{
					fc.NewAccount(acs[2].AccountId, model.MustAddress(acs[2].AccountId).Account(), []model.PublicKey{acs[2].Pubkey}, 1, 30000, ""),
					fc.NewAccount(acs[1].AccountId, model.MustAddress(acs[1].AccountId).Account(), []model.PublicKey{acs[1].Pubkey}, 1, 20000, ""),
				})).
		AddBalance(genesisRootId, acs[2].AccountId, 10000).
		Build()
	testIncentiveExecuteProsl(t, "./test_yaml/test_2.yaml", fc, rp, conf, expTx)

	expTx = fc.NewTxBuilder().
		AddBalance(genesisRootId, acs[1].AccountId, 10000).
		Build()
	testIncentiveExecuteProsl(t, "./test_yaml/test_2.yaml", fc, rp, conf, expTx)
}
