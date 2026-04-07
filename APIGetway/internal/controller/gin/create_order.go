package gin

import (
	"net/http"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/domain/order"
	"github.com/gin-gonic/gin"
)

func (api *GinAPI) CreateOrderHandler(c *gin.Context) {
	var order order.OrderInfo
	err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	uid, exist := c.Get("x-user-id")
	if !exist {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "user_id missing",
		})
		return
	}
	user_id, ok := uid.(string)
	if !ok {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "error casting user_id to string type",
		})
		return
	}
	order.User_id = user_id

	output, err := api.service.CreateOrder(c.Request.Context(), order)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"order_id":     output.Order_id,
		"order_status": output.Order_status,
	})
}
