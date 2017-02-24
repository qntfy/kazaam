package transform

import (
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
)

// Extract returns the specified path as the top-level object.
func Extract(spec *Config, data *simplejson.Json) error {
	outPath, ok := (*spec.Spec)["path"]
	if !ok {
		return &Error{ErrMsg: fmt.Sprintf("Unable to get path"), ErrType: SpecError}
	}
	tmp, err := getJSONPath(data, outPath.(string), spec.Require)
	if err != nil {
		return err
	}
	*data = *tmp
	return nil
}
