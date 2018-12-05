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
	if err := wsv.Query(cmd.GetTargetId(), srcAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotFoundSrcAccountId, err.Error())
	}
	if err := wsv.Query(transfer.GetDestAccountId(), destAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotFoundDestAccountId, err.Error())
	}
	if srcAccount.GetBalance()-transfer.GetBalance() < 0 {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotEnoughSrcAccountBalance,
			fmt.Errorf("srcAccount Amount: %d, transfer Acmount: %d", srcAccount.GetBalance(), transfer.GetBalance()).Error())
	}
	newSrcAccount := c.factory.NewAccount(
		srcAccount.GetAccountId(),
		srcAccount.GetAccountName(),
		srcAccount.GetPublicKeys(),
		srcAccount.GetBalance()-transfer.GetBalance(),
	)
	newDestAccount := c.factory.NewAccount(
		destAccount.GetAccountId(),
		destAccount.GetAccountName(),
		destAccount.GetPublicKeys(),
		destAccount.GetBalance()+transfer.GetBalance(),
	)
	if err := wsv.Append(newSrcAccount.GetAccountId(), newSrcAccount); err != nil {
		return err
	}
	if err := wsv.Append(newDestAccount.GetAccountId(), newDestAccount); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) CreateAccount(wsv model.ObjectFinder, cmd model.Command) error {
	newAccount := c.factory.NewAccount(cmd.GetTargetId(), cmd.GetTargetId(), make([]model.PublicKey, 0), 0)
	ac := c.factory.NewEmptyAccount()
	if err := wsv.Query(cmd.GetTargetId(), ac); err == nil {
		if ac.GetAccountId() == cmd.GetTargetId() {
			return errors.Wrap(core.ErrCommandExecutorCreateAccountAlreadyExistAccount,
				fmt.Errorf("already exist accountId : %s", cmd.GetTargetId()).Error())
		}
	}
	if err := wsv.Append(cmd.GetTargetId(), newAccount); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) AddBalance(wsv model.ObjectFinder, cmd model.Command) error {
	aa := cmd.GetAddBalance()
	ac := c.factory.NewEmptyAccount()
	if err := wsv.Query(cmd.GetTargetId(), ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddBalanceNotExistAccount, err.Error())
	}
	if ac.GetAccountId() != cmd.GetTargetId() {
		return core.ErrCommandExecutorAddBalanceNotExistAccount
	}
	newAc := c.factory.NewAccount(
		ac.GetAccountId(),
		ac.GetAccountName(),
		ac.GetPublicKeys(),
		ac.GetBalance()+aa.GetBalance(),
	)
	if err := wsv.Append(newAc.GetAccountId(), newAc); err != nil {
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
	if err := wsv.Query(cmd.GetTargetId(), ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddPublicKeyNotExistAccount, err.Error())
	}
	if ac.GetAccountId() != cmd.GetTargetId() {
		return core.ErrCommandExecutorAddPublicKeyNotExistAccount
	}
	if containsPublicKey(ac.GetPublicKeys(), ap.GetPublicKeys()[0]) {
		return errors.Wrapf(core.ErrCommandExecutorAddPublicKeyDuplicatePubkey,
			"duplicate key : %x", ap.GetPublicKeys())
	}
	newAc := c.factory.NewAccount(
		ac.GetAccountId(),
		ac.GetAccountName(),
		append(ac.GetPublicKeys(), ap.GetPublicKeys()[0]),
		ac.GetBalance(),
	)
	if err := wsv.Append(newAc.GetAccountId(), newAc); err != nil {
		return err
	}
	return nil
}
