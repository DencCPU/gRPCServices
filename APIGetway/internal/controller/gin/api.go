package gin

import (
	"github.com/gin-gonic/gin"
)

type GinAPI struct {
	r *gin.Engine
}

func NewGinAPI() GinAPI {
	api := GinAPI{}
	api.r = gin.Default()
	api.endpoints()
	return api
}

func (api *GinAPI) endpoints() {
	api.r.POST("/order", api.CreateOrderHandler)
	api.r.GET("/status/{id}")
	api.r.GET("/realtime_status/{id}")
}

func (api *GinAPI) Router() *gin.Engine {
	return api.r
}
