package core

import "github.com/proskenion/proskenion/core/model"

type Prosl interface {
	ConvertFromYaml(yaml []byte) error
	Validate() error
	Execute() (model.Object, map[string]model.Object, error)
	ExecuteWithParams(map[string]model.Object) (model.Object, map[string]model.Object, error)
	model.Modelor
}
