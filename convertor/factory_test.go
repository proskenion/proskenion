package convertor_test

import (
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestBlockFactory(t *testing.T) {
	for _, c := range []struct {
		name                  string
		expectedHeight        int64
		expectedPreBlockHash  model.Hash
		expectedCreatedTime   int64
		expectedWSVHash       model.Hash
		expectedTxHistoryHash model.Hash
		expectedTxsHash       model.Hash
		expectedRound         int32
	}{
		{
			"case 1",
			10,
			model.Hash("preBlockHash"),
			5,
			model.Hash("WSVHash"),
			model.Hash("TxHistoryHash"),
			model.Hash("txHash"),
			1,
		},
		{
			"case 2",
			999999999999,
			model.Hash("preBlockHash"),
			0,
			model.Hash("WSVHash"),
			model.Hash("TxHistoryHash"),
			model.Hash("txHash"),
			1,
		},
		{
			"hash nil case no problem",
			0,
			nil,
			999999999999,
			nil,
			nil,
			nil,
			0,
		},
		{
			"minus number is no problem case",
			-1,
			nil,
			-1,
			nil,
			nil,
			nil,
			-1111111,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			block := RandomFactory().
				NewBlockBuilder().
				Height(c.expectedHeight).
				PreBlockHash(c.expectedPreBlockHash).
				CreatedTime(c.expectedCreatedTime).
				WSVHash(c.expectedWSVHash).
				TxHistoryHash(c.expectedTxHistoryHash).
				TxsHash(c.expectedTxsHash).
				Round(c.expectedRound).
				Build()
			assert.Equal(t, c.expectedHeight, block.GetPayload().GetHeight())
			assert.Equal(t, c.expectedPreBlockHash, block.GetPayload().GetPreBlockHash())
			assert.Equal(t, c.expectedCreatedTime, block.GetPayload().GetCreatedTime())
			assert.Equal(t, c.expectedWSVHash, block.GetPayload().GetWSVHash())
			assert.Equal(t, c.expectedTxHistoryHash, block.GetPayload().GetTxHistoryHash())
			assert.Equal(t, c.expectedTxsHash, block.GetPayload().GetTxsHash())
			assert.Equal(t, c.expectedRound, block.GetPayload().GetRound())
		})
	}
}

func TestSignatureFactory(t *testing.T) {
	for _, c := range []struct {
		name        string
		expectedPub model.PublicKey
		expectedSig []byte
	}{
		{
			"case 1",
			RandomByte(),
			RandomByte(),
		},
		{
			"case 2",
			RandomByte(),
			RandomByte(),
		},
		{
			"case 3",
			RandomByte(),
			RandomByte(),
		},
		{
			"case 4",
			RandomByte(),
			RandomByte(),
		},
		{
			"case 5",
			RandomByte(),
			RandomByte(),
		},
		{
			"pub nil case no problem",
			nil,
			RandomByte(),
		},
		{
			"sig nil case no problem",
			RandomByte(),
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			sig := RandomFactory().NewSignature(c.expectedPub, c.expectedSig)
			assert.Equal(t, c.expectedPub, sig.GetPublicKey())
			assert.Equal(t, c.expectedSig, sig.GetSignature())
		})
	}
}

