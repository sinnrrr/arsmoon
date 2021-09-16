 package main
 
import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// The request body representation.
type Body struct {
	Action  string   `json:"action" validate:"required,oneof=subscribe unsubscribe"`
	Symbols []string `json:"symbols"`
}

 func handler(c *gin.Context) {
	var upgrader = websocket.Upgrader{}

	clientConnection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Error upgrading connection: ", err)
	}

	// Infinite loop of reading messages and parsing them.
	for {
		_, msg, err := clientConnection.ReadMessage()
		if err != nil {
			log.Fatal("Error reading message: ", err)
		}

		var body Body
		if err := parseJsonAndValidate(msg, &body); err != nil {
			log.Fatal("Error parsing JSON: ", err)
		}

		switch body.Action {
			case Subscribe: bitmex.Subscribe(clientConnection)
			case Unsubscribe: bitmex.Unsubscribe(clientConnection)
		}
	}
}