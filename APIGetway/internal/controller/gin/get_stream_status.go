package gin

import (
	"context"
	"encoding/json"
	"net/http"

	orderdto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/order_service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (api *GinAPI) GetStreamStatus(c *gin.Context) {
	var orderInfo orderdto.GetInput

	err := c.ShouldBindJSON(&orderInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
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
	orderInfo.User_id = user_id

	ws, err := api.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer ws.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	msgChan := make(chan orderdto.StreamOutput, 10)
	errChan := make(chan error, 1)

	go func() {
		err := api.service.GetStreamStatus(ctx, orderInfo, msgChan)
		if err != nil {
			errChan <- err
		}
		close(msgChan)
	}()

	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				ws.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "stream finished"))
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				errorMsg, _ := json.Marshal(map[string]string{
					"error": err.Error(),
				})
				ws.WriteMessage(websocket.CloseMessage, errorMsg)
				return
			}

			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				errorMsg, _ := json.Marshal(map[string]string{
					"error": err.Error(),
				})
				ws.WriteMessage(websocket.CloseMessage, errorMsg)
				cancel()
				return
			}

		case err := <-errChan:
			errorMsg, _ := json.Marshal(map[string]string{
				"error": err.Error(),
			})
			ws.WriteMessage(websocket.CloseMessage, errorMsg)
			return

		case <-ctx.Done():

			ws.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseGoingAway, "connection closed"))
			return
		}
	}
}
