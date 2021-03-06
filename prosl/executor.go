package prosl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/convertor"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"github.com/proskenion/proskenion/query"
	"github.com/satellitex/pagerank"
	"sort"
	"strings"
)

var (
	ErrProslExecuteSentence              = fmt.Errorf("Faile Prosl Execute Sentence error")
	ErrProslExecuteInternal              = fmt.Errorf("Failed Prosl Execute internal error")
	ErrProslExecuteUnImplemented         = fmt.Errorf("Failed Prosl Execute unimplemented error")
	ErrProslExecuteAssertation           = fmt.Errorf("Failed Prosl Execute assertation error")
	ErrProslExecuteQueryVerify           = fmt.Errorf("Failed Prosl Execute query verify error")
	ErrProslExecuteQueryValidate         = fmt.Errorf("Failed Prosl Execute query validate error")
	ErrProslExecuteType                  = fmt.Errorf("Failed Prosl Execute type error")
	ErrProslExecuteNotEnoughArgument     = fmt.Errorf("Failed Prosl Execute not enough argument")
	ErrProslExecuteFailedOperate         = fmt.Errorf("Failed Prosl Execute failed operate")
	ErrProslExecuteUnExpectedReturnValue = fmt.Errorf("Failed Prosl Execute unexpected return value")
	ErrProslExecuteOutOfRange            = fmt.Errorf("Failed Prosl Execute out of range")
	ErrProslExecuteUndefined             = fmt.Errorf("Failed Prosl EXecute undefined")
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

type ProslConstState struct {
	Variables map[string]model.Object
	Fc        model.ModelFactory
	Wsv       model.ObjectFinder
	Qc        core.Querycutor
	C         core.Cryptor
}

type ProslStateValue struct {
	*ProslConstState
	ReturnObject model.Object
	St           OperatorState
	ErrCode      proskenion.ErrCode
	Err          error
}

func InitProslStateValue(fc model.ModelFactory, wsv model.ObjectFinder, top model.Block, c core.Cryptor, conf *config.Config) *ProslStateValue {
	qc := struct {
		core.QueryProcessor
		core.QueryValidator
		core.QueryVerifier
	}{query.NewQueryProcessor(fc, conf), query.NewQueryValidator(fc, conf), query.NewQueryVerifier()}
	variables := make(map[string]model.Object)
	if top != nil {
		variables["top"] = fc.NewObjectBuilder().Block(top)
	}
	return &ProslStateValue{
		ProslConstState: &ProslConstState{
			Fc:        fc,
			Wsv:       wsv,
			Qc:        qc,
			C:         c,
			Variables: variables,
		},
		ReturnObject: nil,
		St:           AnotherOperator_State,
		ErrCode:      proskenion.ErrCode_NoErr,
		Err:          nil,
	}
}

func InitProslStateValueWithPrams(fc model.ModelFactory, wsv model.ObjectFinder, top model.Block, c core.Cryptor, conf *config.Config, params map[string]model.Object) *ProslStateValue {
	qc := struct {
		core.QueryProcessor
		core.QueryValidator
		core.QueryVerifier
	}{query.NewQueryProcessor(fc, conf), query.NewQueryValidator(fc, conf), query.NewQueryVerifier()}
	variables := make(map[string]model.Object)
	// params setting
	for key, value := range params {
		variables[key] = value
	}
	if top != nil {
		variables["top"] = fc.NewObjectBuilder().Block(top)
	}
	return &ProslStateValue{
		ProslConstState: &ProslConstState{
			Fc:        fc,
			Wsv:       wsv,
			Qc:        qc,
			C:         c,
			Variables: variables,
		},
		ReturnObject: nil,
		St:           AnotherOperator_State,
		ErrCode:      proskenion.ErrCode_NoErr,
		Err:          nil,
	}
}

func ReturnOpProslStateValue(state *ProslStateValue, st OperatorState) *ProslStateValue {
	if state.St == ReturnOperator_State {
		return state
	}
	return &ProslStateValue{
		ProslConstState: state.ProslConstState,
		ReturnObject:    nil,
		St:              st,
		ErrCode:         proskenion.ErrCode_NoErr,
		Err:             nil,
	}
}

func ReturnProslStateValue(state *ProslStateValue, value model.Object) *ProslStateValue {
	return &ProslStateValue{
		ProslConstState: state.ProslConstState,
		ReturnObject:    value,
		St:              AnotherOperator_State,
		ErrCode:         proskenion.ErrCode_NoErr,
		Err:             nil,
	}
}

func ReturnTxProslStateValue(state *ProslStateValue, value model.Transaction) *ProslStateValue {
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Transaction(value))
}

