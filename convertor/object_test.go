package convertor_test

import (
	. "github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount_GetFromKey(t *testing.T) {
	keys := []model.PublicKey{
		RandomPublicKey(),
		RandomPublicKey(),
	}
	exKeys := PublicKeysToListObject(keys, RandomCryptor())
	a := RandomFactory().NewAccountBuilder().
		AccountId("account@domain").
		AccountName("account").
		PublicKeys(keys).
		Balance(111).
		Quorum(1).
		DelegatePeerId("peer@domain").
		Build()

	assert.Equal(t, "account@domain", a.GetFromKey("account_id").GetAddress())
	assert.Equal(t, "account@domain", a.GetFromKey("id").GetAddress())

	assert.Equal(t, "account", a.GetFromKey("account_name").GetStr())
	assert.Equal(t, "account", a.GetFromKey("name").GetStr())

	assert.Equal(t, exKeys.Hash(), a.GetFromKey("keys").Hash())
	assert.Equal(t, exKeys.Hash(), a.GetFromKey("public_keys").Hash())

	assert.Equal(t, int64(111), a.GetFromKey("balance").GetI64())

	assert.Equal(t, int32(1), a.GetFromKey("quorum").GetI32())

	assert.Equal(t, "peer@domain", a.GetFromKey("peer_id").GetAddress())
	assert.Equal(t, "peer@domain", a.GetFromKey("delegate_peer_id").GetAddress())
}

func TestPeer_GetFromKey(t *testing.T) {
	key := RandomPublicKey()
	p := RandomFactory().NewPeer("peer@domain", "1.1.1.1:0000", key)

	assert.Equal(t, "peer@domain", p.GetFromKey("id").GetAddress())
	assert.Equal(t, "peer@domain", p.GetFromKey("peer_id").GetAddress())

	assert.Equal(t, "1.1.1.1:0000", p.GetFromKey("address").GetStr())
	assert.Equal(t, "1.1.1.1:0000", p.GetFromKey("ip").GetStr())

	assert.Equal(t, key, model.PublicKey(p.GetFromKey("public_key").GetData()))
	assert.Equal(t, key, model.PublicKey(p.GetFromKey("key").GetData()))
}

func TestStorage_GetFromKey(t *testing.T) {
	st := RandomFactory().
		NewStorageBuilder().
		Address("address", "account@com").
		Str("str", "mojimoji").
		Int64("int64", 111).
		Build()

	assert.Equal(t, "account@com", st.GetFromKey("address").GetAddress())

	assert.Equal(t, "mojimoji", st.GetFromKey("str").GetStr())

	assert.Equal(t, int64(111), st.GetFromKey("int64").GetI64())
}

func TestMapConvertor(t *testing.T) {
	binary := ConvertYamlFileToProtoBinary(t, RandomConfig().Prosl.Incentive.Path)
	for i := 0; i < 10; i++ {
		binary2 := ConvertYamlFileToProtoBinary(t, RandomConfig().Prosl.Incentive.Path)
		assert.Equal(t, binary, binary2)
	}
}
