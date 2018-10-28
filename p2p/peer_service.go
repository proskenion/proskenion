package p2p

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type PeerService struct {
	peers []model.Peer
}

// TODO arrangement peers list
func NewPeerService(peers []model.Peer) core.PeerService {
	return &PeerService{peers}
}

func (s *PeerService) List() []model.Peer {
	return s.peers
}
