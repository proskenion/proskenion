package prosl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"strings"
)

var (
	ErrProslExecuteNotExpectedType         = fmt.Errorf("Failed Prosl Execute not expected type")
	ErrProslExecuteUnExpectedCastType      = fmt.Errorf("Failed Prosl Execute Cast not expected type")
	ErrProslExecuteArgumentSize            = fmt.Errorf("Failed Prosl Execute argument size")
	ErrProslExecuteUnExpectedOperationName = fmt.Errorf("Failed Prosl Execute operation")
	ErrProslExecuteInternalErr             = fmt.Errorf("Failed Prosl Execute Internal")
	ErrProslExecuteUnknownObjectCode       = fmt.Errorf("Failed Prosl Execute Unklnown object code")
	ErrProslExecuteQueryOperatorArgument   = fmt.Errorf("Failed Prosl not enough query operator arguments")

	ErrProslExecuteSentence      = fmt.Errorf("Faile Prosl Execute Sentence error")
	ErrProslExecuteInternal      = fmt.Errorf("Failed Prosl Execute internal error")
	ErrProslExecuteUnImplemented = fmt.Errorf("Failed Prosl Execute unimplemented error")
	ErrProslExecuteAssertation   = fmt.Errorf("Failed Prosl Execute assertation error")
	ErrProslExecuteQueryVerify   = fmt.Errorf("Failed Prosl Execute query verify error")
	ErrProslExecuteQueryValidate = fmt.Errorf("Failed Prosl Execute query validate error")
)

type OperatorState int

const (
	AnotherOperator_State OperatorState = iota
	IfOperatorTrue_State
	IfOperatorFalse_State
	ElifOperatorTrue_State
	ElifOperatorFalse_State
	ReturnOperator_State
	AssertOperator_State
)

type ReturnObject struct {
	model.Object
	model.Transaction
	model.Command
}

type ProslStateValue struct {
	Variables    map[string]*ReturnObject
	ReturnObject *ReturnObject
	St           OperatorState
	ErrCode      proskenion.ErrCode
	Err          error
	Wsv          model.ObjectFinder
	Fc           model.ModelFactory
	Qc           core.Querycutor
}

func ReturnOpProslStateValue(state *ProslStateValue, st OperatorState) *ProslStateValue {
	return &ProslStateValue{
		Variables:    state.Variables,
		ReturnObject: nil,
		St:           st,
		ErrCode:      proskenion.ErrCode_NoErr,
		Err:          nil,
		Wsv:          state.Wsv,
		Fc:           state.Fc,
		Qc:           state.Qc,
	}
}

func ReturnValueProslStateValue(state *ProslStateValue, value model.Object) *ProslStateValue {
	return &ProslStateValue{
		Variables:    state.Variables,
		ReturnObject: &ReturnObject{Object: value},
		St:           AnotherOperator_State,
		ErrCode:      proskenion.ErrCode_NoErr,
		Err:          nil,
		Wsv:          state.Wsv,
		Fc:           state.Fc,
		Qc:           state.Qc,
	}
}

func ReturnTxProslStateValue(state *ProslStateValue, value model.Transaction) *ProslStateValue {
	return &ProslStateValue{
		Variables:    state.Variables,
		ReturnObject: &ReturnObject{Transaction: value},
		St:           AnotherOperator_State,
		ErrCode:      proskenion.ErrCode_NoErr,
		Err:          nil,
		Wsv:          state.Wsv,
		Fc:           state.Fc,
		Qc:           state.Qc,
	}
}

func ReturnCmdProslStateValue(state *ProslStateValue, value model.Command) *ProslStateValue {
	return &ProslStateValue{
		Variables:    state.Variables,
		ReturnObject: &ReturnObject{Command: value},
		St:           AnotherOperator_State,
		ErrCode:      proskenion.ErrCode_NoErr,
		Err:          nil,
		Wsv:          state.Wsv,
		Fc:           state.Fc,
		Qc:           state.Qc,
	}
}

func ReturnErrorProslStateValue(state *ProslStateValue, code proskenion.ErrCode, message string) *ProslStateValue {
	var err error
	switch code {
	case proskenion.ErrCode_Sentence:
		err = errors.Wrap(ErrProslExecuteSentence, message)
	case proskenion.ErrCode_UnImplemented:
		err = errors.Wrap(ErrProslExecuteUnImplemented, message)
	default:
		err = errors.Wrap(ErrProslExecuteInternal, message)
	}
	return &ProslStateValue{
		Variables:    state.Variables,
		ReturnObject: nil,
		St:           AnotherOperator_State,
		ErrCode:      code,
		Err:          err,
		Wsv:          state.Wsv,
		Fc:           state.Fc,
		Qc:           state.Qc,
	}
}

