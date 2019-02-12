package prosl

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"go.uber.org/multierr"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

var (
	ErrConvertYamlToMap          = fmt.Errorf("Failed Convert yaml to map")
	ErrConvertMapToProtobuf      = fmt.Errorf("Failed Convert Map to protobuf")
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

func ProslParseCastError(ex interface{}, ac interface{}, yaml interface{}) error {
	return errors.Wrapf(ErrProslParseUnExpectedCastType, "expected: %T, actual: %T, %#v", ex, ac, yaml)
}

func ProslParseUnknownCastError(ac interface{}, yaml interface{}) error {
	return errors.Wrapf(ErrProslParseUnExpectedCastType, "unknown type : %T, %#v, %#v", ac, ac, yaml)
}

func ProslParseArgumentError(ex int, ac int, yaml interface{}) error {
	return errors.Wrapf(ErrProslParseArgumentSize, "expected: %d, actual: %d, %#v", ex, ac, yaml)
}

func ProslParseArgumentErrorMin(minEx int, ac int, yaml interface{}) error {
	return errors.Wrapf(ErrProslParseArgumentSize, "expected moret than: %d, actual: %d, %#v", minEx, ac, yaml)
}

func ProslParseErrOperation(op interface{}, yaml interface{}) error {
	return errors.Wrapf(ErrProslParseUnExpectedOperationName, "unknown operation: %#v\n%#v", op, yaml)
}

func ProslParseErrCode(code string) proskenion.ErrCode {
	switch code {
	case "AnythingErrCode":
		return proskenion.ErrCode_Anything
	}
	return proskenion.ErrCode_Anything
}

func ProslParsePrimitiveObject(value interface{}) (*proskenion.Object, error) {
	switch s := value.(type) {
	case int:
		return &proskenion.Object{
			Type:   proskenion.ObjectCode_Int32ObjectCode,
			Object: &proskenion.Object_I32{int32(s)},
		}, nil
	case string:
		// case : int64, suffix has "LL" or "ll"
		if strings.HasSuffix(s, "ll") || strings.HasSuffix(s, "LL") {
			ret, err := strconv.ParseInt(s[:len(s)-2], 10, 64)
			if err == nil {
				return &proskenion.Object{
					Type:   proskenion.ObjectCode_Int64ObjectCode,
					Object: &proskenion.Object_I64{ret},
				}, nil
			}
		}
		// case : int32, can convert to int
		if ret, err := strconv.ParseInt(s, 10, 32); err == nil {
			return &proskenion.Object{
				Type:   proskenion.ObjectCode_Int32ObjectCode,
				Object: &proskenion.Object_I32{int32(ret)},
			}, nil
		}
		// case : binary, prefix is 0x
		if strings.HasPrefix(s, "0x") {
			data, err := hex.DecodeString(s[2:])
			if err == nil {
				return &proskenion.Object{
					Type:   proskenion.ObjectCode_BytesObjectCode,
					Object: &proskenion.Object_Data{data},
				}, nil
			}
		}
		// case : address, can convert address
		if _, err := model.NewAddress(s); err == nil {
			return &proskenion.Object{
				Type:   proskenion.ObjectCode_AddressObjectCode,
				Object: &proskenion.Object_Address{s},
			}, nil
		}
		// case : string, another case
		return &proskenion.Object{
			Type:   proskenion.ObjectCode_StringObjectCode,
			Object: &proskenion.Object_Str{s},
		}, nil
	case bool:
		return &proskenion.Object{
			Type:   proskenion.ObjectCode_BoolObjectCode,
			Object: &proskenion.Object_Boolean{s},
		}, nil
	}
	return nil, ProslParseUnknownCastError(value, value)
}

func ProslParseObjectCode(yaml interface{}) (proskenion.ObjectCode, error) {
	s, ok := yaml.(string)
	if !ok {
		return 0, ProslParseCastError("", yaml, yaml)
	}
	code := strings.ToLower(s)
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
	case "storage":
		return proskenion.ObjectCode_StorageObjectCode, nil
	case "command":
		return proskenion.ObjectCode_CommandObjectCode, nil
	case "transaction":
		return proskenion.ObjectCode_TransactionObjectCode, nil
	case "block":
		return proskenion.ObjectCode_BlockObjectCode, nil
	}
	return 0, errors.Wrapf(ErrProslParseUnknownObjectCode, "unknown type : %s", code)
}