func TestTxModelBuilder(t *testing.T) {
	t.Run("case transfer", func(t *testing.T) {
		txBuilder := RandomFactory().NewTxBuilder()
		tx := txBuilder.CreatedTime(10).
			TransferBalance("au", "a", "b", 10).                                           // [0]
			CreateAccount("x", "y", []model.PublicKey{[]byte{1, 2, 3}}, 1).                // [1]
			AddBalance("au", "w", 10).                                                     // [2]
			AddPublicKeys("auth", "ac", []model.PublicKey{[]byte{1, 2, 3}}).               // [3]
			TransferBalance("au", "src", "dest", 200).                                     // [4]
			AddPublicKeys("authorizer", "account", []model.PublicKey{[]byte{4, 5, 6}}).    // [5]
			RemovePublicKeys("authorizer", "account", []model.PublicKey{[]byte{4, 5, 6}}). // [6]
			SetQuorum("authorizer", "account", 2).                                         // [7]
			DefineStorage("authorizer", "account",
															RandomFactory().NewStorageBuilder().Int32("int", 32).Build()). // [8]
			CreateStorage("authorizer", "wallet_id").                                                   // [9]
			UpdateObject("authorizer", "wallet_id", "key", RandomFactory().NewEmptyObject()).           // [10]
			AddObject("authorizer", "wallet_id", "key", RandomFactory().NewEmptyObject()).              // [11]
			TransferObject("authorizer", "wallet_id", "dest", "key", RandomFactory().NewEmptyObject()). // [12]
			AddPeer("authorizer", "account", "localhost", model.PublicKey{2, 2, 2}).                    // [13]
			ActivatePeer("authorizer", "peer").                                                         // [14]
			SuspendPeer("authorizer", "peer").                                                          // [15]
			BanPeer("authorizer", "peer").                                                              // [16]
			Consign("authorizer", "account", "peer").                                                   // [17]
			CheckAndCommitProsl("authorizer", "a@c/p",
				map[string]model.Object{"key": RandomFactory().NewObjectBuilder().Str("yyy")}). //[18]
			Build()
		assert.Equal(t, int64(10), tx.GetPayload().GetCreatedTime())

		// transfer balance
		assert.Equal(t, "au", tx.GetPayload().GetCommands()[0].GetAuthorizerId())
		assert.Equal(t, "a", tx.GetPayload().GetCommands()[0].GetTargetId())
		assert.Equal(t, "b", tx.GetPayload().GetCommands()[0].GetTransferBalance().GetDestAccountId())
		assert.Equal(t, int64(10), tx.GetPayload().GetCommands()[0].GetTransferBalance().GetBalance())

		// create account
		assert.Equal(t, "x", tx.GetPayload().GetCommands()[1].GetAuthorizerId())
		assert.Equal(t, "y", tx.GetPayload().GetCommands()[1].GetTargetId())
		assert.Equal(t, []model.PublicKey{[]byte{1, 2, 3}}, tx.GetPayload().GetCommands()[1].GetCreateAccount().GetPublicKeys())
		assert.Equal(t, int32(1), tx.GetPayload().GetCommands()[1].GetCreateAccount().GetQuorum())

		// add balance
		assert.Equal(t, "au", tx.GetPayload().GetCommands()[2].GetAuthorizerId())
		assert.Equal(t, "w", tx.GetPayload().GetCommands()[2].GetTargetId())
		assert.Equal(t, int64(10), tx.GetPayload().GetCommands()[2].GetAddBalance().GetBalance())

		// add public keys
		assert.Equal(t, "auth", tx.GetPayload().GetCommands()[3].GetAuthorizerId())
		assert.Equal(t, "ac", tx.GetPayload().GetCommands()[3].GetTargetId())
		assert.Equal(t, model.PublicKey{1, 2, 3}, tx.GetPayload().GetCommands()[3].GetAddPublicKeys().GetPublicKeys()[0])

		// transfer balance
		assert.Equal(t, "au", tx.GetPayload().GetCommands()[4].GetAuthorizerId())
		assert.Equal(t, "src", tx.GetPayload().GetCommands()[4].GetTargetId())
		assert.Equal(t, "dest", tx.GetPayload().GetCommands()[4].GetTransferBalance().GetDestAccountId())
		assert.Equal(t, int64(200), tx.GetPayload().GetCommands()[4].GetTransferBalance().GetBalance())

		// add publicKeys
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[5].GetAuthorizerId())
		assert.Equal(t, "account", tx.GetPayload().GetCommands()[5].GetTargetId())
		assert.EqualValues(t, []model.PublicKey{[]byte{4, 5, 6}}, tx.GetPayload().GetCommands()[5].GetAddPublicKeys().GetPublicKeys())

		// remove publicKeys
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[6].GetAuthorizerId())
		assert.Equal(t, "account", tx.GetPayload().GetCommands()[6].GetTargetId())
		assert.EqualValues(t, []model.PublicKey{[]byte{4, 5, 6}}, tx.GetPayload().GetCommands()[6].GetRemovePublicKeys().GetPublicKeys())

		// set quorum
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[7].GetAuthorizerId())
		assert.Equal(t, "account", tx.GetPayload().GetCommands()[7].GetTargetId())
		assert.Equal(t, int32(2), tx.GetPayload().GetCommands()[7].GetSetQuorum().GetQuorum())

		// define storage
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[8].GetAuthorizerId())
		assert.Equal(t, "account", tx.GetPayload().GetCommands()[8].GetTargetId())
		assert.Equal(t, model.Int32ObjectCode, model.ObjectCode(tx.GetPayload().GetCommands()[8].GetDefineStorage().GetStorage().GetObject()["int"].GetType()))
		assert.Equal(t, int32(32), tx.GetPayload().GetCommands()[8].GetDefineStorage().GetStorage().GetObject()["int"].GetI32())

		// create storage
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[9].GetAuthorizerId())
		assert.Equal(t, "wallet_id", tx.GetPayload().GetCommands()[9].GetTargetId())

		// update object
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[10].GetAuthorizerId())
		assert.Equal(t, "wallet_id", tx.GetPayload().GetCommands()[10].GetTargetId())
		assert.Equal(t, RandomFactory().NewEmptyObject().GetType(), tx.GetPayload().GetCommands()[10].GetUpdateObject().GetObject().GetType())

		// add object
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[11].GetAuthorizerId())
		assert.Equal(t, "wallet_id", tx.GetPayload().GetCommands()[11].GetTargetId())
		assert.Equal(t, RandomFactory().NewEmptyObject().GetType(), tx.GetPayload().GetCommands()[11].GetAddObject().GetObject().GetType())

		// transferObject
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[12].GetAuthorizerId())
		assert.Equal(t, "wallet_id", tx.GetPayload().GetCommands()[12].GetTargetId())
		assert.Equal(t, "dest", tx.GetPayload().GetCommands()[12].GetTransferObject().GetDestAccountId())
		assert.Equal(t, "key", tx.GetPayload().GetCommands()[12].GetTransferObject().GetKey())
		assert.Equal(t, RandomFactory().NewEmptyObject().GetType(), tx.GetPayload().GetCommands()[12].GetTransferObject().GetObject().GetType())

		// add peer
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[13].GetAuthorizerId())
		assert.Equal(t, "account", tx.GetPayload().GetCommands()[13].GetTargetId())
		assert.Equal(t, "localhost", tx.GetPayload().GetCommands()[13].GetAddPeer().GetAddress())
		assert.EqualValues(t, model.PublicKey{2, 2, 2}, tx.GetPayload().GetCommands()[13].GetAddPeer().GetPublicKey())

		// activate peer
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[14].GetAuthorizerId())
		assert.Equal(t, "peer", tx.GetPayload().GetCommands()[14].GetTargetId())
		assert.IsType(t, &proskenion.Command_ActivatePeer{}, tx.GetPayload().GetCommands()[14].(*convertor.Command).GetCommand())

		// suspend peer
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[15].GetAuthorizerId())
		assert.Equal(t, "peer", tx.GetPayload().GetCommands()[15].GetTargetId())
		assert.IsType(t, &proskenion.Command_SuspendPeer{}, tx.GetPayload().GetCommands()[15].(*convertor.Command).GetCommand())

		// ban peer
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[16].GetAuthorizerId())
		assert.Equal(t, "peer", tx.GetPayload().GetCommands()[16].GetTargetId())
		assert.IsType(t, &proskenion.Command_BanPeer{}, tx.GetPayload().GetCommands()[16].(*convertor.Command).GetCommand())

		// consign
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[17].GetAuthorizerId())
		assert.Equal(t, "account", tx.GetPayload().GetCommands()[17].GetTargetId())
		assert.Equal(t, "peer", tx.GetPayload().GetCommands()[17].GetConsign().GetPeerId())

		// check and commit prosl
		assert.Equal(t, "authorizer", tx.GetPayload().GetCommands()[18].GetAuthorizerId())
		assert.Equal(t, "a@c/p", tx.GetPayload().GetCommands()[18].GetTargetId())
		assert.Equal(t, "yyy", tx.GetPayload().GetCommands()[18].GetCheckAndCommitProsl().GetVariables()["key"].GetStr())
	})
}

