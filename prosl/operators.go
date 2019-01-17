package prosl

import (
	"github.com/proskenion/proskenion/core/model"
)

// ========================= ValueOp =========================
func ExecutePlus(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() + b.GetI32())
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() + b.GetI64())
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() + b.GetU32())
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() + b.GetU64())
	case model.StringObjectCode:
		return builder.Str(a.GetStr() + b.GetStr())
	case model.AddressObjectCode:
		return builder.Address(a.GetAddress() + b.GetAddress())
	}
	return nil
}

func ExecuteMinus(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() - b.GetI32())
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() - b.GetI64())
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() - b.GetU32())
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() - b.GetU64())
	}
	return nil
}
func ExecuteMul(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() * b.GetI32())
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() * b.GetI64())
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() * b.GetU32())
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() * b.GetU64())
	}
	return nil
}
func ExecuteDiv(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() / b.GetI32())
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() / b.GetI64())
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() / b.GetU32())
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() / b.GetU64())
	}
	return nil
}
func ExecuteMod(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() % b.GetI32())
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() % b.GetI64())
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() % b.GetU32())
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() % b.GetU64())
	}
	return nil
}
func ExecuteOr(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() | b.GetI32())
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() | b.GetI64())
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() | b.GetU32())
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() | b.GetU64())
	}
	return nil
}
func ExecuteAnd(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() & b.GetI32())
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() & b.GetI64())
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() & b.GetU32())
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() & b.GetU64())
	}
	return nil
}

func ExecuteXor(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() ^ b.GetI32())
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() ^ b.GetI64())
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() ^ b.GetU32())
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() ^ b.GetU64())
	}
	return nil
}

func ExecuteConcat(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.StringObjectCode:
		return builder.Str(a.GetStr() + b.GetStr())
	case model.ListObjectCode:
		return builder.List(append(a.GetList(), b.GetList()...))
	}
	return nil
}

// ========================= Cond =========================
func ExecuteCondOr(a, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	if a.GetType() == model.BoolObjectCode {
		return builder.Bool(a.GetBoolean() && b.GetBoolean())
	}
	return nil
}

func ExecuteCondAnd(a, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil || a.GetType() != b.GetType() {
		return nil
	}
	builder := fc.NewObjectBuilder()
	if a.GetType() == model.BoolObjectCode {
		return builder.Bool(a.GetBoolean() && b.GetBoolean())
	}
	return nil
}

func ExecuteCondNot(a model.Object, fc model.ModelFactory) model.Object {
	if a == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	if a.GetType() == model.BoolObjectCode {
		return builder.Bool(!a.GetBoolean())
	}
	return nil
}

func ExecuteCondEq(os []model.Object, fc model.ModelFactory) model.Object {
	if len(os) < 2 {
		return nil
	}
	pr := os[0]
	for _, o := range os[1:] {
		if !model.ObjectEq(pr, o) {
			return fc.NewObjectBuilder().Bool(false)
		}
		pr = o
	}
	return fc.NewObjectBuilder().Bool(true)
}

func ExecuteCondNe(os []model.Object, fc model.ModelFactory) model.Object {
	if len(os) < 2 {
		return nil
	}
	st := make(map[string]struct{})
	for _, o := range os {
		if _, ok := st[string(o.Hash())]; ok {
			return fc.NewObjectBuilder().Bool(false)
		}
		st[string(o.Hash())] = struct{}{}
	}
	return fc.NewObjectBuilder().Bool(true)
}

func ExecuteCondGt(a, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	return fc.NewObjectBuilder().Bool(model.ObjectLess(b, a))
}

func ExecuteCondGe(a, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	return fc.NewObjectBuilder().Bool(model.ObjectLess(b, a) || model.ObjectEq(b, a))
}

func ExecuteCondLt(a, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	return fc.NewObjectBuilder().Bool(model.ObjectLess(a, b))
}

func ExecuteCondLe(a, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	return fc.NewObjectBuilder().Bool(model.ObjectLess(a, b) || model.ObjectEq(a, b))
}