func ReturnUnExpectedRetrunValue(state *ProslStateValue, exp Stringer, act Stringer, op Stringer) *ProslStateValue {
	return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "Return Object type expected: %s, but %s\n %s", exp.String(), act.String(), op.String())
}

func ReturnErrObjectCodeRetrunValue(state *ProslStateValue, exp Stringer, act Stringer, op Stringer) *ProslStateValue {
	return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type, "Expected type: %s, but %s\n%s", exp.String(), act.String(), op.String())
}

func ReturnErrorProslStateValue(state *ProslStateValue, code proskenion.ErrCode, format string, a ...interface{}) *ProslStateValue {
	message := fmt.Sprintf(format, a...)
	var err error
	switch code {
	case proskenion.ErrCode_Sentence:
		err = errors.Wrap(ErrProslExecuteSentence, message)
	case proskenion.ErrCode_UnImplemented:
		err = errors.Wrap(ErrProslExecuteUnImplemented, message)
	case proskenion.ErrCode_Assertation:
		err = errors.Wrap(ErrProslExecuteAssertation, message)
	case proskenion.ErrCode_QueryVerify:
		err = errors.Wrap(ErrProslExecuteQueryVerify, message)
	case proskenion.ErrCode_QueryValidate:
		err = errors.Wrap(ErrProslExecuteQueryValidate, message)
	case proskenion.ErrCode_Type:
		err = errors.Wrap(ErrProslExecuteType, message)
	case proskenion.ErrCode_NotEnoughArgument:
		err = errors.Wrap(ErrProslExecuteNotEnoughArgument, message)
	case proskenion.ErrCode_FailedOperate:
		err = errors.Wrap(ErrProslExecuteFailedOperate, message)
	case proskenion.ErrCode_UnExpectedReturnValue:
		err = errors.Wrap(ErrProslExecuteUnExpectedReturnValue, message)
	case proskenion.ErrCode_OutOfRange:
		err = errors.Wrap(ErrProslExecuteOutOfRange, message)
	case proskenion.ErrCode_Undefined:
		err = errors.Wrap(ErrProslExecuteUndefined, message)
	default:
		err = errors.Wrap(ErrProslExecuteInternal, message)
	}
	return &ProslStateValue{
		ProslConstState: state.ProslConstState,
		ReturnObject:    nil,
		St:              AnotherOperator_State,
		ErrCode:         code,
		Err:             err,
	}
}

func ReturnAssertProslStateValue(state *ProslStateValue, message string) *ProslStateValue {
	return &ProslStateValue{
		ProslConstState: state.ProslConstState,
		St:              AssertOperator_State,
		ErrCode:         proskenion.ErrCode_Assertation,
		Err:             errors.Wrap(ErrProslExecuteAssertation, message),
	}
}

func ReturnReturnProslStateValue(state *ProslStateValue) *ProslStateValue {
	return &ProslStateValue{
		ProslConstState: state.ProslConstState,
		ReturnObject:    state.ReturnObject,
		St:              ReturnOperator_State,
	}
}

type Stringer interface {
	String() string
}

