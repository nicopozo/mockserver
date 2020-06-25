package mockserrors

type RuleNotFoundError struct {
	Message string
}

func (e RuleNotFoundError) Error() string {
	return e.Message
}

type InvalidRulesError struct {
	Message string
}

func (e InvalidRulesError) Error() string {
	return e.Message
}