func ParseProsl(yalist []interface{}) (*proskenion.Prosl, error) {
	ops := make([]*proskenion.ProslOperator, 0, len(yalist))
	for _, ya := range yalist {
		switch v := ya.(type) {
		case map[interface{}]interface{}:
			pop, err := ParseProslOperator(v)
			if err != nil {
				return nil, err
			}
			ops = append(ops, pop)
		default:
			return nil, errors.Wrap(ErrProslParseNotExpectedType, fmt.Sprintf("%T, %#v", v, yalist))
		}
	}
	return &proskenion.Prosl{Ops: ops}, nil
}

// ParseProslOperator
func ParseProslOperator(yamap map[interface{}]interface{}) (*proskenion.ProslOperator, error) {
	if len(yamap) != 1 {
		return nil, ProslParseArgumentError(1, len(yamap), yamap)
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
		case "return":
			op, err := ParseReturnOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_ReturnOp{ReturnOp: op}}, nil
		case "each":
			op, err := ParseEachOperator(value)
			if err != nil {
				return nil, err
			}
			return &proskenion.ProslOperator{Op: &proskenion.ProslOperator_EachOp{EachOp: op}}, nil
		default:
			return nil, ProslParseErrOperation(key, yamap)
		}
	}
	return nil, ErrProslParseInternalErr
}

// set:
// 	- variableName (string)
//  - valueOperator (interface{})
func ParseSetOperator(yaml interface{}) (*proskenion.SetOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) != 2 {
			return nil, ProslParseArgumentError(2, len(yalist), yaml)
		}
		ret := &proskenion.SetOperator{}

		// variableName
		if v, ok := yalist[0].(string); ok {
			ret.VariableName = v
		} else {
			return nil, ProslParseCastError("", v, yaml)
		}

		// valueOperator
		if vop, err := ParseValueOperator(yalist[1].(interface{})); err != nil {
			return nil, err
		} else {
			ret.Value = vop
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseIfOperator(yaml interface{}) (*proskenion.IfOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) < 2 {
			return nil, ProslParseArgumentErrorMin(2, len(yalist), yaml)
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
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseElifOperator(yaml interface{}) (*proskenion.ElifOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) < 2 {
			return nil, ProslParseArgumentErrorMin(2, len(yalist), yaml)
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
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseElseOperator(yaml interface{}) (*proskenion.ElseOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) < 1 {
			return nil, ProslParseArgumentErrorMin(1, len(yalist), yaml)
		}
		ret := &proskenion.ElseOperator{}
		prosl, err := ParseProsl(yalist)
		if err != nil {
			return nil, err
		}
		ret.Prosl = prosl
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseErrCatchOperator(yaml interface{}) (*proskenion.ErrCatchOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) < 2 {
			return nil, ProslParseArgumentErrorMin(2, len(yalist), yaml)
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
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
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

func ParseReturnOperator(yaml interface{}) (*proskenion.ReturnOperator, error) {
	op, err := ParseValueOperator(yaml)
	if err != nil {
		return nil, err
	}
	ret := &proskenion.ReturnOperator{Op: op}
	return ret, nil
}

func ParseEachOperator(yaml interface{}) (*proskenion.EachOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) < 2 {
			return nil, ProslParseArgumentErrorMin(3, len(yalist), yaml)
		}
		ret := &proskenion.EachOperator{}
		op, err := ParseValueOperator(yalist[0])
		if err != nil {
			return nil, err
		}
		ret.List = op

		// variableName
		if v, ok := yalist[1].(string); ok {
			ret.VariableName = v
		} else {
			return nil, ProslParseCastError("", v, yaml)
		}

		prosl, err := ParseProsl(yalist[2:])
		if err != nil {
			return nil, err
		}
		ret.Do = prosl
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseValueOperator(yaml interface{}) (*proskenion.ValueOperator, error) {
	if v, ok := yaml.(map[interface{}]interface{}); ok {
		if len(v) != 1 {
			return nil, ProslParseArgumentError(1, len(v), yaml)
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
			case "storage":
				op, err := ParseStorageOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_StorageOp{op}}, nil
			case "map":
				op, err := ParseMapOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_MapOp{op}}, nil
			case "list":
				op, err := ParseListOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_ListOp{op}}, nil
			case "plus":
				op, err := ParsePlusOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_PlusOp{op}}, nil
			case "minus":
				op, err := ParseMinusOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_MinusOp{op}}, nil
			case "mult":
				op, err := ParseMultipleOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_MulOp{op}}, nil
			case "div":
				op, err := ParseDivisionOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_DivOp{op}}, nil
			case "mod":
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
			case "concat":
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
			case "variable", "var":
				op, err := ParseVariableOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_VariableOp{op}}, nil
			case "cast":
				op, err := ParseCastOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_CastOp{op}}, nil
			case "list_comprehension", "list_comp", "comprehension", "comp":
				op, err := ParseListComprehensionOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_ListComprehensionOp{op}}, nil
			case "sort":
				op, err := ParseSortOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_SortOp{op}}, nil
			case "slice":
				op, err := ParseSliceOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_SliceOp{op}}, nil
			case "is_defined":
				op, err := ParseIsDefinedOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_IsDefinedOp{op}}, nil
			case "verify":
				op, err := ParseVerifyOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_VerifyOp{op}}, nil
			case "pagerank":
				op, err := ParsePageRankOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_PageRankOp{op}}, nil
			case "len":
				op, err := ParseLenOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_LenOp{op}}, nil
			default: // another case, all command
				op, err := ParseCommandOperator(v)
				if err != nil {
					return nil, err
				}
				return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_CmdOp{op}}, nil
			}
		}
	} else if _, ok := yaml.([]interface{}); ok { // list op
		op, err := ParseListOperator(yaml)
		if err != nil {
			return nil, err
		}
		return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_ListOp{op}}, nil
	}
	ob, err := ProslParsePrimitiveObject(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.ValueOperator{Op: &proskenion.ValueOperator_Object{Object: ob}}, nil
}

