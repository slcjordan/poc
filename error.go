package poc

// An ErrorCategory must be at least detailed enough to only correspond to a
// single http status code, but may be broken down further (one-to-many).
type ErrorCategory uint8

// supported error categories
const (
	_ ErrorCategory = iota
	SemanticError
	MalformedError
	UnavailableError
	UnimplementedError
	NotFoundError
	UnknownError
)

// Error wraps an existing error and assigns it a category.
type Error struct {
	Actual   error
	Category ErrorCategory
}

func (e Error) Unwrap() error {
	return e.Actual
}

func (e Error) Error() string {
	return e.Actual.Error()
}
