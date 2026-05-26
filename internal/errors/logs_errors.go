package mockserrors

import "fmt"

type LogEntryNotFoundError struct {
	ID string
}

func (e LogEntryNotFoundError) Error() string {
	return fmt.Sprintf("log entry not found: %s", e.ID)
}

func NewLogEntryNotFoundError(id string) error {
	return LogEntryNotFoundError{ID: id}
}
