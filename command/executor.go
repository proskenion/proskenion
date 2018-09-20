package command

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type CommandExecutor struct {
	factory model.ModelFactory
}

func NewCommandExecutor(factory model.ModelFactory) core.CommandExecutor {
	return &CommandExecutor{factory}
}

func (c *CommandExecutor) Transfer(wsv model.ObjectFinder, transfer model.Transfer) error {
	return nil
}
