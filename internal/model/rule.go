package model

import (
	"io"

	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

type Rule struct {
	Key         string     `json:"key"`
	Application string     `json:"application"`
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Strategy    string     `json:"strategy"`
	Method      string     `json:"method"`
	Status      string     `json:"status"`
	Responses   []Response `json:"responses"`
}

type RuleList struct {
	Paging  Paging
	Results []*Rule
}

func UnmarshalRule(body io.Reader) (*Rule, error) {
	task := &Rule{}
	err := jsonutils.Unmarshal(body, task)

	return task, err
}