func ReturnAssertProslStateValue(state *ProslStateValue, message string) *ProslStateValue {
	return &ProslStateValue{
		St:      AssertOperator_State,
		ErrCode: proskenion.ErrCode_Assertation,
		Err:     errors.Wrap(ErrProslExecuteAssertation, message),
		Wsv:     state.Wsv,
		Fc:      state.Fc,
		Qc:      state.Qc,
	}
}

func ReturnReturnProslStateValue(state *ProslStateValue) *ProslStateValue {
	return &ProslStateValue{
		ReturnObject: state.ReturnObject,
		St:           ReturnOperator_State,
		Wsv:          state.Wsv,
		Fc:           state.Fc,
		Qc:           state.Qc,
	}
}

type Stringer interface {
	String() string
}

func ExecuteAssertType(op Stringer, o model.Object, expectedType model.ObjectCode, state *ProslStateValue) *ProslStateValue {
	if o != nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type,
			fmt.Sprintf("expected type is %s, but nil, %s", expectedType.String(), op.String()))
	}
	if o.GetType() != expectedType {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type,
			fmt.Sprintf("expected type is %s, but %s, %s", expectedType.String(), o.GetType().String(), op.String()))
	}
	return state
}

func ExecuteProsl(prosl *proskenion.Prosl, state *ProslStateValue) *ProslStateValue {
	ops := prosl.GetOps()
	for _, op := range ops {
		state = ExecuteProslOpFormula(op, state)
		if state.Err != nil {
			return state
		}
	}
	return state
}

func ExecuteProslOpFormula(op *proskenion.ProslOperator, state *ProslStateValue) *ProslStateValue {
	switch op.GetOp().(type) {
	case *proskenion.ProslOperator_SetOp:
		state = ExecuteProslSetOperator(op.GetSetOp(), state)
	case *proskenion.ProslOperator_IfOp:
		state = ExecuteProslIfOperator(op.GetIfOp(), state)
	case *proskenion.ProslOperator_ElifOp:
		state = ExecuteProslElifOperator(op.GetElifOp(), state)
	case *proskenion.ProslOperator_ElseOp:
		state = ExecuteProslElseOperator(op.GetElseOp(), state)
	case *proskenion.ProslOperator_ErrOp:
		state = ExecuteProslErrOperator(op.GetErrOp(), state)
	case *proskenion.ProslOperator_RequireOp:
		state = ExecuteProslRequireOperator(op.GetRequireOp(), state)
	case *proskenion.ProslOperator_AssertOp:
		state = ExecuteProslAssertOperator(op.GetAssertOp(), state)
	case *proskenion.ProslOperator_ReturnOp:
		state = ExecuteProslReturnOperator(op.GetReturnOp(), state)
	default:
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented, "unimlemented operator")
	}
	return state
}

func ExecuteProslSetOperator(op *proskenion.SetOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetValue(), state)
	if state.Err != nil {
		return state
	}
	state.Variables[op.GetVariableName()] = state.ReturnObject
	return ReturnOpProslStateValue(state, AnotherOperator_State)
}

