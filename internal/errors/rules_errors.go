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

type RuleAlreadyCreatedError struct {
	Message string
}

func (e RuleAlreadyCreatedError) Error() string {
	return e.Message
}
