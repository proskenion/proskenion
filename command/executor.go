package command

import (
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/core"
)

type CommandExecutor struct {}

func NewCommandExecutor() core.CommandExecutor {
	return &CommandExecutor{}
}

// WIP
func (c *CommandExecutor) Transfer(transfer model.Transfer) error {
	return nil
}

