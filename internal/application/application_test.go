package application_test

import (
	"bytes"
	"encoding/json"
	"github.com/larkovsasha/sprint1/internal/application"
	"github.com/larkovsasha/sprint1/pkg/calculation"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testExpression struct {
	Expression string `json:"expression"`
}

func TestApplication_RunServer(t *testing.T) {
	tests := []struct {
		expression string
		statusCode int
		result     float64
	}{
		{"2+2*2", http.StatusOK, 6},
		{"2*(2+2)", http.StatusOK, 8},
		{"2-(1-6)+(8+2)*3", http.StatusOK, 37},
		{"7*(1/(2+8))", http.StatusOK, 0.7},
	}

	for _, tc := range tests {
		t.Run(tc.expression, func(t *testing.T) {
			data := testExpression{Expression: tc.expression}
			jsonData, _ := json.Marshal(data)
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(application.CalcHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.statusCode)
			}

			body, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			var responseBody map[string]interface{}
			err = json.Unmarshal(body, &responseBody)
			if err != nil {
				t.Fatal(err)
			}

			eps := 0.0000001
			if responseBody["result"].(float64)-tc.result > eps {
				t.Errorf("handler returned wrong result: got %v want %v", responseBody["result"], tc.result)
			}
		})
	}

	uncorrect_tests := []struct {
		expression  string
		statusCode  int
		expectedErr string
	}{
		{
			statusCode:  http.StatusUnprocessableEntity,
			expectedErr: calculation.ErrInvalidExpression.Error(),
			expression:  "1+1*",
		},

		{
			statusCode:  http.StatusUnprocessableEntity,
			expression:  "",
			expectedErr: calculation.ErrInvalidExpression.Error(),
		},

		{
			statusCode:  http.StatusUnprocessableEntity,
			expression:  "1+*1",
			expectedErr: calculation.ErrInvalidExpression.Error(),
		},
		{
			statusCode:  http.StatusUnprocessableEntity,
			expression:  "1+(1+5))",
			expectedErr: calculation.ErrInvalidExpression.Error(),
		},
		{
			statusCode:  http.StatusUnprocessableEntity,
			expression:  "1+((1+5)",
			expectedErr: calculation.ErrInvalidExpression.Error(),
		},
		{
			statusCode:  http.StatusUnprocessableEntity,
			expression:  "25/(6-6)+(8+2)*15",
			expectedErr: calculation.ErrInvalidExpression.Error(),
		},
		{
			statusCode:  http.StatusUnprocessableEntity,
			expression:  "25/(6-6)+(d+2)*15",
			expectedErr: calculation.ErrInvalidExpression.Error(),
		},
	}

	for _, tc := range uncorrect_tests {
		t.Run(tc.expression, func(t *testing.T) {
			data := testExpression{Expression: tc.expression}
			jsonData, _ := json.Marshal(data)
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := application.CalcHandler
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.statusCode)
			}

			body, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			var responseBody map[string]interface{}
			err = json.Unmarshal(body, &responseBody)
			if err != nil {
				t.Fatal(err)
			}

			if responseBody["error"].(string) != tc.expectedErr {
				t.Errorf("handler returned wrong result: got %v want %v", responseBody["result"], tc.expectedErr)
			}
		})
	}
}
