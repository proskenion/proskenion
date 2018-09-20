package command

import (
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

func (c *CommandExecutor) Transfer(wsv model.ObjectFinder, transfer model.Transfer) error {
	return nil
}