func ExecuteAssertType(op Stringer, o model.Object, expectedType model.ObjectCode, state *ProslStateValue) *ProslStateValue {
	if o == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type,
			"expected type is %s, but nil, %s", expectedType.String(), op.String())
	}
	if o.GetType() != expectedType {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type,
			"expected type is %s, but %s, %s", expectedType.String(), o.GetType().String(), op.String())
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
		if state.St == ReturnOperator_State {
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
	case *proskenion.ProslOperator_EachOp:
		state = ExecuteProslEachOperator(op.GetEachOp(), state)
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

func ExecuteProslEachOperator(op *proskenion.EachOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetList(), state)
	if state.Err != nil {
		return state
	}

	list := state.ReturnObject.GetList()
	if list == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "unexpected return type, expected List. %+v", op)
	}

	for _, o := range list {
		state.Variables[op.VariableName] = o
		state = ExecuteProsl(op.GetDo(), state)
		if state.Err != nil {
			return state
		}
	}
	return ReturnOpProslStateValue(state, AnotherOperator_State)
}

func ExecuteProslValueOperator(op *proskenion.ValueOperator, state *ProslStateValue) *ProslStateValue {
	switch op.GetOp().(type) {
	case *proskenion.ValueOperator_QueryOp:
		state = ExecuteProslQueryOperator(op.GetQueryOp(), state)
	case *proskenion.ValueOperator_TxOp:
		state = ExecuteProslTxOperator(op.GetTxOp(), state)
	case *proskenion.ValueOperator_CmdOp:
		state = ExecuteProslCmdOperator(op.GetCmdOp(), state)
	case *proskenion.ValueOperator_StorageOp:
		state = ExecuteProslStorageOperator(op.GetStorageOp(), state)
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
	case *proskenion.ValueOperator_PageRankOp:
		state = ExecuteProslPageRankOperator(op.GetPageRankOp(), state)
	case *proskenion.ValueOperator_ListOp:
		state = ExecuteProslListOperator(op.GetListOp(), state)
	case *proskenion.ValueOperator_MapOp:
		state = ExecuteProslMapOperator(op.GetMapOp(), state)
	case *proskenion.ValueOperator_CastOp:
		state = ExecuteProslCastOperator(op.GetCastOp(), state)
	case *proskenion.ValueOperator_ListComprehensionOp:
		state = ExecuteProslListComprehensionOperator(op.GetListComprehensionOp(), state)
	case *proskenion.ValueOperator_SortOp:
		state = ExecuteProslSortOperator(op.GetSortOp(), state)
	case *proskenion.ValueOperator_SliceOp:
		state = ExecuteProslSliceOperator(op.GetSliceOp(), state)
	case *proskenion.ValueOperator_LenOp:
		state = ExecuteProslLenOperator(op.GetLenOp(), state)
	default:
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented, "unimlemented value operator, %s", op.String())
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
	builder = builder.FromId(state.ReturnObject.GetAddress())

	// authorizer
	state = ExecuteProslValueOperator(op.GetAuthorizerId(), state)
	if state.Err != nil {
		return state
	}
	builder = builder.AuthorizerId(state.ReturnObject.GetAddress())

	// where
	if op.GetWhere() != nil {
		state = ExecuteProslValueOperator(op.GetWhere(), state)
		if state.Err != nil {
			return state
		}
		builder = builder.Where(state.ReturnObject.GetStr())
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
	// WIP : no validate
	/*
		if err := state.Qc.Validate(query); err != nil {
			return ReturnErrorProslStateValue(state, proskenion.ErrCode_QueryValidate, err.Error())
		}
	*/
	ret, err := state.Qc.Query(state.Wsv, query)
	if err != nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_Internal, err.Error())
	}
	if ret.GetObject().GetType() != model.ObjectCode(op.GetType()) {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type,
			fmt.Sprintf("unexpected type, expected: %s, actual: %s", op.GetType().String(), ret.GetObject().GetType().String()))
	}
	return ReturnProslStateValue(state, ret.GetObject())
}

func ExecuteProslTxOperator(op *proskenion.TxOperator, state *ProslStateValue) *ProslStateValue {
	builder := state.Fc.NewTxBuilder()
	state = ExecuteProslValueOperator(op.GetCommands(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetType() != model.ListObjectCode {
		return ReturnUnExpectedRetrunValue(state, model.ListObjectCode, state.ReturnObject.GetType(), op)
	}
	list := state.ReturnObject.GetList()
	for _, o := range list {
		if o.GetType() != model.CommandObjectCode {
			return ReturnErrObjectCodeRetrunValue(state, model.CommandObjectCode, state.ReturnObject.GetType(), op)
		}
		builder = builder.AppendCommand(o.GetCommand())
	}
	builder = builder.CreatedTime(0)
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
	case "activatepeer":
		return ExecuteProslActivatePeer(op.GetParams(), state)
	case "forceupdate", "forceupdatestorage", "updatestorage":
		return ExecuteProslForceUpdateStorage(op.GetParams(), state)
	default:
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented, fmt.Sprintf("unimplemented command : %s, %s", op.GetCommandName(), op.String()))
	}
}

