package prosl

import (
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
	"github.com/satellitex/protobuf/proto"
)

type Prosl struct {
	prosl *proskenion.Prosl
	fc    model.ModelFactory
	c     core.Cryptor
	conf  *config.Config
}

func NewProsl(fc model.ModelFactory, c core.Cryptor, conf *config.Config) core.Prosl {
	return &Prosl{&proskenion.Prosl{}, fc, c, conf}
}

func (p *Prosl) ConvertFromYaml(yaml []byte) error {
	prosl, err := ConvertYamlToProtobuf(yaml)
	if err != nil {
		return err
	}
	p.prosl = prosl
	return nil
}

func (p *Prosl) Validate() error {
	return nil
}

func (p *Prosl) Execute(wsv model.ObjectFinder, top model.Block) (model.Object, map[string]model.Object, error) {
	if p.prosl == nil {
		return nil, nil, errors.Errorf("Must be prosl setting, from yaml or protobuf binary")
	}
	state := ExecuteProsl(p.prosl, InitProslStateValue(p.fc, wsv, top, p.c, p.conf))
	if state.Err != nil {
		return nil, state.Variables, state.Err
	}
	return state.ReturnObject, state.Variables, nil
}

func (p *Prosl) ExecuteWithParams(wsv model.ObjectFinder, top model.Block, params map[string]model.Object) (model.Object, map[string]model.Object, error) {
	if p.prosl == nil {
		return nil, nil, errors.Errorf("Must be prosl setting, from yaml or protobuf binary")
	}
	state := ExecuteProsl(p.prosl, InitProslStateValueWithPrams(p.fc, wsv, top, p.c, p.conf, params))
	if state.Err != nil {
		return nil, state.Variables, state.Err
	}
	return state.ReturnObject, state.Variables, nil
}

func (p *Prosl) Unmarshal(proslData []byte) error {
	err := proto.Unmarshal(proslData, p.prosl)
	if err != nil {
		return err
	}
	return nil
}

func (p *Prosl) Marshal() ([]byte, error) {
	return proto.Marshal(p.prosl)
}

func (p *Prosl) Hash() model.Hash {
	return p.c.Hash(p)
}
