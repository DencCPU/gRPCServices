package main

import (
	"log"

	"github.com/DencCPU/gRPCServices/APIGetway/pkg/apprunner"
)

func main() {
	err := apprunner.Apprunner()
	if err != nil {
		log.Fatal(err)
	}
}
