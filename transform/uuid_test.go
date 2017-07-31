package transform

import (
	"testing"

	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
)

func TestUUIDV4(t *testing.T) {
	spec := `{"a.uuid":{"version": 4}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}`

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
	jsonIn := `{"a":{"author":"jason","id":2323223}`

	cfg := getConfig(spec, false)
	kazaamOut, _ := getTransformTestWrapper(UUID, cfg, jsonIn)
	out, _, _, _ := jsonparser.Get([]byte(kazaamOut), "a", "uuid")
	_, err := uuid.FromString(string(out))

	if err == nil {
		t.Error("Shouldn't be able to generate UUID due to incompatible version")
		t.FailNow()
	}
}

func TestUUIDNoVersionErro(t *testing.T) {
	spec := `{"a.uuid":{}`
	jsonIn := `{"a":{"author":"jason","id":2323223}`

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
	jsonIn := `{"a":{"author":"jason","id":2323223}`

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

	spec := `{"a.uuid": {"version": 3, "names": [{"path": "a.id", "default": "test"}, {"path": "a.author", "default": "test"}]}}`

	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error("Should fail on missing namespace")
		t.FailNow()
	}

}

func TestUUIDV3InvalidNamespace(t *testing.T) {

	spec := `{"a.uuid": {"version": 3, "nameSpace": "test", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`

	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error("Should fail on missing namespace")
		t.FailNow()
	}

}

func TestUUIDV3URLNameSpace(t *testing.T) {

	spec := `{"a.uuid": {"version": 3, "nameSpace": "URL", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"cad3ae9e-7a89-3b0b-8ef4-9c22ead9c6eb"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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

	spec := `{"a.uuid": {"version": 3, "nameSpace": "DNS", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"9a77a459-b1a3-32bc-b758-e5569a667a61"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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

	spec := `{"a.uuid": {"version": 3, "nameSpace": "OID", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"c01bef62-619d-3524-a36c-4bdcf263e0cb"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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

	spec := `{"a.uuid": {"version": 3, "nameSpace": "X500", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"b7b0fc51-085c-35a3-9b1b-e1b5dcef128b"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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

func TestUUIDWithCustomeNameSpace(t *testing.T) {

	spec := `{"a.uuid": {"version": 3, "nameSpace": "04536ac7-c030-4f10-811b-451bcc4c8ef5", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"49121a9c-2d58-30aa-8eed-02eb1b61a0e1"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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

	spec := `{"a.uuid":{"version": 3, "nameSpace": "URL", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`
	jsonIn := `{"a":{"id":2323223}}`
	jsonOut := `{"a":{"id":2323223,"uuid":"9a4d8062-cefd-35d5-907e-b2da04873d95"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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

	spec := `{"a.uuid": {"version": 3, "nameSpace": "not a uuid", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`
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

	spec := `{"a.uuid": {"version": 5, "nameSpace": "URL", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "test"}]}}`
	jsonIn := `{"a":{"author":"jason","id":2323223}}`
	jsonOut := `{"a":{"author":"jason","id":2323223,"uuid":"7e7c9ede-828f-39f7-9cc0-55c4a259d8f4"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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

	spec := `{"a.uuid": {"version": 5, "nameSpace": "URL", "names": [{"path": "a.id", "default": "test"},
	{"path": "a.author", "default": "go lang rules!"}]}}`
	jsonIn := `{"a":{"id":2323223}}`
	jsonOut := `{"a":{"id":2323223,"uuid":"015747dc-eef7-36ab-b22f-1c851ef3118e"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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

func TestUUIDV5NoNames(t *testing.T) {

	spec := `{"a.uuid": {"version": 5, "nameSpace": "URL", "names": []},
	{"path": "a.author", "default": "go lang rules!"}]}}`
	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error(err)
		t.FailNow()
	}

}

func TestUUIDV5NoName(t *testing.T) {

	spec := `{"a.uuid": {"version": 5, "nameSpace": "URL"},
	{"path": "a.author", "default": "go lang rules!"}]}}`
	jsonIn := `{"a":{"id":2323223}}`

	cfg := getConfig(spec, false)
	_, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err == nil {
		t.Error(err)
		t.FailNow()
	}

}

func TestUUIDV5ArrayIndex(t *testing.T) {
	spec := `{"a.uuid": {"version": 5, "nameSpace": "URL", "names": [{"path": "a.tags[0]", "default": "test"},
	{"path": "a.author", "default": "go lang rules!"}]}}`
	jsonIn := `{"a":{"id":2323223, "tags": ["tag1", "tag2"]}}`
	jsonOut := `{"a":{"id":2323223, "tags": ["tag1", "tag2"],"uuid":"49ad2943-37a6-3baa-96ca-980861b80191"}}`

	cfg := getConfig(spec, false)
	kazaamOut, err := getTransformTestWrapper(UUID, cfg, jsonIn)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if kazaamOut != jsonOut {
		t.Error("Transformed data does not match expectation.")
		t.Log("Expected:   ", jsonOut)
		t.Log("Actual:     ", kazaamOut)
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
