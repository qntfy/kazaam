package converter

import (
	"encoding/json"
	"errors"
	"github.com/mbordner/kazaam/transform"
	"regexp"
	"strings"
)

type Split struct {
	ConverterBase
}

// |substr start end , end is optional, and will be the last char sliced's index + 1,
// start is the start index and required
func (c *Split) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	newValue = value

	var jsonValue, argsValue *transform.JSONValue
	jsonValue, err = c.NewJSONValue(value)
	if err != nil {
		return
	}

	argsValue, err = transform.NewJSONValue(args)
	if err != nil {
		return
	}

	if jsonValue.IsString() == false || argsValue.IsString() == false {
		err = errors.New("invalid value or arguments for split converter")
		return
	}

	var re *regexp.Regexp
	re, err = regexp.Compile(`(?Us)^(?:\s*)(.+)(?:\s*)$`)
	if err != nil {
		return
	}

	argsString := argsValue.GetStringValue()
	origValue := jsonValue.GetStringValue()

	newValue = []byte("null")

	if matches := re.FindStringSubmatch(argsString); matches != nil {

		vals := strings.Split(origValue,matches[1])

		newValue, err = json.Marshal(vals)

	}


	return
}
