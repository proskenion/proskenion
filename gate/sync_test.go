package gate_test

import (
	"fmt"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/gate"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Top(rp core.Repository) model.Block {
	b, _ := rp.Top()
	return b
}

func TestNewSyncGate(t *testing.T) {
	fc, _, _, c, rp, _, conf := NewTestFactories()
	sg := NewSyncGate(rp, fc, c, RandomLogger(), conf)

	require.NoError(t, rp.GenesisCommit(RandomGenesisTxList(t)))
	for i := 0; i < conf.Sync.Limits*2; i++ {
		RandomCommitableBlockAndTxList(t, rp)
	}

	newRp := RandomRepository()
	require.NoError(t, newRp.GenesisCommit(RandomGenesisTxList(t)))
	var newBlock model.Block
	var newTxList core.TxList
	for i := 0; i < conf.Sync.Limits*2; i += conf.Sync.Limits {
		blockHash := MusTop(newRp).Hash()
		blockChan := make(chan model.Block)
		txListChan := make(chan core.TxList)
		errChan := make(chan error)
		go func(blockHash model.Hash) {
			defer close(blockChan)
			defer close(txListChan)
			defer close(errChan)
			err := sg.Sync(blockHash, blockChan, txListChan)
			require.NoError(t, err)
			errChan <- err
		}(blockHash)
		for {
			select {
			case newBlock = <-blockChan:
				fmt.Println("innewBlock")
			case newTxList = <-txListChan:
				fmt.Println("innewTxList")
				fmt.Println(newBlock, newTxList)
				err := newRp.Commit(newBlock, newTxList)
				require.NoError(t, err)
			case err := <-errChan:
				require.NoError(t, err)
				goto afterFor
			}
		}
	afterFor:
	}
	newTop, ok := newRp.Top()
	require.True(t, ok)
	assert.Equal(t, newTop.Hash(), MusTop(rp).Hash())
}
