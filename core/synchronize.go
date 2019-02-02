package core

import (
	. "github.com/proskenion/proskenion/core/model"
)

type Synchronizer interface {
	Sync(peer Peer) error
}
