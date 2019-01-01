package converter

import (
	"encoding/json"
	"errors"
	"github.com/qntfy/kazaam/transform"
	"regexp"
)

type regexSpec struct {
	Match   *string `json:"match"`
	Replace *string `json:"replace"`
}

type regexSpecsStruct map[string]regexSpec

var regexSpecs regexSpecsStruct

type Regex struct {
	ConverterBase
}

func (r *Regex) Init(config []byte) (err error) {
	if err := json.Unmarshal(config, &regexSpecs); err != nil {
		regexSpecs = make(regexSpecsStruct)
	}
	return
}

func (r *Regex) Convert(jsonData []byte, value []byte, args []byte) (newValue []byte, err error) {

	var jsonValue, argsValue *transform.JSONValue
	jsonValue, err = transform.NewJSONValue(value)
	if err != nil {
		return
	}

	argsValue, err = transform.NewJSONValue(args)
	if err != nil {
		return
	}

	if jsonValue.IsString() == false || argsValue.IsString() == false {
		err = errors.New("invalid value or arguments for regex converter")
		return
	}

	reName := argsValue.GetStringValue()
	if spec, ok := regexSpecs[reName]; ok {
		var re *regexp.Regexp
		re, err = regexp.Compile(*spec.Match)
		if err != nil {
			return
		}

		newValue = re.ReplaceAll(jsonValue.GetData(), []byte(*spec.Replace))

	} else {
		err = errors.New("regex not defined")
		return
	}

	return
}
