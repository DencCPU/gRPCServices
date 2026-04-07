package apprunner

import (
	"fmt"
	"log"
	"net/http"

	orderclient "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/order_client"
	spotclient "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/spot_client"
	userclient "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/user_client"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/controller/gin"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/usecase"
	"github.com/DencCPU/gRPCServices/Shared/logger"
	"go.uber.org/zap"
)

func Apprunner() error {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	user_client, err := userclient.NewClient()
	if err != nil {
		logger.Error("user_client creation error",
			zap.Error(err),
		)
		return err
	}
	order_client, err := orderclient.NewClient()
	if err != nil {
		logger.Error("order_client creation error",
			zap.Error(err),
		)
		return err
	}

	spot_client, err := spotclient.NewClient()
	if err != nil {
		logger.Error("spot_client creation error",
			zap.Error(err),
		)
		return err
	}

	service := usecase.NewService(user_client, order_client, spot_client, logger)
	api := gin.NewGinAPI(service)
	fmt.Println("Server is running")
	err = http.ListenAndServe(":8082", api.Router())
	if err != nil {
		logger.Error("server error:",
			zap.Error(err),
		)
		return err
	}
	return nil
}
