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
		return c.GetTransfer().Execute(wsv)
	default:
		return fmt.Errorf("Command has unexpected type %T", x)
	}
}

func (c *Command) Validate(wsv model.ObjectFinder) error {
	switch x := c.GetCommand().(type) {
	case *proskenion.Command_Transfer:
		return c.GetTransfer().Validate(wsv)
	default:
		return fmt.Errorf("Command has unexpected type %T", x)
	}
}

func (c *Command) GetTransfer() model.Transfer {
	if c.Command != nil {
		return &Transfer{c.Command.GetTransfer(), c.executor, c.validator}
	}
	return &Transfer{nil, c.executor, c.validator}
}

type Transfer struct {
	*proskenion.Transfer
	executor  core.CommandExecutor
	validator core.CommandValidator
}

func (c *Transfer) Execute(wsv model.ObjectFinder) error {
	return c.executor.Transfer(wsv, c)
}

func (c *Transfer) Validate(wsv model.ObjectFinder) error {
	return c.validator.Transfer(wsv, c)
}
