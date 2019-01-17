package convertor

import (
	"fmt"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Command struct {
	*proskenion.Command
	executor  core.CommandExecutor
	validator core.CommandValidator
	cryptor   core.Cryptor
}

func (c *Command) Execute(wsv model.ObjectFinder) error {
	switch x := c.GetCommand().(type) {
	case *proskenion.Command_TransferBalance:
		return c.executor.TransferBalance(wsv, c)
	case *proskenion.Command_AddBalance:
		return c.executor.AddBalance(wsv, c)
	case *proskenion.Command_CreateAccount:
		return c.executor.CreateAccount(wsv, c)
	case *proskenion.Command_AddPublicKeys:
		return c.executor.AddPublicKeys(wsv, c)
	case *proskenion.Command_DefineStorage:
		return c.executor.DefineStorage(wsv, c)
	case *proskenion.Command_CreateStorage:
		return c.executor.CreateStorage(wsv, c)
	case *proskenion.Command_UpdateObject:
		return c.executor.UpdateObject(wsv, c)
	case *proskenion.Command_AddObject:
		return c.executor.AddObject(wsv, c)
	case *proskenion.Command_TransferObject:
		return c.executor.TransferObject(wsv, c)
	case *proskenion.Command_AddPeer:
		return c.executor.AddPeer(wsv, c)
	case *proskenion.Command_Consign:
		return c.executor.Consign(wsv, c)
	case *proskenion.Command_CheckAndCommitProsl:
		return c.executor.CheckAndCommitProsl(wsv, c)
	default:
		return fmt.Errorf("Command has unexpected type %T", x)
	}
}

func (c *Command) Validate(wsv model.ObjectFinder) error {
	switch x := c.GetCommand().(type) {
	case *proskenion.Command_TransferBalance:
		return c.validator.TransferBalance(wsv, c)
	case *proskenion.Command_AddBalance:
		return c.validator.AddBalance(wsv, c)
	case *proskenion.Command_CreateAccount:
		return c.validator.CreateAccount(wsv, c)
	case *proskenion.Command_AddPublicKeys:
		return c.validator.AddPublicKeys(wsv, c)
	case *proskenion.Command_DefineStorage:
		return c.validator.DefineStorage(wsv, c)
	case *proskenion.Command_CreateStorage:
		return c.validator.CreateStorage(wsv, c)
	case *proskenion.Command_UpdateObject:
		return c.validator.AddObject(wsv, c)
	case *proskenion.Command_TransferObject:
		return c.validator.TransferObject(wsv, c)
	case *proskenion.Command_AddPeer:
		return c.validator.AddPeer(wsv, c)
	case *proskenion.Command_Consign:
		return c.validator.Consign(wsv, c)
	case *proskenion.Command_CheckAndCommitProsl:
		return c.validator.CheckAndCommitProsl(wsv, c)

	default:
		return fmt.Errorf("Command has unexpected type %T", x)
	}
}

func (c *Command) GetTransferBalance() model.TransferBalance {
	return c.Command.GetTransferBalance()
}

func (c *Command) GetCreateAccount() model.CreateAccount {
	return &CreateAccount{c.Command.GetCreateAccount()}
}

func (c *Command) GetAddBalance() model.AddBalance {
	return c.Command.GetAddBalance()
}

func (c *Command) GetAddPublicKeys() model.AddPublicKeys {
	return &AddPublicKeys{c.Command.GetAddPublicKeys()}
}

func (c *Command) GetRemovePublicKeys() model.RemovePublicKeys {
	return &RemovePublicKeys{c.Command.GetRemovePublicKeys()}
}

func (c *Command) GetSetQuorum() model.SetQuroum {
	return c.Command.GetSetQurum()
}

type CreateAccount struct {
	*proskenion.CreateAccount
}

func (c *CreateAccount) GetPublicKeys() []model.PublicKey {
	if c.CreateAccount == nil {
		return nil
	}
	return model.PublicKeysFromBytesSlice(c.CreateAccount.GetPublicKeys())
}

type AddPublicKeys struct {
	*proskenion.AddPublicKeys
}

func (c *AddPublicKeys) GetPublicKeys() []model.PublicKey {
	if c.AddPublicKeys == nil {
		return nil
	}
	return model.PublicKeysFromBytesSlice(c.AddPublicKeys.GetPublicKeys())
}

type RemovePublicKeys struct {
	*proskenion.RemovePublicKeys
}

func (c *RemovePublicKeys) GetPublicKeys() []model.PublicKey {
	if c.RemovePublicKeys == nil {
		return nil
	}
	return model.PublicKeysFromBytesSlice(c.RemovePublicKeys.GetPublicKeys())
}

type DefineStorage struct {
	c core.Cryptor
	e core.CommandExecutor
	v core.CommandValidator
	*proskenion.DefineStorage
}

func (c *DefineStorage) GetStorage() model.Storage {
	return &Storage{c.c, c.e, c.v, c.Storage}
}

func (c *Command) GetDefineStorage() model.DefineStorage {
	return &DefineStorage{c.cryptor, c.executor, c.validator, c.Command.GetDefineStorage()}
}

func (c *Command) GetCreateStorage() model.CreateStorage {
	return c.Command.GetCreateStorage()
}

type UpdateObject struct {
	c core.Cryptor
	e core.CommandExecutor
	v core.CommandValidator
	*proskenion.UpdateObject
}

func (c *UpdateObject) GetObject() model.Object {
	return &Object{c.c, c.e, c.v, c.Object}
}

func (c *Command) GetUpdateObject() model.UpdateObject {
	return &UpdateObject{c.cryptor, c.executor, c.validator, c.Command.GetUpdateObject()}
}

type AddObject struct {
	c core.Cryptor
	e core.CommandExecutor
	v core.CommandValidator
	*proskenion.AddObject
}

func (c *AddObject) GetObject() model.Object {
	return &Object{c.c, c.e, c.v, c.AddObject.Object}
}

func (c *Command) GetAddObject() model.AddObject {
	return &AddObject{c.cryptor, c.executor, c.validator, c.Command.GetAddObject()}
}

type TransferObject struct {
	c core.Cryptor
	e core.CommandExecutor
	v core.CommandValidator
	*proskenion.TransferObject
}

func (c *TransferObject) GetObject() model.Object {
	return &Object{c.c, c.e, c.v, c.TransferObject.Object}
}

func (c *Command) GetTransferObject() model.TransferObject {
	return &TransferObject{c.cryptor, c.executor, c.validator, c.Command.GetTransferObject()}
}

func (c *Command) GetAddPeer() model.AddPeer {
	return c.Command.GetAddPeer()
}

func (c *Command) GetConsign() model.Consign {
	return c.Command.GetConsign()
}

type CheckAndCommitProsl struct {
	c core.Cryptor
	e core.CommandExecutor
	v core.CommandValidator
	*proskenion.CheckAndCommitProsl
}

func (c *Command) GetCheckAndCommitProsl() model.CheckAndCommitProsl {
	return &CheckAndCommitProsl{}
}

func (c *CheckAndCommitProsl) GetVariables() map[string]model.Object {
	return ObjectMapsFromProslObjectMaps(c.c, c.e, c.v, c.CheckAndCommitProsl.GetVariables())
}
