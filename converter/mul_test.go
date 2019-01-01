package converter

import (
	"github.com/qntfy/kazaam/registry"
	"strconv"
	"testing"
)

func TestMul_Convert(t *testing.T) {
	registry.RegisterConverter("mul", &Mul{})
	c := registry.GetConverter("mul")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`5`, `1`, `5`,},
		{`5`, `2`, `10`,},
		{`5`, `2.5`, `12.5`},
		{`10`, `.5`, `5`},
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
