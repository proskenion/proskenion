package grpc_test

import "github.com/proskenion/proskenion/core/model"

func MakeEdgeStorageFromObjects(fc model.ModelFactory, objs []model.Object) model.Storage {
	objList := make([]model.Object, 0)
	for _, o := range objs {
		objList = append(objList, o)
	}
	st := fc.NewStorageBuilder().
		List(TrustStorage, objList).
		Build()
	return st
}
