package config

import (
	"encoding/hex"
	"github.com/proskenion/proskenion/core/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	DB                DBConfig     `yaml:"db"`
	ProposalTxsLimits int          `yaml:"proposal_txs_limits"`
	Commit            CommitConfig `yaml:"commit"`
	Peer              PeerConfig   `yaml:"peer"`
	Prosl             ProslConfig  `yaml:"prosl"`
	Root              RootConfig   `yaml:"root"`
}

type DBConfig struct {
	Path string `yaml:"path"`
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}

type CommitConfig struct {
	NumTxInBlock int `yaml:"num_tx_in_block"`
}

type PeerConfig struct {
	PublicKey  string `yaml:"public_key"`
	PrivateKey string `yaml:"private_key"`
	Port       string `yaml:"port"`
}

type ProslConfig struct {
	Id        string             `yaml:"id"`
	Genesis   DefaultProslConfig `yaml:genesis`
	Incentive DefaultProslConfig `yaml:"incentive"`
	Consensus DefaultProslConfig `yaml:"consensus"`
	Update    DefaultProslConfig `yaml:"update"`
}

type DefaultProslConfig struct {
	Path string `yaml:"path"`
	Id   string `yaml:"id"`
}

type RootConfig struct {
	Id string `yaml:"id" default:"root@root"`
}

func (c PeerConfig) PublicKeyBytes() model.PublicKey {
	pub, err := hex.DecodeString(c.PublicKey)
	if err != nil {
		panic(err)
	}
	return pub
}

func (c PeerConfig) PrivateKeyBytes() model.PrivateKey {
	pri, err := hex.DecodeString(c.PrivateKey)
	if err != nil {
		panic(err)
	}
	return pri
}

func NewConfig(configPath string) *Config {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config := &Config{Root: RootConfig{Id: "root@root"}}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}
	return config
}
