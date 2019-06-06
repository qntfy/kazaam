package converter

import (
	"errors"
	"github.com/mbordner/kazaam/transform"
	"regexp"
	"strconv"
	"strings"
)

type Splitn struct {
	ConverterBase
}

// |substr start end , end is optional, and will be the last char sliced's index + 1,
// start is the start index and required
func (c *Splitn) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {
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
		err = errors.New("invalid value or arguments for splintn converter")
		return
	}

	var re *regexp.Regexp
	re, err = regexp.Compile(`(?Us)^(?:\s*)(.+)(?:\s*)(\d+)*(?:\s*)$`)
	if err != nil {
		return
	}

	argsString := argsValue.GetStringValue()
	origValue := jsonValue.GetStringValue()

	var n int64

	newValue = []byte("null")

	if matches := re.FindStringSubmatch(argsString); matches != nil {
		n, err = strconv.ParseInt(matches[2], 10, 64)
		if err != nil {
			return
		}

		vals := strings.SplitN(origValue,matches[1],int(n)+1)

		if len(vals) >= int(n) {
			newValue = []byte(strconv.Quote(vals[int(n)-1]))
		}

	}


	return
}
