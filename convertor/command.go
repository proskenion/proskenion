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
	case *proskenion.Command_TransferBalance:
		return c.executor.TransferBalance(wsv, c)
	case *proskenion.Command_AddBalance:
		return c.executor.AddBalance(wsv, c)
	case *proskenion.Command_CreateAccount:
		return c.executor.CreateAccount(wsv, c)
	case *proskenion.Command_AddPublicKeys:
		return c.executor.AddPublicKeys(wsv, c)
	default:
		return fmt.Errorf("Command has unexpected type %T", x)
	}
}

func (c *Command) Validate(wsv model.ObjectFinder) error {
	switch x := c.GetCommand().(type) {
	case *proskenion.Command_TransferBalance:
		return c.validator.TransferBalance(wsv, c)
	case *proskenion.Command_AddBalance:
		return c.validator.AddBalance(wsv, c)
	case *proskenion.Command_CreateAccount:
		return c.validator.CreateAccount(wsv, c)
	case *proskenion.Command_AddPublicKeys:
		return c.validator.AddPublicKeys(wsv, c)
	default:
		return fmt.Errorf("Command has unexpected type %T", x)
	}
}

func (c *Command) GetTransferBalance() model.TransferBalance {
	return c.Command.GetTransferBalance()
}

func (c *Command) GetCreateAccount() model.CreateAccount {
	return c.Command.GetCreateAccount()
}

func (c *Command) GetAddBalance() model.AddBalance {
	return c.Command.GetAddBalance()
}

func (c *Command) GetAddPublicKeys() model.AddPublicKeys {
	return c.Command.GetAddPublicKeys()
}
