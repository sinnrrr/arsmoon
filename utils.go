package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

var validate = validator.New()

func parseJsonAndValidate(data []byte, s interface{}) error {
	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("failed unmarshaling JSON: %s", err)
	}
	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("failed validating unmarshaled JSON: %s", err)
	}

	return nil
}

func readMessages(
	conn *websocket.Conn,
	statusMessageId uint,
) error {
	var i uint = 0
	for {
		i++

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed reading message from remote server: %s", err)
		}

		if statusMessageId > 0 {
			if i == statusMessageId {
				// Second message resproduces the subscription status
				var subscriptionStatus SubscriptionStatus
				if err := parseJsonAndValidate(msg, &subscriptionStatus); err != nil {
					return fmt.Errorf("error validating subscription message: %s", err)
				}
	
				// Checking if successfully subscribed to channel
				if !subscriptionStatus.Success {
					return fmt.Errorf("error while subscribing to Bitmex")
				}

				statusMessageId = 0
			}
		} else {
			log.Printf(string(msg))
			panic("asd")
			// var instrumentResponse InstrumentResponse
			// if err := parseJsonAndValidate(msg, &instrumentResponse); err != nil {
			// 	return fmt.Errorf("error validating instrument message: %s", err)
			// }

			// if err = conn.WriteJSON(instrumentResponse); err != nil {
			// 	log.Fatal(err)
			// }
		}
	}
}
