package main

import (
	"Academy/gRPCServices/APIGetway/internal/controller/gin"
	"fmt"
	"log"
	"net/http"
)

func main() {
	api := gin.NewGinAPI()
	fmt.Println("Сервер запущен")
	err := http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