func ExecuteProslIfOperator(op *proskenion.IfOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslConditionalFormula(op.GetOp(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetBoolean() {
		state = ExecuteProsl(op.GetProsl(), state)
		if state.Err != nil {
			return state
		}
		return ReturnOpProslStateValue(state, IfOperatorTrue_State)
	}
	return ReturnOpProslStateValue(state, IfOperatorFalse_State)
}

func ExecuteProslElifOperator(op *proskenion.ElifOperator, state *ProslStateValue) *ProslStateValue {
	switch state.St {
	case IfOperatorFalse_State, ElifOperatorFalse_State:
		state = ExecuteProslConditionalFormula(op.GetOp(), state)
		if state.Err != nil {
			return state
		}
		if state.ReturnObject.GetBoolean() {
			state = ExecuteProsl(op.GetProsl(), state)
			if state.Err != nil {
				return state
			}
			return ReturnOpProslStateValue(state, ElifOperatorTrue_State)
		}
		return ReturnOpProslStateValue(state, ElifOperatorFalse_State)
	case IfOperatorTrue_State, ElifOperatorTrue_State:
		return state
	}
	return ReturnErrorProslStateValue(state, proskenion.ErrCode_Sentence,
		"elif operator must have previous operator that is if or elif operator")
}

func ExecuteProslElseOperator(op *proskenion.ElseOperator, state *ProslStateValue) *ProslStateValue {
	switch state.St {
	case IfOperatorFalse_State, ElifOperatorFalse_State:
		state = ExecuteProsl(op.GetProsl(), state)
		if state.Err != nil {
			return state
		}
		return ReturnOpProslStateValue(state, AnotherOperator_State)
	case IfOperatorTrue_State, ElifOperatorTrue_State:
		return state
	}
	return ReturnErrorProslStateValue(state, proskenion.ErrCode_Sentence,
		"else operator must have previous operator that is if or elif operator")
}

func ExecuteProslErrOperator(op *proskenion.ErrCatchOperator, state *ProslStateValue) *ProslStateValue {
	if op.GetCode() == state.ErrCode {
		state = ExecuteProsl(op.GetProsl(), state)
		if state.Err != nil {
			return state
		}
	}
	return ReturnOpProslStateValue(state, AnotherOperator_State)
}

func ExecuteProslRequireOperator(op *proskenion.RequireOperator, state *ProslStateValue) *ProslStateValue {
	return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented, "Require operator is unimplemented, yet")
}

func ExecuteProslAssertOperator(op *proskenion.AssertOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslConditionalFormula(op.GetOp(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetBoolean() {
		return ReturnAssertProslStateValue(state, fmt.Sprintf("%#v", op.GetOp()))
	}
	return ReturnOpProslStateValue(state, AnotherOperator_State)
}

func ExecuteProslReturnOperator(op *proskenion.ReturnOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetOp(), state)
	if state.Err != nil {
		return state
	}
	return ReturnReturnProslStateValue(state)
}

func ExecuteProslValueOperator(op *proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	switch op.GetOp().(type) {
	case *proskenion.ValueOperator_QueryOp:
		state = ExecuteProslQueryOperator(op.GetQueryOp(), state)
	case *proskenion.ValueOperator_TxOp:
		state = ExecuteProslTxOperator(op.GetTxOp(), state)
	case *proskenion.ValueOperator_CmdOp:
		state = ExecuteProslCmdOperator(op.GetCmdOp(), state)
	case *proskenion.ValueOperator_PlusOp:
		state = ExecuteProslPlusOperator(op.GetPlusOp(), state)
	case *proskenion.ValueOperator_MinusOp:
		state = ExecuteProslMinusOperator(op.GetMinusOp(), state)
	case *proskenion.ValueOperator_MulOp:
		state = ExecuteProslMulOperator(op.GetMulOp(), state)
	case *proskenion.ValueOperator_DivOp:
		state = ExecuteProslDivOperator(op.GetDivOp(), state)
	case *proskenion.ValueOperator_ModOp:
		state = ExecuteProslModOperator(op.GetModOp(), state)
	case *proskenion.ValueOperator_OrOp:
		state = ExecuteProslOrOperator(op.GetOrOp(), state)
	case *proskenion.ValueOperator_AndOp:
		state = ExecuteProslAndOperator(op.GetAndOp(), state)
	case *proskenion.ValueOperator_XorOp:
		state = ExecuteProslXorOperator(op.GetXorOp(), state)
	case *proskenion.ValueOperator_ConcatOp:
		state = ExecuteProslConcatOperator(op.GetConcatOp(), state)
	case *proskenion.ValueOperator_ValuedOp:
		state = ExecuteProslValuedOperator(op.GetValuedOp(), state)
	case *proskenion.ValueOperator_IndexedOp:
		state = ExecuteProslIndexedOperator(op.GetIndexedOp(), state)
	case *proskenion.ValueOperator_VariableOp:
		state = ExecuteProslVariableOperator(op.GetVariableOp(), state)
	case *proskenion.ValueOperator_Object:
		state = ExecuteProslObjectOperator(op.GetObject(), state)
	case *proskenion.ValueOperator_IsDefinedOp:
		state = ExecuteProslIsDefinedOperator(op.GetIsDefinedOp(), state)
	case *proskenion.ValueOperator_VerifyOp:
		state = ExecuteProslVerifyOperator(op.GetVerifyOp(), state)
	default:
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented, "unimlemented value operator")
	}
	return state
}

func ExecuteProslQueryOperator(op *proskenion.QueryOperator, state *ProslStateValue) *ProslStateValue {
	// must : select, request code
	builder := state.Fc.NewQueryBuilder().
		Select(op.GetSelect()).
		RequestCode(model.ObjectCode(op.GetType()))
	// from
	state = ExecuteProslValueOperator(op.GetFrom(), state)
	if state.Err != nil {
		return state
	}
	builder = builder.FromId(state.ReturnObject.GetStr())

	// optional
	// authorizerId
	if op.GetAuthorizerId() != nil {
		state = ExecuteProslValueOperator(op.GetAuthorizerId(), state)
		if state.Err != nil {
			return state
		}
		builder = builder.AuthorizerId(state.ReturnObject.GetStr())
	}
	// where
	if op.GetWhere() != nil {
		state = ExecuteProslConditionalFormula(op.GetWhere(), state)
		if state.Err != nil {
			return state
		}
		b, err := state.ReturnObject.Object.Marshal()//TODO CAUTION:::!!!!! has bug
		if err != nil {
			return ReturnErrorProslStateValue(state, proskenion.ErrCode_Internal, "where conditional formula marshal error")
		}
		builder = builder.Where(b)
	}
	// order_by
	if op.GetOrderBy() != nil {
		builder = builder.OrderBy(op.GetOrderBy().GetKey(),
			model.OrderCode(op.GetOrderBy().GetOrder()))
	}
	// limit
	builder = builder.Limit(op.GetLimit())

	// query -> object
	query := builder.Build()
	if err := state.Qc.Verify(query); err != nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_QueryVerify, err.Error())
	}
	if err := state.Qc.Validate(query); err != nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_QueryValidate, err.Error())
	}
	ret, err := state.Qc.Query(query)
	if err != nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_Internal, err.Error())
	}
	if ret.GetObject().GetType() != model.ObjectCode(op.GetType()) {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type,
			fmt.Sprintf("unexpected type, expected: %d, actual: %d", op.GetType(), ret.GetObject().GetType()))
	}
	return ReturnValueProslStateValue(state, ret.GetObject())
}

