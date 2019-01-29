package convertor_test

import (
	"github.com/pkg/errors"
	. "github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_GetHash(t *testing.T) {
	txs := make([]model.Transaction, 50)
	for id, _ := range txs {
		txs[id] = RandomValidTx()
	}
	for id, a := range txs {
		for jd, b := range txs {
			if id != jd {
				assert.NotEqual(t, a.Hash(), b.Hash())
			} else {
				assert.Equal(t, a.Hash(), b.Hash())
			}
		}
	}
}

func TestTransaction_Verfy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tx := RandomValidTx()
		assert.NoError(t, tx.Verify())
	})
	t.Run("success w valid sign ", func(t *testing.T) {
		tx := RandomTx()
		vpu1, vpr1 := RandomKeyPairs()
		vpu2, vpr2 := RandomKeyPairs()
		tx.Sign(vpu1, vpr1)
		tx.Sign(vpu2, vpr2)
		assert.NoError(t, tx.Verify())
	})
	t.Run("failed twice invalid sign ", func(t *testing.T) {
		tx := RandomTx()
		vpu1, vpr1 := RandomKeyPairs()
		_, vpr2 := RandomKeyPairs()
		tx.Sign(vpu1, vpr1)
		tx.Sign(vpu1, vpr2)
		assert.EqualError(t, errors.Cause(tx.Verify()), core.ErrCryptorVerify.Error())
	})
	t.Run("failed invalid signature", func(t *testing.T) {
		tx := RandomInvalidTx(t)
		assert.EqualError(t, errors.Cause(tx.Verify()), core.ErrCryptorVerify.Error())
	})
	t.Run("failed not signed", func(t *testing.T) {
		tx := RandomTx()
		assert.EqualError(t, errors.Cause(tx.Verify()), ErrInvalidSignatures.Error())
	})
	t.Run("failed nil signature", func(t *testing.T) {
		tx := RandomTx()
		tx.(*Transaction).Signatures = make([]*proskenion.Signature, 5)
		assert.EqualError(t, errors.Cause(tx.Verify()), model.ErrInvalidSignature.Error())
	})
	t.Run("failed nil transaction", func(t *testing.T) {
		tx := RandomTx()
		tx.(*Transaction).Transaction = nil
		assert.EqualError(t, errors.Cause(tx.Verify()), ErrInvalidSignatures.Error())
	})
}
