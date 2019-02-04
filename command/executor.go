package command

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
)

type CommandExecutor struct {
	factory model.ModelFactory
	prosl   core.Prosl
	conf    *config.Config
}

func NewCommandExecutor(conf *config.Config) core.CommandExecutor {
	return &CommandExecutor{conf: conf}
}

func (c *CommandExecutor) SetField(factory model.ModelFactory, prosl core.Prosl) {
	c.factory = factory
	c.prosl = prosl
}

func (c *CommandExecutor) TransferBalance(wsv model.ObjectFinder, cmd model.Command) error {
	transfer := cmd.GetTransferBalance()
	srcAccount := c.factory.NewEmptyAccount()
	destAccount := c.factory.NewEmptyAccount()
	srcId := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	destId := model.MustAddress(model.MustAddress(cmd.GetTransferBalance().GetDestAccountId()).AccountId())
	if err := wsv.Query(srcId, srcAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotFoundSrcAccountId, err.Error())
	}
	if err := wsv.Query(destId, destAccount); err != nil {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotFoundDestAccountId, err.Error())
	}
	if srcAccount.GetBalance()-transfer.GetBalance() < 0 {
		return errors.Wrap(core.ErrCommandExecutorTransferBalanceNotEnoughSrcAccountBalance,
			fmt.Errorf("srcAccount Amount: %d, transfer Acmount: %d", srcAccount.GetBalance(), transfer.GetBalance()).Error())
	}
	newSrcAccount := c.factory.NewAccountBuilder().
		From(srcAccount).
		Balance(srcAccount.GetBalance() - transfer.GetBalance()).
		Build()
	newDestAccount := c.factory.NewAccountBuilder().
		From(destAccount).
		Balance(destAccount.GetBalance() + transfer.GetBalance()).
		Build()
	if err := wsv.Append(srcId, newSrcAccount); err != nil {
		return err
	}
	if err := wsv.Append(destId, newDestAccount); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) CreateAccount(wsv model.ObjectFinder, cmd model.Command) error {
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	newAccount := c.factory.NewAccountBuilder().
		AccountId(id.Account() + "@" + id.Domain()).
		AccountName(id.Account()).
		PublicKeys(cmd.GetCreateAccount().GetPublicKeys()).
		Quorum(cmd.GetCreateAccount().GetQuorum()).
		Build()

	/*
		ac := c.factory.NewEmptyAccount()
		if err := wsv.Query(id, ac); err == nil {
			if ac.GetAccountId() == cmd.GetTargetId() {
				return errors.Wrap(core.ErrCommandExecutorCreateAccountAlreadyExistAccount,
					fmt.Errorf("already exist accountId : %s", id.AccountId()).Error())
			}
		}
	*/
	if err := wsv.Append(id, newAccount); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) SetQuorum(wsv model.ObjectFinder, cmd model.Command) error {
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	sq := cmd.GetSetQuorum()
	ac := c.factory.NewEmptyAccount()
	if err := wsv.Query(id, ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddBalanceNotExistAccount, err.Error())
	}
	newAc := c.factory.NewAccountBuilder().
		From(ac).
		Quorum(sq.GetQuorum()).
		Build()
	if err := wsv.Append(id, newAc); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) AddBalance(wsv model.ObjectFinder, cmd model.Command) error {
	aa := cmd.GetAddBalance()
	ac := c.factory.NewEmptyAccount()
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	if err := wsv.Query(id, ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddBalanceNotExistAccount, err.Error())
	}
	newAc := c.factory.NewAccountBuilder().
		From(ac).
		Balance(ac.GetBalance() + aa.GetBalance()).
		Build()
	if err := wsv.Append(id, newAc); err != nil {
		return err
	}
	return nil
}

func containsPublicKey(keys []model.PublicKey, key model.PublicKey) bool {
	for _, k := range keys {
		if bytes.Equal(k, key) {
			return true
		}
	}
	return false
}