func ParseOrderBy(yaml interface{}) (*proskenion.OrderBy, error) {
	if value, ok := yaml.([]interface{}); ok {
		orderBy := &proskenion.OrderBy{}

		// 0-index key
		if s, ok := value[0].(string); ok {
			orderBy.Key = s
		} else {
			return nil, ProslParseCastError("", s, yaml)
		}

		// 1-index order
		switch value[1] {
		case "DESC":
			orderBy.Order = proskenion.OrderCode_DESC
		default:
			orderBy.Order = proskenion.OrderCode_ASC
		}
		return orderBy, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseQueryOperator(yaml interface{}) (*proskenion.QueryOperator, error) {
	if v, ok := yaml.(map[interface{}]interface{}); ok {
		if len(v) < 2 {
			return nil, ProslParseArgumentErrorMin(2, len(v), yaml)
		}
		ret := &proskenion.QueryOperator{}
		mustFlags := 0
		for key, value := range v {
			switch key {
			case "authorizer", "authorizer_id":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.AuthorizerId = op
			case "select":
				if s, ok := value.(string); ok {
					ret.Select = s
				} else {
					return nil, ProslParseCastError("", value, yaml)
				}
				mustFlags |= 1
			case "type":
				code, err := ProslParseObjectCode(value)
				if err != nil {
					return nil, err
				}
				ret.Type = code
				mustFlags |= 2
			case "from", "from_id":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.From = op
				mustFlags |= 4
			case "where":
				op, err := ParseValueOperator(value)
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
				// TODO : int64未対応
				if s, ok := value.(int); ok {
					ret.Limit = int32(s)
				} else {
					return nil, ProslParseCastError(int32(0), value, yaml)
				}
			default:
				return nil, ProslParseErrOperation(key, yaml)
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
				if mustFlags&1 == 0 {
					err = multierr.Append(err, e.err)
				}
				mustFlags >>= 1
			}
			return nil, errors.Wrapf(ErrProslParseQueryOperatorArgument, err.Error())
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml, yaml)
}

func ParseCommandOperator(yaml interface{}) (*proskenion.CommandOperator, error) {
	if yamap, ok := yaml.(map[interface{}]interface{}); ok {
		if len(yamap) != 1 {
			return nil, ProslParseArgumentError(1, len(yamap), yaml)
		}
		ret := &proskenion.CommandOperator{Params: make(map[string]*proskenion.ValueOperator)}
		for key, value := range yamap {
			if ks, ok := key.(string); ok {
				ret.CommandName = ks
			} else {
				return nil, ProslParseCastError("", key, yaml)
			}
			if yam, ok := value.(map[interface{}]interface{}); ok {
				for k, v := range yam {
					s, ok := k.(string)
					if !ok {
						return nil, ProslParseCastError("", k, yaml)
					}
					vop, err := ParseValueOperator(v)
					if err != nil {
						return nil, err
					}
					ret.Params[s] = vop
				}
			}
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml, yaml)
}

func ParseCommandsOperator(yaml interface{}) ([]*proskenion.ValueOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			cmd, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ret = append(ret, cmd)
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseTxOperator(yaml interface{}) (*proskenion.TxOperator, error) {
	if yamap, ok := yaml.(map[interface{}]interface{}); ok {
		ret := &proskenion.TxOperator{}
		for key, value := range yamap {
			switch key {
			case "commands":
				v, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.Commands = v
			default:
				return nil, ProslParseUnknownCastError(key, yaml)
			}
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml, yaml)
}

func ParseStorageOperator(yaml interface{}) (*proskenion.StorageOperator, error) {
	ret, err := ParseMapOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.StorageOperator{Object: ret}, nil
}

func ParsePlusOperator(yaml interface{}) (*proskenion.PlusOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.PlusOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseMinusOperator(yaml interface{}) (*proskenion.MinusOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.MinusOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseMultipleOperator(yaml interface{}) (*proskenion.MultipleOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.MultipleOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)

}

func ParseDivisionOperator(yaml interface{}) (*proskenion.DivisionOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.DivisionOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)

}

func ParseModOperator(yaml interface{}) (*proskenion.ModOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.ModOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)

}

func ParseOrOperator(yaml interface{}) (*proskenion.OrOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.OrOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)

}

func ParseAndOperator(yaml interface{}) (*proskenion.AndOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.AndOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)

}

func ParseXorOperator(yaml interface{}) (*proskenion.XorOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.XorOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)

}

func ParseConcatOperator(yaml interface{}) (*proskenion.ConcatOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ret := &proskenion.ConcatOperator{}
		ops := make([]*proskenion.ValueOperator, 0, len(yalist))
		for _, value := range yalist {
			v, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, v)
		}
		ret.Ops = ops
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)

}

