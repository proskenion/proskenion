package prosl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/proto"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v2"
	"strings"
)

var (
	ErrConvertYamlToMap          = fmt.Errorf("Failed Convert yaml to map")
	ErrProslParseNotExpectedType = fmt.Errorf("Failed Prosl Parse not expected type")

	ErrProslParseUnExpectedCastType      = fmt.Errorf("Failed Prosl Parse Cast not expected type")
	ErrProslParseArgumentSize            = fmt.Errorf("Failed Prosl Parse argument size")
	ErrProslParseUnExpectedOperationName = fmt.Errorf("Failed Prosl Parse operation")
	ErrProslParseInternalErr             = fmt.Errorf("Failed Prosl Parse Internal")
	ErrProslParseUnknownObjectCode       = fmt.Errorf("Failed Prosl Parse Unklnown object code")
	ErrProslParseQueryOperatorArgument   = fmt.Errorf("Failed Prosl not enough query operator arguments")
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
	return errors.Wrapf(ErrProslParseUnExpectedCastType, "expected: %T, actual: %T", ex, ac)
}

func ProslParseUnknownCastError(ac interface{}) error {
	return errors.Wrapf(ErrProslParseUnExpectedCastType, "unknown type : %T", ac)
}
func ProslParseArgumentError(ex int, ac int) error {
	return errors.Wrapf(ErrProslParseArgumentSize, "expected: %d, actual: %d", ex, ac)
}

func ProslParseArgumentErrorMin(minEx int, ac int) error {
	return errors.Wrapf(ErrProslParseArgumentSize, "expected moret than: %d, actual: %d", minEx, ac)
}

func ProslParseErrOperation(op string) error {
	return errors.Wrapf(ErrProslParseUnExpectedOperationName, "unknown operation: %s", op)
}

func ProslParseErrCode(code string) proskenion.ErrCode {
	switch code {
	case "AnythingErrCode":
		return proskenion.ErrCode_AnythingErrCode
	}
	return proskenion.ErrCode_AnythingErrCode
}

func ProslParsePrimitiveObject(value interface{}) (*proskenion.Object, error) {
	switch v := value.(type) {
	case int64:
		return &proskenion.Object{
			Type:   proskenion.ObjectCode_Int64ObjectCode,
			Object: &proskenion.Object_I64{v},
		}, nil
	case int:
		return &proskenion.Object{
			Type:   proskenion.ObjectCode_Int32ObjectCode,
			Object: &proskenion.Object_I32{int32(v)},
		}, nil
	case string:
		return &proskenion.Object{
			Type:   proskenion.ObjectCode_StringObjectCode,
			Object: &proskenion.Object_Str{v},
		}, nil
	}
	return nil, ProslParseUnknownCastError(value)
}

func ProslParseObjectCode(code string) (proskenion.ObjectCode, error) {
	code = strings.ToLower(code)
	switch code {
	case "bool":
		return proskenion.ObjectCode_BoolObjectCode, nil
	case "int32":
		return proskenion.ObjectCode_Int32ObjectCode, nil
	case "int64":
		return proskenion.ObjectCode_Int64ObjectCode, nil
	case "uint32":
		return proskenion.ObjectCode_Uint32ObjectCode, nil
	case "uint64":
		return proskenion.ObjectCode_Uint64ObjectCode, nil
	case "string":
		return proskenion.ObjectCode_StringObjectCode, nil
	case "bytes":
		return proskenion.ObjectCode_BytesObjectCode, nil
	case "address":
		return proskenion.ObjectCode_AddressObjectCode, nil
	case "signature":
		return proskenion.ObjectCode_SignatureObjectCode, nil
	case "account":
		return proskenion.ObjectCode_AccountObjectCode, nil
	case "peer":
		return proskenion.ObjectCode_PeerObjectCode, nil
	case "list":
		return proskenion.ObjectCode_ListObjectCode, nil
	case "dict":
		return proskenion.ObjectCode_DictObjectCode, nil
	}
	return 0, errors.Wrapf(ErrProslParseUnknownObjectCode, "unknown type : %s", code)
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
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) > 1 {
			return nil, ProslParseArgumentErrorMin(2, len(yalist))
		}
		ret := &proskenion.IfOperator{}

		if v, ok := yalist[0].(interface{}); ok {
			op, err := ParseConditionalFormula(v)
			if err != nil {
				return nil, err
			}
			ret.Op = op
		}
		prosl, err := ParseProsl(yalist[1:])
		if err != nil {
			return nil, err
		}
		ret.Prosl = prosl
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml)
}

func ParseElifOperator(yaml interface{}) (*proskenion.ElifOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) > 1 {
			return nil, ProslParseArgumentErrorMin(2, len(yalist))
		}
		ret := &proskenion.ElifOperator{}

		if v, ok := yalist[0].(interface{}); ok {
			op, err := ParseConditionalFormula(v)
			if err != nil {
				return nil, err
			}
			ret.Op = op
		}
		prosl, err := ParseProsl(yalist[1:])
		if err != nil {
			return nil, err
		}
		ret.Prosl = prosl
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml)
}

func ParseElseOperator(yaml interface{}) (*proskenion.ElseOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) > 0 {
			return nil, ProslParseArgumentErrorMin(1, len(yalist))
		}
		ret := &proskenion.ElseOperator{}
		prosl, err := ParseProsl(yalist)
		if err != nil {
			return nil, err
		}
		ret.Prosl = prosl
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml)
}

func ParseErrCatchOperator(yaml interface{}) (*proskenion.ErrCatchOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) > 1 {
			return nil, ProslParseArgumentErrorMin(2, len(yalist))
		}
		ret := &proskenion.ErrCatchOperator{}

		if v, ok := yalist[0].(string); ok {
			ret.Code = ProslParseErrCode(v)
		}
		prosl, err := ParseProsl(yalist[1:])
		if err != nil {
			return nil, err
		}
		ret.Prosl = prosl
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml)
}

