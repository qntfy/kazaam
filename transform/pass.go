package transform

import simplejson "github.com/bitly/go-simplejson"

// Pass performs no manipulation of the passed-in data. It is useful
// for testing/default behavior.
func Pass(spec *Config, data *simplejson.Json) error {
	return nil
}
