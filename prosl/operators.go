package prosl

import (
	"github.com/proskenion/proskenion/core/model"
)

func ExecutePlus(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() + b.GetI32()).Build()
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() + b.GetI64()).Build()
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() + b.GetU32()).Build()
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() + b.GetU64()).Build()
	case model.StringObjectCode:
		return builder.Str(a.GetStr() + b.GetStr()).Build()
	case model.AddressObjectCode:
		return builder.Address(a.GetAddress() + b.GetAddress()).Build()
	}
	return nil
}

func ExecuteMinus(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() - b.GetI32()).Build()
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() - b.GetI64()).Build()
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() - b.GetU32()).Build()
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() - b.GetU64()).Build()
	}
	return nil
}
func ExecuteMul(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() * b.GetI32()).Build()
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() * b.GetI64()).Build()
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() * b.GetU32()).Build()
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() * b.GetU64()).Build()
	}
	return nil
}
func ExecuteDiv(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() / b.GetI32()).Build()
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() / b.GetI64()).Build()
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() / b.GetU32()).Build()
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() / b.GetU64()).Build()
	}
	return nil
}
func ExecuteMod(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() % b.GetI32()).Build()
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() % b.GetI64()).Build()
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() % b.GetU32()).Build()
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() % b.GetU64()).Build()
	}
	return nil
}
func ExecuteOr(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() | b.GetI32()).Build()
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() | b.GetI64()).Build()
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() | b.GetU32()).Build()
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() | b.GetU64()).Build()
	}
	return nil
}
func ExecuteAnd(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() & b.GetI32()).Build()
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() & b.GetI64()).Build()
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() & b.GetU32()).Build()
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() & b.GetU64()).Build()
	}
	return nil
}

func ExecuteXor(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.Int32ObjectCode:
		return builder.Int32(a.GetI32() ^ b.GetI32()).Build()
	case model.Int64ObjectCode:
		return builder.Int64(a.GetI64() ^ b.GetI64()).Build()
	case model.Uint32ObjectCode:
		return builder.Uint32(a.GetU32() ^ b.GetU32()).Build()
	case model.Uint64ObjectCode:
		return builder.Uint64(a.GetU64() ^ b.GetU64()).Build()
	}
	return nil
}

func ExecuteConcat(a model.Object, b model.Object, fc model.ModelFactory) model.Object {
	if a == nil || b == nil {
		return nil
	}
	builder := fc.NewObjectBuilder()
	switch a.GetType() {
	case model.StringObjectCode:
		return builder.Str(a.GetStr() + b.GetStr()).Build()
	case model.ListObjectCode:
		return builder.List(append(a.GetList(), b.GetList()...)).Build()
	}
	return nil
}
