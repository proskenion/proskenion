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

func testProposalBlockQueue(t *testing.T, queue core.ProposalBlockQueue) {
	t.Run("Success, random valid block push", func(t *testing.T) {
		// Invalid Block no problem
		blocks := []model.Block{
			RandomBlock(),
			RandomBlock(),
			RandomBlock(),
			RandomBlock(),
			RandomBlock(),
			RandomBlock(),
		}

		for _, block := range blocks {
			err := queue.Push(block)
			assert.NoError(t, err)
		}

		for _, block := range blocks {
			front, ok := queue.Pop()
			assert.True(t, ok)
			assert.Equal(t, block, front)
		}
	})

	t.Run("Failed, nil block push", func(t *testing.T) {
		err := queue.Push(nil)
		assert.EqualError(t, errors.Cause(err), core.ErrProposalQueuePushNil.Error())
	})

	t.Run("Empty pop, return nil", func(t *testing.T) {
		front, ok := queue.Pop()
		assert.False(t, ok)
		assert.Equal(t, nil, front)
	})

	NumWaitChan := func(num int) {
		for i := 0; i < num; i++ {
			queue.WaitPush()
		}
		assert.True(t, true)
	}

	t.Run("Erase test", func(t *testing.T) {
		blocks := []model.Block{
			RandomBlock(), RandomBlock(), RandomBlock(), RandomBlock(), RandomBlock(),
		}
		go NumWaitChan(len(blocks))
		for _, block := range blocks {
			err := queue.Push(block)
			require.NoError(t, err)
		}
		require.NoError(t, queue.Erase(blocks[0].Hash()))
		require.NoError(t, queue.Erase(blocks[2].Hash()))
		require.NoError(t, queue.Erase(blocks[4].Hash()))

		block1, ok := queue.Pop()
		assert.True(t, ok)
		assert.Equal(t, block1, blocks[1])

		block2, ok := queue.Pop()
		assert.True(t, ok)
		assert.Equal(t, block2, blocks[3])

		_, ok = queue.Pop()
		assert.False(t, ok)
	})

	t.Run("Failed, over limits random valid block push", func(t *testing.T) {
		limit := RandomConfig().Queue.BlockLimits
		blocks := make([]model.Block, 0, limit)
		go NumWaitChan(limit)
		for i := 0; i < limit; i++ {
			blocks = append(blocks, RandomBlock())
			require.NoError(t, queue.Push(blocks[i]))
		}

		err := queue.Push(RandomBlock())
		assert.EqualError(t, errors.Cause(err), core.ErrProposalQueueLimits.Error())

		for i := 0; i < limit; i++ {
			front, ok := queue.Pop()
			assert.True(t, ok)
			assert.Equal(t, blocks[i], front)
		}

		_, ok := queue.Pop()
		assert.False(t, ok)
	})

	t.Run("Failed, alrady exist block", func(t *testing.T) {
		block := RandomBlock()
		err := queue.Push(block)
		require.NoError(t, err)
		queue.WaitPush()

		err = queue.Push(block)
		assert.EqualError(t, errors.Cause(err), core.ErrProposalQueueAlreadyExist.Error())

		exBlock, ok := queue.Pop()
		require.True(t, ok)
		assert.Equal(t, block, exBlock)
	})

	t.Run("Failed, over limits random block push and erase", func(t *testing.T) {
		limit := RandomConfig().Queue.BlockLimits
		blocks := make([]model.Block, 0, limit)
		blocks2 := make([]model.Block, 0, limit/2)
		for i := 0; i < limit; i++ {
			blocks = append(blocks, RandomBlock())
			require.NoError(t, queue.Push(blocks[i]))
			queue.WaitPush()
		}

		err := queue.Push(RandomBlock())
		assert.EqualError(t, errors.Cause(err), core.ErrProposalQueueLimits.Error())

		for i := 0; i < limit; i += 2 {
			blocks2 = append(blocks2, RandomBlock())
			assert.NoError(t, queue.Erase(blocks[i].Hash()))
			require.NoError(t, queue.Push(blocks2[i/2]))
			queue.WaitPush()
		}

		for i := 0; i < limit; i++ {
			front, ok := queue.Pop()
			assert.True(t, ok)
			if i < limit/2 {
				assert.Equal(t, blocks[i*2+1], front)
			} else {
				assert.Equal(t, blocks2[i-limit/2], front)
			}
		}

		_, ok := queue.Pop()
		assert.False(t, ok)
	})

	t.Run("Failed, nothing hash erase", func(t *testing.T) {
		block := RandomBlock()
		require.NoError(t, queue.Push(block))
		queue.WaitPush()

		block2 := RandomBlock()
		err := queue.Erase(block2.Hash())
		assert.EqualError(t, core.ErrProposalQueueEraseUnexist, errors.Cause(err).Error())

		require.NoError(t, queue.Push(block2))
		queue.WaitPush()
		require.NoError(t, queue.Erase(block2.Hash()))

		exBlock, ok := queue.Pop()
		require.True(t, ok)
		assert.Equal(t, block, exBlock)

		err = queue.Erase(block.Hash())
		assert.EqualError(t, core.ErrProposalQueueEraseUnexist, errors.Cause(err).Error())

		_, ok = queue.Pop()
		assert.False(t, ok)
	})
}

func TestProposalBlockQueueOnMemory(t *testing.T) {
	queue := NewProposalBlockQueueOnMemory(RandomConfig())
	testProposalBlockQueue(t, queue)
}
