package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Body struct {
	Action  string   `json:"action" validate:"required,oneof=subscribe unsubscribe"`
	Symbols []string `json:"symbols"`
}

var upgrader = websocket.Upgrader{}

func handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Error upgrading connection: ", err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if (err != nil) {
			log.Fatal("Error reading message: ", err)
		}

		var body Body;
		if err := parseJsonAndValidate(msg, &body); err != nil {
			log.Fatal("Error parsing JSON: ", err)
		}

		// Connecting in handler not to waste connection resources
		bitmex.Connect()

		switch body.Action {
			case Subscribe: bitmex.Subscribe()
			case Unsubscribe: bitmex.Unsubscribe()
		}

		fmt.Printf("%+v\n", <-instrumentMessages)
	}
}