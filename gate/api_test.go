package gate_test
/*
import (
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/require"
	"testing"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

func genesisCommit(t *testing.T, rp core.Repository) {
	txList := repository.NewTxList(RandomCryptor())
	require.NoError(t, txList.Push(
		NewTestFactory().NewTxBuilder().
			CreateAccount("root", "authorizer@com").
			CreateAccount("root", "target1@com").
			CreateAccount("root", "target2@com").
			Build()))
	require.NoError(t, rp.GenesisCommit(txList))
}

func createAccount(t *testing.T, authorizer string, target string, pub model.PublicKey, pri model.PrivateKey) model.Transaction {
	tx := NewTestFactory().NewTxBuilder().
		CreateAccount(authorizer, target).
		Build()
	require.NoError(t,tx.Sign(pub, pri))
	return tx
}
/*
func TestAPIGate_WriteAndRead(t *testing.T) {
	fc := NewTestFactory()
	rp := repository.NewRepository(RandomDBA(), RandomCryptor(), fc)
	queue := RandomQueue()
	logger := log15.New(context.TODO())
	qp := query.NewQueryProcessor(rp, fc)
	api := NewAPIGate(queue, logger, qp)

	// genesis Commit
	genesisCommit(t, rp)

	txs := []model.Transaction {
		createAccount(t, "authorizer@com", "target3@com")
	}
}
*/