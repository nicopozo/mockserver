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

type ESPatternList struct {
	Source *PatternList `json:"_source"`
}

func UnmarshalESPatternList(body io.Reader) (*PatternList, error) {
	patternList := &ESPatternList{}
	err := jsonutils.Unmarshal(body, patternList)

	return patternList.Source, err
}