func ExecuteProslStorageOperator(op *proskenion.StorageOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslMapOperator(op.GetObject(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetDict() == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "expected dict, but %s, %s", state.ReturnObject.GetType(), op.String())
	}
	storage := state.Fc.NewStorageBuilder().FromMap(state.ReturnObject.GetDict()).Build()
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Storage(storage))
}

type GetOpser interface {
	GetOps() []*proskenion.ValueOperator
	String() string
}

func ExecutePolynomiaValueOperator(op GetOpser, f func(model.Object, model.Object, model.ModelFactory) model.Object, symbol string, state *ProslStateValue) *ProslStateValue {
	if len(op.GetOps()) < 2 {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_NotEnoughArgument, fmt.Sprintf("%s Operator minimum number of argument is 2, %s", symbol, op.String()))
	}
	state = ExecuteProslValueOperator(op.GetOps()[0], state)
	if state.Err != nil {
		return state
	}
	ret := state.ReturnObject
	for _, o := range op.GetOps()[1:] {
		state = ExecuteProslValueOperator(o, state)
		if state.Err != nil {
			return state
		}
		ret = f(ret, state.ReturnObject, state.Fc)
	}
	if ret == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_FailedOperate, "symbol: %s, %s", symbol, op.String())
	}
	return ReturnProslStateValue(state, ret)
}

func ExecuteProslPlusOperator(op *proskenion.PlusOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecutePlus, "+", state)
}

func ExecuteProslMinusOperator(op *proskenion.MinusOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteMinus, "-", state)
}

func ExecuteProslMulOperator(op *proskenion.MultipleOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteMul, "*", state)
}

func ExecuteProslDivOperator(op *proskenion.DivisionOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteDiv, "/", state)
}

func ExecuteProslModOperator(op *proskenion.ModOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteMod, "%", state)
}

func ExecuteProslOrOperator(op *proskenion.OrOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteOr, "or", state)
}

func ExecuteProslAndOperator(op *proskenion.AndOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteAnd, "and", state)
}

func ExecuteProslXorOperator(op *proskenion.XorOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteXor, "xor", state)
}

func ExecuteProslConcatOperator(op *proskenion.ConcatOperator, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteConcat, "concat", state)
}

func ExecuteProslValuedOperator(op *proskenion.ValuedOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetObject(), state)
	if state.Err != nil {
		return state
	}
	object := state.ReturnObject
	if object == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, fmt.Sprintf("expected return Object, but: %#v, %s", state.ReturnObject, op.String()))
	}

	var ret model.Object
	switch object.GetType() {
	case model.StorageObjectCode:
		ret = object.GetStorage().GetFromKey(op.GetKey())
		state = ExecuteAssertType(op, ret, model.ObjectCode(op.GetType()), state)
		if state.Err != nil {
			return state
		}
	case model.DictObjectCode:
		ret = object.GetDict()[op.GetKey()]
		state = ExecuteAssertType(op, ret, model.ObjectCode(op.GetType()), state)
		if state.Err != nil {
			return state
		}
	case model.AccountObjectCode:
		ret = object.GetAccount().GetFromKey(op.GetKey())
		state = ExecuteAssertType(op, ret, model.ObjectCode(op.GetType()), state)
		if state.Err != nil {
			return state
		}
	case model.PeerObjectCode:
		ret = object.GetPeer().GetFromKey(op.GetKey())
		state = ExecuteAssertType(op, ret, model.ObjectCode(op.GetType()), state)
		if state.Err != nil {
			return state
		}
	case model.BlockObjectCode:
		ret = object.GetBlock().GetFromKey(op.GetKey())
		state = ExecuteAssertType(op, ret, model.ObjectCode(op.GetType()), state)
		if state.Err != nil {
			return state
		}
	default:
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented,
			fmt.Sprintf("unimplemented valued type: %s, %s", object.GetType().String(), op.String()))
	}
	// already asserted ret == nil check
	return ReturnProslStateValue(state, ret)
}

