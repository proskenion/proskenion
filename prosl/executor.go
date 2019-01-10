package prosl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
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

type ProslExecutor struct {
	Variables map[string]model.Object
}

type ProslStateValue struct {
	Variables    map[string]model.Object
	ReturnObject model.Object
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
		ReturnObject: value,
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
		b, err := state.ReturnObject.Marshal()
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

// TODO
func ExecuteProslTxOperator(op *proskenion.TxOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslCmdOperator(op *proskenion.CommandOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslPlusOperator(op *proskenion.PlusOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslMinusOperator(op *proskenion.MinusOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslMulOperator(op *proskenion.MultipleOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslDivOperator(op *proskenion.DivisionOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslModOperator(op *proskenion.ModOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslOrOperator(op *proskenion.OrOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslAndOperator(op *proskenion.AndOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslXorOperator(op *proskenion.XorOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslConcatOperator(op *proskenion.ConcatOperator, state *ProslStateValue) *ProslStateValue {
	return nil
}

func ExecuteProslValuedOperator(op *proskenion.ValuedOperator, state *ProslStateValue) *ProslStateValue {
	return nil
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
