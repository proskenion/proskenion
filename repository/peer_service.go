package repository

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type PeerService struct {
	c         core.Cryptor
	peers     []model.Peer
	cacheHash model.Hash
}

// TODO arrangement peers list
func NewPeerService(c core.Cryptor) core.PeerService {
	return &PeerService{c, make([]model.Peer, 0), nil}
}

func (s *PeerService) Set(peers []model.Peer) {
	s.cacheHash = nil
	s.peers = peers
}

func (s *PeerService) List() []model.Peer {
	return s.peers
}
func (s *PeerService) Hash() model.Hash {
	if s.cacheHash == nil {
		hashes := make([]model.Hash, 0, len(s.peers))
		for _, peer := range s.peers {
			hashes = append(hashes, peer.Hash())
		}
		s.cacheHash = s.c.ConcatHash(hashes...)
	}
	return s.cacheHash
}
