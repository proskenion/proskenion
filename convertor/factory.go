package convertor

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type ModelFactory struct {
	cryptor          core.Cryptor
	executor         core.CommandExecutor
	commandValidator core.CommandValidator
	queryValidator   core.QueryValidator
}

func NewModelFactory(cryptor core.Cryptor,
	executor core.CommandExecutor,
	cmdValidator core.CommandValidator, queryValidator core.QueryValidator) model.ModelFactory {
	return &ModelFactory{cryptor, executor, cmdValidator, queryValidator}
}

func (f *ModelFactory) NewEmptyBlock() model.Block {
	return &Block{
		&proskenion.Block{
			Payload:   &proskenion.Block_Payload{},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
	}
}

func (f *ModelFactory) NewEmptyAccount() model.Account {
	return &Account{
		f.cryptor,
		&proskenion.Account{},
	}
}

func (f *ModelFactory) NewEmptyPeer() model.Peer {
	return &Peer{
		f.cryptor,
		&proskenion.Peer{},
	}
}

func (f *ModelFactory) NewEmptyTx() model.Transaction {
	return f.NewTxBuilder().Build()
}

func (f *ModelFactory) NewEmptyQuery() model.Query {
	return f.NewQueryBuilder().Build()
}

func (f *ModelFactory) NewEmptyQueryResponse() model.QueryResponse {
	return f.NewQueryResponseBuilder().Build()
}

func (f *ModelFactory) NewBlock(height int64,
	preBlockHash model.Hash, createdTime int64,
	merkleHash model.Hash, txsHash model.Hash,
	round int32) model.Block {
	return &Block{
		&proskenion.Block{
			Payload: &proskenion.Block_Payload{
				Height:       height,
				PreBlockHash: preBlockHash,
				CreatedTime:  createdTime,
				MerkleHash:   merkleHash,
				TxsHash:      txsHash,
				Round:        round,
			},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
	}
}

func (f *ModelFactory) NewSignature(pubkey model.PublicKey, signature []byte) model.Signature {
	return &Signature{
		&proskenion.Signature{
			PublicKey: []byte(pubkey),
			Signature: signature,
		},
	}
}

func (f *ModelFactory) NewAccount(accountId string, accountName string, publicKeys []model.PublicKey, amount int64) model.Account {
	return &Account{
		f.cryptor,
		&proskenion.Account{
			AccountId:   accountId,
			AccountName: accountName,
			PublicKeys:  model.BytesListFromPublicKeys(publicKeys),
			Amount:      amount,
		},
	}
}

func (f *ModelFactory) NewPeer(address string, pubkey model.PublicKey) model.Peer {
	return &Peer{
		f.cryptor,
		&proskenion.Peer{
			Address:   address,
			PublicKey: []byte(pubkey),
		},
	}
}

func (f *ModelFactory) NewTxBuilder() model.TxBuilder {
	return &TxBuilder{
		&proskenion.Transaction{
			Payload:    &proskenion.Transaction_Payload{},
			Signatures: make([]*proskenion.Signature, 0),
		},
		f.cryptor,
		f.executor,
		f.commandValidator,
	}
}

func (f *ModelFactory) NewQueryBuilder() model.QueryBuilder {
	return &QueryBuilder{
		&proskenion.Query{
			Payload:   &proskenion.Query_Payload{},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
		f.queryValidator,
	}
}

func (f *ModelFactory) NewQueryResponseBuilder() model.QueryResponseBuilder {
	return &QueryResponseBuilder{
		&proskenion.QueryResponse{
			Payload:   &proskenion.QueryResponse_Payload{},
			Signature: &proskenion.Signature{},
		},
		f.cryptor,
	}
}

type TxBuilder struct {
	*proskenion.Transaction
	cryptor   core.Cryptor
	executor  core.CommandExecutor
	validator core.CommandValidator
}

func (t *TxBuilder) CreatedTime(time int64) model.TxBuilder {
	t.Payload.CreatedTime = time
	return t
}

func (t *TxBuilder) Transfer(srcAccountId string, destAccountId string, amount int64) model.TxBuilder {
	t.Payload.Commands = append(t.Payload.Commands,
		&proskenion.Command{
			Command: &proskenion.Command_Transfer{
				Transfer: &proskenion.Transfer{
					SrcAccountId:  srcAccountId,
					DestAccountId: destAccountId,
					Amount:        amount,
				},
			},
		})
	return t
}

func (t *TxBuilder) Build() model.Transaction {
	return &Transaction{t.Transaction,
		t.cryptor, t.executor, t.validator}
}

type QueryBuilder struct {
	*proskenion.Query
	cryptor   core.Cryptor
	validator core.QueryValidator
}

func (q *QueryBuilder) AuthorizerId(authorizerId string) model.QueryBuilder {
	q.Query.Payload.AuthorizerId = authorizerId
	return q
}

func (q *QueryBuilder) TargetId(targetId string) model.QueryBuilder {
	q.Query.Payload.TargetId = targetId
	return q
}

func (q *QueryBuilder) CreatedTime(time int64) model.QueryBuilder {
	q.Query.Payload.CreatedTime = time
	return q
}

func (q *QueryBuilder) RequestCode(code model.ObjectCode) model.QueryBuilder {
	q.Query.Payload.RequstCode = proskenion.ObjectCode(code)
	return q
}

func (q *QueryBuilder) Build() model.Query {
	return &Query{q.Query, q.cryptor, q.validator}
}

type QueryResponseBuilder struct {
	*proskenion.QueryResponse
	cryptor core.Cryptor
}

func (q *QueryResponseBuilder) Account(ac model.Account) model.QueryResponseBuilder {
	q.QueryResponse.Payload.Object = &proskenion.QueryResponse_Payload_Account{
		Account: &proskenion.Account{
			AccountId:   ac.GetAccountId(),
			AccountName: ac.GetAccountName(),
			PublicKeys:  model.BytesListFromPublicKeys(ac.GetPublicKeys()),
			Amount:      ac.GetAmount(),
		},
	}
	q.QueryResponse.Payload.ResponseCode = proskenion.ObjectCode_AccountObjectCode
	return q
}

func (q *QueryResponseBuilder) Build() model.QueryResponse {
	return &QueryResponse{q.QueryResponse, q.cryptor}
}
