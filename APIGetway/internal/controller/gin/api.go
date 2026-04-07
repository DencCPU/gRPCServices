package gin

import (
	"context"
	"net/http"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	spotservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/spot_service"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Service interface {
	RegistrationUser(ctx context.Context, newUser userdomain.User) (tokens.PairToken, error)
	Validation(ctx context.Context, accessToken string) (userservicedto.Output, error)
	UpdateTokens(ctx context.Context, accessToken, refreshToken string) (tokens.PairToken, error)
	CreateOrder(ctx context.Context, order order.OrderInfo) (orderdto.Output, error)
	ViewEnableMarkets(ctx context.Context, role string) ([]spotservicedto.Output, error)
	GetOrderStatus(ctx context.Context, input orderdto.GetInput) (orderdto.GetOutput, error)
	GetStreamStatus(ctx context.Context, input orderdto.GetInput, msgChan chan orderdto.StreamOutput) error
}
type GinAPI struct {
	r              *gin.Engine
	service        Service
	exeptionalPath map[string]bool
	upgrader       websocket.Upgrader
}

func NewGinAPI(service Service) GinAPI {
	api := GinAPI{}
	api.r = gin.Default()
	api.service = service

	api.exeptionalPath = map[string]bool{
		"/user/reg": true,
	}

	api.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	api.endpoints()
	return api
}

func (api *GinAPI) endpoints() {
	//OrderSevice
	api.r.Use(api.Middleware())
	api.r.POST("/order", api.CreateOrderHandler)             //Create order
	api.r.POST("/order/status", api.GetOrderStatus)          //Get status order
	api.r.GET("/order/realtime_status", api.GetStreamStatus) //Get stream status order

	//SpotService
	api.r.POST("/markets", api.ViewEnableMarkets) //Get list available markets

	//UserService
	api.r.POST("/user/reg", api.RegistrationUser) //Registration a new user
}

func (api *GinAPI) Router() *gin.Engine {
	return api.r
}
