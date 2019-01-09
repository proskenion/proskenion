package prosl

import (
	"fmt"
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
)

type ProslReturnValue struct {
	model.Object
}

type ProslExecutor struct {
	Variables []map[string]model.Object   // stack
	PreOp     []*proskenion.ProslOperator // stack
	// 順次
	// 判断
	// 繰り返し
}

func (e *ProslExecutor) Prosl(prosl *proskenion.Prosl) (*ProslReturnValue, error) {
	ops := prosl.GetOps()
	for _, op := range ops {
		_, err := e.ProslOperator(op)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (e *ProslExecutor) ProslOperator(op *proskenion.ProslOperator) (*ProslReturnValue, error) {
	switch op.GetOp().(type) {
	case *proskenion.ProslOperator_SetOp:
	case *proskenion.ProslOperator_IfOp:
	case *proskenion.ProslOperator_ElifOp:
	case *proskenion.ProslOperator_ElseOp:
	case *proskenion.ProslOperator_ErrOp:
	case *proskenion.ProslOperator_RequireOp:
	case *proskenion.ProslOperator_AssertOp:
	case *proskenion.ProslOperator_VerifyOp:
	case *proskenion.ProslOperator_ReturnOp:
	default:
	}
	return nil, nil
}

func (e *ProslExecutor) ProslSetOperator(op *proskenion.SetOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslIfOperator(op *proskenion.IfOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProsEliffOperator(op *proskenion.ElifOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslElseOperator(op *proskenion.ElseOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslErrOperator(op *proskenion.ErrCatchOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslRequireOperator(op *proskenion.RequireOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslAssertOperator(op *proskenion.AssertOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslVerifyOperator(op *proskenion.VerifyOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslReturnOperator(op *proskenion.ReturnOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ValueOperator(op *proskenion.ValueOperator) (*ProslReturnValue, error) {
	switch op.GetOp().(type) {
	case *proskenion.ValueOperator_QueryOp:
	case *proskenion.ValueOperator_TxOp:
	case *proskenion.ValueOperator_CmdOp:
	case *proskenion.ValueOperator_PlusOp:
	case *proskenion.ValueOperator_MinusOp:
	case *proskenion.ValueOperator_MulOp:
	case *proskenion.ValueOperator_DivOp:
	case *proskenion.ValueOperator_ModOp:
	case *proskenion.ValueOperator_OrOp:
	case *proskenion.ValueOperator_AndOp:
	case *proskenion.ValueOperator_XorOp:
	case *proskenion.ValueOperator_ConcatOp:
	case *proskenion.ValueOperator_ValuedOp:
	case *proskenion.ValueOperator_IndexedOp:
	case *proskenion.ValueOperator_VariableOp:
	case *proskenion.ValueOperator_Object:
	default:
	}
	return nil, nil
}

func (e *ProslExecutor) ProslQueryOperator(op *proskenion.QueryOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslTxOperator(op *proskenion.TxOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslCmdOperator(op *proskenion.CommandOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslPlusOperator(op *proskenion.PlusOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslMinusOperator(op *proskenion.MinusOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslMulOperator(op *proskenion.MultipleOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslDivOperator(op *proskenion.DivisionOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslModOperator(op *proskenion.ModOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslOrOperator(op *proskenion.OrOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslAndOperator(op *proskenion.AndOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslXorOperator(op *proskenion.XorOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslConcatOperator(op *proskenion.ConcatOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslValuedOperator(op *proskenion.ValuedOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslIndexedOperator(op *proskenion.IndexedOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslVariableOperator(op *proskenion.VariableOperator) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslObjectOperator(op *proskenion.Object) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ConditionalFormula(op *proskenion.ConditionalFormula) (*ProslReturnValue, error) {
	switch op.GetOp().(type) {
	case *proskenion.ConditionalFormula_Or:
	case *proskenion.ConditionalFormula_And:
	case *proskenion.ConditionalFormula_Not:
	case *proskenion.ConditionalFormula_Eq:
	case *proskenion.ConditionalFormula_Ne:
	case *proskenion.ConditionalFormula_Gt:
	case *proskenion.ConditionalFormula_Ge:
	case *proskenion.ConditionalFormula_Lt:
	case *proskenion.ConditionalFormula_Le:
	default:
	}
	return nil, nil
}

func (e *ProslExecutor) ProslOrFormula(op *proskenion.OrFormula) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslAndFormula(op *proskenion.AndFormula) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslNotFormula(op *proskenion.NotFormula) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslEqFormula(op *proskenion.EqFormula) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslNeFormula(op *proskenion.NeFormula) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslGtFormula(op *proskenion.GtFormula) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslGeFormula(op *proskenion.GeFormula) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslLtFormula(op *proskenion.LtFormula) (*ProslReturnValue, error) {
	return nil, nil
}

func (e *ProslExecutor) ProslLeFormula(op *proskenion.LeFormula) (*ProslReturnValue, error) {
	return nil, nil
}