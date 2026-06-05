package gin

import (
	"net/http"

	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"github.com/gin-gonic/gin"
)

func (api *GinAPI) RegistrationUser(c *gin.Context) {
	var newUser userdomain.User
	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"serialization error": err.Error(),
		})
		return
	}

	pairToken, err := api.service.RegistrationUser(c.Request.Context(), newUser)
	if err != nil {

		c.JSON(http.StatusBadGateway, gin.H{
			"error:": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access token":  pairToken.AccessToken,
		"refresh token": pairToken.RefreshToken,
		"expire_at":     pairToken.ExpireAt,
	})
}
