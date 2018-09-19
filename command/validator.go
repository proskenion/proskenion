package command

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type CommandValidator struct {}

func NewCommandValidator() core.CommandValidator {
	return &CommandValidator{}
}

// WIP
func (c *CommandValidator) Transfer(transfer model.Transfer) error {
	return nil
}