func ExecuteProslTxOperator(op *proskenion.TxOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	for _, cmd := range op.GetCommands() {
		state = ExecuteProslValueOperator(cmd, state)
		if state.Err != nil {
			return state
		}
		if state.ReturnObject.Command == nil {
			return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type, "expected type: command, but not command object")
		}
		builder = builder.AppendCommand(state.ReturnObject.Command)
	}
	return ReturnTxProslStateValue(state, builder.Build())
}

func ExecuteProslCmdOperator(op *proskenion.CommandOperator, state *ProslStateValue) *ProslStateValue {
	cmdName := strings.Replace(strings.ToLower(op.CommandName), "_", "", -1)
	switch cmdName {
	case "createaccount":
		return ExecuteProslCreateAccount(op.GetParams(), state)
	case "addbalance":
		return ExecuteProslAddBalance(op.GetParams(), state)
	case "transferbalance":
		return ExecuteProslTransferBalance(op.GetParams(), state)
	case "addpublickeys":
		return ExecuteProslAddPublicKeys(op.GetParams(), state)
	case "removepublickeys":
		return ExecuteProslRemovePublicKeys(op.GetParams(), state)
	case "setqurum":
		return ExecuteProslSetQuorum(op.GetParams(), state)
	case "definestorage":
		return ExecuteProslDefineStorage(op.GetParams(), state)
	case "createstorage":
		return ExecuteProslCreateStorage(op.GetParams(), state)
	case "updateobject":
		return ExecuteProslUpdateObject(op.GetParams(), state)
	case "addobject":
		return ExecuteProslAddObject(op.GetParams(), state)
	case "transferobject":
		return ExecuteProslTransferObject(op.GetParams(), state)
	case "addpeer":
		return ExecuteProslAddPeer(op.GetParams(), state)
	case "consign":
		return ExecuteProslConsign(op.GetParams(), state)
	default:
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented, fmt.Sprintf("unimplemented command : %s", op.GetCommandName()))
	}
	return ReturnErrorProslStateValue(state, proskenion.ErrCode_Internal, "internal error")
}

type GetOpser interface {
	GetOps() []*proskenion.ValueOperator
	String() string
}

