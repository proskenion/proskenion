package prosl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/proto"
	"gopkg.in/yaml.v2"
	"reflect"
)

var (
	ErrConvertYamlToMap          = fmt.Errorf("Failed Convert yaml to map")
	ErrProslParseNotExpectedType = fmt.Errorf("Failed Prosl Parse not expected type")

	ErrProslParseUnExpectedCastType      = fmt.Errorf("Failed Prosl Parse Cast not expected type")
	ErrProslParseArgumentSize            = fmt.Errorf("Failed Prosl Parse argument size")
	ErrProslParseUnExpectedOperationName = fmt.Errorf("Failed Prosl Parse operation")
	ErrProslParseInternalErr             = fmt.Errorf("Failed Prosl Parse Internal")
)

func ConvertYamlToMap(yamlBytes []byte) ([]interface{}, error) {
	yamap := make([]interface{}, 0)
	err := yaml.Unmarshal(yamlBytes, &yamap)
	if err != nil {
		return nil, err
	}
	return yamap, nil
}

func ProslParseCastError(ex interface{}, ac interface{}) error {
	return errors.Wrapf(ErrProslParseArgumentSize, "expected: %T, actual: %T", ex, ac)
}

func ProslParseArgumentError(ex int, ac int) error {
	return errors.Wrapf(ErrProslParseArgumentSize, "expected: %d, actual: %d", ex, ac)
}

func ProslParseErrOperation(op string) error {
	return errors.Wrapf(ErrProslParseUnExpectedOperationName, "unknown operation: %s", op)
}

func ParseProsl(yalist []interface{}) (*proskenion.Prosl, error) {
	ops := make([]*proskenion.ProslOperator, len(yalist))
	for _, ya := range yalist {
		switch v := ya.(type) {
		case map[interface{}]interface{}:
			pop, err := ParseProslOperator(v)
			if err != nil {
				return nil, err
			}
			ops = append(ops, pop)
		default:
			return nil, errors.Wrap(ErrProslParseNotExpectedType, fmt.Sprintf("%T", v))
		}
	}
	return &proskenion.Prosl{Ops: ops}, nil
}

// ParseProslOperator
func ParseProslOperator(yamap map[interface{}]interface{}) (*proskenion.ProslOperator, error) {
	if len(yamap) != 1 {
		return nil, ProslParseArgumentError(1, len(yamap))
	}
	for key, value := range yamap {
		switch key {
		case "set":
			op, err := ParseSetOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_SetOp{SetOp: op}}, nil
		case "if":
			op, err := ParseIfOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_IfOp{IfOp: op}}, nil
		case "elif":
			op, err := ParseElifOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_ElifOp{ElifOp: op}}, nil
		case "else":
			op, err := ParseElseOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_ElseOp{ElseOp: op}}, nil
		case "err":
			op, err := ParseErrCatchOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_ErrOp{ErrOp: op}}, nil
		case "require":
			op, err := ParseRequireOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_RequireOp{RequireOp: op}}, nil
		case "assert":
			op, err := ParseAssertOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_AssertOp{AssertOp: op}}, nil
		case "verify":
			op, err := ParseVerifyOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_VerifyOp{VerifyOp: op}}, nil
		case "return":
			op, err := ParseReturnOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_ReturnOp{ReturnOp: op}}, nil
		default:
			return nil, ProslParseErrOperation(value.(string))
		}
	}
	return nil, ErrProslParseInternalErr
}

// set:
// 	- variableName (string)
//  - valueOperator (interface{})
func ParseSetOperator(yaml interface{}) (*proskenion.SetOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) == 2 {
			return nil, ProslParseArgumentError(2, len(yalist))
		}
		ret := &proskenion.SetOperator{}

		// variableName
		if v, ok := yalist[0].(string); ok {
			ret.VariableName = v
		} else {
			return nil, ProslParseCastError("", v)
		}

		// valueOperator
		if vop, err := ParseValueOperator(yalist[1].(interface{})); err != nil {
			return nil, err
		} else {
			ret.Value = vop
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml)
}

func ParseIfOperator(yaml interface{}) (*proskenion.IfOperator, error) {
	return nil, nil
}

func ParseElifOperator(yaml interface{}) (*proskenion.ElifOperator, error) {
	return nil, nil
}

func ParseElseOperator(yaml interface{}) (*proskenion.ElseOperator, error) {
	return nil, nil
}

