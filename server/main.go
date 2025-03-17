package main

import (
	"chat-server/server/model"
	"chat-server/server/tasks"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func subscribeToChannel(channel string) {
	pubsub := rdb.Subscribe(ctx, channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Failed to subscribe to channel: %v", err)
	}

	ch := pubsub.Channel()
	go func() {
		for msg := range ch {
			fmt.Printf("Received message from channel %s: %s\n", msg.Channel, msg.Payload)
			var message model.Message
			err := json.Unmarshal([]byte(msg.Payload), &message)
			if err != nil {
				log.Fatalf("Failed to unmarshal message: %v", err)
			}
			sendConnection, ok := onlineUsers.Load(message.To)
			if !ok {
				fmt.Println("User not found")
				continue
			}
			if conn, ok := sendConnection.(*websocket.Conn); ok {
				conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			} else {
				fmt.Println("Invalid connection type")
			}
		}
	}()
}

func publishToChannel(channel string, message model.Message) {
	msg, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	err = rdb.Publish(ctx, channel, msg).Err()
	if err != nil {
		log.Fatalf("Failed to publish message to channel: %v", err)
	}
}

func init() {
	subscribeToChannel("chat-channel")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// var onlineUsers = make(map[string]*websocket.Conn)
var onlineUsers sync.Map

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	uid := c.Request.URL.Query().Get("uid")
	onlineUsers.Store(uid, conn)

	fmt.Print(uid)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}
	defer conn.Close()

	var message model.Message
	// go subscribeToChannel("chat-channel")
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error: ", err)
			break
		}
		fmt.Printf("send message: %\n", string(p))
		json.Unmarshal(p, &message)
		publishToChannel("chat-channel", message)
	}
}

func main() {

	var port string = os.Args[1]

	r := gin.Default()
	r.LoadHTMLGlob("html/*")
	r.GET("/ws", handleWebSocket)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{})
	})
	go tasks.CheckOnlineUsers(&onlineUsers)
	r.Run("0.0.0.0:" + port)
}
