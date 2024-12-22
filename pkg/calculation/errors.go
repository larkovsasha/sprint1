package calculation

import "errors"

var (
	ErrWrongBracketsSequence = errors.New("Wrong brackets sequence")
	ErrTwoSignsInRow         = errors.New("Two signs in a row")
	ErrExpEndsWithSign       = errors.New("Expressions ends with sign")
	ErrExpIsEmpty            = errors.New("Expressions is empty")
	ErrDivisionByZero        = errors.New("Division by zero")
	ErrInternalServer        = errors.New("internal server error")
)

var ErrInvalidExpression = errors.New("Expression is not valid")
var ErrInvalidChars = errors.New("Invalid chars")
