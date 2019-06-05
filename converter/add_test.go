package converter

import (
	"github.com/mbordner/kazaam/registry"
	"strconv"
	"testing"
)

func TestAdd_Convert(t *testing.T) {

	registry.RegisterConverter("add", &Add{})
	c := registry.GetConverter("add")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`5`, `1`, `6`,},
		{`5`, `-1`, `4`,},
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
