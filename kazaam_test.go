package kazaam

import (
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/qntfy/kazaam/transform"
)

func TestDefaultKazaamGetUnknownTransform(t *testing.T) {
	_, err := NewKazaam(`[{"operation": "doesnt-exist"}]`)

	if err == nil {
		t.Error("Should have thrown error for unknown transform")
	}
}

func TestKazaamWithRegisteredTransform(t *testing.T) {
	kc := NewDefaultConfig()
	kc.RegisterTransform("3rd-party", func(spec *transform.Config, data *simplejson.Json) error {
		data.Set("doesnt-exist", "does-exist")
		return nil
	})
	_, err := New(`[{"operation": "3rd-party"}]`, kc)
	if err != nil {
		t.Error("Shouldn't have thrown error for registered 3rd-party transform")
	}
}

func TestReregisterKazaamTransform(t *testing.T) {
	kc := NewDefaultConfig()
	err := kc.RegisterTransform("shift", nil)

	if err == nil {
		t.Error("Should have thrown error for duplicated transform name")
	}
}
