package main

import (
	"log"

	"github.com/DencCPU/gRPCServices/UserService/pkg/apprunner"
)

func main() {
	app, err := apprunner.FxAppRunner()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
