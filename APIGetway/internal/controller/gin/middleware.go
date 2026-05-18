package gin

import (
	"fmt"
	"net/http"
	"strings"

	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
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

		accessToken := c.GetHeader("access-token")
		refreshToken := c.GetHeader("refresh-token")

		if accessToken == "" || refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "unfind tokens",
			})
			return
		}
		var user userservicedto.Output
		user, err := api.service.Validation(c.Request.Context(), accessToken)
		if err != nil {

			if status.Code(err) == codes.Unauthenticated && strings.Contains(status.Convert(err).Message(), sharederrors.ExpiredToken.Error()) {
				pairToken, err := api.service.UpdateTokens(c.Request.Context(), accessToken, refreshToken)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}
				layout := "2006-01-02 15:04:05"

				c.Header("new-access-token", pairToken.AccessToken)
				c.Header("new-refresh-token", pairToken.RefreshToken)
				c.Header("new-expire_at", pairToken.ExpireAt.Format(layout))

				user, err = api.service.Validation(c.Request.Context(), pairToken.AccessToken)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}

			}
		}
		fmt.Println(user)
		c.Set("x-user-id", user.User_id)
		c.Set("x-user-role", user.Role)
		c.Next()
	}
}