func ExecutePolynomialOperator(op GetOpser, f func(model.Object, model.Object, model.ModelFactory) model.Object, symbol string, state *ProslStateValue) *ProslStateValue {
	if len(op.GetOps()) < 2 {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_NotEnoughArgument, fmt.Sprintf("%s Operator minimum number of argument is 2, %s", symbol, op.String()))
	}
	state = ExecuteProslValueOperator(op.GetOps()[0], state)
	if state.Err != nil {
		return state
	}
	ret := state.ReturnObject.Object
	for _, o := range op.GetOps()[1:] {
		state = ExecuteProslValueOperator(o, state)
		if state.Err != nil {
			return state
		}
		ret = f(ret, state.ReturnObject.Object, state.Fc)
	}
	if ret == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_FailedOperate, op.String())
	}
	return ReturnValueProslStateValue(state, ret)
}

func ExecuteProslPlusOperator(op *proskenion.PlusOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecutePlus, "+", state)
}

func ExecuteProslMinusOperator(op *proskenion.MinusOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecuteMinus, "-", state)
}

func ExecuteProslMulOperator(op *proskenion.MultipleOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecuteMul, "*", state)
}

func ExecuteProslDivOperator(op *proskenion.DivisionOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecuteDiv, "/", state)
}

func ExecuteProslModOperator(op *proskenion.ModOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecuteMod, "%", state)
}

func ExecuteProslOrOperator(op *proskenion.OrOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecuteOr, "|", state)
}

func ExecuteProslAndOperator(op *proskenion.AndOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecuteAnd, "&", state)
}

func ExecuteProslXorOperator(op *proskenion.XorOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecuteXor, "^", state)
}

func ExecuteProslConcatOperator(op *proskenion.ConcatOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialOperator(op, ExecuteConcat, "concat", state)
}

func ExecuteProslValuedOperator(op *proskenion.ValuedOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetObject(), state)
	if state.Err != nil {
		return state
	}
	object := state.ReturnObject.Object
	switch object.GetType() {
	case model.StorageObjectCode:
		ret := object.GetStorage().GetObject()[op.GetKey()]
		state = ExecuteAssertType(op, ret, model.ObjectCode(op.GetType()), state)
		if state.Err != nil {
			return state
		}
	case model.DictObjectCode:
		ret := object.GetDict()[op.GetKey()]
		state = ExecuteAssertType(op, ret, model.ObjectCode(op.GetType()), state)
		if state.Err != nil {
			return state
		}
	case model.AccountObjectCode:
		// TODO model.Account has model.Account.Get(key string)
	case model.PeerObjectCode:
		// TODO
	default:
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented,
			fmt.Sprintf("unimplemented valued type: %s, %s", object.GetType().String(), op.String()))
	}
	return ReturnValueProslStateValue(state, state.ReturnObject.Object)
}

func ExecuteProslIndexedOperator(op *proskenion.IndexedOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslVariableOperator(op *proskenion.VariableOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslObjectOperator(op *proskenion.Object, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslIsDefinedOperator(op *proskenion.IsDefinedOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslVerifyOperator(op *proskenion.VerifyOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslConditionalFormula(op *proskenion.ConditionalFormula, state *ProslStateValue) *ProslStateValue {
	switch op.GetOp().(type) {
	case *proskenion.ConditionalFormula_Or:
		state = ExecuteProslOrFormula(op.GetOr(), state)
	case *proskenion.ConditionalFormula_And:
		state = ExecuteProslAndFormula(op.GetAnd(), state)
	case *proskenion.ConditionalFormula_Not:
		state = ExecuteProslNotFormula(op.GetNot(), state)
	case *proskenion.ConditionalFormula_Eq:
		state = ExecuteProslEqFormula(op.GetEq(), state)
	case *proskenion.ConditionalFormula_Ne:
		state = ExecuteProslNeFormula(op.GetNe(), state)
	case *proskenion.ConditionalFormula_Gt:
		state = ExecuteProslGtFormula(op.GetGt(), state)
	case *proskenion.ConditionalFormula_Ge:
		state = ExecuteProslGeFormula(op.GetGe(), state)
	case *proskenion.ConditionalFormula_Lt:
		state = ExecuteProslLtFormula(op.GetLt(), state)
	case *proskenion.ConditionalFormula_Le:
		state = ExecuteProslLeFormula(op.GetLe(), state)
	default:
	}
	return state
}

func ExecuteProslOrFormula(op *proskenion.OrFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslAndFormula(op *proskenion.AndFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslNotFormula(op *proskenion.NotFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslEqFormula(op *proskenion.EqFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslNeFormula(op *proskenion.NeFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslGtFormula(op *proskenion.GtFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslGeFormula(op *proskenion.GeFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslLtFormula(op *proskenion.LtFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslLeFormula(op *proskenion.LeFormula, state *ProslStateValue) *ProslStateValue {
	return nil
}
