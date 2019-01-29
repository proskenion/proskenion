package test_utils

import (
	"github.com/proskenion/proskenion/config"
)

func RandomConfig() *config.Config {
	config := config.NewConfig("../config/config.yaml")

	config.DB.Path = "../database"
	config.Prosl.Genesis.Path = "../test_utils/genesis.yaml"
	config.Prosl.Update.Path = "../test_utils/update.yaml"
	config.Prosl.Incentive.Path = "../test_utils/incentive.yaml"
	config.Prosl.Consensus.Path = "../test_utils/consensus.yaml"
	return config
}
