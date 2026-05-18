package gin

import (
	"net/http"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"github.com/gin-gonic/gin"
)

func (api *GinAPI) GetOrderStatus(c *gin.Context) {
	var orderInfo orderdto.GetInput

	orderInfo.OrderId = c.Query("id")
	if orderInfo.OrderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unfinde orderID in query row",
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
	orderInfo.UserId = user_id
	output, err := api.service.GetOrderStatus(c.Request.Context(), orderInfo)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"order_id":     output.OrderId,
		"order_status": output.Status,
		"price":        output.Price,
		"quantity":     output.Quantity,
		"market_name":  output.MarketName,
	})
}