func ExecuteProslIndexedOperator(op *proskenion.IndexedOperator, state *ProslStateValue) *ProslStateValue {
	// list
	state = ExecuteProslValueOperator(op.GetObject(), state)
	if state.Err != nil {
		return state
	}
	object := state.ReturnObject
	if object == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, fmt.Sprintf("expected return Object, but: %#v, %s", state.ReturnObject, op.String()))
	}
	if object.GetType() != model.ListObjectCode {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented,
			fmt.Sprintf("unimplemented indexed type: %s, %s", object.GetType().String(), op.String()))
	}

	// int32
	state = ExecuteProslValueOperator(op.GetIndex(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetType() != model.Int32ObjectCode {
		ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, fmt.Sprintf("expected return Int32, but: %#v, %s", state.ReturnObject, op.String()))
	}
	index := state.ReturnObject.GetI32()
	if len(object.GetList()) <= int(index) {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_OutOfRange,
			"list object length is %d, but index is %d, %s", len(object.GetList()), index, op)
	}

	ret := object.GetList()[index]
	state = ExecuteAssertType(op, ret, model.ObjectCode(op.GetType()), state)
	if state.Err != nil {
		return state
	}
	return ReturnProslStateValue(state, ret)
}

func ExecuteProslVariableOperator(op *proskenion.VariableOperator, state *ProslStateValue) *ProslStateValue {
	if ret, ok := state.Variables[op.GetVariableName()]; ok {
		return ReturnProslStateValue(state, ret)
	}
	return ReturnErrorProslStateValue(state, proskenion.ErrCode_Undefined,
		fmt.Sprintf("undefined variable name: %s, %s", op.GetVariableName(), op.String()))
}

func ExecuteProslObjectOperator(op *proskenion.Object, state *ProslStateValue) *ProslStateValue {
	object := state.Fc.NewObjectBuilder().Build()
	object.(*convertor.Object).Object = op
	return ReturnProslStateValue(state, object)
}

func ExecuteProslIsDefinedOperator(op *proskenion.IsDefinedOperator, state *ProslStateValue) *ProslStateValue {
	if _, ok := state.Variables[op.GetVariableName()]; ok {
		return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Bool(true))
	}
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Bool(false))
}

func ExecuteProslVerifyOperator(op *proskenion.VerifyOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetSig(), state)
	if state.Err != nil {
		return state
	}
	sig := state.ReturnObject.GetSig()
	if sig == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "unexpected return object. expected signature type.")
	}
	state = ExecuteProslValueOperator(op.GetHash(), state)
	if state.Err != nil {
		return state
	}
	hasher := state.ReturnObject
	if hasher == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "unexpected return object. expected bytes data.")
	}
	if err := state.C.Verify(sig.GetPublicKey(), hasher, sig.GetSignature()); err != nil {
		return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Bool(false))
	}
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Bool(true))
}

