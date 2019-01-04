package command

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type CommandExecutor struct {
	factory model.ModelFactory
}

func NewCommandExecutor() core.CommandExecutor {
	return &CommandExecutor{}
}

func (c *CommandExecutor) SetFactory(factory model.ModelFactory) {
	c.factory = factory
}

func (c *CommandExecutor) TransferBalance(wsv model.ObjectFinder, cmd model.Command) error {
	transfer := cmd.GetTransferBalance()
	srcAccount := c.factory.NewEmptyAccount()
	destAccount := c.factory.NewEmptyAccount()
	srcId := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	destId := model.MustAddress(model.MustAddress(cmd.GetTransferBalance().GetDestAccountId()).AccountId())
	if err := wsv.Query(srcId, srcAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotFoundSrcAccountId, err.Error())
	}
	if err := wsv.Query(destId, destAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotFoundDestAccountId, err.Error())
	}
	if srcAccount.GetBalance()-transfer.GetBalance() < 0 {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotEnoughSrcAccountBalance,
			fmt.Errorf("srcAccount Amount: %d, transfer Acmount: %d", srcAccount.GetBalance(), transfer.GetBalance()).Error())
	}
	newSrcAccount := c.factory.NewAccountBuilder().
		From(srcAccount).
		Balance(srcAccount.GetBalance() - transfer.GetBalance()).
		Build()
	newDestAccount := c.factory.NewAccountBuilder().
		From(destAccount).
		Balance(destAccount.GetBalance() + transfer.GetBalance()).
		Build()
	if err := wsv.Append(srcId, newSrcAccount); err != nil {
		return err
	}
	if err := wsv.Append(destId, newDestAccount); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) CreateAccount(wsv model.ObjectFinder, cmd model.Command) error {
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	newAccount := c.factory.NewAccountBuilder().
		AccountId(id.Account() + "@" + id.Domain()).
		AccountName(id.Account()).
		PublicKeys(cmd.GetCreateAccount().GetPublicKeys()).
		Quorum(cmd.GetCreateAccount().GetQuorum()).
		Build()

	ac := c.factory.NewEmptyAccount()
	if err := wsv.Query(id, ac); err == nil {
		if ac.GetAccountId() == cmd.GetTargetId() {
			return errors.Wrap(core.ErrCommandExecutorCreateAccountAlreadyExistAccount,
				fmt.Errorf("already exist accountId : %s", id.AccountId()).Error())
		}
	}
	if err := wsv.Append(id, newAccount); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) AddBalance(wsv model.ObjectFinder, cmd model.Command) error {
	aa := cmd.GetAddBalance()
	ac := c.factory.NewEmptyAccount()
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	if err := wsv.Query(id, ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddBalanceNotExistAccount, err.Error())
	}
	newAc := c.factory.NewAccountBuilder().
		From(ac).
		Balance(ac.GetBalance() + aa.GetBalance()).
		Build()
	if err := wsv.Append(id, newAc); err != nil {
		return err
	}
	return nil
}

func containsPublicKey(keys []model.PublicKey, key model.PublicKey) bool {
	for _, k := range keys {
		if bytes.Equal(k, key) {
			return true
		}
	}
	return false
}

func (c *CommandExecutor) AddPublicKeys(wsv model.ObjectFinder, cmd model.Command) error {
	ap := cmd.GetAddPublicKeys()
	ac := c.factory.NewEmptyAccount()
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	if err := wsv.Query(id, ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddPublicKeyNotExistAccount, err.Error())
	}
	if containsPublicKey(ac.GetPublicKeys(), ap.GetPublicKeys()[0]) {
		return errors.Wrapf(core.ErrCommandExecutorAddPublicKeyDuplicatePubkey,
			"duplicate key : %x", ap.GetPublicKeys())
	}
	newAc := c.factory.NewAccountBuilder().
		From(ac).
		PublicKeys(append(ac.GetPublicKeys(), ap.GetPublicKeys()[0])).
		Build()
	if err := wsv.Append(id, newAc); err != nil {
		return err
	}
	return nil
}
