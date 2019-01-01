package converter

import (
	"github.com/qntfy/kazaam/registry"
	"strconv"
	"testing"
)

func TestUpper_Convert(t *testing.T) {
	registry.RegisterConverter("upper", &Upper{})
	c := registry.GetConverter("upper")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`"this is a test"`, ``, `"THIS IS A TEST"`,},
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