func ExecuteProslPageRankOperator(op *proskenion.PageRankOperator, state *ProslStateValue) *ProslStateValue {
	// storage
	state = ExecuteProslValueOperator(op.GetStorages(), state)
	if state.Err != nil {
		return state
	}
	storages := state.ReturnObject.GetList()
	if storages == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "unexpected return object. expected List type.")
	}

	// toKey
	state = ExecuteProslValueOperator(op.GetToKey(), state)
	if state.Err != nil {
		return state
	}
	toKey := state.ReturnObject.GetStr()
	if toKey == "" {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "unexpected return object. expected String type.")
	}

	// outName
	state = ExecuteProslValueOperator(op.GetOutName(), state)
	if state.Err != nil {
		return state
	}
	outName := state.ReturnObject.GetStr()
	if outName == "" {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "unexpected return object. expected String type.")
	}

	graph := pagerank.New()
	for _, o := range storages {
		st := o.GetStorage()
		if st == nil {
			return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "unexpected return object. expected storage type. : %s", op.String())
		}
		edges := st.GetFromKey(toKey).GetList()
		if len(edges) == 0 {
			return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type, "edge: \"%s\", unexpected type. expected list type. : %s", toKey, op.String())
		}
		stId := model.MustAddress(st.GetId())
		for _, o := range edges {
			to := o.GetAddress()
			if to == "" {
				return ReturnErrorProslStateValue(state, proskenion.ErrCode_Type, "edge: \"%s\", unexpected type. expected address type. : %s", toKey, op.String())
			}
			toId := model.MustAddress(to)
			graph.Link(fmt.Sprintf("%s@%s", stId.Account(), stId.Domain()),
				fmt.Sprintf("%s@%s", toId.Account(), toId.Domain()))
		}
	}
	res := make([]model.Object, 0, len(storages))
	graph.Rank(85*pagerank.Dot2ONE, 6*pagerank.DotONE, func(label string, rank int64) {
		st := state.Fc.NewStorageBuilder().
			Id(fmt.Sprintf("%s/%s", label, outName)).
			Address("account_id", label).
			Int64("rank", rank).
			Build()
		res = append(res, state.Fc.NewObjectBuilder().Storage(st))
	})
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().List(res))
}

func ExecuteProslListOperator(op *proskenion.ListOperator, state *ProslStateValue) *ProslStateValue {
	obs := make([]model.Object, 0, len(op.GetObject()))
	for _, v := range op.GetObject() {
		state = ExecuteProslValueOperator(v, state)
		if state.Err != nil {
			return state
		}
		obs = append(obs, state.ReturnObject)
	}
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().List(obs))
}

func ExecuteProslMapOperator(op *proskenion.MapOperator, state *ProslStateValue) *ProslStateValue {
	obs := make(map[string]model.Object)
	for key, value := range op.GetObject() {
		state = ExecuteProslValueOperator(value, state)
		if state.Err != nil {
			return state
		}
		obs[key] = state.ReturnObject
	}
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Dict(obs))
}

func ExecuteProslCastOperator(op *proskenion.CastOperator, state *ProslStateValue) *ProslStateValue {
	code := op.GetType()
	state = ExecuteProslValueOperator(op.GetObject(), state)
	if state.Err != nil {
		return state
	}
	object := state.ReturnObject
	if object == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "Return Object is nil, %s", op.String())
	}
	ret, ok := object.Cast(model.ObjectCode(code))
	if !ok {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_CastType, "Can not cast %s to %s, %s", object.GetType().String(), op.GetType().String(), op.String())
	}
	return ReturnProslStateValue(state, ret)
}

func ExecuteProslListComprehensionOperator(op *proskenion.ListComprehensionOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetList(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetType() != model.ListObjectCode {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue, "Return Object type expected: %s,but actual: %s\n%s",
			model.ListObjectCode.String(), state.ReturnObject.GetType(), op.String())
	}
	list := state.ReturnObject.GetList()

	ret := make([]model.Object, 0)
	for _, o := range list {
		state.Variables[op.GetVariableName()] = o
		if op.GetIf() != nil {
			state = ExecuteProslConditionalFormula(op.GetIf(), state)
			if state.Err != nil {
				return state
			}
			if !state.ReturnObject.GetBoolean() {
				break
			}
		}
		state = ExecuteProslValueOperator(op.GetElement(), state)
		if state.Err != nil {
			return state
		}
		ret = append(ret, state.ReturnObject)
	}
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().List(ret))
}

type ComparedObjects struct {
	objs []model.Object
	code model.ObjectCode
	key  string
}

func (a *ComparedObjects) Len() int {
	return len(a.objs)
}

type GetFromKeyer interface {
	GetFromKey(string) model.Object
}

func lessGetFromKey(l GetFromKeyer, r GetFromKeyer, key string) bool {
	if l == nil || r == nil {
		return false
	}
	return model.ObjectLess(l.GetFromKey(key), r.GetFromKey(key))
}

