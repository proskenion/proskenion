package synchronize

import (
	"fmt"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"io"
)

type Synchronizer struct {
	rp core.Repository
	cf core.ClientFactory
	fc model.ModelFactory
}

func NewSynchronizer(rp core.Repository, cf core.ClientFactory, fc model.ModelFactory) core.Synchronizer {
	return &Synchronizer{rp, cf, fc}
}

func (s *Synchronizer) activate(peer model.Peer) error {
	client, err := s.cf.APIClient(peer)
	if err != nil {
		return err
	}
	me := s.rp.Me()
	in := s.fc.NewTxBuilder().
		ActivatePeer(me.GetPeerId(), me.GetPeerId()).
		Build()
	if err := in.Sign(me.GetPublicKey(), me.GetPrivateKey()); err != nil {
		return err
	}
	return client.Write(in)
}

func (s *Synchronizer) Sync(peer model.Peer) error {
	top, ok := s.rp.Top()
	if !ok {
		return fmt.Errorf("Failed Sync top block nil error.")
	}
	blockHash := top.Hash()

	client, err := s.cf.SyncClient(peer)
	if err != nil {
		return err
	}

	blockChan := make(chan model.Block)
	txListChan := make(chan core.TxList)
	errChan := make(chan error)

	retErrChan := make(chan error)
	defer close(retErrChan)
	go func() {
		defer close(blockChan)
		defer close(txListChan)
		defer close(errChan)
		err := client.Sync(blockHash, blockChan, txListChan, errChan)
		retErrChan <- err
	}()

	var newBlock model.Block
	var newTxList core.TxList
	for {
		select {
		case newBlock = <-blockChan:
		case newTxList = <-txListChan:
			err := s.rp.Commit(newBlock, newTxList)
			if s.rp.Me().GetActive() {
				errChan <- io.EOF
			} else {
				errChan <- err
			}
		case err := <-retErrChan:
			if err != nil && err != io.EOF {
				return err
			}
			goto afterSync
		}
	}

afterSync:
	return s.activate(peer)
}
