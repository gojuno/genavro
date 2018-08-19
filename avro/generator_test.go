package avro

import (
	"testing"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mkorolyov/astparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	cfg := astparser.Config{
		InputDir: "fixtures_test",
	}
	sources, err := astparser.Load(cfg)
	require.NoError(t, err)

	protocols := Generate(sources, "junolab.net")
	for name, protocol := range protocols {
		got, err := json.Marshal(protocol)
		require.NoError(t, err)
		want, err := ioutil.ReadFile(
			fmt.Sprintf("fixtures_test/%s.avpr", name))
		require.NoError(t, err)
		assert.Equal(t, want, got)
	}
}
