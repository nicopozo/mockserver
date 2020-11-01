package jsonutils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func Unmarshal(body io.Reader, model interface{}) error {
	var jsonReq []byte

	var err error
	if jsonReq, err = ioutil.ReadAll(body); err != nil {
		return fmt.Errorf("error reading body, %w", err)
	}

	return json.Unmarshal(jsonReq, &model)
}

func Marshal(model interface{}) string {
	result := ""

	out, err := json.Marshal(model)
	if err == nil {
		result = string(out)
	}

	return result
}
