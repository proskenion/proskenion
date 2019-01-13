package prosl

import (
	"github.com/proskenion/proskenion/core"
	"github.com/proskenion/proskenion/query"
	"github.com/proskenion/proskenion/repository"
	. "github.com/proskenion/proskenion/test_utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestExecuteProsl(t *testing.T) {
	buf, err := ioutil.ReadFile("./example.yaml")
	require.NoError(t, err)

	prosl, err := ConvertYamlToProtobuf(buf)
	require.NoError(t, err)

	dba := RandomDBA()
	cryptor := RandomCryptor()
	fc := NewTestFactory()
	rp := repository.NewRepository(dba, cryptor, fc)
	rtx, err := rp.Begin()
	require.NoError(t, err)

	conf := NewTestConfig()
	wsv, err := rtx.WSV(nil)
	require.NoError(t, err)

	qp := query.NewQueryProcessor(rp, fc, conf)
	qv := query.NewQueryValidator(rp, fc, conf)
	qc := query.NewQueryVerifier()

	ExecuteProsl(prosl, InitProslStateValue(wsv, fc, struct {
		core.QueryProcessor
		core.QueryValidator
		core.QueryVerifier
	}{qp, qv, qc}))
}
