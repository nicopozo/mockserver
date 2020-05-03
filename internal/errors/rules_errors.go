package mockserrors

type RuleNotFoundError struct {
	Message string
}

func (e RuleNotFoundError) Error() string {
	return e.Message
}

type InvalidRulesErrorError struct {
	Message string
}

func (e InvalidRulesErrorError) Error() string {
	return e.Message
}