func (a *ComparedObjects) Less(i, j int) bool {
	if a.key == "" {
		return model.ObjectLess(a.objs[i], a.objs[j])
	}
	switch a.code {
	case model.AccountObjectCode:
		return lessGetFromKey(a.objs[i].GetAccount(), a.objs[j].GetAccount(), a.key)
	case model.PeerObjectCode:
		return lessGetFromKey(a.objs[i].GetPeer(), a.objs[j].GetPeer(), a.key)
	case model.StorageObjectCode:
		return lessGetFromKey(a.objs[i].GetStorage(), a.objs[j].GetStorage(), a.key)
	default:
		return model.ObjectLess(a.objs[i], a.objs[j])
	}
}

func (a *ComparedObjects) Swap(i, j int) {
	a.objs[i], a.objs[j] = a.objs[j], a.objs[i]
}

func ExecuteProslSortOperator(op *proskenion.SortOperator, state *ProslStateValue) *ProslStateValue {
	// list : getList
	state = ExecuteProslValueOperator(op.GetList(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetType() != model.ListObjectCode {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnExpectedReturnValue,
			"Return Object type expected: %s,but actual: %s\n%s", model.ListObjectCode.String(), state.ReturnObject.GetType().String(), op.String())
	}
	list := state.ReturnObject.GetList()

	// Execute sort with list
	ret := &ComparedObjects{list, model.ObjectCode(op.GetType()), op.GetOrderBy().GetKey()}
	if op.GetOrderBy().GetOrder() == model.DESC {
		sort.Sort(sort.Reverse(ret))
	} else {
		sort.Sort(ret)
	}

	// limit check
	if op.GetLimit() == nil {
		return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().List(ret.objs))
	}
	state = ExecuteProslValueOperator(op.GetLimit(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetType() != model.Int32ObjectCode {
		return ReturnUnExpectedRetrunValue(state, model.Int32ObjectCode, state.ReturnObject.GetType(), op)
	}
	limit := state.ReturnObject.GetI32()
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().List(ret.objs[:int(limit)]))
}

func ExecuteProslSliceOperator(op *proskenion.SliceOperator, state *ProslStateValue) *ProslStateValue {
	// list : getList
	state = ExecuteProslValueOperator(op.GetList(), state)
	if state.Err != nil {
		return state
	}
	if state.ReturnObject.GetType() != model.ListObjectCode {
		return ReturnUnExpectedRetrunValue(state, model.ListObjectCode, state.ReturnObject.GetType(), op)
	}
	list := state.ReturnObject.GetList()

	// left [left, right)
	left := 0
	if op.GetLeft() != nil {
		state = ExecuteProslValueOperator(op.GetLeft(), state)
		if state.Err != nil {
			return state
		}
		if state.ReturnObject.GetType() != model.Int32ObjectCode {
			return ReturnUnExpectedRetrunValue(state, model.Int32ObjectCode, state.ReturnObject.GetType(), op)
		}
		left = int(state.ReturnObject.GetI32())
	}
	if left < 0 {
		left = 0
	}

	// right [left, right)
	right := len(list)
	if op.GetRight() != nil {
		state = ExecuteProslValueOperator(op.GetRight(), state)
		if state.Err != nil {
			return state
		}
		if state.ReturnObject.GetType() != model.Int32ObjectCode {
			return ReturnUnExpectedRetrunValue(state, model.Int32ObjectCode, state.ReturnObject.GetType(), op)
		}
		right = int(state.ReturnObject.GetI32())
	}
	if right > len(list) {
		right = len(list)
	}

	if len(list) <= left || left > right {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_OutOfRange,
			"SliceOperator invalid range left: %d, right: %d, len(list): %d\n%s", left, right, len(list), op.String())
	}
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().List(list[left:right]))
}

func ExecuteProslLenOperator(op *proskenion.LenOperator, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetList(), state)
	if state.Err != nil {
		return state
	}
	list := state.ReturnObject.GetList()
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Int32(int32(len(list))))
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
	case *proskenion.ConditionalFormula_VerifyOp:
		state = ExecuteProslVerifyOperator(op.GetVerifyOp(), state)
	default:
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_UnImplemented, "undefined forumula: %s", op.String())
	}
	return state
}

