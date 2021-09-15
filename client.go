package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
	Instrument  = "instrument"
)

// Initializing global bitmex client.
var bitmex = NewBitmex()

// The message, which send when you need to
// subscribe or unsubscribe to services on Bitmex channel.
type SubscriptionRequestInfo struct {
	Op   string `json:"op" validate:"oneof=subscribe unsubscribe"`
	Args string `json:"args" validate:"eq=instrument"`
}

// The message, which comes from Bitmex channel
// to check if subscription request was successful.
type SubscriptionStatus struct {
	Success   bool                    `json:"success" validate:"required"`
	Subscribe string                  `json:"subscribe" validate:"eq=instrument"`
	Request   SubscriptionRequestInfo `json:"request" validate:"required,dive,required"`
}

// Form of the message, which sends to client from this channel.
// Used for transforming to from Bitmex instrument message.
type InstrumentMessage struct {
	Timestamp uint16  `json:"timestamp"`
	Symbol    string  `json:"symbol"`
	Price     float32 `json:"price"`
}

type InstrumentInfo struct {
	Symbol             string  `json:"symbol"`
	LastPriceProtected float32 `json:"lastPriceProtected"`
	BidPrice           float32 `json:"bidPrice"`
	MidPrice           float32 `json:"midPrice"`
	AskPrice           float32 `json:"askPrice"`
	ImpactBidPrice     float32 `json:"impactBidPrice"`
	ImpactMidPrice     float32 `json:"impactMidPrice"`
	ImpactAskPrice     float32 `json:"impactAskPrice"`
	Timestamp          string  `json:"timestamp"`
}

type InstrumentResponse struct {
	Table  string           `json:"table" validate:"eq=instrument"`
	Action string           `json:"action" validate:"eq=update"`
	Data   []InstrumentInfo `json:"data" validate:"required,dive,required"`
}

// Bitmex client.
type Bitmex struct {
	URL        *url.URL
	Connection *websocket.Conn
	Channel    chan InstrumentMessage
}

// Constructor for Bitmex client.
func NewBitmex() *Bitmex {
	return &Bitmex{
		URL: &url.URL{
			Scheme: "wss",
			Host:   "www.bitmex.com",
			Path:   "realtime",
		},
	}
}

// The subscription request constructor.
func NewSubscriptionRequestInfo(
	subscribe bool,
) *SubscriptionRequestInfo {
	var op string

	if subscribe {
		op = Subscribe
	} else {
		op = Unsubscribe
	}

	return &SubscriptionRequestInfo{
		Op:   op,
		Args: Instrument,
	}
}

// Helper function, which reads messages from the Bitmex channel.
func (bitmex *Bitmex) ReadMessages(
	/*
		The ID of the message, which represents the status
	 	of the subscrition, to check if everything got fine.

		All the messages before are ignored.
	*/
	statusMessageId uint,
) error {
	var i uint = 0
	for {
		i++

		_, msg, err := bitmex.Connection.ReadMessage()
		if err != nil {
			log.Fatal("failed reading message from remote server: %s", err)
		}

		// If ID is set:
		if statusMessageId > 0 {
			// If the message ID equals status message ID:
			if i == statusMessageId {
				// Transforming and validating status message from Bitmex.
				var subscriptionStatus SubscriptionStatus
				if err := parseJsonAndValidate(msg, &subscriptionStatus); err != nil {
					log.Fatal("error validating subscription message: %s", err)
				}

				// Checking if successfully subscribed to channel.
				if !subscriptionStatus.Success {
					log.Fatal("error while subscribing to Bitmex")
				}

				// Reseting status message to zero
				// in order to continue looping through messages.
				statusMessageId = 0
			}
		} else {
			panic(string(msg))
			// var instrumentResponse InstrumentResponse
			// if err := parseJsonAndValidate(msg, &instrumentResponse); err != nil {
			// 	log.Fatal("error validating instrument message: %s", err)
			// }

			// if err = conn.WriteJSON(instrumentResponse); err != nil {
			// 	log.Fatal(err)
			// }
		}
	}
}


// Establish connection to Bitmex websocket server.
func (bitmex *Bitmex) Connect() {
	conn, _, err := websocket.DefaultDialer.Dial(bitmex.URL.String(), nil)
	if err != nil {
		log.Fatal("Error dialing remote server: ", err)
	}

	// Passing connection to struct in order other handlers
	// to have easy idiomatic access to websocket connection.
	bitmex.Connection = conn
}

// Subscribe to instruments service.
func (bitmex *Bitmex) Subscribe() {
	// Sending subscription request message.
	err := bitmex.Connection.WriteJSON(NewSubscriptionRequestInfo(true))
	if err != nil {
		log.Fatal("Error writing subscription message: ", err)
	}

	if err := bitmex.ReadMessages(2); err != nil {
		log.Fatal(err)
	}
}

// Unsubscribe from instruments service.
func (bitmex *Bitmex) Unsubscribe() {
	// Sending unsubscribtion request message.
	err := bitmex.Connection.WriteJSON(NewSubscriptionRequestInfo(false))
	if err != nil {
		log.Fatal("Error writing subscription message: ", err)
	}

	if err := bitmex.ReadMessages(0); err != nil {
		log.Fatal(err)
	}
}
