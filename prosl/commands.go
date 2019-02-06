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

func ReturnCmdProslStateValue(state *ProslStateValue, cmd model.Command) *ProslStateValue {
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Command(cmd))
}

func ExecuteProslCreateAccount(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId string
	var targetId string
	var publicKeys []model.PublicKey
	var quorum int32
	for key, value := range params {
		switch key {
		case "authorizer_id", "authorizer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "account_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()
		case "public_keys", "keys":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			publicKeys = ObjectListToPublicKeys(state.ReturnObject.GetList())
		case "quorum":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			quorum = state.ReturnObject.GetI32()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.CreateAccount(authorizerId, targetId, publicKeys, quorum).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslTransferBalance(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId, destAccountId string
	var balance int64
	for key, value := range params {
		switch key {
		case "authorizer_id", "authorizer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "account_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()

		case "dest_account_id", "dest", "dest_account":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			destAccountId = state.ReturnObject.GetStr()
		case "balance":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			balance = state.ReturnObject.GetI64()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.TransferBalance(authorizerId, targetId, destAccountId, balance).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslAddBalance(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId string
	var balance int64
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "account_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()
		case "balance":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			balance = state.ReturnObject.GetI64()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.AddBalance(authorizerId, targetId, balance).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslAddPublicKeys(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId string
	var publicKeys []model.PublicKey
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "account_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()
		case "public_keys", "keys":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			publicKeys = ObjectListToPublicKeys(state.ReturnObject.GetList())
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.AddPublicKeys(authorizerId, targetId, publicKeys).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslRemovePublicKeys(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId string
	var publicKeys []model.PublicKey
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "account_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()

		case "public_keys", "keys":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			publicKeys = ObjectListToPublicKeys(state.ReturnObject.GetList())
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.RemovePublicKeys(authorizerId, targetId, publicKeys).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslSetQuorum(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId string
	var quorum int32
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "account_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()

		case "quorum":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			quorum = state.ReturnObject.GetI32()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.SetQuorum(authorizerId, targetId, quorum).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslDefineStorage(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId string
	var storage model.Storage
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "storage_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()

		case "storage":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			storage = state.ReturnObject.GetStorage()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.DefineStorage(authorizerId, targetId, storage).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslCreateStorage(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId string
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "wallet_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.CreateStorage(authorizerId, targetId).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslUpdateObject(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId, k string
	var object model.Object
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "wallet_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()
		case "key":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			k = state.ReturnObject.GetStr()
		case "object":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			object = state.ReturnObject
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.UpdateObject(authorizerId, targetId, k, object).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslAddObject(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId, key string
	var object model.Object
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "wallet_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()

		case "key":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			key = state.ReturnObject.GetStr()

		case "object":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			object = state.ReturnObject
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.AddObject(authorizerId, targetId, key, object).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslTransferObject(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId, key, destAccountId string
	var object model.Object
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "src_wallet_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()

		case "key":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			key = state.ReturnObject.GetStr()

		case "dest_wallet_id", "dest", "dest_wallet":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			destAccountId = state.ReturnObject.GetStr()

		case "object":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			object = state.ReturnObject
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.TransferObject(authorizerId, targetId, key, destAccountId, object).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslAddPeer(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId, address string
	var publicKey model.PublicKey
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()

		case "peer_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()

		case "address", "ip":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			address = state.ReturnObject.GetStr()

		case "public_key", "key":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			publicKey = state.ReturnObject.GetData()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.AddPeer(authorizerId, targetId, address, publicKey).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslConsign(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId, peerId string
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()
		case "account_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()

		case "peer_id", "peer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			peerId = state.ReturnObject.GetAddress()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.Consign(authorizerId, targetId, peerId).Build().GetPayload().GetCommands()[0])
}

func ExecuteProslActivatePeer(params map[string]*proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	var authorizerId, targetId string
	for key, value := range params {
		switch key {
		case "authorizer_id", "authoirzer":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			authorizerId = state.ReturnObject.GetAddress()

		case "peer_id", "target_id", "target":
			state = ExecuteProslValueOperator(value, state)
			if state.Err != nil {
				return state
			}
			targetId = state.ReturnObject.GetAddress()
		}
	}
	return ReturnCmdProslStateValue(state,
		builder.ActivatePeer(authorizerId, targetId).Build().GetPayload().GetCommands()[0])
}