func ExecutePolynomialCondOperator(op GetOpser, f func(model.Object, model.Object, model.ModelFactory) model.Object, symbol string, state *ProslStateValue) *ProslStateValue {
	if len(op.GetOps()) < 2 {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_NotEnoughArgument, fmt.Sprintf("%s Operator minimum number of argument is 2, %s", symbol, op.String()))
	}
	state = ExecuteProslValueOperator(op.GetOps()[0], state)
	if state.Err != nil {
		return state
	}
	pr := state.ReturnObject
	for _, o := range op.GetOps()[1:] {
		state = ExecuteProslValueOperator(o, state)
		if state.Err != nil {
			return state
		}
		ret := f(pr, state.ReturnObject, state.Fc)
		if ret == nil {
			return ReturnErrorProslStateValue(state, proskenion.ErrCode_FailedOperate, "symbol: %s, %s", symbol, op.String())
		}
		if !ret.GetBoolean() {
			return ReturnProslStateValue(state, ret)
		}
	}
	return ReturnProslStateValue(state, state.Fc.NewObjectBuilder().Bool(true))
}

func ExecuteProslOrFormula(op *proskenion.OrFormula, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteCondOr, "cond-or", state)
}

func ExecuteProslAndFormula(op *proskenion.AndFormula, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomiaValueOperator(op, ExecuteCondAnd, "cond-and", state)
}

func ExecuteProslNotFormula(op *proskenion.NotFormula, state *ProslStateValue) *ProslStateValue {
	state = ExecuteProslValueOperator(op.GetOp(), state)
	if state.Err != nil {
		return state
	}
	ret := ExecuteCondNot(state.ReturnObject, state.Fc)
	if ret == nil {
		return ReturnErrorProslStateValue(state, proskenion.ErrCode_FailedOperate, op.String())
	}
	return ReturnProslStateValue(state, ret)
}

func ExecuteValueOpeatorsToObjectList(op GetOpser, symbol string, state *ProslStateValue) ([]model.Object, *ProslStateValue) {
	if len(op.GetOps()) < 2 {
		return nil, ReturnErrorProslStateValue(state, proskenion.ErrCode_NotEnoughArgument,
			fmt.Sprintf("%s Operator minimum number of argument is 2, %s", symbol, op.String()))
	}
	os := make([]model.Object, 0, len(op.GetOps()))
	for _, o := range op.GetOps() {
		state = ExecuteProslValueOperator(o, state)
		if state.Err != nil {
			return nil, state
		}
		os = append(os, state.ReturnObject)
	}
	return os, state
}

func ExecuteProslEqFormula(op *proskenion.EqFormula, state *ProslStateValue) *ProslStateValue {
	os, state := ExecuteValueOpeatorsToObjectList(op, "eq(==)", state)
	if state.Err != nil {
		return state
	}
	ret := ExecuteCondEq(os, state.Fc)
	return ReturnProslStateValue(state, ret)
}

func ExecuteProslNeFormula(op *proskenion.NeFormula, state *ProslStateValue) *ProslStateValue {
	os, state := ExecuteValueOpeatorsToObjectList(op, "ne(!=)", state)
	if state.Err != nil {
		return state
	}
	ret := ExecuteCondNe(os, state.Fc)
	return ReturnProslStateValue(state, ret)
}

func ExecuteProslGtFormula(op *proskenion.GtFormula, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialCondOperator(op, ExecuteCondGt, "gt(>)", state)
}

func ExecuteProslGeFormula(op *proskenion.GeFormula, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialCondOperator(op, ExecuteCondGe, "ge(>=)", state)
}

func ExecuteProslLtFormula(op *proskenion.LtFormula, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialCondOperator(op, ExecuteCondLt, "lt(<)", state)
}

func ExecuteProslLeFormula(op *proskenion.LeFormula, state *ProslStateValue) *ProslStateValue {
	return ExecutePolynomialCondOperator(op, ExecuteCondLe, "le(<=)", state)
}
