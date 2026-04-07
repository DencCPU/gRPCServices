package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *GinAPI) ViewEnableMarkets(c *gin.Context) {
	r, exist := c.Get("x-user-role")
	if !exist {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "user role missing",
		})
		c.Abort()
		return
	}
	role, ok := r.(string)
	if !ok {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "error casting user_role to string type",
		})
		return
	}
	fmt.Println(role)
	markets, err := api.service.ViewEnableMarkets(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "error getting available markets",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"available markets": markets,
	})
}
