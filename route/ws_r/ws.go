package ws_r

import (
	"fmt"
	"git.hilo.cn/hilo-common/mycontext"
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

func WsHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	userId, exists := c.Get(mycontext.USERID)
	if !exists {
		c.Writer.Write([]byte("not Authorization"))
		return
	}
	clientMap.Store(userId, ws)
	defer func() {
		_ = ws.Close()
		clientMap.Delete(userId)
	}()
	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		//If client message is ping will return pong
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func WsTest(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Request.URL.Query().Get("uid"), 10, 64)
	if ws, ok := clientMap.Load(userId); ok {
		msg := c.Request.URL.Query().Get("msg")
		ws.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(msg))
	}
}
