package convertor

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Command struct {
	*proskenion.Command
	executor  core.CommandExecutor
	validator core.CommandValidator
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

func (c *Transfer) Execute() error {
	return c.executor.Transfer(c)
}

func (c *Transfer) Validate() error {
	return c.validator.Transfer(c)
}
