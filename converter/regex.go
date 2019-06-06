package converter

import (
	"encoding/json"
	"errors"
	"github.com/mbordner/kazaam/transform"
	"regexp"
	"strconv"
)

type regexSpecStruct struct {
	Match   *string `json:"match"`
	Replace *string `json:"replace"`
}

type regexSpecsList []regexSpecStruct

func (rsl *regexSpecsList) UnmarshalJSON(config []byte) (err error) {
	if config[0] == '{' {
		spec := regexSpecStruct{}
		err = json.Unmarshal(config, &spec)
		if err != nil {
			return err
		}
		*rsl = append(*rsl, spec)
		return nil
	}

	var list []regexSpecStruct
	json.Unmarshal(config, &list)

	*rsl = regexSpecsList(list)

	return
}

type regexSpecs map[string]regexSpecsList

var specs regexSpecs

type Regex struct {
	ConverterBase
}

func (r *Regex) Init(config []byte) (err error) {
	specs = make(regexSpecs)
	err = json.Unmarshal(config, &specs)
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
	if specs, ok := specs[reName]; ok {
		for _, spec := range specs {
			var re *regexp.Regexp
			re, err = regexp.Compile(*spec.Match)
			if err != nil {
				return
			}

			src := jsonValue.GetStringValue()

			if re.Match([]byte(src)) {
				newValue = re.ReplaceAll([]byte(src), []byte(*spec.Replace))

				newValue = []byte(strconv.Quote(string(newValue)))
				break
			} else {
				newValue = []byte(strconv.Quote(string(src)))
			}
		}
	} else {
		err = errors.New("regex not defined")
		return
	}

	return
}
