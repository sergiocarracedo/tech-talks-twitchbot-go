package notifications

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func lastNotificationHandler() func(*gin.Context) {
	return func(ctx *gin.Context) {
		//Upgrade get request to webSocket protocol
		ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("error get connection")
			log.Fatal(err)
		}
		defer ws.Close()
		clients[ws] = true

		for {
			// Grab the next message from the broadcast channel
			msg := <-broadcast
			// Send it out to every client that is currently connected
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func SendMessages () {
	ticker := time.NewTicker(time.Second * 1)
	defer func() {
		ticker.Stop()
	}()
	i := 0

	for {
		<- ticker.C
		i = i + 1
		log.Println("Broadcast", i)
		broadcast <- fmt.Sprintf("%d", i)
	}

	log.Println("End send messages")
}