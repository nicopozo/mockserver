package model

import (
	"io"

	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

type Pattern struct {
	PathExpression string   `json:"path_expression"`
	RuleKeys       []string `json:"rule_keys"`
}

type PatternList struct {
	Patterns []*Pattern `json:"patterns"`
}

func UnmarshalPatternList(body io.Reader) (*PatternList, error) {
	patternList := &PatternList{}
	err := jsonutils.Unmarshal(body, patternList)

	return patternList, err
}
