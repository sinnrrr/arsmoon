package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/go-playground/validator/v10"
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


func subscriptionHandler(
	clientConnection *websocket.Conn,
	serverConnection *websocket.Conn,
	subscriptionStatusMessageId uint,
	isSubscribeAction bool,
) {
	// Sending unsubscribtion request message.
	err := clientConnection.WriteJSON(NewSubscriptionRequestInfo(isSubscribeAction))
	if err != nil {
		log.Fatalf("error writing subscription message: %s", err)
	}

	var i uint = 0

	for {
		i++

		_, msg, err := clientConnection.ReadMessage()
		if err != nil {
			log.Fatalf("failed reading message from remote server: %s", err)
		}

		// If the message ID equals status message ID:
		if i == subscriptionStatusMessageId {
			// Transforming and validating status message from Bitmex.
			var subscriptionStatus SubscriptionStatus
			if err := parseJsonAndValidate(msg, &subscriptionStatus); err != nil {
				log.Fatalf("error validating subscription message: %s", err)
			}

			// Checking if successfully subscribed to channel.
			if !subscriptionStatus.Success {
				log.Fatalf("error while subscribing to Bitmex: %+v", subscriptionStatus)
			}
		} else if i > subscriptionStatusMessageId && isSubscribeAction {
			var instumentResponse InstrumentResponse
			if err := parseJsonAndValidate(msg, &instumentResponse); err == nil {
				for _, element := range instumentResponse.Data {
					if err = clientConnection.WriteJSON(InstrumentMessage{
						Timestamp: element.Timestamp,
						Symbol:    element.Symbol,
					}); err != nil {
						log.Fatalf("Error writing messages JSON to channel: %s", err)
					}
				}
			}
		}
	}
}