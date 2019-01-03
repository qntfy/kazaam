package converter

import (
	"errors"
	"github.com/qntfy/kazaam/transform"
	"regexp"
	"strconv"
)

type Substr struct {
	ConverterBase
}

// |substr start end , end is optional, and will be the last char sliced's index + 1,
// start is the start index and required
func (c *Substr) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
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
		err = errors.New("invalid value or arguments for substr converter")
		return
	}

	var re *regexp.Regexp
	re, err = regexp.Compile(`(?:\s*)(\d+)(?:\s*)(\d+)*(?:\s*)`)
	if err != nil {
		return
	}

	argsString := argsValue.GetStringValue()
	origValue := jsonValue.GetStringValue()

	var start, end int64
	end = int64(len(origValue))

	if matches := re.FindStringSubmatch(argsString); matches != nil {
		start, err = strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return
		}
		if len(matches) > 2 {
			end, err = strconv.ParseInt(matches[2], 10, 64)
			if err != nil {
				return
			}
		}
	}

	newValue = []byte(strconv.Quote(origValue[start:end]))

	return
}
