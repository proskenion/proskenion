package prosl

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/proskenion/proskenion/config"
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/core/model"
	"github.com/proskenion/proskenion/proto"
)

type Prosl struct {
	prosl *proskenion.Prosl
	fc    model.ModelFactory
	rp    core.Repository
	c     core.Cryptor
	conf  *config.Config
}

func NewProsl(fc model.ModelFactory, rp core.Repository, c core.Cryptor, conf *config.Config) core.Prosl {
	return &Prosl{&proskenion.Prosl{}, fc, rp, c, conf}
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

func (p *Prosl) Execute() (model.Object, error) {
	if p.prosl == nil {
		return nil, errors.Errorf("Must be prosl setting, from yaml or protobuf binary")
	}
	state := ExecuteProsl(p.prosl, InitProslStateValue(p.fc, p.rp, p.conf))
	if state.Err != nil {
		return nil, state.Err
	}
	return state.ReturnObject, nil
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
