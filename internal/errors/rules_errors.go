package ruleserrors

type RuleNotFoundError struct {
	Message string
}

func (e RuleNotFoundError) Error() string {
	return e.Message
}
