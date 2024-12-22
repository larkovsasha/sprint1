package main

import (
	"github.com/larkovsasha/sprint1/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
