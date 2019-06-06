package converter

import (
	"github.com/mbordner/kazaam/registry"
	"strconv"
	"testing"
)

func TestFloat_Convert(t *testing.T) {

	registry.RegisterConverter("float", &Float{})
	c := registry.GetConverter("float")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`5`, `1`, `5.0`,},
		{`5.01`, `2`, `5.01`,},
		{`5.012`,`1`,`5.0`},
		{`7.77`,`1`,`7.8`},
		{`500.01`,`1`,`500.0`},
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