func TestNewAccount(t *testing.T) {
	for _, c := range []struct {
		name        string
		accountId   string
		accountName string
		pubkeys     []model.PublicKey
		quorum      int32
		amount      int64
		peerId      string
	}{
		{
			"case 1",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{RandomByte()},
			rand.Int31(),
			rand.Int63(),
			RandomStr(),
		}, {
			"case 2",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{RandomByte()},
			rand.Int31(),
			rand.Int63(),
			RandomStr(),
		}, {
			"case 3",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{RandomByte(), RandomByte(), RandomByte(), RandomByte()},
			rand.Int31(),
			rand.Int63(),
			RandomStr(),
		}, {
			"case 4",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{},
			rand.Int31(),
			rand.Int63(),
			RandomStr(),
		},
		{
			"case 5",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{},
			rand.Int31(),
			rand.Int63(),
			RandomStr(),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ac := RandomFactory().NewAccount(c.accountId, c.accountName, c.pubkeys, c.quorum, c.amount, c.peerId)
			assert.Equal(t, c.accountId, ac.GetAccountId())
			assert.Equal(t, c.accountName, ac.GetAccountName())
			assert.Equal(t, c.pubkeys, ac.GetPublicKeys())
			assert.Equal(t, c.quorum, ac.GetQuorum())
			assert.Equal(t, c.amount, ac.GetBalance())
			assert.Equal(t, c.peerId, ac.GetDelegatePeerId())

			ac2 := RandomFactory().NewAccountBuilder().From(ac).Build()
			assert.Equal(t, c.accountId, ac2.GetAccountId())
			assert.Equal(t, c.accountName, ac2.GetAccountName())
			assert.Equal(t, c.pubkeys, ac2.GetPublicKeys())
			assert.Equal(t, c.quorum, ac2.GetQuorum())
			assert.Equal(t, c.amount, ac2.GetBalance())
			assert.Equal(t, c.peerId, ac2.GetDelegatePeerId())

			ac3 := RandomFactory().NewAccountBuilder().
				AccountId(c.accountId).
				AccountName(c.accountName).
				Balance(c.amount).
				PublicKeys(c.pubkeys).
				Quorum(c.quorum).
				DelegatePeerId(c.peerId).
				Build()

			assert.Equal(t, c.accountId, ac3.GetAccountId())
			assert.Equal(t, c.accountName, ac3.GetAccountName())
			assert.Equal(t, c.pubkeys, ac3.GetPublicKeys())
			assert.Equal(t, c.quorum, ac3.GetQuorum())
			assert.Equal(t, c.amount, ac3.GetBalance())
			assert.Equal(t, c.peerId, ac3.GetDelegatePeerId())
		})
	}
}

