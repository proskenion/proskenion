package commit

import "github.com/proskenion/proskenion/core/model"

type CommitProperty struct {
	NumTxInBlock int
	PublicKey model.PublicKey
	PrivateKey model.PrivateKey
}
