package kazaam

import (
	"fmt"
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

func TestDefaultTransformsSetCardinarily(t *testing.T) {
	if len(validSpecTypes) != 6 {
		t.Error("Unexpected number of default transforms. Missing tests?")
	}
}

func ExampleNewKazaam() {
	k, _ := NewKazaam(`[{"operation": "shift", "spec": {"output": "input"}}]`)
	kazaamOut, _ := k.TransformJSONStringToString(`{"input":"input value"}`)

	fmt.Println(kazaamOut)
	// Output:
	// {"output":"input value"}
}

func ExampleNew() {
	// Initialize a default Kazaam instance (i.e. same as NewKazaam(spec string))
	k, _ := New(`[{"operation": "shift", "spec": {"output": "input"}}]`, NewDefaultConfig())
	kazaamOut, _ := k.TransformJSONStringToString(`{"input":"input value"}`)

	fmt.Println(kazaamOut)
	// Output:
	// {"output":"input value"}
}

func ExampleConfig_RegisterTransform() {
	// use the default config to have access to built-in kazaam transforms
	kc := NewDefaultConfig()

	// register the new custom transform called "copy" which supports copying the
	// value of a top-level key to another top-level key
	kc.RegisterTransform("copy", func(spec *transform.Config, data *simplejson.Json) error {
		// the internal `Spec` will contain a mapping of source and target keys
		for targetField, sourceFieldInt := range *spec.Spec {
			sourceField := sourceFieldInt.(string)
			data.Set(targetField, data.Get(sourceField).Interface())
		}
		return nil
	})

	k, _ := New(`[{"operation": "copy", "spec": {"output": "input"}}]`, kc)
	kazaamOut, _ := k.TransformJSONStringToString(`{"input":"input value"}`)

	fmt.Println(kazaamOut)
	// Output:
	// {"input":"input value","output":"input value"}
}
