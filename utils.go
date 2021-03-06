package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

// A method, which handles the subsctiption part of the service.
// It makes subsctiption request, validates the result of it and outputs the results.
func subscriptionHandler(
	// Bitmex connection.
	clientConnection *websocket.Conn,
	// Server connection (gorilla webscoket).
	serverConnection *websocket.Conn,
	isSubscribeAction bool,
) {
	const SubscriptionStatusMessageId = 2

	// Sending unsubscribtion request message.
	err := clientConnection.WriteJSON(NewSubscriptionRequestInfo(isSubscribeAction))
	if err != nil {
		log.Fatalf("error writing subscription message: %s", err)
	}

	var i uint = 0

	go func() {
		for {
			i++

			_, msg, err := clientConnection.ReadMessage()
			if err != nil {
				break
			}

			// If the message ID equals status message ID:
			if SubscriptionStatusMessageId > 0 && i == SubscriptionStatusMessageId {
				// Transforming and validating status message from Bitmex.
				var subscriptionStatus SubscriptionStatus

				if err := parseJsonAndValidate(msg, &subscriptionStatus); err != nil {
					log.Printf("error validating subscription message: %s", err)
					continue
				}

				if err = serverConnection.WriteJSON(SubscriptionStatusMessage{
					Success: subscriptionStatus.Success,
				}); err != nil {
					log.Printf("Error writing messages JSON to channel: %s", err)
					continue
				}
			}

			if i > SubscriptionStatusMessageId && isSubscribeAction {
				var instumentResponse InstrumentResponse
				if err := parseJsonAndValidate(msg, &instumentResponse); err != nil {
					continue
				}

				for _, element := range instumentResponse.Data {
					if err = serverConnection.WriteJSON(InstrumentMessage{
						Timestamp: element.Timestamp,
						Symbol:    element.Symbol,
						Price: maxOutOfSlice([]float32{
							element.ImpactAskPrice,
							element.LastPrice,
							element.MarkPrice,
						}),
					}); err != nil {
						log.Fatalf("Error writing messages JSON to channel: %s", err)
						continue
					}
				}
			}
		}
	}()

}

func generateApiSignature(apiSecret, method, path, expires, data string) string {
	hash := hmac.New(sha256.New, []byte(apiSecret))
	hash.Write([]byte(method + path + expires + data))

	return hex.EncodeToString(hash.Sum(nil))
}

func maxOutOfSlice(input []float32) float32 {
	var max = input[0]
	for _, value := range input {
		if value > max {
			max = value
		}
	}

	return max
}
