package convertor

import (
	"fmt"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Command struct {
	*proskenion.Command
	executor  core.CommandExecutor
	validator core.CommandValidator
}

func (c *Command) Execute(wsv model.ObjectFinder) error {
	switch x := c.GetCommand().(type) {
	case *proskenion.Command_Transfer:
		return c.executor.Transfer(wsv, c)
	default:
		return fmt.Errorf("Command has unexpected type %T", x)
	}
}

func (c *Command) Validate(wsv model.ObjectFinder) error {
	switch x := c.GetCommand().(type) {
	case *proskenion.Command_Transfer:
		return c.validator.Transfer(wsv, c)
	default:
		return fmt.Errorf("Command has unexpected type %T", x)
	}
}

func (c *Command) GetTransfer() model.Transfer {
	return c.Command.GetTransfer()
}

func (c *Command) GetCreateAccount() model.CreateAccount {
	return c.Command.GetCreateAccount()
}

func (c *Command) GetAddAsset() model.AddAsset {
	return c.Command.GetAddAsset()
}
