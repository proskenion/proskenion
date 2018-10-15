package command

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type CommandValidator struct{}

func NewCommandValidator() core.CommandValidator {
	return &CommandValidator{}
}

func (c *CommandValidator) Transfer(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}

func (c *CommandValidator) CreateAccount(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}
func (c *CommandValidator) AddAsset(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}
func (c *CommandValidator) AddPublicKey(wsv model.ObjectFinder, cmd model.Command) error {
	return nil
}