func (c *CommandExecutor) AddPublicKeys(wsv model.ObjectFinder, cmd model.Command) error {
	ap := cmd.GetAddPublicKeys()
	ac := c.factory.NewEmptyAccount()
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	if err := wsv.Query(id, ac); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddPublicKeyNotExistAccount, err.Error())
	}
	if containsPublicKey(ac.GetPublicKeys(), ap.GetPublicKeys()[0]) {
		return errors.Wrapf(core.ErrCommandExecutorAddPublicKeyDuplicatePubkey,
			"duplicate key : %x", ap.GetPublicKeys())
	}
	newAc := c.factory.NewAccountBuilder().
		From(ac).
		PublicKeys(append(ac.GetPublicKeys(), ap.GetPublicKeys()[0])).
		Build()
	if err := wsv.Append(id, newAc); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) DefineStorage(wsv model.ObjectFinder, cmd model.Command) error {
	ds := cmd.GetDefineStorage()
	id := model.MustAddress(cmd.GetTargetId())

	newStorage := c.factory.NewStorageBuilder().
		From(ds.GetStorage()).
		Build()
	if err := wsv.Append(id, newStorage); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) CreateStorage(wsv model.ObjectFinder, cmd model.Command) error {
	id := model.MustAddress(cmd.GetTargetId())
	mtSt := c.factory.NewEmptyStorage()
	if err := wsv.Query(model.MustAddress("/"+id.Storage()), mtSt); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorCreateStorageNotDefinedStorage, err.Error())
	}
	if err := wsv.Append(id, mtSt); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) UpdateObject(wsv model.ObjectFinder, cmd model.Command) error {
	uo := cmd.GetUpdateObject()
	id := model.MustAddress(cmd.GetTargetId())
	mtSt := c.factory.NewEmptyStorage()
	if err := wsv.Query(id, mtSt); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorUpdateObjectNotExistWallet, err.Error())
	}
	newSt := c.factory.NewStorageBuilder().
		From(mtSt).
		Set(uo.GetKey(), uo.GetObject()).
		Build()
	if err := wsv.Append(id, newSt); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) AddObject(wsv model.ObjectFinder, cmd model.Command) error {
	ao := cmd.GetAddObject()
	id := model.MustAddress(cmd.GetTargetId())
	mtSt := c.factory.NewEmptyStorage()
	if err := wsv.Query(id, mtSt); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorAddObjectNotExistWallet, err.Error())
	}
	if o, ok := mtSt.GetObject()[ao.GetKey()]; ok {
		var newSt model.Storage
		switch o.GetType() {
		case model.ListObjectCode:
			objectList := o.GetList()
			if objectList == nil {
				objectList = make([]model.Object, 0, 1)
			}
			objectList = append(objectList, ao.GetObject())

			newSt = c.factory.NewStorageBuilder().
				From(mtSt).
				List(ao.GetKey(), objectList).
				Build()
		case model.DictObjectCode:
			return errors.Errorf("Failed AddObject UnImplement")
		default:
			return errors.Errorf("Failed AddObject type is not dict or list key: %s", ao.GetKey())
		}
		if err := wsv.Append(id, newSt); err != nil {
			return err
		}
	} else {
		return errors.Errorf("Failed AddObject is not key: %s", ao.GetKey())
	}
	return nil
}

func (c *CommandExecutor) TransferObject(wsv model.ObjectFinder, cmd model.Command) error {
	to := cmd.GetTransferObject()
	srcId := model.MustAddress(cmd.GetTargetId())
	destId := model.MustAddress(to.GetDestAccountId())
	srcSt := c.factory.NewEmptyStorage()
	destSt := c.factory.NewEmptyStorage()
	if err := wsv.Query(srcId, srcSt); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorTransferObjectNotExistSrcWallet, err.Error())
	}
	if err := wsv.Query(destId, destSt); err != nil {
		return errors.Wrapf(core.ErrCommandExecutorTransferObjectNotExistDestWallet, err.Error())
	}
	srco, ok1 := srcSt.GetObject()[to.GetKey()]
	desto, ok2 := destSt.GetObject()[to.GetKey()]
	if !ok1 || !ok2 {
		return errors.Errorf("Failed TransferObject is not key: %s", to.GetKey())
	}

	if srco.GetType() == model.ListObjectCode &&
		desto.GetType() == model.ListObjectCode {
		srcList := srco.GetList()
		destList := desto.GetList()

		if destList == nil {
			destList = make([]model.Object, 0, 1)
		}
		f := false
		for i, o := range srcList {
			if bytes.Equal(o.Hash(), to.GetObject().Hash()) {
				destList = append(destList, to.GetObject())
				srcList = append(srcList[0:i], srcList[i+1:]...)
				f = true
			}
		}
		if !f {
			return errors.Errorf("Failed TranferObject is not found object: %x", to.GetObject().Hash())
		}

		newSrcSt := c.factory.NewStorageBuilder().
			From(srcSt).
			List(to.GetKey(), srcList).
			Build()
		newDestSt := c.factory.NewStorageBuilder().
			From(destSt).
			List(to.GetKey(), destList).
			Build()

		if err := wsv.Append(srcId, newSrcSt); err != nil {
			return err
		}
		if err := wsv.Append(destId, newDestSt); err != nil {
			return err
		}
	} else {
		return errors.Errorf("Failed TranferObject type is not dict or list key: %s", to.GetKey())
	}
	return nil
}

