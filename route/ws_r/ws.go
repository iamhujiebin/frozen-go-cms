package ws_r

import (
	"fmt"
	"frozen-go-cms/_const/enum/ws_e"
	"frozen-go-cms/req/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clientMap sync.Map

func StoreClient(token string, ws *websocket.Conn) error {
	claim, err := jwt.ParseToken(token)
	if err != nil {
		return err
	}
	userId := claim.UserId
	if data, ok := clientMap.Load(userId); ok {
		data.(*sync.Map).Store(token, ws)
	} else {
		var userMap = new(sync.Map)
		userMap.Store(token, ws)
		clientMap.Store(userId, userMap)
	}
	return nil
}

func RemoveClient(token string) {
	claim, err := jwt.ParseToken(token)
	if err != nil {
		return
	}
	userId := claim.UserId
	if data, ok := clientMap.Load(userId); ok {
		data.(*sync.Map).Range(func(key, value interface{}) bool {
			if key.(string) == token {
				data.(*sync.Map).Delete(token)
				return false // stop range
			}
			return true
		})
	} else {
	}
	return
}

func WsHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	token := c.Param("token")
	if err := StoreClient(token, ws); err != nil {
		return
	}

	defer func() {
		_ = ws.Close()
		RemoveClient(token)
	}()
	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		//If client message is ping will return pong
		if string(message) == "ping" && mt == websocket.TextMessage {
			//redirect to other actions
			message = []byte("pong")
			err = ws.WriteMessage(mt, message)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}

func WsTest(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Request.URL.Query().Get("uid"), 10, 64)
	if userMap, ok := clientMap.Load(userId); ok {
		userMap.(*sync.Map).Range(func(key, value interface{}) bool {
			if ws, ok := value.(*websocket.Conn); ok {
				msg := c.Request.URL.Query().Get("msg")
				ws.WriteMessage(websocket.TextMessage, []byte(msg))
			}
			return true
		})
	}
}

// 发信息到客户端
func SendToClient(userId uint64, cmd ws_e.CMD) {
	if userMap, ok := clientMap.Load(userId); ok {
		userMap.(*sync.Map).Range(func(key, value interface{}) bool {
			if ws, ok := value.(*websocket.Conn); ok {
				if err := ws.WriteMessage(websocket.TextMessage, []byte(cmd)); err != nil {
					fmt.Println(err)
				}
			}
			return true
		})
	}
}
