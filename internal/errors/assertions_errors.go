package mockserrors

import jsonutils "github.com/nicopozo/mockserver/internal/utils/json"

type AssertionError struct {
	Errors []string
}

func (e AssertionError) Error() string {
	return jsonutils.Marshal(e.Errors)
}