func TestNewPeer(t *testing.T) {
	for _, c := range []struct {
		name    string
		id      string
		address string
		pubkey  model.PublicKey
	}{
		{
			"case 1",
			"peer@com.pr/peer",
			"111.111.111.111",
			RandomByte(),
		},
		{
			"case 2",
			"peer@com.pr/peer",
			RandomStr(),
			RandomByte(),
		},
		{
			"case 3",
			"peer@com.pr/peer",
			"localhost",
			nil,
		},
		{
			"case 4",
			"peer@com.pr/peer",
			"",
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			peer := RandomFactory().NewPeer(c.id, c.address, c.pubkey)
			assert.Equal(t, c.address, peer.GetAddress())
			assert.Equal(t, c.pubkey, peer.GetPublicKey())
		})
	}
}

func TestModelFactory_NewQueryBuilder(t *testing.T) {
	t.Run("case 1 account query", func(t *testing.T) {
		builder := RandomFactory().NewQueryBuilder()
		query := builder.CreatedTime(1).
			FromId("a").
			AuthorizerId("b").
			Select("*").
			OrderBy("key", model.DESC).
			Where("key > 10").
			Limit(10).
			RequestCode(model.AccountObjectCode).
			Build()
		assert.Equal(t, int64(1), query.GetPayload().GetCreatedTime())
		assert.Equal(t, "a", query.GetPayload().GetFromId())
		assert.Equal(t, "b", query.GetPayload().GetAuthorizerId())
		assert.Equal(t, "*", query.GetPayload().GetSelect())
		assert.Equal(t, "key", query.GetPayload().GetOrderBy().GetKey())
		assert.Equal(t, model.OrderCode(model.DESC), query.GetPayload().GetOrderBy().GetOrder())
		assert.Equal(t, "key > 10", query.GetPayload().GetWhere())
		assert.Equal(t, int32(10), query.GetPayload().GetLimit())
		assert.Equal(t, model.AccountObjectCode, query.GetPayload().GetRequestCode())
	})
}