func ParseValuedOperator(yaml interface{}) (*proskenion.ValuedOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) != 3 {
			return nil, ProslParseArgumentError(3, len(yalist), yaml)
		}
		ret := &proskenion.ValuedOperator{}

		// 0 - value operator
		value, err := ParseValueOperator(yalist[0])
		if err != nil {
			return nil, err
		}
		ret.Object = value

		// 1 - type operator
		t, err := ProslParseObjectCode(yalist[1])
		if err != nil {
			return nil, err
		}
		ret.Type = t

		// 2 - key operator
		keys, ok := yalist[2].(string)
		if !ok {
			return nil, ProslParseCastError("", yalist, yaml)
		}
		ret.Key = keys
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseIndexedOperator(yaml interface{}) (*proskenion.IndexedOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) != 3 {
			return nil, ProslParseArgumentError(3, len(yalist), yaml)
		}
		ret := &proskenion.IndexedOperator{}

		// 0 - value operator
		value, err := ParseValueOperator(yalist[0])
		if err != nil {
			return nil, err
		}
		ret.Object = value

		// 1 - type operator
		t, err := ProslParseObjectCode(yalist[1])
		if err != nil {
			return nil, err
		}
		ret.Type = t

		// 2 - index operator
		index, err := ParseValueOperator(yalist[2])
		if err != nil {
			return nil, err
		}
		ret.Index = index
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseVariableOperator(yaml interface{}) (*proskenion.VariableOperator, error) {
	if s, ok := yaml.(string); ok {
		return &proskenion.VariableOperator{VariableName: s}, nil
	} else {
		return nil, ProslParseCastError("", yaml, yaml)
	}
}

