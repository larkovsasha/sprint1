package calculation_test

import (
	"errors"
	"github.com/larkovsasha/sprint1/pkg/calculation"
	"testing"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
		{
			name:           "complicated",
			expression:     "2-(1-6)+(8+2)*3",
			expectedResult: 37,
		},
		{
			name:           "complicated",
			expression:     "25*(0-5)+(8+2)*15",
			expectedResult: 25,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculation.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "ends with sign",
			expectedErr: calculation.ErrExpEndsWithSign,
			expression:  "1+1*",
		},

		{
			name:        "empty string",
			expression:  "",
			expectedErr: calculation.ErrExpIsEmpty,
		},

		{
			name:        "two signs in row",
			expression:  "1+*1",
			expectedErr: calculation.ErrTwoSignsInRow,
		},
		{
			name:        "wrong brackets",
			expression:  "1+(1+5))",
			expectedErr: calculation.ErrWrongBracketsSequence,
		},
		{
			name:        "wrong brackets",
			expression:  "1+((1+5)",
			expectedErr: calculation.ErrWrongBracketsSequence,
		},
		{
			name:        "invalid expression",
			expression:  "25(0-6)+(8+2)*15",
			expectedErr: calculation.ErrInternalServer,
		},
		{
			name:        "division by zero",
			expression:  "25/(6-6)+(8+2)*15",
			expectedErr: calculation.ErrDivisionByZero,
		},
		{
			name:        "internal server",
			expression:  "25/(6-6)+(d+2)*15",
			expectedErr: calculation.ErrInvalidChars,
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := calculation.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("test must be failed")
			}
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("error should be %s but is %s", testCase.expectedErr, err)
			}
		})
	}
}