func TestModelFactory_NewQueryResponseBuilder(t *testing.T) {
	t.Run("case 1 account query", func(t *testing.T) {
		expAc := RandomFactory().NewAccount(RandomStr(), RandomStr(), []model.PublicKey{RandomByte(), RandomByte()}, rand.Int31(), rand.Int63(), RandomStr())
		builder := RandomFactory().NewQueryResponseBuilder()
		res := builder.Account(expAc).Build()
		actAc := res.GetObject().GetAccount()
		assert.Equal(t, expAc.GetAccountId(), actAc.GetAccountId())
		assert.Equal(t, expAc.GetAccountName(), actAc.GetAccountName())
		assert.Equal(t, expAc.GetPublicKeys(), actAc.GetPublicKeys())
		assert.Equal(t, expAc.GetBalance(), actAc.GetBalance())
	})

	t.Run("case 2 peer query", func(t *testing.T) {
		pub, _ := RandomCryptor().NewKeyPairs()
		expPeer := RandomFactory().NewPeer(RandomStr(), "address:50051", pub)
		res := RandomFactory().NewQueryResponseBuilder().
			Peer(expPeer).
			Build()
		actPeer := res.GetObject().GetPeer()
		assert.Equal(t, expPeer.GetPublicKey(), actPeer.GetPublicKey())
		assert.Equal(t, expPeer.GetAddress(), actPeer.GetAddress())
	})

}

func TestNewObjectFactory_NewObjectBuilder(t *testing.T) {
	fc := RandomFactory()
	t.Run("case 1 object builder", func(t *testing.T) {
		dict := fc.NewObjectBuilder().Dict(map[string]model.Object{"key": fc.NewEmptyObject()})
		list := fc.NewObjectBuilder().List([]model.Object{fc.NewEmptyObject(), fc.NewEmptyObject()})
		account := fc.NewObjectBuilder().Account(fc.NewEmptyAccount())
		sig := fc.NewObjectBuilder().Sig(fc.NewEmptySignature())
		address := fc.NewObjectBuilder().Address("target@account.com")
		data := fc.NewObjectBuilder().Data([]byte("aaaa"))
		str := fc.NewObjectBuilder().Str("str")
		peer := fc.NewObjectBuilder().Peer(fc.NewEmptyPeer())
		i32 := fc.NewObjectBuilder().Int32(32)
		i64 := fc.NewObjectBuilder().Int64(64)
		u32 := fc.NewObjectBuilder().Uint32(1)
		u64 := fc.NewObjectBuilder().Uint64(2)
		b := fc.NewObjectBuilder().Bool(true)

		expTx := fc.NewTxBuilder().CreateAccount("authorizer@com", "account@com", nil, 0).Build()
		cmd := fc.NewObjectBuilder().Command(expTx.GetPayload().GetCommands()[0])
		tx := fc.NewObjectBuilder().Transaction(expTx)

		storage := fc.NewObjectBuilder().Storage(RandomFactory().NewStorageBuilder().Int32("int32", 1).Build())

		assert.Equal(t, map[string]model.Object{"key": fc.NewEmptyObject()}, dict.GetDict())
		assert.Equal(t, model.DictObjectCode, dict.GetType())

		assert.Equal(t, []model.Object{fc.NewEmptyObject(), fc.NewEmptyObject()}, list.GetList())
		assert.Equal(t, model.ListObjectCode, list.GetType())

		assert.Equal(t, fc.NewEmptyAccount(), account.GetAccount())
		assert.Equal(t, model.AccountObjectCode, account.GetType())

		assert.Equal(t, fc.NewEmptySignature(), sig.GetSig())
		assert.Equal(t, model.SignatureObjectCode, sig.GetType())

		assert.Equal(t, "target@account.com", address.GetAddress())
		assert.Equal(t, model.AddressObjectCode, address.GetType())

		assert.Equal(t, []byte("aaaa"), data.GetData())
		assert.Equal(t, model.BytesObjectCode, data.GetType())

		assert.Equal(t, "str", str.GetStr())
		assert.Equal(t, model.StringObjectCode, str.GetType())

		assert.Equal(t, fc.NewEmptyPeer(), peer.GetPeer())
		assert.Equal(t, model.PeerObjectCode, peer.GetType())

		assert.Equal(t, int32(32), i32.GetI32())
		assert.Equal(t, model.Int32ObjectCode, i32.GetType())

		assert.Equal(t, int64(64), i64.GetI64())
		assert.Equal(t, model.Int64ObjectCode, i64.GetType())

		assert.Equal(t, uint32(1), u32.GetU32())
		assert.Equal(t, model.Uint32ObjectCode, u32.GetType())

		assert.Equal(t, uint64(2), u64.GetU64())
		assert.Equal(t, model.Uint64ObjectCode, u64.GetType())

		assert.Equal(t, RandomFactory().NewStorageBuilder().Int32("int32", 1).Build().Hash(), storage.GetStorage().Hash())
		assert.Equal(t, model.StorageObjectCode, storage.GetType())

		assert.Equal(t, true, b.GetBoolean())
		assert.Equal(t, model.BoolObjectCode, b.GetType())

		assert.Equal(t, model.CommandObjectCode, cmd.GetType())
		assert.Equal(t, expTx.GetPayload().GetCommands()[0].GetCreateAccount(),
			cmd.GetCommand().GetCreateAccount())

		assert.Equal(t, model.TransactionObjectCode, tx.GetType())
		assert.Equal(t, expTx.Hash(), tx.GetTransaction().Hash())
	})
}

