package converter

import (
	"strconv"
	"testing"
	"github.com/mbordner/kazaam/registry"
)

func TestNot_Convert(t *testing.T) {

	registry.RegisterConverter("not", &Not{})
	c := registry.GetConverter("not")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`true"`, ``, `false`,},
		{`false`,``,`true`,},
		{`42"`, ``, `false`,},
		{`"42"`, ``, `false`,},
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