func ParseCastOperator(yaml interface{}) (*proskenion.CastOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		if len(yalist) != 2 {
			return nil, ProslParseArgumentError(2, len(yalist), yaml)
		}
		ret := &proskenion.CastOperator{}

		// 0 - type operator
		types, err := ProslParseObjectCode(yalist[0])
		if err != nil {
			return nil, err
		}
		ret.Type = types

		// 1 - value operator
		value, err := ParseValueOperator(yalist[1])
		if err != nil {
			return nil, err
		}
		ret.Object = value
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseListComprehensionOperator(yaml interface{}) (*proskenion.ListComprehensionOperator, error) {
	if yamap, ok := yaml.(map[interface{}]interface{}); ok {
		if len(yamap) < 3 {
			return nil, ProslParseArgumentErrorMin(3, len(yamap), yaml)
		}
		ret := &proskenion.ListComprehensionOperator{}
		for key, value := range yamap {
			switch key {
			case "list":
				list, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.List = list
			case "var", "variable", "variable_name":
				variable, ok := value.(string)
				if !ok {
					return nil, ProslParseCastError("", value, yaml)
				}
				ret.VariableName = variable
			case "if":
				condIF, err := ParseConditionalFormula(value)
				if err != nil {
					return nil, err
				}
				ret.If = condIF
			case "element", "elem", "value":
				element, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.Element = element
			}
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseSortOperator(yaml interface{}) (*proskenion.SortOperator, error) {
	if v, ok := yaml.(map[interface{}]interface{}); ok {
		if len(v) < 1 {
			return nil, ProslParseArgumentErrorMin(1, len(v), yaml)
		}
		ret := &proskenion.SortOperator{}
		mustFlags := 0
		for key, value := range v {
			switch key {
			case "list":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.List = op
				mustFlags |= 1
			case "order_by":
				op, err := ParseOrderBy(value)
				if err != nil {
					return nil, err
				}
				ret.OrderBy = op
			case "obj_code", "object_code", "code", "type":
				op, err := ProslParseObjectCode(value)
				if err != nil {
					return nil, err
				}
				ret.Type = op
			case "limit":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.Limit = op
			}
		}
		if mustFlags == 0 {
			return nil, errors.Wrapf(ErrProslParseQueryOperatorArgument, "sort operator must be list. %#v", yaml)
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml, yaml)
}

func ParseSliceOperator(yaml interface{}) (*proskenion.SliceOperator, error) {
	if v, ok := yaml.(map[interface{}]interface{}); ok {
		if len(v) < 1 {
			return nil, ProslParseArgumentErrorMin(1, len(v), yaml)
		}
		ret := &proskenion.SliceOperator{}
		mustFlags := 0
		for key, value := range v {
			switch key {
			case "list":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.List = op
				mustFlags |= 1
			case "left", "l":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.Left = op
			case "right", "r":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.Right = op
			}
		}
		if mustFlags == 0 {
			return nil, errors.Wrapf(ErrProslParseQueryOperatorArgument, "sort operator must be list. %#v", yaml)
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml, yaml)
}

func ParseIsDefinedOperator(yaml interface{}) (*proskenion.IsDefinedOperator, error) {
	if variableName, ok := yaml.(string); ok {
		return &proskenion.IsDefinedOperator{VariableName: variableName}, nil
	}
	return nil, ProslParseCastError("", yaml, yaml)
}

func ParseVerifyOperator(yaml interface{}) (*proskenion.VerifyOperator, error) {
	if yamap, ok := yaml.(map[interface{}]interface{}); ok {
		ret := &proskenion.VerifyOperator{}
		for key, value := range yamap {
			switch key {
			case "sig", "signature":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.Sig = op
			case "hasher", "hash":
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.Hash = op
			}
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml, yaml)
}

func ParsePageRankOperator(yaml interface{}) (*proskenion.PageRankOperator, error) {
	if yamap, ok := yaml.(map[interface{}]interface{}); ok {
		ret := &proskenion.PageRankOperator{}
		for key, value := range yamap {
			switch key {
			case "storages", "edges":
				v, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.Storages = v
			case "to_key", "toKey", "tokey", "key":
				v, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.ToKey = v
			case "out_name", "outName", "outname", "name", "out":
				v, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				ret.OutName = v
			}
		}
		return ret, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml, yaml)
}

func ParseLenOperator(yaml interface{}) (*proskenion.LenOperator, error) {
	v, err := ParseValueOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.LenOperator{List: v}, nil
}

func ParseListOperator(yaml interface{}) (*proskenion.ListOperator, error) {
	vops := make([]*proskenion.ValueOperator, 0)
	if list, ok := yaml.([]interface{}); ok {
		for _, value := range list {
			op, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			vops = append(vops, op)
		}
		return &proskenion.ListOperator{Object: vops}, nil
	} else if s, ok := yaml.(string); ok {
		s = strings.ToLower(s)
		if s == "nil" {
			return &proskenion.ListOperator{Object: vops}, nil
		}
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseMapOperator(yaml interface{}) (*proskenion.MapOperator, error) {
	vops := make(map[string]*proskenion.ValueOperator)
	if yamap, ok := yaml.(map[interface{}]interface{}); ok {
		for key, value := range yamap {
			if s, ok := key.(string); ok {
				op, err := ParseValueOperator(value)
				if err != nil {
					return nil, err
				}
				vops[s] = op
			} else {
				return nil, ProslParseCastError("", yaml, yaml)
			}
		}
		return &proskenion.MapOperator{Object: vops}, nil
	}
	return nil, ProslParseCastError(make(map[interface{}]interface{}), yaml, yaml)
}

func ParseConditionalFormula(yaml interface{}) (*proskenion.ConditionalFormula, error) {
	if yamap, ok := yaml.(map[interface{}]interface{}); ok {
		if len(yamap) != 1 {
			return nil, ProslParseCastError("", yaml, yaml)
		}
		for key, value := range yamap {
			switch key {
			case "or":
				op, err := ParseOrFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_Or{Or: op}}, nil
			case "and":
				op, err := ParseAndFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_And{And: op}}, nil
			case "not":
				op, err := ParseNotFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_Not{Not: op}}, nil
			case "eq":
				op, err := ParseEqFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_Eq{Eq: op}}, nil
			case "ne":
				op, err := ParseNeFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_Ne{Ne: op}}, nil
			case "gt":
				op, err := ParseGtFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_Gt{Gt: op}}, nil
			case "ge":
				op, err := ParseGeFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_Ge{Ge: op}}, nil
			case "lt":
				op, err := ParseLtFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_Lt{Lt: op}}, nil
			case "le":
				op, err := ParseLeFormula(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_Le{Le: op}}, nil
			case "verify":
				op, err := ParseVerifyOperator(value)
				if err != nil {
					return nil, err
				}
				return &proskenion.ConditionalFormula{Op: &proskenion.ConditionalFormula_VerifyOp{VerifyOp: op}}, nil
			}
		}
	}
	return nil, ProslParseErrOperation(yaml, yaml)
}

func ParsePolynomialOperator(yaml interface{}) ([]*proskenion.ValueOperator, error) {
	if yalist, ok := yaml.([]interface{}); ok {
		ops := make([]*proskenion.ValueOperator, 0)
		for _, value := range yalist {
			op, err := ParseValueOperator(value)
			if err != nil {
				return nil, err
			}
			ops = append(ops, op)
		}
		return ops, nil
	}
	return nil, ProslParseCastError(make([]interface{}, 0), yaml, yaml)
}

func ParseOrFormula(yaml interface{}) (*proskenion.OrFormula, error) {
	ops, err := ParsePolynomialOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.OrFormula{Ops: ops}, nil
}

func ParseAndFormula(yaml interface{}) (*proskenion.AndFormula, error) {
	ops, err := ParsePolynomialOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.AndFormula{Ops: ops}, nil
}

func ParseNotFormula(yaml interface{}) (*proskenion.NotFormula, error) {
	op, err := ParseValueOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.NotFormula{Op: op}, nil
}

func ParseEqFormula(yaml interface{}) (*proskenion.EqFormula, error) {
	ops, err := ParsePolynomialOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.EqFormula{Ops: ops}, nil
}

func ParseNeFormula(yaml interface{}) (*proskenion.NeFormula, error) {
	ops, err := ParsePolynomialOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.NeFormula{Ops: ops}, nil
}

func ParseGtFormula(yaml interface{}) (*proskenion.GtFormula, error) {
	ops, err := ParsePolynomialOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.GtFormula{Ops: ops}, nil
}
func ParseGeFormula(yaml interface{}) (*proskenion.GeFormula, error) {
	ops, err := ParsePolynomialOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.GeFormula{Ops: ops}, nil
}

func ParseLtFormula(yaml interface{}) (*proskenion.LtFormula, error) {
	ops, err := ParsePolynomialOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.LtFormula{Ops: ops}, nil
}
func ParseLeFormula(yaml interface{}) (*proskenion.LeFormula, error) {
	ops, err := ParsePolynomialOperator(yaml)
	if err != nil {
		return nil, err
	}
	return &proskenion.LeFormula{Ops: ops}, nil
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
	if err != nil {
		return nil, errors.Wrap(ErrConvertMapToProtobuf, err.Error())
	}
	return prosl, nil
}
