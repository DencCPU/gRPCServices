package gin

import (
	"fmt"
	"net/http"
	"strings"

	sharederrors "github.com/DencCPU/gRPCServices/Shared/errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (api *GinAPI) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//exceptional path
		if api.exeptionalPath[c.Request.URL.Path] {
			c.Next()
			return
		}

		//Cheack headers
		//Get access token from header
		accessToken := c.GetHeader("access-token")
		if accessToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the access token header is missing ",
			})
			c.Abort()
			return
		}

		//Get refresh token from header
		refreshToken := c.GetHeader("refresh-token")
		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the refresh token header is missing ",
			})
			c.Abort()
			return
		}

		user, err := api.service.Validation(c.Request.Context(), accessToken)
		fmt.Println("ID:", user.User_id, "Role:", user.Role)

		if err != nil {
			if status.Code(err) == codes.Unauthenticated && strings.Contains(status.Convert(err).Message(), sharederrors.EXPIRED_TOKEN.Error()) {
				fmt.Println("Код здесь")
				pairToken, err := api.service.UpdateTokens(c.Request.Context(), accessToken, refreshToken)
				if err != nil {
					c.JSON(http.StatusBadGateway, gin.H{
						"error": err.Error(),
					})
					c.Abort()
					return
				}

				fmt.Println("Установак новых заголовков")
				layout := "2006-01-02 15:04:05"
				c.Header("new-access-token", pairToken.AccessToken)
				c.Header("new-refresh-token", pairToken.RefreshToken)
				c.Header("new-expire_at", pairToken.Expire_at.Format(layout))

				fmt.Println("Новый токен доступа:", c.GetHeader("new-access-token"))
				user, err = api.service.Validation(c.Request.Context(), pairToken.AccessToken)
				if err != nil {
					c.JSON(http.StatusBadGateway, gin.H{
						"error": err.Error(),
					})
					c.Abort()
					return
				}
			}
		}

		c.Set("x-user-id", user.User_id)
		c.Set("x-user-role", user.Role)
		c.Next()
	}
}
