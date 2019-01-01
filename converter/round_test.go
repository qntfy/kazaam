package converter

import (
	"github.com/qntfy/kazaam/registry"
	"strconv"
	"testing"
)

func TestRound_Convert(t *testing.T) {
	registry.RegisterConverter("round", &Round{})
	c := registry.GetConverter("round")

	table := []struct {
		value     string
		arguments string
		expected  string
	}{
		{`5.1`, ``, `5`,},
		{`5.6`, ``, `6`,},
		{`0.01`, ``, `0`},
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