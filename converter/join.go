package converter

import (
	"encoding/json"
	"errors"
	"github.com/mbordner/kazaam/transform"
	"regexp"
	"strconv"
	"strings"
)

type Join struct {
	ConverterBase
}

// |substr start end , end is optional, and will be the last char sliced's index + 1,
// start is the start index and required
func (c *Join) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
	newValue = value

	var stringValues []string

	err = json.Unmarshal(value,&stringValues)
	if err != nil {
		return
	}

	argsValue, err := transform.NewJSONValue(args)
	if err != nil {
		return
	}

	if argsValue.IsString() == false {
		err = errors.New("invalid value or arguments for join converter")
		return
	}

	var re *regexp.Regexp
	re, err = regexp.Compile(`(?Us)^(?:\s*)(.+)(?:\s*)$`)
	if err != nil {
		return
	}

	argsString := argsValue.GetStringValue()

	newValue = []byte("null")

	if matches := re.FindStringSubmatch(argsString); matches != nil {

		val := strings.Join(stringValues,matches[1])

		newValue = []byte(strconv.Quote(val))

	}


	return
}
