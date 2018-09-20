package command

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type CommandValidator struct{}

func NewCommandValidator() core.CommandValidator {
	return &CommandValidator{}
}

func (c *CommandValidator) Transfer(wsv model.ObjectFinder, transfer model.Transfer) error {
	return nil
}
