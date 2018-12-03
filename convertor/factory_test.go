package convertor_test

import (
	"github.com/proskenion/proskenion/core/model"
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
			block := NewTestFactory().
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
			sig := NewTestFactory().NewSignature(c.expectedPub, c.expectedSig)
			assert.Equal(t, c.expectedPub, sig.GetPublicKey())
			assert.Equal(t, c.expectedSig, sig.GetSignature())
		})
	}
}

func TestTxModelBuilder(t *testing.T) {
	t.Run("case transfer", func(t *testing.T) {
		txBuilder := NewTestFactory().NewTxBuilder()
		tx := txBuilder.CreatedTime(10).
			Transfer("a", "b", 10).
			CreateAccount("x", "y").
			AddAsset("w", 10).
			AddPublicKey("auth", "ac", []byte{1, 2, 3}).
			Build()
		assert.Equal(t, int64(10), tx.GetPayload().GetCreatedTime())

		assert.Equal(t, "a", tx.GetPayload().GetCommands()[0].GetAuthorizerId())
		assert.Equal(t, "a", tx.GetPayload().GetCommands()[0].GetTargetId())
		assert.Equal(t, "b", tx.GetPayload().GetCommands()[0].GetTransfer().GetDestAccountId())
		assert.Equal(t, int64(10), tx.GetPayload().GetCommands()[0].GetTransfer().GetBalance())

		assert.Equal(t, "x", tx.GetPayload().GetCommands()[1].GetAuthorizerId())
		assert.Equal(t, "y", tx.GetPayload().GetCommands()[1].GetTargetId())

		assert.Equal(t, "w", tx.GetPayload().GetCommands()[2].GetAuthorizerId())
		assert.Equal(t, "w", tx.GetPayload().GetCommands()[2].GetTargetId())
		assert.Equal(t, int64(10), tx.GetPayload().GetCommands()[2].GetAddAsset().GetBalance())

		assert.Equal(t, "auth", tx.GetPayload().GetCommands()[3].GetAuthorizerId())
		assert.Equal(t, "ac", tx.GetPayload().GetCommands()[3].GetTargetId())
		assert.Equal(t, []byte{1, 2, 3}, tx.GetPayload().GetCommands()[3].GetAddPublicKey().GetPublicKey())
	})
}

func TestNewAccount(t *testing.T) {
	for _, c := range []struct {
		name        string
		accountId   string
		accountName string
		pubkeys     []model.PublicKey
		amount      int64
	}{
		{
			"case 1",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{RandomByte()},
			rand.Int63(),
		}, {
			"case 2",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{RandomByte()},
			rand.Int63(),
		}, {
			"case 3",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{RandomByte(), RandomByte(), RandomByte(), RandomByte()},
			rand.Int63(),
		}, {
			"case 4",
			RandomStr(),
			RandomStr(),
			[]model.PublicKey{},
			rand.Int63(),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ac := NewTestFactory().NewAccount(c.accountId, c.accountName, c.pubkeys, c.amount)
			assert.Equal(t, c.accountId, ac.GetAccountId())
			assert.Equal(t, c.accountName, ac.GetAccountName())
			assert.Equal(t, c.pubkeys, ac.GetPublicKeys())
			assert.Equal(t, c.amount, ac.GetBalance())
		})
	}
}

func TestNewPeer(t *testing.T) {
	for _, c := range []struct {
		name    string
		address string
		pubkey  model.PublicKey
	}{
		{
			"case 1",
			"111.111.111.111",
			RandomByte(),
		},
		{
			"case 2",
			RandomStr(),
			RandomByte(),
		},
		{
			"case 3",
			"localhost",
			nil,
		},
		{
			"case 4",
			"",
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			peer := NewTestFactory().NewPeer(c.address, c.pubkey)
			assert.Equal(t, c.address, peer.GetAddress())
			assert.Equal(t, c.pubkey, peer.GetPublicKey())
		})
	}
}

func TestModelFactory_NewQueryBuilder(t *testing.T) {
	t.Run("case 1 account query", func(t *testing.T) {
		builder := NewTestFactory().NewQueryBuilder()
		query := builder.CreatedTime(1).
			TargetId("a").
			AuthorizerId("b").
			RequestCode(model.AccountObjectCode).
			Build()
		assert.Equal(t, int64(1), query.GetPayload().GetCreatedTime())
		assert.Equal(t, "a", query.GetPayload().GetTargetId())
		assert.Equal(t, "b", query.GetPayload().GetAuthorizerId())
		assert.Equal(t, model.ObjectCode(model.AccountObjectCode), query.GetPayload().GetRequestCode())
	})
}

