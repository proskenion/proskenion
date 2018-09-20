package core

import . "github.com/proskenion/proskenion/core/model"

type Executor interface {
	Execute(ObjectFinder) error
}

type Validator interface {
	Validate(ObjectFinder) error
}

type CommandExecutor interface {
	SetFactory(factory ModelFactory)
	Transfer(ObjectFinder, Transfer) error
}

type CommandValidator interface {
	Transfer(ObjectFinder, Transfer) error
}
