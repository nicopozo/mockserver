package jsonutils

import (
	"encoding/json"
	"fmt"
	"io"
)

func Unmarshal(body io.Reader, model any) error {
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

func Marshal(model any) string {
	if model == nil {
		return ""
	}

	if s, ok := model.(string); ok {
		return s
	}

	out, err := json.Marshal(model)
	if err == nil {
		return string(out)
	}

	return ""
}
