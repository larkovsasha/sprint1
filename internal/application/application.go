package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/larkovsasha/sprint1/pkg/calculation"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr string
}

func GetConfig() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: GetConfig(),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type SuccessResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func calcHandlerFunction(w http.ResponseWriter, r *http.Request) {
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		log.Println(err)
		if errors.Is(err, calculation.ErrInternalServer) {
			writeErrorResponse(w, http.StatusInternalServerError, calculation.ErrInternalServer.Error())
		} else {
			writeErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		}
		return
	}
	writeSuccessResponse(w, SuccessResponse{result})
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeSuccessResponse(w http.ResponseWriter, response SuccessResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("%s", err)
			}
		}()
		log.Printf("%s", r.Method)
		next.ServeHTTP(w, r)
	}
}

func ParseRequestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Wrong method")
			writeErrorResponse(w, http.StatusInternalServerError, calculation.ErrInternalServer.Error())
			return
		}
		next.ServeHTTP(w, r)
	}
}

var CalcHandler = ParseRequestMiddleware(LoggingMiddleware(calcHandlerFunction))

func (a *Application) RunServer() error {
	fmt.Printf("Starting server at %s\n", a.config.Addr)
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
