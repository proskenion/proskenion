package prosl

import (
	"fmt"
	"github.com/pkg/errors"
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

	ErrProslExecuteSentenceErr = fmt.Errorf("Faile Prosl Execute Sentence error")
)

type OperatorState int

const (
	AnotherOperator_State OperatorState = iota
	IfOperatorTrue_State
	IfOperatorFalse_State
	ElifOperatorTrue_State
	ElifOperatorFalse_State
	ReturnOperator_State
)

type ProslStateValue struct {
	Variables    map[string]model.Object
	ReturnObject model.Object
	St           OperatorState
}

func ReturnOpProslStateValue(state *ProslStateValue, st OperatorState) *ProslStateValue {
	return &ProslStateValue{
		Variables:    state.Variables,
		ReturnObject: nil,
		St:           st,
	}
}

func ReturnValueProslStateValue(state *ProslStateValue, value model.Object) *ProslStateValue {
	return &ProslStateValue{
		Variables:    state.Variables,
		ReturnObject: value,
		St:           AnotherOperator_State,
	}
}

func ExecuteProsl(prosl *proskenion.Prosl, state *ProslStateValue) (*ProslStateValue, error) {
	ops := prosl.GetOps()
	var err error
	for _, op := range ops {
		state, err = ExecuteProslOpFormula(op, state)
		if err != nil {
			return nil, err
		}
	}
	return state, nil
}

func ExecuteProslOpFormula(op *proskenion.ProslOperator, state *ProslStateValue) (*ProslStateValue, error) {
	var err error
	switch op.GetOp().(type) {
	case *proskenion.ProslOperator_SetOp:
		state, err = ExecuteProslSetOperator(op.GetSetOp(), state)
		state, err = ExecuteProslSetOperator(op.GetSetOp(), state)
	case *proskenion.ProslOperator_IfOp:
		state, err = ExecuteProslIfOperator(op.GetIfOp(), state)
	case *proskenion.ProslOperator_ElifOp:
		state, err = ExecuteProslElifOperator(op.GetElifOp(), state)
	case *proskenion.ProslOperator_ElseOp:
		state, err = ExecuteProslElseOperator(op.GetElseOp(), state)
	case *proskenion.ProslOperator_ErrOp:
		state, err = ExecuteProslErrOperator(op.GetErrOp(), state)
	case *proskenion.ProslOperator_RequireOp:
		state, err = ExecuteProslRequireOperator(op.GetRequireOp(), state)
	case *proskenion.ProslOperator_AssertOp:
		state, err = ExecuteProslAssertOperator(op.GetAssertOp(), state)
	case *proskenion.ProslOperator_VerifyOp:
		state, err = ExecuteProslVerifyOperator(op.GetVerifyOp(), state)
	case *proskenion.ProslOperator_ReturnOp:
		state, err = ExecuteProslReturnOperator(op.GetReturnOp(), state)
	default:
	}
	return state, err
}

func ExecuteProslSetOperator(op *proskenion.SetOperator, state *ProslStateValue) (*ProslStateValue, error) {
	var err error
	state, err = ExecuteProslValueOperator(op.GetValue(), state)
	if err != nil {
		return nil, err
	}
	state.Variables[op.GetVariableName()] = state.ReturnObject
	return ReturnOpProslStateValue(state, AnotherOperator_State), nil
}

func ExecuteProslIfOperator(op *proskenion.IfOperator, state *ProslStateValue) (*ProslStateValue, error) {
	var err error
	state, err = ExecuteProslConditionalFormula(op.GetOp(), state)
	if err != nil {
		return nil, err
	}
	if state.ReturnObject.GetBoolean() {
		state, err = ExecuteProsl(op.GetProsl(), state)
		if err != nil {
			return nil, err
		}
		return ReturnOpProslStateValue(state, IfOperatorTrue_State), nil
	}
	return ReturnOpProslStateValue(state, IfOperatorFalse_State), nil
}

func ExecuteProslElifOperator(op *proskenion.ElifOperator, state *ProslStateValue) (*ProslStateValue, error) {
	var err error
	switch state.St {
	case IfOperatorFalse_State, ElifOperatorFalse_State:
		state, err = ExecuteProslConditionalFormula(op.GetOp(), state)
		if err != nil {
			return nil, err
		}
		if state.ReturnObject.GetBoolean() {
			state, err = ExecuteProsl(op.GetProsl(), state)
			if err != nil {
				return nil, err
			}
			return ReturnOpProslStateValue(state, ElifOperatorTrue_State), nil
		}
		return ReturnOpProslStateValue(state, ElifOperatorFalse_State), nil
	case IfOperatorTrue_State, ElifOperatorTrue_State:
		return state, nil
	}
	return nil, errors.Wrapf(ErrProslExecuteSentenceErr, "elif operator must have previous operator that is if or elif operator")
}

func ExecuteProslElseOperator(op *proskenion.ElseOperator, state *ProslStateValue) (*ProslStateValue, error) {
	var err error
	switch state.St {
	case IfOperatorFalse_State, ElifOperatorFalse_State:
		state, err = ExecuteProsl(op.GetProsl(), state)
		if err != nil {
			return nil, err
		}
		return ReturnOpProslStateValue(state, AnotherOperator_State), nil
	case IfOperatorTrue_State, ElifOperatorTrue_State:
		return state, nil
	}
	return nil, errors.Wrapf(ErrProslExecuteSentenceErr, "else operator must have previous operator that is if or elif operator")
}

