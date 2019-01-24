package repository

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/prosl"
	"io/ioutil"
)

func rollBackTx(tx core.RepositoryTx, mtErr error) error {
	if err := tx.Rollback(); err != nil {
		return errors.Wrap(err, mtErr.Error())
	}
	return mtErr
}

func commitTx(tx core.RepositoryTx) error {
	if err := tx.Commit(); err != nil {
		return rollBackTx(tx, err)
	}
	return nil
}

type Repository struct {
	dba     core.DBA
	cryptor core.Cryptor
	fc      model.ModelFactory
	conf    *config.Config

	TopBlock model.Block
	Height   int64
}

func NewRepository(dba core.DBA, cryptor core.Cryptor, fc model.ModelFactory, conf *config.Config) core.Repository {
	return &Repository{dba, cryptor, fc, conf, nil, 0}
}

func (r *Repository) Begin() (core.RepositoryTx, error) {
	tx, err := r.dba.Begin()
	if err != nil {
		return nil, err
	}
	return &RepositoryTx{tx, r.cryptor, r.fc}, nil
}

func (r *Repository) Top() (model.Block, bool) {
	if r.TopBlock == nil {
		return nil, false
	}
	return r.TopBlock, true
}

func (r *Repository) GetDelegatedAccounts() ([]model.Account, error) {
	top, ok := r.Top()
	if !ok {
		panic("Failed Repository error empty top")
	}
	rtx, err := r.Begin()
	if err != nil {
		return nil, err
	}
	wsv, err := rtx.WSV(top.GetPayload().GetWSVHash())
	if err != nil {
		return nil, err
	}
	st := ProslStorage(r.fc, r.conf)
	id := model.MustAddress(r.conf.Prosl.Consensus.Id)
	if err := wsv.Query(id, st); err != nil {
		return nil, err
	}
	if err := commitTx(rtx); err != nil {
		return nil, err
	}

	prData := st.GetFromKey(core.ProslKey).GetData()
	pr := prosl.NewProsl(r.fc, r, r.cryptor, r.conf)
	if err := pr.Unmarshal(prData); err != nil {
		return nil, err
	}
	ret, vars, err := pr.Execute()
	if err != nil {
		return nil, fmt.Errorf("errors: %s\nvariables: %+v\n", err.Error(), vars)
	}
	list := ret.GetList()
	acs := make([]model.Account, 0, len(list))
	for _, ac := range list {
		acs = append(acs, ac.GetAccount())
	}
	return acs, nil
}

func (r *Repository) loadMPTrees(dtx core.RepositoryTx, preBlock model.Block, preBlockHash model.Hash) (core.Blockchain, core.WSV, core.TxHistory, model.Block, error) {
	bc, err := dtx.Blockchain(preBlockHash)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(core.ErrRepositoryCommitLoadPreBlock, err.Error())
	}
	if preBlock == nil {
		preBlock, err = bc.Get(preBlockHash)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(core.ErrRepositoryCommitLoadPreBlock, err.Error())
		}
	}
	wsv, err := dtx.WSV(preBlock.GetPayload().GetWSVHash())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(core.ErrRepositoryCommitLoadWSV, err.Error())
	}
	txHistory, err := dtx.TxHistory(preBlock.GetPayload().GetTxHistoryHash())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(core.ErrRepositoryCommitLoadTxHistory, err.Error())
	}
	return bc, wsv, txHistory, preBlock, nil
}

// Incentive Prosl exeucute (fource execute)
func (r *Repository) executeProslIncentive(wsv core.WSV, txHistory core.TxHistory) error {
	// 1. get prosl
	proSt := r.fc.NewEmptyStorage()
	if err := wsv.Query(model.MustAddress(r.conf.Prosl.Incentive.Id), proSt); err != nil {
		return nil
	}
	pr := prosl.NewProsl(r.fc, r, r.cryptor, r.conf)
	proslByte := proSt.GetFromKey(core.ProslKey).GetData()
	if err := pr.Unmarshal(proslByte); err != nil {
		return nil
	}
	// 2. execute incentive prosl
	ret, vars, err := pr.Execute()
	if err != nil {
		fmt.Printf("Incentive Prosl Error\nvariables: %+v, error: %s\n", vars, err.Error())
	} else if ret == nil || ret.GetTransaction() == nil {
		fmt.Printf("Incentive Prosl Error\nvariables: %+v, error: empty incentive tx\n", vars)
	} else {
		// 3. execute incentive tx
		for _, cmd := range ret.GetTransaction().GetPayload().GetCommands() {
			if err := cmd.Execute(wsv); err != nil {
				return fmt.Errorf("Incentive Tx Error\n%s\n%+v", err.Error(), ret.GetTransaction())
			}
		}
		if err := txHistory.Append(ret.GetTransaction()); err != nil {
			return fmt.Errorf("Incentive Tx Append Error\n%s\n%+v", err.Error(), ret.GetTransaction())
		}
	}
	return nil
}

