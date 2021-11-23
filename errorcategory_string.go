// Code generated by "stringer -type=ErrorCategory"; DO NOT EDIT.

package poc

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SemanticError-1]
	_ = x[MalformedError-2]
	_ = x[UnavailableError-3]
	_ = x[UnimplementedError-4]
	_ = x[NotFoundError-5]
	_ = x[UnknownError-6]
}

const _ErrorCategory_name = "SemanticErrorMalformedErrorUnavailableErrorUnimplementedErrorNotFoundErrorUnknownError"

var _ErrorCategory_index = [...]uint8{0, 13, 27, 43, 61, 74, 86}

func (i ErrorCategory) String() string {
	i -= 1
	if i >= ErrorCategory(len(_ErrorCategory_index)-1) {
		return "ErrorCategory(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ErrorCategory_name[_ErrorCategory_index[i]:_ErrorCategory_index[i+1]]
}