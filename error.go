package poc

type ErrorCategory uint8

const (
	_ ErrorCategory = iota
	SemanticError
	MalformedError
	UnavailableError
	UnimplementedError
	NotFoundError
	UnknownError
)

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