func (c *CommandExecutor) AddPeer(wsv model.ObjectFinder, cmd model.Command) error {
	ap := cmd.GetAddPeer()
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).PeerId())
	newPeer := c.factory.NewPeer(id.Account()+"@"+id.Domain(), ap.GetAddress(), ap.GetPublicKey())

	if err := wsv.Append(id, newPeer); err != nil {
		return err
	}
	return nil
}

// TODO implent sync
func (c *CommandExecutor) ActivatePeer(wsv model.ObjectFinder, cmd model.Command) error {
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).PeerId())
	peer := c.factory.NewEmptyPeer()
	if err := wsv.Query(id, peer); err != nil {
		return err
	}
	peer.Activate()
	if err := wsv.Append(id, peer); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) SuspendPeer(wsv model.ObjectFinder, cmd model.Command) error {
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).PeerId())
	peer := c.factory.NewEmptyPeer()
	if err := wsv.Query(id, peer); err != nil {
		return err
	}
	peer.Suspend()
	if err := wsv.Append(id, peer); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) BanPeer(wsv model.ObjectFinder, cmd model.Command) error {
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).PeerId())
	peer := c.factory.NewEmptyPeer()
	if err := wsv.Query(id, peer); err != nil {
		return err
	}
	peer.Ban()
	if err := wsv.Append(id, peer); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) Consign(wsv model.ObjectFinder, cmd model.Command) error {
	cc := cmd.GetConsign()
	id := model.MustAddress(model.MustAddress(cmd.GetTargetId()).AccountId())
	ac := c.factory.NewEmptyAccount()
	if err := wsv.Query(id, ac); err != nil {
		return errors.Wrap(core.ErrCommandExecutorConsignNotFoundAccount, err.Error())
	}
	newAc := c.factory.NewAccountBuilder().
		From(ac).
		DelegatePeerId(cc.GetPeerId()).
		Build()
	if err := wsv.Append(id, newAc); err != nil {
		return err
	}
	return nil
}

func (c *CommandExecutor) CheckAndCommitProsl(wsv model.ObjectFinder, cmd model.Command) error {
	cc := cmd.GetCheckAndCommitProsl()
	targetId := model.MustAddress(cmd.GetTargetId())

	// 1. get update prosl
	updateId := model.MustAddress(c.conf.Prosl.Update.Id)
	updateSt := c.factory.NewEmptyStorage()
	if err := wsv.Query(updateId, updateSt); err != nil {
		return err
	}
	buf := updateSt.GetFromKey(core.ProslKey).GetData()
	if err := c.prosl.Unmarshal(buf); err != nil {
		return err
	}

	// 2. update prosl execute with prams + ["target_id"] = target_id
	params := cc.GetVariables()
	params[core.TargetIdKey] = c.factory.NewObjectBuilder().Address(cmd.GetTargetId())
	if check, variables, err := c.prosl.ExecuteWithParams(wsv, nil, params); err != nil {
		return errors.Errorf("variales: %+v, error: %s", variables, err.Error())
	} else if !check.GetBoolean() {
		return errors.Wrapf(core.ErrCommandExecutorCheckAndCommitProslInvalid,
			"variables: %+v", variables)
	}

	// 3. if true, targetId 's prosl setting to dest incentive or consensus or update
	proSt := c.factory.NewEmptyStorage()
	if err := wsv.Query(targetId, proSt); err != nil {
		return errors.Wrap(core.ErrCommandExecutorCheckAndCommitProslNotFound, err.Error())
	}
	t := proSt.GetFromKey(core.ProslTypeKey).GetStr()
	var destId model.Address
	switch t {
	case core.IncentiveKey:
		destId = model.MustAddress(c.conf.Prosl.Incentive.Id)
	case core.ConsensusKey:
		destId = model.MustAddress(c.conf.Prosl.Consensus.Id)
	case core.UpdateKey:
		destId = model.MustAddress(c.conf.Prosl.Update.Id)
	default:
		return errors.Errorf("not found key %s, or unexpected value: %s", core.ProslKey, t)
	}
	if err := wsv.Append(destId, proSt); err != nil {
		return err
	}
	return nil
}
