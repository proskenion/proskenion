package core

type Consensus interface {
	Boot()
	Receiver()
	Patrol()
}

