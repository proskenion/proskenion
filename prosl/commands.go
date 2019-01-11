package prosl

import (
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

func ObjectListToPublicKeys(list []model.Object) []model.PublicKey {
	ret := make([]model.PublicKey, 0, len(list))
	for _, o := range list {
		ret = append(ret, o.GetData())
	}
	return ret
}

func ExecuteProslCreateAccount(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["account_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["public_keys"], state)
	if state.Err != nil {
		return state
	}
	publickeys := ObjectListToPublicKeys(state.ReturnObject.GetList())

	state = ExecuteProslValueOperator(params["quorum"], state)
	if state.Err != nil {
		return state
	}
	quorum := state.ReturnObject.GetI32()
	return ReturnCmdProslStateValue(state,
		builder.CreateAccount(authorizerId, targetId, publickeys, quorum).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslTransferBalance(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["account_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["dest_account_id"], state)
	if state.Err != nil {
		return state
	}
	destAccountId := state.ReturnObject.GetStr()
	state = ExecuteProslValueOperator(params["balance"], state)
	if state.Err != nil {
		return state
	}
	balance := state.ReturnObject.GetI64()

	return ReturnCmdProslStateValue(state,
		builder.TransferBalance(authorizerId, targetId, destAccountId, balance).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslAddBalance(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["account_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["balance"], state)
	if state.Err != nil {
		return state
	}
	balance := state.ReturnObject.GetI64()
	return ReturnCmdProslStateValue(state,
		builder.AddBalance(authorizerId, targetId, balance).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslAddPublicKeys(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["account_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["public_keys"], state)
	if state.Err != nil {
		return state
	}
	publicKeys := ObjectListToPublicKeys(state.ReturnObject.GetList())
	return ReturnCmdProslStateValue(state,
		builder.AddPublicKeys(authorizerId, targetId, publicKeys).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslRemovePublicKeys(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["account_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["public_keys"], state)
	if state.Err != nil {
		return state
	}
	publicKeys := ObjectListToPublicKeys(state.ReturnObject.GetList())
	return ReturnCmdProslStateValue(state,
		builder.RemovePublicKeys(authorizerId, targetId, publicKeys).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslSetQuorum(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["account_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["quorum"], state)
	if state.Err != nil {
		return state
	}
	quorum := state.ReturnObject.GetI32()
	return ReturnCmdProslStateValue(state,
		builder.SetQuorum(authorizerId, targetId, quorum).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslDefineStorage(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["storage_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["storage"], state)
	if state.Err != nil {
		return state
	}
	storage := state.ReturnObject.GetStorage()
	return ReturnCmdProslStateValue(state,
		builder.DefineStorage(authorizerId, targetId, storage).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslCreateStorage(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["wallet_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	return ReturnCmdProslStateValue(state,
		builder.CreateStorage(authorizerId, targetId).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslUpdateObject(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["wallet_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["key"], state)
	if state.Err != nil {
		return state
	}
	key := state.ReturnObject.GetStr()

	state = ExecuteProslValueOperator(params["object"], state)
	if state.Err != nil {
		return state
	}
	object := state.ReturnObject.Object

	return ReturnCmdProslStateValue(state,
		builder.UpdateObject(authorizerId, targetId, key, object).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslAddObject(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["wallet_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["key"], state)
	if state.Err != nil {
		return state
	}
	key := state.ReturnObject.GetStr()

	state = ExecuteProslValueOperator(params["object"], state)
	if state.Err != nil {
		return state
	}
	object := state.ReturnObject.Object

	return ReturnCmdProslStateValue(state,
		builder.AddObject(authorizerId, targetId, key, object).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslTransferObject(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["src_wallet_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["key"], state)
	if state.Err != nil {
		return state
	}
	key := state.ReturnObject.GetStr()

	state = ExecuteProslValueOperator(params["dest_wallet_id"], state)
	if state.Err != nil {
		return state
	}
	destAccountId := state.ReturnObject.GetStr()

	state = ExecuteProslValueOperator(params["object"], state)
	if state.Err != nil {
		return state
	}
	object := state.ReturnObject.Object

	return ReturnCmdProslStateValue(state,
		builder.TransferObject(authorizerId, targetId, key, destAccountId, object).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslAddPeer(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["peer_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["address"], state)
	if state.Err != nil {
		return state
	}
	address := state.ReturnObject.GetStr()

	state = ExecuteProslValueOperator(params["public_key"], state)
	if state.Err != nil {
		return state
	}
	publicKey := state.ReturnObject.GetData()

	return ReturnCmdProslStateValue(state,
		builder.AddPeer(authorizerId, targetId, address, publicKey).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslConsign(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(params["authorizer_id"], state)
	if state.Err != nil {
		return state
	}
	authorizerId := state.ReturnObject.GetAddress()
	state = ExecuteProslValueOperator(params["account_id"], state)
	if state.Err != nil {
		return state
	}
	targetId := state.ReturnObject.GetAddress()

	state = ExecuteProslValueOperator(params["peer_id"], state)
	if state.Err != nil {
		return state
	}
	peerId := state.ReturnObject.GetStr()
	return ReturnCmdProslStateValue(state,
		builder.Consign(authorizerId, targetId, peerId).Build().GetPayload().GetCommands()[0])
}
