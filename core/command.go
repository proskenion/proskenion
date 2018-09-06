package core

import . "github.com/proskenion/proskenion/core/model"

type Executor interface {
	Execute() error
}

type Validator interface {
	Validate() error
}

type CommandExecutor interface {
	Transfer(transfer Transfer) error
}

type CommandValidator interface {
	Transfer(transfer Transfer) error
}