func TestModelFactory_NewQueryResponseBuilder(t *testing.T) {
	t.Run("case 1 account query", func(t *testing.T) {
		expAc := NewTestFactory().NewAccount(RandomStr(), RandomStr(), []model.PublicKey{RandomByte(), RandomByte()}, rand.Int63())
		builder := NewTestFactory().NewQueryResponseBuilder()
		res := builder.Account(expAc).Build()
		actAc := res.GetPayload().GetAccount()
		assert.Equal(t, expAc.GetAccountId(), actAc.GetAccountId())
		assert.Equal(t, expAc.GetAccountName(), actAc.GetAccountName())
		assert.Equal(t, expAc.GetPublicKeys(), actAc.GetPublicKeys())
		assert.Equal(t, expAc.GetBalance(), actAc.GetBalance())
	})

	t.Run("case 2 peer query", func(t *testing.T) {
		pub, _ := RandomCryptor().NewKeyPairs()
		expPeer := NewTestFactory().NewPeer("address:50051", pub)
		res := NewTestFactory().NewQueryResponseBuilder().
			Peer(expPeer).
			Build()
		actPeer := res.GetPayload().GetPeer()
		assert.Equal(t, expPeer.GetPublicKey(), actPeer.GetPublicKey())
		assert.Equal(t, expPeer.GetAddress(), actPeer.GetAddress())
	})

}

func NewObjectFactory_NewStorageBuilder(t *testing.T) {
	fc := NewTestFactory()
	t.Run("case 1 storage builder", func(t *testing.T) {
		builder := fc.NewStorageBuilder()
		storage := builder.Dict("dict", map[string]model.Object{"key": fc.NewEmptyObject()}).
			List("list", []model.Object{fc.NewEmptyObject(), fc.NewEmptyObject()}).
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
		assert.Equal(t, map[string]model.Object{"key": fc.NewEmptyObject()}, dict["dict"].GetDict())
		assert.Equal(t, model.ObjectCode(model.DictObjectCode), dict["dict"].GetType())

		assert.Equal(t, []model.Object{fc.NewEmptyObject(), fc.NewEmptyObject()}, dict["list"].GetList())
		assert.Equal(t, model.ObjectCode(model.ListObjectCode), dict["list"].GetType())

		assert.Equal(t, fc.NewEmptyAccount(), dict["account"].GetAccount())
		assert.Equal(t, model.ObjectCode(model.AccountObjectCode), dict["account"].GetType())

		assert.Equal(t, fc.NewEmptySignature(), dict["sig"].GetSig())
		assert.Equal(t, model.ObjectCode(model.SignatureObjectCode), dict["sig"].GetType())

		assert.Equal(t, "target@account.com", dict["address"].GetAddress())
		assert.Equal(t, model.ObjectCode(model.AddressObjectCode), dict["address"].GetType())

		assert.Equal(t, []byte("aaaa"), dict["data"].GetAddress())
		assert.Equal(t, model.ObjectCode(model.BytesObjectCode), dict["data"].GetType())

		assert.Equal(t, "str", dict["str"].GetStr())
		assert.Equal(t, model.ObjectCode(model.StringObjectCode), dict["str"].GetType())

		assert.Equal(t, fc.NewEmptyPeer(), dict["peer"].GetPeer())
		assert.Equal(t, model.ObjectCode(model.PeerObjectCode), dict["peer"].GetType())

		assert.Equal(t, 32, dict["int32"].GetI32())
		assert.Equal(t, model.ObjectCode(model.Int32ObjectCode), dict["int32"].GetType())

		assert.Equal(t, 64, dict["int64"].GetI64())
		assert.Equal(t, model.ObjectCode(model.Int64ObjectCode), dict["int64"].GetType())

		assert.Equal(t, 1, dict["uint64"].GetU32())
		assert.Equal(t, model.ObjectCode(model.Uint32ObjectCode), dict["uint32"].GetType())

		assert.Equal(t, 2, dict["uint64"].GetU64())
		assert.Equal(t, model.ObjectCode(model.Uint64ObjectCode), dict["uint64"].GetType())
	})
}
