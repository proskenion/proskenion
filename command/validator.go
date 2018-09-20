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

func (c *CommandValidator) CreateAccount(wsv model.ObjectFinder, ca model.CreateAccount) error {
	return nil
}
func (c *CommandValidator) AddAsset(wsv model.ObjectFinder, aa model.AddAsset) error {
	return nil
}