func (r *Repository) appendAndUpdateBlock(bc core.Blockchain, block model.Block) error {
	// block を追加・
	if err := bc.Append(block); err != nil {
		return err
	}
	// top ブロックを更新
	if r.Height < block.GetPayload().GetHeight() {
		r.Height = block.GetPayload().GetHeight()
		r.TopBlock = block
	}
	return nil
}

func (r *Repository) CreateBlock(queue core.ProposalTxQueue, round int32, now int64) (model.Block, core.TxList, error) {
	dtx, err := r.Begin()
	if err != nil {
		return nil, nil, err
	}
	preBlock, ok := r.Top()
	if !ok {
		return nil, nil, errors.Errorf("Failed CreateBlock internal error, after execute genesis block")
	}
	// load state
	bc, wsv, txHistory, _, err := r.loadMPTrees(dtx, preBlock, preBlock.Hash())
	if err != nil {
		return nil, nil, err
	}

	// execute incentive prosl transaction. (fource execute)
	if err := r.executeProslIncentive(wsv, txHistory); err != nil {
		return nil, nil, rollBackTx(dtx, err)
	}

	txList := NewTxList(r.cryptor)
	// ProposalTxQueue から valid な Tx をとってきて hoge る
	for txList.Size() < r.conf.Commit.NumTxInBlock {
		tx, ok := queue.Pop()
		if !ok {
			break
		}

		// tx を構築
		if err := tx.Validate(wsv, txHistory); err != nil {
			goto txskip
		}
		for _, cmd := range tx.GetPayload().GetCommands() {
			if err := cmd.Validate(wsv); err != nil {
				goto txskip
			}
			if err := cmd.Execute(wsv); err != nil {
				goto txskip // WIP : 要考
				//return nil, nil, rollBackTx(dtx, err)
			}
		}
		if err := txHistory.Append(tx); err != nil {
			return nil, nil, rollBackTx(dtx, err)
		}
		if err := txList.Push(tx); err != nil {
			return nil, nil, rollBackTx(dtx, err)
		}

	txskip:
	}

	newBlock := r.fc.NewBlockBuilder().
		Round(round).
		TxsHash(txList.Hash()).
		TxHistoryHash(txHistory.Hash()).
		WSVHash(wsv.Hash()).
		CreatedTime(now).
		Height(preBlock.GetPayload().GetHeight() + 1).
		PreBlockHash(preBlock.Hash()).
		Build()
	if err = newBlock.Sign(r.conf.Peer.PublicKeyBytes(), r.conf.Peer.PrivateKeyBytes()); err != nil {
		return nil, nil, rollBackTx(dtx, err)
	}

	// append Block and repository state update
	if err := r.appendAndUpdateBlock(bc, newBlock); err != nil {
		return nil, nil, rollBackTx(dtx, err)
	}
	if err := commitTx(dtx); err != nil {
		return nil, nil, err
	}
	return newBlock, txList, nil
}

func (r *Repository) Commit(block model.Block, txList core.TxList) (err error) {
	dtx, err := r.Begin()
	if err != nil {
		return err
	}

	// load state
	preBlockHash := block.GetPayload().GetPreBlockHash()
	bc, wsv, txHistory, _, err := r.loadMPTrees(dtx, nil, preBlockHash)
	if err != nil {
		return err
	}

	// transactions execute
	for _, tx := range txList.List() {
		if err := tx.Validate(wsv, txHistory); err != nil {
			return rollBackTx(dtx, err)
		}
		for _, cmd := range tx.GetPayload().GetCommands() {
			if err := cmd.Validate(wsv); err != nil {
				return rollBackTx(dtx, err)
			}
			if err := cmd.Execute(wsv); err != nil {
				return rollBackTx(dtx, err)
			}
		}
		if err := txHistory.Append(tx); err != nil {
			return rollBackTx(dtx, err)
		}
	}

	// Incentive Prosl exeucute (fource execute)
	if err := r.executeProslIncentive(wsv, txHistory); err != nil {
		return rollBackTx(dtx, err)
	}

	// hash check
	if !bytes.Equal(block.GetPayload().GetWSVHash(), wsv.Hash()) {
		return rollBackTx(dtx, errors.Errorf("not equaled wsv Hash : %x", wsv.Hash()))
	}
	if !bytes.Equal(block.GetPayload().GetTxHistoryHash(), txHistory.Hash()) {
		return rollBackTx(dtx, errors.Errorf("not equaled txHistory Hash : %x", txHistory.Hash()))
	}

	// append Block and repository state update
	if err := r.appendAndUpdateBlock(bc, block); err != nil {
		return rollBackTx(dtx, err)
	}
	return commitTx(dtx)
}

func ProslStorage(fc model.ModelFactory, conf *config.Config) model.Storage {
	return fc.NewStorageBuilder().
		Data(core.ProslKey, nil).
		Str(core.ProslTypeKey, "none").
		Build()
}

func (r *Repository) getProslBytes(filename string, pr core.Prosl) ([]byte, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := pr.ConvertFromYaml(buf); err != nil {
		return nil, err
	}
	return pr.Marshal()
}

