package gate_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	. "github.com/proskenion/proskenion/gate"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func newRandomConsensusGate() core.ConsensusGate {
	fc, _, _, c, _, _, conf := NewTestFactories()
	return NewConsensusGate(fc, c,
		repository.NewProposalTxQueueOnMemory(conf), repository.NewTxListCache(conf),
		repository.NewProposalBlockQueueOnMemory(conf), conf)
}

func TestConsensus_PropagateTx(t *testing.T) {
	cg := newRandomConsensusGate()
	t.Run("case 1 : correct", func(t *testing.T) {
		tx := RandomSignedTx(t)
		assert.NoError(t, cg.PropagateTx(tx))
	})

	t.Run("case 2 : err verify", func(t *testing.T) {
		tx := RandomTx()
		err := cg.PropagateTx(tx)
		assert.EqualError(t, errors.Cause(err), core.ErrConsensusGatePropagateTxVerifyError.Error())
	})
}

func TestConsensus_PropagateBlockAck(t *testing.T) {
	cg := newRandomConsensusGate()
	t.Run("case 1 : correct", func(t *testing.T) {
		block := RandomSignedBlock(t)
		sig, err := cg.PropagateBlockAck(block)
		assert.NoError(t, err)
		assert.NoError(t, RandomCryptor().Verify(sig.GetPublicKey(), block, sig.GetSignature()))
	})

	t.Run("case 2 : no error", func(t *testing.T) {
		_, err := cg.PropagateBlockAck(RandomBlock())
		assert.Equal(t, errors.Cause(err), core.ErrConsensusGatePropagateBlockVerifyError)
	})
}

func TestConsensus_PropagateBlockStreamTx(t *testing.T) {
	cg := newRandomConsensusGate()
	t.Run("case 1 : correct", func(t *testing.T) {
		txList := RandomTxList()
		block := RandomFactory().NewBlockBuilder().TxsHash(txList.Hash()).Build()

		fmt.Println(txList.Size())
		txChan := make(chan model.Transaction)
		errChan := make(chan error)

		go func(txList core.TxList) {
			defer close(txChan)
			defer close(errChan)

			for _, tx := range txList.List() {
				txChan <- tx
			}
			errChan <- io.EOF
		}(txList)

		err := cg.PropagateBlockStreamTx(block, txChan, errChan)
		require.NoError(t, err)
	})

	t.Run("case 2 : different error", func(t *testing.T) {
		txList := RandomTxList()
		block := RandomBlock()

		txChan := make(chan model.Transaction)
		errChan := make(chan error)

		go func(txList core.TxList) {
			defer close(txChan)
			defer close(errChan)
			for _, tx := range txList.List() {
				txChan <- tx
			}
			errChan <- nil
		}(txList)

		err := cg.PropagateBlockStreamTx(block, txChan, errChan)
		assert.Error(t, errors.Cause(err), core.ErrConsensusGatePropagateBlockDifferentHash)
	})

	t.Run("case 2 : error in the middle", func(t *testing.T) {
		txList := RandomTxList()
		block := RandomBlock()

		txChan := make(chan model.Transaction)
		errChan := make(chan error)

		go func(txList core.TxList) {
			defer close(txChan)
			defer close(errChan)
			for i, tx := range txList.List() {
				if i == 10 {
					errChan <- fmt.Errorf("expected error")
					return
				}
				txChan <- tx
			}
		}(txList)

		err := cg.PropagateBlockStreamTx(block, txChan, errChan)
		assert.Error(t, errors.Cause(err), "expected error")
	})

}
