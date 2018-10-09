package command

import (
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

func (c *CommandExecutor) Transfer(wsv model.ObjectFinder, cmd model.Command) error {
	transfer := cmd.GetTransfer()
	srcAccount := c.factory.NewEmptyAccount()
	destAccount := c.factory.NewEmptyAccount()
	if err := wsv.Query(cmd.GetTargetId(), srcAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferNotFoundSrcAccountId, err.Error())
	}
	if err := wsv.Query(transfer.GetDestAccountId(), destAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferNotFoundDestAccountId, err.Error())
	}
	if srcAccount.GetAmount()-transfer.GetAmount() < 0 {
		return errors.Wrap(core.ErrCommandExecutorTransferNotEnoughSrcAccountAmount,
			fmt.Errorf("srcAccount Amount: %d, transfer Acmount: %d", srcAccount.GetAmount(), transfer.GetAmount()).Error())
	}
	newSrcAccount := c.factory.NewAccount(
		srcAccount.GetAccountId(),
		srcAccount.GetAccountName(),
		srcAccount.GetPublicKeys(),
		srcAccount.GetAmount()-transfer.GetAmount(),
	)
	newDestAccount := c.factory.NewAccount(
		destAccount.GetAccountId(),
		destAccount.GetAccountName(),
		destAccount.GetPublicKeys(),
		destAccount.GetAmount()+transfer.GetAmount(),
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

func (c *CommandExecutor) AddAsset(wsv model.ObjectFinder, cmd model.Command) error {
	aa := cmd.GetAddAsset()
	ac := c.factory.NewEmptyAccount()
	if err := wsv.Query(cmd.GetTargetId(), ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddAssetNotExistAccount, err.Error())
	}
	if ac.GetAccountId() != cmd.GetTargetId() {
		return core.ErrCommandExecutorAddAssetNotExistAccount
	}
	newAc := c.factory.NewAccount(
		ac.GetAccountId(),
		ac.GetAccountName(),
		ac.GetPublicKeys(),
		ac.GetAmount()+aa.GetAmount(),
	)
	if err := wsv.Append(newAc.GetAccountId(), newAc); err != nil {
		return err
	}
	return nil
}
