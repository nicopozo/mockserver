package model

import (
	"fmt"
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
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, %w", err)
	}

	return patternList.Source, nil
}