func ParseRequireOperator(yaml interface{}) (*proskenion.RequireOperator, error) {
	ret := &proskenion.RequireOperator{}
	op, err := ParseConditionalFormula(yaml)
	if err != nil {
		return nil, err
	}
	ret.Op = op
	return ret, nil
}

func ParseAssertOperator(yaml interface{}) (*proskenion.AssertOperator, error) {
	op, err := ParseConditionalFormula(yaml)
	if err != nil {
		return nil, err
	}
	ret := &proskenion.AssertOperator{Op: op}
	return ret, nil
}

func ParseVerifyOperator(yaml interface{}) (*proskenion.VerifyOperator, error) {
	op, err := ParseValueOperator(yaml)
	if err != nil {
		return nil, err
	}
	ret := &proskenion.VerifyOperator{Op: op}
	return ret, nil
}

func ParseReturnOperator(yaml interface{}) (*proskenion.ReturnOperator, error) {
	op, err := ParseValueOperator(yaml)
	if err != nil {
		return nil, err
	}
	ret := &proskenion.ReturnOperator{Op: op}
	return ret, nil
}

func ParseValueOperator(yaml interface{}) (*proskenion.ValueOperator, error) {
	if v, ok := yaml.(map[interface{}]interface{}); ok {
		if len(v) != 1 {
			return nil, ProslParseArgumentError(1, len(v))
		}
		for key, value := range v {
			switch key {
			case "query":
				op, err := ParseQueryOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_QueryOp{op}}, nil
			case "transaction":
				op, err := ParseTxOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_TxOp{op}}, nil
			case "command":
				op, err := ParseCommandOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_CmdOp{op}}, nil
			case "+":
				op, err := ParsePlusOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_PlusOp{op}}, nil
			case "-":
				op, err := ParseMinusOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_MinusOp{op}}, nil
			case "*":
				op, err := ParseMultipleOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_MulOp{op}}, nil
			case "/":
				op, err := ParseDivisionOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_DivOp{op}}, nil
			case "%":
				op, err := ParseModOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_ModOp{op}}, nil
			case "or":
				op, err := ParseOrOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_OrOp{op}}, nil
			case "and":
				op, err := ParseAndOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_AndOp{op}}, nil
			case "xor":
				op, err := ParseXorOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_XorOp{op}}, nil
			case "cocnat":
				op, err := ParseConcatOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_ConcatOp{op}}, nil
			case "valued":
				op, err := ParseValuedOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_ValuedOp{op}}, nil
			case "indexed":
				op, err := ParseIndexedOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_IndexedOp{op}}, nil
			}
			return nil, ProslParseErrOperation(key.(string))
		}
	}
	ob, err := ProslParsePrimitiveObject(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_Object{Object: ob}}, nil
}

func ParseOrderBy(yaml interface{}) (*proskenion.QueryOperator_OrderBy, error) {
	if value, ok := yaml.(map[interface{}]interface{}); ok {
		orderBy := &proskenion.QueryOperator_OrderBy{}
		for k, v := range value {
			switch k {
			case "key":
				if s, ok := v.(string); ok {
					orderBy.Key = s
				} else {
					return nil, ProslParseCastError("", s)
				}
			case "order":
				switch v {
				case "DESC":
					orderBy.Order = proskenion.QueryOperator_DESC
				default:
					orderBy.Order = proskenion.QueryOperator_ASC
				}
			default:
				return nil, ProslParseErrOperation(k)
			}
		}
		return orderBy, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml)
}

func ParseQueryOperator(yaml interface{}) (*proskenion.QueryOperator, error) {
	if v, ok := yaml.(map[interface{}]interface{}); ok {
		if len(v) > 2 {
			return nil, ProslParseArgumentErrorMin(2, len(v))
		}
		ret := &proskenion.QueryOperator{}
		mustFlags := 0
		for key, value := range v {
			switch key {
			case "select":
				if s, ok := value.(string); ok {
					ret.Select = s
				} else {
					return nil, ProslParseCastError("", value)
				}
				mustFlags |= 1
			case "type":
				if s, ok := value.(string); ok {
					code, err := ProslParseObjectCode(s)
					if err != nil {
						return nil, err
					}
					ret.Type = code
				} else {
					return nil, ProslParseCastError("", value)
				}
				mustFlags |= 2
			case "from":
				if s, ok := value.(string); ok {
					ret.From = s
				} else {
					return nil, ProslParseCastError("", value)
				}
				mustFlags |= 4
			case "where":
				op, err := ParseConditionalFormula(value)
				if err != nil {
					return nil, err
				}
				ret.Where = op
			case "order_by":
				op, err := ParseOrderBy(value)
				if err != nil {
					return nil, err
				}
				ret.OrderBy = op
			case "limit":
				if s, ok := value.(int32); ok {
					ret.Limit = s
				} else {
					return nil, ProslParseCastError(int32(0), value)
				}
			default:
				return nil, ProslParseErrOperation(key.(string))
			}
		}
		if mustFlags != 7 {
			var err error
			for _, e := range []struct {
				err error
			}{
				{fmt.Errorf("Must be select operand")},
				{fmt.Errorf("Must be type operand")},
				{fmt.Errorf("Must be from operand")},
			} {
				if mustFlags&1 == 1 {
					err = multierr.Append(err, e.err)
				}
				mustFlags >>= 1
			}
			return nil, errors.Wrapf(ErrProslParseQueryOperatorArgument, err.Error())
		}
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml)
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
