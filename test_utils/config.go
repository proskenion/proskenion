package test_utils

import "github.com/proskenion/proskenion/config"

func NewTestConfig() *config.Config {
	config := config.NewConfig("../config/config.yaml")
	config.DB.Path = "../database"
	return config
}
