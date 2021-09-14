package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type SubscriptionRequestInfo struct {
	Op string `json:"op" validate:"oneof=subscribe unsubscribe"`
	Args string `json:"args" validate:"eq=instrument"`
}

type SubscriptionStatus struct {
	Success   bool   `json:"success" validate:"required"`
	Subscribe string `json:"subscribe" validate:"eq=instrument"`
	Request SubscriptionRequestInfo `json:"request" validate:"required,dive,required"`
}

type InstrumentResponse struct {
	Timestamp uint16  `json:"timestamp"`
	Symbol    string  `json:"symbol"`
	Price     float32 `json:"price"`
}

type Bitmex struct {
	URL        *url.URL
	Connection *websocket.Conn
	Channel    chan InstrumentResponse
}

func NewBitmex() *Bitmex {
	return &Bitmex{
		URL: &url.URL{
			Scheme:   "wss",
			Host:     "www.bitmex.com",
			Path:     "realtime",
			RawQuery: "subscribe=instrument",
		},
	}
}

func (bitmex *Bitmex) Connect() {
	conn, _, err := websocket.DefaultDialer.Dial(bitmex.URL.String(), nil)
	if err != nil {
		log.Fatal("Error dialing remote server: ", err)
	}

	bitmex.Connection = conn
}

func (bitmex *Bitmex) Subscribe() {
	// var messages chan InstrumentResponse

	go func() {
		i := 0;
		for {
			i++;

			_, msg, err := bitmex.Connection.ReadMessage()
			if err != nil {
				log.Fatal("Failed reading message from remote server: ", err)
			}

			// First two messages are
			// just welcome ones (omiting push to the channel)
			if (i == 2) {
				// Second message resproduces the subscription status
				var subscriptionStatus SubscriptionStatus;
				if err := parseJsonAndValidate(msg, &subscriptionStatus); err != nil {
					log.Fatal("Error validating request: ", err)
				}

				if !subscriptionStatus.Success {
					log.Fatal("Error while subscribing to Bitmex: ", )
				}
			} else if (i > 2) {
				// todo
			}
		}
	}()
}
