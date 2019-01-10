package prosl_test

import (
	. "github.com/proskenion/proskenion/prosl"
	"github.com/proskenion/proskenion/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestConvertYamlToMap(t *testing.T) {
	buf, err := ioutil.ReadFile("./example.yaml")
	require.NoError(t, err)

	yamap, err := ConvertYamlToMap(buf)
	require.NoError(t, err)

	assert.Equal(t, 2, len(yamap))
	{
		yamap_0 := yamap[0].(map[interface{}]interface{})
		assert.Equal(t, 1, len(yamap_0))
		assert.Contains(t, yamap_0, "set")

		setmap := yamap_0["set"].([]interface{})
		assert.Equal(t, 2, len(setmap))
		{
			setmap_0 := setmap[0].(string)
			assert.Equal(t, setmap_0, "peers")
		}
		{
			setmap_1 := setmap[1].(map[interface{}]interface{})
			assert.Equal(t, 1, len(setmap_1))
			assert.Contains(t, setmap[1], "query")

			querymap := setmap_1["query"].(map[interface{}]interface{})
			assert.Equal(t, 5, len(querymap))
			{
				assert.Contains(t, querymap, "select")
				selectmap := querymap["select"].(string)
				assert.Equal(t, selectmap, "peer")

				assert.Contains(t, querymap, "type")
				typemap := querymap["type"].(string)
				assert.Equal(t, typemap, "Peer")

				assert.Contains(t, querymap, "from")
				frommap := querymap["from"].(string)
				assert.Equal(t, frommap, "domain.com/peer")

				assert.Contains(t, querymap, "order_by")
				ordermap := querymap["order_by"].([]interface{})
				assert.Equal(t, 2, len(ordermap))
				{
					ordermap_0 := ordermap[0].(string)
					assert.Equal(t, ordermap_0, "fav")

					ordermap_1 := ordermap[1].(string)
					assert.Equal(t, ordermap_1, "DESC")
				}

				assert.Contains(t, querymap, "limit")
				limitmap := querymap["limit"].(int)
				assert.Equal(t, limitmap, 20)
			}
		}
	}
	{
		yamap_1 := yamap[1].(map[interface{}]interface{})
		assert.Equal(t, 1, len(yamap_1))
		assert.Contains(t, yamap_1, "return")

		returnmap := yamap_1["return"].(map[interface{}]interface{})
		{
			variablemap := returnmap["variable"]
			assert.Equal(t, "peers", variablemap.(string))
		}
	}
}

func TestConvertYamlToProbuf(t *testing.T) {
	buf, err := ioutil.ReadFile("./example.yaml")
	require.NoError(t, err)

	prosl, err := ConvertYamlToProtobuf(buf)
	require.NoError(t, err)

	setOp := prosl.GetOps()[0].GetSetOp()
	{
		// setOp variableName
		assert.Equal(t, "peers", setOp.GetVariableName())
		// setOp value is queryOp
		queryOp := setOp.GetValue().GetQueryOp()
		{
			assert.Equal(t, "peer", queryOp.GetSelect())
			assert.Equal(t, proskenion.ObjectCode_PeerObjectCode, queryOp.GetType())
			assert.Equal(t, "domain.com/peer", queryOp.GetFrom().GetObject().GetStr())
			assert.Equal(t, "fav", queryOp.GetOrderBy().GetKey())
			assert.Equal(t, proskenion.QueryOperator_DESC, queryOp.GetOrderBy().GetOrder())
			assert.Equal(t, int32(20), queryOp.GetLimit())
		}
	}

	returnOp := prosl.GetOps()[1].GetReturnOp()
	{
		// return Op is variableOp
		variableOp := returnOp.GetOp().GetVariableOp()
		{
			assert.Equal(t, "peers", variableOp.GetVariableName())
		}
	}
	returnOp.GetOp()
}
