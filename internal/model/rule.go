package model

import (
	"fmt"
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
	RuleStrategyScene      = "scene"
)

type Rule struct {
	Key             string            `json:"key" example:"payments_get_556032950"`
	Group           string            `json:"group" example:"payments"`
	Name            string            `json:"name" example:"get payment"`
	Path            string            `json:"path" example:"/v1/payments/{payment_id}"`
	Strategy        string            `json:"strategy" example:"normal"`
	Method          string            `json:"method" example:"GET"`
	Status          string            `json:"status" example:"enabled"`
	Responses       []Response        `json:"responses"`
	Variables       []*Variable       `json:"variables"`
	AssertionGroups []*AssertionGroup `json:"assertion_groups"`
}

type RuleList struct {
	Paging  Paging  `json:"paging"`
	Results []*Rule `json:"results"`
}

type RuleStatus struct {
	Status string `json:"status" example:"enabled"`
}

type ESRule struct {
	Source *Rule `json:"_source"` //nolint:tagliatelle
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
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, %w", err)
	}

	return rule, nil
}

func UnmarshalRules(body io.Reader) ([]Rule, error) {
	rules := new([]Rule)

	err := jsonutils.Unmarshal(body, rules)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, %w", err)
	}

	return *rules, nil
}

func UnmarshalRuleStatus(body io.Reader) (*RuleStatus, error) {
	status := &RuleStatus{}

	err := jsonutils.Unmarshal(body, status)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, %w", err)
	}

	return status, nil
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
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, %w", err)
	}

	return rule.Source, nil
}

func UnmarshalSearchESRule(body io.Reader) (*ESSearchResult, error) {
	result := &ESSearchResult{}

	err := jsonutils.Unmarshal(body, result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, %w", err)
	}

	return result, nil
}
