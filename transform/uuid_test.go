package transform

import (
	"testing"

	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
)

func TestUUIDV4(t *testing.T) {
	spec := `{"a.uuid":{"version": 4}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(UUID, cfg, jsonIn)
	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err := uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.FailNow()
	}
}

func TestUUIDVersionError(t *testing.T) {
	spec := `{"a.uuid":{"version": 6}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(UUID, cfg, jsonIn)
	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err := uuid.FromString(string(out))

	if err == nil {
		t.Error("Shouldn't be able to generate UUID due to incompatible version")
		t.FailNow()
	}
}

func TestUUIDNoVersionError(t *testing.T) {
	spec := `{"a.uuid":{}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(UUID, cfg, jsonIn)
	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err := uuid.FromString(string(out))

	if err == nil {
		t.Error("Shouldn't be able to generate UUID due to incompatible version")
		t.FailNow()
	}
}

func TestUUIDInValidSpec(t *testing.T) {
	spec := `{}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(UUID, cfg, jsonIn)
	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err := uuid.FromString(string(out))

	if err == nil {
		t.Error("Shouldn't be able to generate UUID due to incompatible version")
		t.FailNow()
	}
}

func TestUUIDV3WithoutNamespace(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error("Should fail on missing namespace")
		t.FailNow()
	}

}

func TestUUIDV3InvalidNamespace(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, "namespace": "test", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error("Should fail on missing namespace")
		t.FailNow()
	}
}

func TestUUIDV3URLNameSpace(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, "namespace": "URL", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"06af7528-f22c-3716-86e0-579192ed244a"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDV3DNSNameSpace(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, "namespace": "DNS", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"83e9b77c-641d-331f-961d-bac57a61e534"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDV3OIDNameSpace(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, "namespace": "OID", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"d26cc082-0ba8-38a7-a738-983f3830f0cf"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDV3X500NameSpace(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, "namespace": "X500", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"438297ee-7562-336d-a913-7e7745455e80"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDWithCustomNameSpace(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, "namespace": "04536ac7-c030-4f10-811b-451bcc4c8ef5", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"82966ef3-a31d-379a-a266-d8f323901397"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDV3UsingDefaultField(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, "namespace": "URL", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"id":2323223}}`
	jsonOut := `{"a":{"id":2323223,"uuid":"a6cb7732-ccc1-3725-a9dc-040e73e0889b"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDBadNameSpace(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 3, "namespace": "not a uuid", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error(err)
		t.Error("Bad namespace should have prevented transform")
		t.FailNow()
	}

}

func TestUUIDV5(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 5, "namespace": "URL", 
			"names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]
		}
	}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"388607bf-b5c1-5f55-b10a-6afbbb91e18e"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDV5UsingDefaultField(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 5, "namespace": "URL", 
			"names": [
				{"path": "a.id", "default": "test"}, 
				{"path": "a.author", "default": "go lang rules!"}
			]
		}
	}`
	jsonIn := `{"a":{"id":2323223}}`
	jsonOut := `{"a":{"id":2323223,"uuid":"fb26c0d3-3cd2-514f-aced-2ecd901a1196"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDV5NoNamespace(t *testing.T) {
	spec := `{"a.uuid": {"version": 5, "names": []}}`
	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDV5NoNames(t *testing.T) {
	spec := `{"a.uuid": {"version": 5, "namespace": "URL"}}`
	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDV5ArrayIndex(t *testing.T) {
	spec := `{
		"a.uuid": {
			"version": 5, "namespace": "URL", 
			"names": [
				{"path": "a.tags[0]", "default": "test"}, 
				{"path": "a.author", "default": "go lang rules!"}
			]
		}
	}`
	jsonIn := `{"a":{"id":2323223, "tags": ["tag1", "tag2"]}}`
	jsonOut := `{
		"a": {
			"id": 2323223, "tags": ["tag1", "tag2"], "uuid": "5ae718bd-0add-5bab-a155-340433391b8c"
		}
	}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	areEqual, _ := checkJSONBytesEqual(kazaamOut, []byte(jsonOut))
	if !areEqual {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", string(kazaamOut))
		t.FailNow()
	}

	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err = uuid.FromString(string(out))

	if err != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.Error(err)
		t.FailNow()
	}
}

func TestUUIDWithMultipleSpecs(t *testing.T) {
	spec := `{"a.uuid": {"version": 4}, "fl.uuid": {"version": 4}}`
	jsonIn := `{"a": {"author": "jason", "id": 2323223}}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(UUID, cfg, jsonIn)
	aOut, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	flOut, _, _, _ := jsonparser.Get([]byte(kazaamOut), "fl", "uuid")
	_, aErr := uuid.FromString(string(aOut))
	_, flErr := uuid.FromString(string(flOut))

	if aErr != nil || flErr != nil {
		t.Error("transformed data didn't contain valid UUID")
		t.FailNow()
	}
}
