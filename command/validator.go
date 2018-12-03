package command

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type CommandValidator struct {
	fc model.ModelFactory
}

func NewCommandValidator() core.CommandValidator {
	return &CommandValidator{}
}

func (c *CommandValidator) SetFactory(factory model.ModelFactory) {
	c.fc = factory
}

func (c *CommandValidator) Transfer(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) CreateAccount(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) AddAsset(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) AddPublicKey(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func containsPublicKeyInSignatures(sigs []model.Signature, key model.PublicKey) bool {
	for _, sig := range sigs {
		if bytes.Equal(sig.GetPublicKey(), key) {
			return true
		}
	}
	return false
}

// Transaction 全体の Stateful Validate
// 1. 既に同一の Transaction Hash を持つ Transaction が存在するか
// 2. Authorizer アカウントが存在するか
// 3. Authorozer の権限を行使するための署名が揃っているか
func (c *CommandValidator) Tx(wsv model.ObjectFinder, txh model.TxFinder, tx model.Transaction) error {
	hash := tx.Hash()
	_, err := txh.Query(hash)
	if errors.Cause(err) != core.ErrTxHistoryNotFound {
		return core.ErrTxValidateAlreadyExist
	}
	for _, cmd := range tx.GetPayload().GetCommands() {
		ac := c.fc.NewEmptyAccount()
		err := wsv.Query(cmd.GetAuthorizerId(), ac)
		if err != nil {
			return errors.Wrapf(core.ErrTxValidateNotFoundAuthorizer,
				"authorizer : %s", cmd.GetAuthorizerId())
		}
		// TODO : sort すれば全体一致判定をO(nlogn)
		for _, key := range ac.GetPublicKeys() {
			if !containsPublicKeyInSignatures(tx.GetSignatures(), key) {
				return errors.Wrapf(core.ErrTxValidateNotSignedAuthorizer,
					"authorizer : %s, expect key : %x",
					cmd.GetAuthorizerId(), key)
			}
		}
	}
	return nil
}
