package converter

import (
	"github.com/mbordner/kazaam/registry"
	"strconv"
	"testing"
)

func TestJoin_Convert(t *testing.T) {

	registry.RegisterConverter("join", &Join{})
	c := registry.GetConverter("join")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`["a","b","c"]`, `|`, `"a|b|c"`,},
	}

	for _, test := range table {
		v, e := c.Convert(getTestData(), []byte(test.value), []byte(strconv.Quote(test.arguments)))

		if e != nil {
			t.Error("error running convert function")
		}

		if string(v) != test.expected {
			t.Error("unexpected result from convert")
			t.Log("Expected: {}", test.expected)
			t.Log("Actual: {}", string(v))
		}
	}

}