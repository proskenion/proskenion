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
	if err := wsv.Query(model.MustAddress(cmd.GetTargetId()), srcAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotFoundSrcAccountId, err.Error())
	}
	if err := wsv.Query(model.MustAddress(transfer.GetDestAccountId()), destAccount); err != nil {
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
	if err := wsv.Append(model.MustAddress(newSrcAccount.GetAccountId()), newSrcAccount); err != nil {
		return err
	}
	if err := wsv.Append(model.MustAddress(newDestAccount.GetAccountId()), newDestAccount); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) CreateAccount(wsv model.ObjectFinder, cmd model.Command) error {
	newAccount := c.factory.NewAccountBuilder().
		AccountId(cmd.GetTargetId()).
		AccountName(cmd.GetTargetId()).
		PublicKeys(make([]model.PublicKey, 0)).
		Build()

	ac := c.factory.NewEmptyAccount()
	if err := wsv.Query(model.MustAddress(cmd.GetTargetId()), ac); err == nil {
		if ac.GetAccountId() == cmd.GetTargetId() {
			return errors.Wrap(core.ErrCommandExecutorCreateAccountAlreadyExistAccount,
				fmt.Errorf("already exist accountId : %s", cmd.GetTargetId()).Error())
		}
	}
	if err := wsv.Append(model.MustAddress(cmd.GetTargetId()), newAccount); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) AddBalance(wsv model.ObjectFinder, cmd model.Command) error {
	aa := cmd.GetAddBalance()
	ac := c.factory.NewEmptyAccount()
	if err := wsv.Query(model.MustAddress(cmd.GetTargetId()), ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddBalanceNotExistAccount, err.Error())
	}
	if ac.GetAccountId() != cmd.GetTargetId() {
		return core.ErrCommandExecutorAddBalanceNotExistAccount
	}
	newAc := c.factory.NewAccountBuilder().
		From(ac).
		Balance(ac.GetBalance() + aa.GetBalance()).
		Build()
	if err := wsv.Append(model.MustAddress(newAc.GetAccountId()), newAc); err != nil {
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
	if err := wsv.Query(model.MustAddress(cmd.GetTargetId()), ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddPublicKeyNotExistAccount, err.Error())
	}
	if ac.GetAccountId() != cmd.GetTargetId() {
		return core.ErrCommandExecutorAddPublicKeyNotExistAccount
	}
	if containsPublicKey(ac.GetPublicKeys(), ap.GetPublicKeys()[0]) {
		return errors.Wrapf(core.ErrCommandExecutorAddPublicKeyDuplicatePubkey,
			"duplicate key : %x", ap.GetPublicKeys())
	}
	newAc := c.factory.NewAccountBuilder().
		From(ac).
		PublicKeys(append(ac.GetPublicKeys(), ap.GetPublicKeys()[0])).
		Build()
	if err := wsv.Append(model.MustAddress(newAc.GetAccountId()), newAc); err != nil {
		return err
	}
	return nil
}