func TestNewObjectFactory_NewStorageBuilder(t *testing.T) {
	fc := RandomFactory()
	t.Run("case 1 storage builder", func(t *testing.T) {
		builder := fc.NewStorageBuilder()
		storage := builder.Dict("dict", map[string]model.Object{"key": fc.NewObjectBuilder().Int32(1)}).
			List("list", []model.Object{fc.NewObjectBuilder().Int32(1), fc.NewObjectBuilder().Int32(2)}).
			Account("account", fc.NewEmptyAccount()).
			Sig("sig", fc.NewEmptySignature()).
			Address("address", "target@account.com").
			Data("data", []byte("aaaa")).
			Str("str", "str").
			Peer("peer", fc.NewEmptyPeer()).
			Int32("int32", 32).
			Int64("int64", 64).
			Uint32("uint32", 1).
			Uint64("uint64", 2).
			Build()

		dict := storage.GetObject()
		assert.Equal(t, fc.NewObjectBuilder().Int32(1).Hash(), dict["dict"].GetDict()["key"].Hash())
		assert.Equal(t, model.DictObjectCode, dict["dict"].GetType())

		for i, o := range []model.Object{fc.NewObjectBuilder().Int32(1), fc.NewObjectBuilder().Int32(2)} {
			assert.Equal(t, o, dict["list"].GetList()[i])
		}
		assert.Equal(t, model.ListObjectCode, dict["list"].GetType())

		assert.Equal(t, fc.NewEmptyAccount(), dict["account"].GetAccount())
		assert.Equal(t, model.AccountObjectCode, dict["account"].GetType())

		assert.Equal(t, fc.NewEmptySignature(), dict["sig"].GetSig())
		assert.Equal(t, model.SignatureObjectCode, dict["sig"].GetType())

		assert.Equal(t, "target@account.com", dict["address"].GetAddress())
		assert.Equal(t, model.AddressObjectCode, dict["address"].GetType())

		assert.Equal(t, []byte("aaaa"), dict["data"].GetData())
		assert.Equal(t, model.BytesObjectCode, dict["data"].GetType())

		assert.Equal(t, "str", dict["str"].GetStr())
		assert.Equal(t, model.StringObjectCode, dict["str"].GetType())

		assert.Equal(t, fc.NewEmptyPeer(), dict["peer"].GetPeer())
		assert.Equal(t, model.PeerObjectCode, dict["peer"].GetType())

		assert.Equal(t, int32(32), dict["int32"].GetI32())
		assert.Equal(t, model.Int32ObjectCode, dict["int32"].GetType())

		assert.Equal(t, int64(64), dict["int64"].GetI64())
		assert.Equal(t, model.Int64ObjectCode, dict["int64"].GetType())

		assert.Equal(t, uint32(1), dict["uint32"].GetU32())
		assert.Equal(t, model.Uint32ObjectCode, dict["uint32"].GetType())

		assert.Equal(t, uint64(2), dict["uint64"].GetU64())
		assert.Equal(t, model.Uint64ObjectCode, dict["uint64"].GetType())
	})
}