func ExecuteProslErrOperator(op *proskenion.ErrCatchOperator, state *ProslStateValue) (*ProslStateValue, error) {
	// TODO
	return nil, nil
}

func ExecuteProslRequireOperator(op *proskenion.RequireOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslAssertOperator(op *proskenion.AssertOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslVerifyOperator(op *proskenion.VerifyOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslReturnOperator(op *proskenion.ReturnOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslValueOperator(op *proskenion.ValueOperator, state *ProslStateValue) (*ProslStateValue, error) {
	var err error
	switch op.GetOp().(type) {
	case *proskenion.ValueOperator_QueryOp:
		state, err = ExecuteProslQueryOperator(op.GetQueryOp(), state)
	case *proskenion.ValueOperator_TxOp:
		state, err = ExecuteProslTxOperator(op.GetTxOp(), state)
	case *proskenion.ValueOperator_CmdOp:
		state, err = ExecuteProslCmdOperator(op.GetCmdOp(), state)
	case *proskenion.ValueOperator_PlusOp:
		state, err = ExecuteProslPlusOperator(op.GetPlusOp(), state)
	case *proskenion.ValueOperator_MinusOp:
		state, err = ExecuteProslMinusOperator(op.GetMinusOp(), state)
	case *proskenion.ValueOperator_MulOp:
		state, err = ExecuteProslMulOperator(op.GetMulOp(), state)
	case *proskenion.ValueOperator_DivOp:
		state, err = ExecuteProslDivOperator(op.GetDivOp(), state)
	case *proskenion.ValueOperator_ModOp:
		state, err = ExecuteProslModOperator(op.GetModOp(), state)
	case *proskenion.ValueOperator_OrOp:
		state, err = ExecuteProslOrOperator(op.GetOrOp(), state)
	case *proskenion.ValueOperator_AndOp:
		state, err = ExecuteProslAndOperator(op.GetAndOp(), state)
	case *proskenion.ValueOperator_XorOp:
		state, err = ExecuteProslXorOperator(op.GetXorOp(), state)
	case *proskenion.ValueOperator_ConcatOp:
		state, err = ExecuteProslConcatOperator(op.GetConcatOp(), state)
	case *proskenion.ValueOperator_ValuedOp:
		state, err = ExecuteProslValuedOperator(op.GetValuedOp(), state)
	case *proskenion.ValueOperator_IndexedOp:
		state, err = ExecuteProslIndexedOperator(op.GetIndexedOp(), state)
	case *proskenion.ValueOperator_VariableOp:
		state, err = ExecuteProslVariableOperator(op.GetVariableOp(), state)
	case *proskenion.ValueOperator_Object:
		state, err = ExecuteProslObjectOperator(op.GetObject(), state)
	default:
	}
	return state, err
}

func ExecuteProslQueryOperator(op *proskenion.QueryOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslTxOperator(op *proskenion.TxOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslCmdOperator(op *proskenion.CommandOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslPlusOperator(op *proskenion.PlusOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslMinusOperator(op *proskenion.MinusOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslMulOperator(op *proskenion.MultipleOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslDivOperator(op *proskenion.DivisionOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslModOperator(op *proskenion.ModOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslOrOperator(op *proskenion.OrOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslAndOperator(op *proskenion.AndOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslXorOperator(op *proskenion.XorOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslConcatOperator(op *proskenion.ConcatOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslValuedOperator(op *proskenion.ValuedOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslIndexedOperator(op *proskenion.IndexedOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslVariableOperator(op *proskenion.VariableOperator, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslObjectOperator(op *proskenion.Object, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslConditionalFormula(op *proskenion.ConditionalFormula, state *ProslStateValue) (*ProslStateValue, error) {
	var err error
	switch op.GetOp().(type) {
	case *proskenion.ConditionalFormula_Or:
		state, err = ExecuteProslOrFormula(op.GetOr(), state)
	case *proskenion.ConditionalFormula_And:
		state, err = ExecuteProslAndFormula(op.GetAnd(), state)
	case *proskenion.ConditionalFormula_Not:
		state, err = ExecuteProslNotFormula(op.GetNot(), state)
	case *proskenion.ConditionalFormula_Eq:
		state, err = ExecuteProslEqFormula(op.GetEq(), state)
	case *proskenion.ConditionalFormula_Ne:
		state, err = ExecuteProslNeFormula(op.GetNe(), state)
	case *proskenion.ConditionalFormula_Gt:
		state, err = ExecuteProslGtFormula(op.GetGt(), state)
	case *proskenion.ConditionalFormula_Ge:
		state, err = ExecuteProslGeFormula(op.GetGe(), state)
	case *proskenion.ConditionalFormula_Lt:
		state, err = ExecuteProslLtFormula(op.GetLt(), state)
	case *proskenion.ConditionalFormula_Le:
		state, err = ExecuteProslLeFormula(op.GetLe(), state)
	default:
	}
	return state, err
}

func ExecuteProslOrFormula(op *proskenion.OrFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslAndFormula(op *proskenion.AndFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslNotFormula(op *proskenion.NotFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslEqFormula(op *proskenion.EqFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslNeFormula(op *proskenion.NeFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslGtFormula(op *proskenion.GtFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslGeFormula(op *proskenion.GeFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslLtFormula(op *proskenion.LtFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}

func ExecuteProslLeFormula(op *proskenion.LeFormula, state *ProslStateValue) (*ProslStateValue, error) {
	return nil, nil
}
