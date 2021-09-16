package main

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
	Instrument  = "instrument"
)

// The message, which send when you need to
// subscribe or unsubscribe to services on Bitmex channel.
type SubscriptionRequestInfo struct {
	Op   string `json:"op" validate:"oneof=subscribe unsubscribe"`
	Args string `json:"args" validate:"eq=instrument"`
}

// The message, which comes from Bitmex channel
// to check if subscription request was successful.
type SubscriptionStatus struct {
	Success     bool                    `json:"success" validate:"required"`
	Subscribe   string                  `json:"subscribe" validate:"eq=instrument"`
	Unsubscribe string                  `json:"unsubscribe" validate:"eq=instrument"`
	Request     SubscriptionRequestInfo `json:"request" validate:"required,dive,required"`
}

// Form of the message, which sends to client from this channel.
// Used for transforming to from Bitmex instrument message.
type InstrumentMessage struct {
	Timestamp string  `json:"timestamp"`
	Symbol    string  `json:"symbol"`
	Price     float32 `json:"price"`
}

type InstrumentInfo struct {
	Symbol         string  `json:"symbol"`
	Timestamp      string  `json:"timestamp"`
	ImpactAskPrice float32 `json:"impactAskPrice"`
	MarkPrice      float32 `json:"markPrice"`
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

// Establish connection to Bitmex websocket server.
func (b *Bitmex) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(b.URL.String(), nil)
	if err != nil {
		return fmt.Errorf("error dialing remote server: %s", err)
	}

	// Passing connection to struct in order other handlers
	// to have easy idiomatic access to websocket connection.
	b.Connection = conn

	return nil
}

// Subscribe to instruments service.
func (b *Bitmex) Subscribe(serverConnection *websocket.Conn) {
	go subscriptionHandler(
		b.Connection,
		serverConnection,
		2,
		true,
	)
}

// Unsubscribe from instruments service.
func (b *Bitmex) Unsubscribe(serverConnection *websocket.Conn) {
	go subscriptionHandler(
		b.Connection,
		serverConnection,
		1,
		false,
	)
}

