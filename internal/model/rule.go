package model

import (
	"io"

	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

const (
	RuleStatusEnabled  = "enabled"
	RuleStatusDisabled = "disabled"

	RuleStrategyNormal     = "normal"
	RuleStrategySequential = "sequential"
	RuleStrategyRandom     = "random"
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

type ESSearchResult struct {
	Hits *ESHits `json:"hits"`
}

type ESHits struct {
	Total *ESTotal  `json:"total"`
	Hits  []*ESRule `json:"hits"`
}

type ESTotal struct {
	Value int `json:"value"`
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

func UnmarshalSearchESRule(body io.Reader) (*ESSearchResult, error) {
	result := &ESSearchResult{}
	err := jsonutils.Unmarshal(body, result)

	return result, err
}
