package gin

import (
	"Academy/gRPCServices/APIGetway/internal/domain/order"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *GinAPI) CreateOrderHandler(c *gin.Context) {
	var order order.OrderInfo
	err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	user_id := c.GetInt64("user_id")
	order.User_id = user_id
	//Достать userID из контектса(middleware)
}
