package test_utils

import "github.com/proskenion/proskenion/commit"

func RandomCommitProperty() *commit.CommitProperty {
	return &commit.CommitProperty{
		NumTxInBlock: 100,
	}
}
