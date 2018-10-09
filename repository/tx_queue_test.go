package repository_test

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	. "github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func testProposalTxQueue(t *testing.T, queue core.ProposalTxQueue) {
	t.Run("Success, random valid tx push", func(t *testing.T) {
		// Invalid Tx no problem
		txs := []model.Transaction{
			RandomTx(),
			RandomTx(),
			RandomTx(),
			RandomTx(),
			RandomTx(),
			RandomTx(),
		}

		for _, tx := range txs {
			err := queue.Push(tx)
			assert.NoError(t, err)
		}

		for _, tx := range txs {
			front, ok := queue.Pop()
			assert.True(t, ok)
			assert.Equal(t, tx, front)
		}
	})

	t.Run("Failed, nil tx push", func(t *testing.T) {
		err := queue.Push(nil)
		assert.EqualError(t, errors.Cause(err), model.ErrInvalidTransaction.Error())
	})

	t.Run("Empty pop, return nil", func(t *testing.T) {
		front, ok := queue.Pop()
		assert.False(t, ok)
		assert.Equal(t, nil, front)
	})

	t.Run("Erase test", func(t *testing.T) {
		txs := []model.Transaction{
			RandomTx(), RandomTx(), RandomTx(), RandomTx(), RandomTx(),
		}
		for _, tx := range txs {
			err := queue.Push(tx)
			require.NoError(t, err)
		}
		queue.Erase(MustHash(txs[0]))
		queue.Erase(MustHash(txs[2]))
		queue.Erase(MustHash(txs[4]))

		tx1, ok := queue.Pop()
		assert.True(t, ok)
		assert.Equal(t, tx1, txs[1])

		tx2, ok := queue.Pop()
		assert.True(t, ok)
		assert.Equal(t, tx2, txs[3])

		_, ok = queue.Pop()
		assert.False(t, ok)
	})

	t.Run("Failed, over limits random valid tx push", func(t *testing.T) {
		limit := NewTestConfig().ProposalTxsLimits
		txs := make([]model.Transaction, 0, limit)
		for i := 0; i < limit; i++ {
			txs = append(txs, RandomTx())
			require.NoError(t, queue.Push(txs[i]))
		}

		err := queue.Push(RandomTx())
		assert.EqualError(t, errors.Cause(err), ErrProposalTxQueueLimits.Error())

		for i := 0; i < limit; i++ {
			front, ok := queue.Pop()
			assert.True(t, ok)
			assert.Equal(t, txs[i], front)
		}

		_, ok := queue.Pop()
		assert.False(t, ok)
	})

	t.Run("Failed, alrady exist tx", func(t *testing.T) {
		tx := RandomTx()
		err := queue.Push(tx)
		require.NoError(t, err)

		err = queue.Push(tx)
		assert.EqualError(t, errors.Cause(err), ErrProposalTxQueueAlreadyExistTx.Error())

		exTx, ok := queue.Pop()
		require.True(t, ok)
		assert.Equal(t, tx, exTx)
	})

	t.Run("Failed, over limits random tx push and erase", func(t *testing.T) {
		limit := NewTestConfig().ProposalTxsLimits
		txs := make([]model.Transaction, 0, limit)
		txs2 := make([]model.Transaction, 0, limit/2)
		for i := 0; i < limit; i++ {
			txs = append(txs, RandomTx())
			require.NoError(t, queue.Push(txs[i]))
		}

		err := queue.Push(RandomTx())
		assert.EqualError(t, errors.Cause(err), ErrProposalTxQueueLimits.Error())

		for i := 0; i < limit; i += 2 {
			txs2 = append(txs2, RandomTx())
			assert.NoError(t, queue.Erase(MustHash(txs[i])))
			require.NoError(t, queue.Push(txs2[i/2]))
		}

		for i := 0; i < limit; i++ {
			front, ok := queue.Pop()
			assert.True(t, ok)
			if i < limit/2 {
				assert.Equal(t, txs[i*2+1], front)
			} else {
				assert.Equal(t, txs2[i-limit/2], front)
			}
		}

		_, ok := queue.Pop()
		assert.False(t, ok)
	})

	t.Run("Failed, nothing hash erase", func(t *testing.T) {
		tx := RandomTx()
		require.NoError(t, queue.Push(tx))

		tx2 := RandomTx()
		err := queue.Erase(MustHash(tx2))
		assert.EqualError(t, ErrProposalTxQueueEraseUnexistTx, errors.Cause(err).Error())

		require.NoError(t, queue.Push(tx2))
		require.NoError(t, queue.Erase(MustHash(tx2)))

		exTx, ok := queue.Pop()
		require.True(t, ok)
		assert.Equal(t, tx, exTx)

		err = queue.Erase(MustHash(tx))
		assert.EqualError(t, ErrProposalTxQueueEraseUnexistTx, errors.Cause(err).Error())

		_, ok = queue.Pop()
		assert.False(t, ok)
	})

}

func TestProposalTxQueueOnMemory(t *testing.T) {
	queue := NewProposalTxQueueOnMemory(NewTestConfig())
	testProposalTxQueue(t, queue)
}
