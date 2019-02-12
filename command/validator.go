package command

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type CommandValidator struct {
	fc    model.ModelFactory
	prosl core.Prosl
	conf  *config.Config
}

func NewCommandValidator(conf *config.Config) core.CommandValidator {
	return &CommandValidator{conf: conf}
}

func (c *CommandValidator) SetField(factory model.ModelFactory, prosl core.Prosl) {
	c.fc = factory
	c.prosl = prosl
}

func (c *CommandValidator) TransferBalance(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) CreateAccount(wsv model.ObjectFinder, cmd model.Command) error {
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	ac := c.fc.NewEmptyAccount()
	if err := wsv.Query(id, ac); err == nil {
		return errors.Wrap(core.ErrCommandExecutorCreateAccountAlreadyExistAccount,
			fmt.Errorf("already exist accountId : %s", id.AccountId()).Error())
	}
	return nil
}

func (c *CommandValidator) SetQuorum(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) AddBalance(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) AddPublicKeys(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) DefineStorage(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) CreateStorage(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) UpdateObject(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) AddObject(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) TransferObject(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) AddPeer(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) ActivatePeer(model.ObjectFinder, model.Command) error {
	return nil
}

func (c *CommandValidator) SuspendPeer(model.ObjectFinder, model.Command) error {
	return nil
}

func (c *CommandValidator) BanPeer(model.ObjectFinder, model.Command) error {
	return nil
}

func (c *CommandValidator) Consign(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) CheckAndCommitProsl(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) ForceUpdateStorage(wsv model.ObjectFinder, cmd model.Command) error {
	return core.ErrCommandValidatorForceUpdateStorageCanNotUsedDefault
}

func containsPublicKeyInSignaturesForQuorum(sigs []model.Signature, keys []model.PublicKey, quorum int32) bool {
	cnt := make(map[string]int)
	for _, sig := range sigs {
		cnt[string(sig.GetPublicKey())] = 1
	}
	for _, key := range keys {
		cnt[string(key)] |= 2
	}
	for _, c := range cnt {
		if c == 3 {
			quorum--
		}
	}
	return quorum <= 0
}

// Transaction 全体の Stateful Validate
// 1. 既に同一の Transaction Hash を持つ Transaction が存在するか
// 2. Authorizer アカウントが存在するか
// 3. Authorozer の権限を行使するための署名が揃っているか
func (c *CommandValidator) Tx(wsv model.ObjectFinder, txh model.TxFinder, tx model.Transaction) error {
	hash := tx.Hash()
	_, err := txh.GetTx(hash)
	if errors.Cause(err) != core.ErrTxHistoryNotFound {
		return core.ErrTxValidateAlreadyExist
	}
	for _, cmd := range tx.GetPayload().GetCommands() {
		ac := c.fc.NewEmptyAccount()
		authorizerId := model.MustAddress(model.MustAddress(cmd.GetAuthorizerId()).AccountId())
		err = wsv.Query(authorizerId, ac)
		if err != nil {
			return errors.Wrapf(core.ErrTxValidateNotFoundAuthorizer,
				"authorizer : %s", cmd.GetAuthorizerId())
		}
		if !containsPublicKeyInSignaturesForQuorum(tx.GetSignatures(), ac.GetPublicKeys(), ac.GetQuorum()) {
			return errors.Wrapf(core.ErrTxValidateNotSignedAuthorizer,
				"authorizer : %s, expect keys : %+v",
				cmd.GetAuthorizerId(), ac.GetPublicKeys())
		}
	}
	return nil
}
