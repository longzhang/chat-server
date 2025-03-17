package tasks

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func CheckOnlineUsers(onlineUsers *sync.Map) {
	for {
		time.Sleep(10 * time.Second)
		onlineUsers.Range(func(key, value interface{}) bool {
			if conn, ok := value.(*websocket.Conn); ok {
				conn.WriteMessage(websocket.TextMessage, []byte("ping"))
			}
			return true
		})
	}
}
