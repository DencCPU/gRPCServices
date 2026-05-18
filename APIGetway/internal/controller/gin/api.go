package gin

import (
	"context"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	spotservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/spot_service"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	orderdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"github.com/gin-gonic/gin"
)

type Service interface {
	RegistrationUser(ctx context.Context, newUser userdomain.User) (tokens.PairToken, error)
	Validation(ctx context.Context, accessToken string) (userservicedto.Output, error)
	UpdateTokens(ctx context.Context, accessToken, refreshToken string) (tokens.PairToken, error)
	CreateOrder(ctx context.Context, order orderdomain.OrderInfo) (orderdto.Output, error)
	ViewEnableMarkets(ctx context.Context, input spotservicedto.Input) ([]spotservicedto.Output, error)
	GetOrderStatus(ctx context.Context, input orderdto.GetInput) (orderdto.GetOutput, error)
	GetStreamStatus(ctx context.Context, input orderdto.GetInput, msgChan chan orderdto.StreamOutput) error
	Authentication(ctx context.Context, email, password string) (tokens.PairToken, error)
}

type GinAPI struct {
	r              *gin.Engine
	service        Service
	exeptionalPath map[string]bool
}

func NewGinAPI(service Service) GinAPI {
	api := GinAPI{}
	api.r = gin.Default()
	api.service = service

	api.exeptionalPath = map[string]bool{
		"/user/reg":  true,
		"/user/auth": true,
	}
	api.r.Use(otelgin.Middleware("API/getway"))
	api.r.Use(api.Middleware())

	api.endpoints()
	return api
}

func (api *GinAPI) endpoints() {

	//OrderSevice
	api.r.POST("/order", api.CreateOrderHandler)             //Create order
	api.r.GET("/order/status", api.GetOrderStatus)           //Get status order
	api.r.GET("/order/realtime_status", api.GetStreamStatus) //Get stream status order

	//SpotService
	api.r.POST("/markets", api.ViewEnableMarkets) //Get list available markets

	//UserService
	api.r.POST("/user/reg", api.RegistrationUser) //Registration a new user
	api.r.POST("/user/auth", api.Authentication)  //Authentication user
}

func (api *GinAPI) Router() *gin.Engine {
	return api.r
}
