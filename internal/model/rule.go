package model

import (
	"io"

	mockserrors "github.com/nicopozo/mockserver/internal/errors"
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
	Key         string     `json:"key" example:"payments_get_556032950"`
	Application string     `json:"application" example:"payments"`
	Name        string     `json:"name" example:"get payment"`
	Path        string     `json:"path" example:"/v1/payments/{payment_id}"`
	Strategy    string     `json:"strategy" example:"normal"`
	Method      string     `json:"method" example:"GET"`
	Status      string     `json:"status" example:"enabled"`
	Responses   []Response `json:"responses"`
}

type RuleList struct {
	Paging  Paging
	Results []*Rule
}

type RuleStatus struct {
	Status string `json:"status" example:"enabled"`
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

func UnmarshalRuleStatus(body io.Reader) (*RuleStatus, error) {
	status := &RuleStatus{}
	err := jsonutils.Unmarshal(body, status)

	return status, err
}

func (status *RuleStatus) Validate() error {
	if status.Status != RuleStatusEnabled && status.Status != RuleStatusDisabled {
		return mockserrors.InvalidRulesError{
			Message: "invalid status: %s, only 'enabled' and 'disabled' are allowed",
		}
	}

	return nil
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
