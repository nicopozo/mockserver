package jsonutils

import (
	"encoding/json"
	"fmt"
	"io"
)

func Unmarshal(body io.Reader, model interface{}) error {
	var jsonReq []byte

	var err error
	if jsonReq, err = io.ReadAll(body); err != nil {
		return fmt.Errorf("error reading body, %w", err)
	}

	err = json.Unmarshal(jsonReq, &model)
	if err != nil {
		return fmt.Errorf("error unmarshalling reader %w", err)
	}

	return nil
}

func Marshal(model interface{}) string {
	result := ""

	if model == nil {
		return result
	}

	out, err := json.Marshal(model)
	if err == nil {
		result = string(out)
	}

	return result
}
