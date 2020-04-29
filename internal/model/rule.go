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

type ESRule struct {
	Source *Rule `json:"_source"`
}

func UnmarshalRule(body io.Reader) (*Rule, error) {
	rule := &Rule{}
	err := jsonutils.Unmarshal(body, rule)

	return rule, err
}

func UnmarshalESRule(body io.Reader) (*Rule, error) {
	rule := &ESRule{}
	err := jsonutils.Unmarshal(body, rule)

	return rule.Source, err
}
