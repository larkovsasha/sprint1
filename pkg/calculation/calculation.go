package calculation

import (
	"strconv"
	"strings"
)

var internalError error = nil

func checkBrackets(expression string) error {
	var count = 0
	for _, char := range expression {
		if string(char) == "(" {
			count++
		}
		if string(char) == ")" {
			count--
		}
		if count < 0 {
			return ErrWrongBracketsSequence
		}
	}
	if count == 0 {
		return nil
	}
	return ErrWrongBracketsSequence
}

func checkSigns(expression string) error {
	for i := 0; i < len(expression)-1; i++ {
		if isSign(expression[i]) && isSign(expression[i+1]) {
			return ErrTwoSignsInRow
		}
	}
	if isSign(expression[len(expression)-1]) {
		return ErrExpEndsWithSign
	}
	return nil
}

func isSign(ch uint8) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/'
}
func isBracket(ch uint8) bool {
	return ch == '(' || ch == ')'
}
func checkChars(expression string) error {
	for i := 0; i < len(expression); i++ {
		if !(isSign(expression[i]) || isBracket(expression[i]) || (expression[i] <= '9' && expression[i] >= '0')) {
			return ErrInvalidChars
		}
	}
	return nil
}

type Stack[T float64 | uint8] struct {
	Values []T
}

func (s *Stack[T]) Push(v T) {
	s.Values = append(s.Values, v)
}
func (s *Stack[T]) Top() T {
	if s.IsEmpty() {
		internalError = ErrInternalServer
		return 0
	}
	return s.Values[len(s.Values)-1]
}

func (s *Stack[T]) Pop() {
	if s.IsEmpty() {
		return
	}
	s.Values = s.Values[:len(s.Values)-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.Values) == 0
}

func doOperation(signs *Stack[uint8], numbers *Stack[float64]) error {
	s := signs.Top()
	signs.Pop()

	v1 := numbers.Top()
	numbers.Pop()
	v2 := numbers.Top()
	numbers.Pop()
	switch s {
	case '+':
		numbers.Push(v1 + v2)
	case '-':
		numbers.Push(v2 - v1)
	case '*':
		numbers.Push(v2 * v1)
	case '/':
		if v1 == 0 {
			return ErrDivisionByZero
		}
		numbers.Push(v2 / v1)
	}
	return nil
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, "()", "")
	if expression == "" {
		return 0, ErrExpIsEmpty
	}

	err := checkBrackets(expression)
	if err != nil {
		return 0, err
	}
	err = checkChars(expression)
	if err != nil {
		return 0, err
	}
	err = checkSigns(expression)
	if err != nil {
		return 0, err
	}

	m := map[uint8]int{'+': 0, '-': 0, '/': 1, '*': 1, '(': 2, ')': 2}

	var signs Stack[uint8]
	var numbers Stack[float64]
	for i := 0; i < len(expression); i++ {
		char := expression[i]
		if isSign(char) || isBracket(char) {
			if char == ')' {
				for signs.Top() != '(' {
					err := doOperation(&signs, &numbers)
					if err != nil {
						return 0, err
					}
				}
				signs.Pop()
			} else if signs.IsEmpty() || m[signs.Top()] < m[char] || signs.Top() == '(' || char == '(' {
				signs.Push(char)
			} else {
				for !signs.IsEmpty() && m[signs.Top()] >= m[char] {
					err := doOperation(&signs, &numbers)
					if err != nil {
						return 0, err
					}
				}
				signs.Push(char)
			}
		} else {
			j := i
			for ; j < len(expression) && !isSign(expression[j]) && !isBracket(expression[j]); j++ {
			}
			n, e := strconv.ParseFloat(expression[i:j], 64)
			if e != nil {
				return 0, e
			}
			numbers.Push(n)
			i = j - 1
		}
	}
	for !signs.IsEmpty() {
		err := doOperation(&signs, &numbers)
		if err != nil {
			return 0, err
		}
	}
	if !signs.IsEmpty() || len(numbers.Values) > 1 {
		return 0, ErrInternalServer
	}
	if internalError != nil {
		return 0, internalError
	}
	return numbers.Top(), nil
}
