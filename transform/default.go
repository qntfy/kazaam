package transform

import (
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

// Default sets specific value(s) in output json.
func Default(spec *Config, data *simplejson.Json) (*simplejson.Json, error) {
	for k, v := range *spec.Spec {
		data.SetPath(strings.Split(k, "."), v)
	}
	return data, nil
}