func ParseErrCatchOperator(yaml interface{}) (*proskenion.ErrCatchOperator, error) {
	return nil, nil
}

func ParseRequireOperator(yaml interface{}) (*proskenion.RequireOperator, error) {
	return nil, nil
}

func ParseAssertOperator(yaml interface{}) (*proskenion.AssertOperator, error) {
	return nil, nil
}

func ParseVerifyOperator(yaml interface{}) (*proskenion.VerifyOperator, error) {
	return nil, nil
}

func ParseReturnOperator(yaml interface{}) (*proskenion.ReturnOperator, error) {
	return nil, nil
}

func ParseValueOperator(yalist interface{}) (*proskenion.ValueOperator, error) {
	return nil, nil
}

func ParseQueryOperator(yalist map[interface{}]interface{}) (*proskenion.QueryOperator, error) {
	return nil, nil
}

func ParseCommandOperator(yalist map[interface{}]interface{}) (*proskenion.CommandOperator, error) {
	return nil, nil
}

func ParseTxOperator(yalist map[interface{}]interface{}) (*proskenion.TxOperator, error) {
	return nil, nil
}

func ParsePlusOperator(yaml interface{}) (*proskenion.PlusOperator, error) {
	return nil, nil
}

func ParseMinusOperator(yaml interface{}) (*proskenion.MinusOperator, error) {
	return nil, nil
}

func ParseMultipleOperator(yaml interface{}) (*proskenion.MultipleOperator, error) {
	return nil, nil
}

func ParseDivisionOperator(yaml interface{}) (*proskenion.DivisionOperator, error) {
	return nil, nil
}

func ParseModOperator(yaml interface{}) (*proskenion.ModOperator, error) {
	return nil, nil
}

func ParseOrOperator(yaml interface{}) (*proskenion.OrOperator, error) {
	return nil, nil
}

func ParseAndOperator(yaml interface{}) (*proskenion.AndOperator, error) {
	return nil, nil
}

func ParseXorOperator(yaml interface{}) (*proskenion.XorOperator, error) {
	return nil, nil
}

func ParseConcatOperator(yaml interface{}) (*proskenion.ConcatOperator, error) {
	return nil, nil
}

func ParseValuedOperator(yaml interface{}) (*proskenion.ValuedOperator, error) {
	return nil, nil
}

func ParseIndexedOperator(yaml interface{}) (*proskenion.ConditionalFormula, error) {
	return nil, nil
}

func ParseConditionalFormula(yamap map[interface{}]interface{}) (*proskenion.ConditionalFormula, error) {
	return nil, nil
}

func ParseOrFormula(yaml interface{}) (*proskenion.OrFormula, error) {
	return nil, nil
}

func ParseAndFormula(yaml interface{}) (*proskenion.AndFormula, error) {
	return nil, nil
}

func ParseNotFormula(yaml interface{}) (*proskenion.NotFormula, error) {
	return nil, nil
}

func ParseEqFormula(yaml interface{}) (*proskenion.EqFormula, error) {
	return nil, nil
}

func ParseNeFormula(yaml interface{}) (*proskenion.NeFormula, error) {
	return nil, nil
}

func ParseGtFormula(yaml interface{}) (*proskenion.GtFormula, error) {
	return nil, nil
}

func ParseGeFormula(yaml interface{}) (*proskenion.GeFormula, error) {
	return nil, nil
}

func ParseLtFormula(yaml interface{}) (*proskenion.LtFormula, error) {
	return nil, nil
}

func ParseLeFormula(yaml interface{}) (*proskenion.LeFormula, error) {
	if len(yalist) == 2 {
		return nil, nil
	}
	return nil, nil
}

func ParseIsDefinedFormula(variableName string) (*proskenion.IsDefinedFormula, error) {
	return &proskenion.IsDefinedFormula{VariableName: variableName}, nil
}

// pattern
// list -> []interface{}
// map -> map[interface{}]interface{}
// string
// int
func ConvertYamlToProtobuf(yamlBytes []byte) (*proskenion.Prosl, error) {
	yamap, err := ConvertYamlToMap(yamlBytes)
	if err != nil {
		return nil, errors.Wrap(ErrConvertYamlToMap, err.Error())
	}

	prosl, err := ParseProsl(yamap)

	return nil, nil
}
