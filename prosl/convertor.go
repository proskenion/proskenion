package prosl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/proto"
	"gopkg.in/yaml.v2"
)

var (
	ErrConvertYamlToMap = fmt.Errorf("Failed Convert yaml to map")
)

func ConvertYamlToMap(yamlBytes []byte) ([]interface{}, error) {
	yamap := make([]interface{},0)
	err := yaml.Unmarshal(yamlBytes, &yamap)
	if err != nil {
		return nil, err
	}
	return yamap, nil
}

// pattern
// list -> []interface{}
// map -> map[string]interface{}
// string
// int
func ConvertYamlToProtobuf(yamlBytes []byte) (*proskenion.Prosl, error) {
	_, err := ConvertYamlToMap(yamlBytes)
	if err != nil {
		return nil, errors.Wrap(ErrConvertYamlToMap, err.Error())
	}

	return nil, nil
}
