package test_utils

import (
	"github.com/proskenion/proskenion/command"
	"github.com/proskenion/proskenion/core"
)

func RandomCommandExecutor() core.CommandExecutor {
	fc := NewTestFactory()
	ex := command.NewCommandExecutor()
	ex.SetFactory(fc)
	return ex
}

func RandomCommandValidator() core.CommandValidator {
	return command.NewCommandValidator()
}