func (r *Repository) genesisProslSetting() (model.Transaction, error) {
	proSt := ProslStorage(r.fc, r.conf)
	pr := prosl.NewProsl(r.fc, r, r.cryptor, r.conf)
	incPr, err := r.getProslBytes(r.conf.Prosl.Incentive.Path, pr)
	if err != nil {
		return nil, err
	}
	conPr, err := r.getProslBytes(r.conf.Prosl.Consensus.Path, pr)
	if err != nil {
		return nil, err
	}
	updPr, err := r.getProslBytes(r.conf.Prosl.Update.Path, pr)
	if err != nil {
		return nil, err
	}
	return r.fc.NewTxBuilder().
		DefineStorage(r.conf.Root.Id, r.conf.Prosl.Incentive.Id, proSt).
		DefineStorage(r.conf.Root.Id, r.conf.Prosl.Consensus.Id, proSt).
		DefineStorage(r.conf.Root.Id, r.conf.Prosl.Update.Id, proSt).
		UpdateObject(r.conf.Root.Id, r.conf.Prosl.Incentive.Id, core.ProslKey,
			r.fc.NewObjectBuilder().Data(incPr)).
		UpdateObject(r.conf.Root.Id, r.conf.Prosl.Consensus.Id, core.ProslKey,
			r.fc.NewObjectBuilder().Data(conPr)).
		UpdateObject(r.conf.Root.Id, r.conf.Prosl.Update.Id, core.ProslKey,
			r.fc.NewObjectBuilder().Data(updPr)).
		UpdateObject(r.conf.Root.Id, r.conf.Prosl.Incentive.Id, core.ProslTypeKey,
			r.fc.NewObjectBuilder().Str(core.IncentiveKey)).
		UpdateObject(r.conf.Root.Id, r.conf.Prosl.Consensus.Id, core.ProslTypeKey,
			r.fc.NewObjectBuilder().Str(core.ConsensusKey)).
		UpdateObject(r.conf.Root.Id, r.conf.Prosl.Update.Id, core.ProslTypeKey,
			r.fc.NewObjectBuilder().Str(core.UpdateKey)).
		CreatedTime(0).
		Build(), nil
}

func (r *Repository) GenesisCommit(txList core.TxList) (err error) {
	dtx, err := r.Begin()
	if err != nil {
		return err
	}

	// load state
	var bc core.Blockchain
	if bc, err = dtx.Blockchain(nil); err != nil {
		return err
	}
	wsv, err := dtx.WSV(nil)
	if err != nil {
		return errors.Wrap(core.ErrRepositoryCommitLoadWSV, err.Error())
	}
	txHistory, err := dtx.TxHistory(nil)
	if err != nil {
		return errors.Wrap(core.ErrRepositoryCommitLoadTxHistory, err.Error())
	}

	// Genesis Commit to add prosl
	genTx, err := r.genesisProslSetting()
	if err != nil {
		return err
	}
	err = txList.Push(genTx)
	if err != nil {
		return err
	}

	// transactions execute (no validate)
	for _, tx := range txList.List() {
		for _, cmd := range tx.GetPayload().GetCommands() {
			if err := cmd.Execute(wsv); err != nil {
				return rollBackTx(dtx, err)
			}
		}
		if err := txHistory.Append(tx); err != nil {
			return rollBackTx(dtx, err)
		}
	}

	// hash check and block 生成
	wsvHash := wsv.Hash()
	if err != nil {
		return rollBackTx(dtx, err)
	}
	txHistoryHash := txHistory.Hash()
	if err != nil {
		return rollBackTx(dtx, err)
	}
	genesisBlock := r.fc.NewBlockBuilder().
		CreatedTime(0).
		TxsHash(txList.Hash()).
		PreBlockHash(nil).
		TxHistoryHash(txHistoryHash).
		WSVHash(wsvHash).
		Round(0).
		Height(0).
		Build()

	// block を追加・
	if err := bc.Append(genesisBlock); err != nil {
		return err
	}
	// top ブロックを更新
	r.Height = genesisBlock.GetPayload().GetHeight()
	r.TopBlock = genesisBlock
	return commitTx(dtx)
}

type RepositoryTx struct {
	tx      core.DBATx
	cryptor core.Cryptor
	fc      model.ModelFactory
}

func (r *RepositoryTx) WSV(hash model.Hash) (core.WSV, error) {
	return NewWSV(r.tx, r.cryptor, r.fc, hash)
}

func (r *RepositoryTx) TxHistory(hash model.Hash) (core.TxHistory, error) {
	return NewTxHistory(r.tx, r.fc, r.cryptor, hash)
}

func (r *RepositoryTx) Blockchain(topBlockHash model.Hash) (core.Blockchain, error) {
	return NewBlockchainFromTopBlock(r.tx, r.fc, r.cryptor, topBlockHash)
}

func (r *RepositoryTx) Top() (model.Block, error) {
	return nil, nil
}

func (r *RepositoryTx) Commit() error {
	return r.tx.Commit()
}

func (r *RepositoryTx) Rollback() error {
	return r.tx.Rollback()
}
