package main

import (
	"fmt"
	"net/url"
	// "os"
	// "net/http"
	// "strconv"
	// "time"

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
	Subscribe   string                  `json:"subscribe"`
	Unsubscribe string                  `json:"unsubscribe"`
	Request     SubscriptionRequestInfo `json:"request" validate:"required,dive,required"`
}

// The type of message, which sends, when 
// the user is informed about subscription result.
type SubscriptionStatusMessage struct {
	Success bool `json:"success"`
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

/// JSON form of the message, which is used 
// for transforming into own form (InstrumentMessage struct).
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
	isSubscribeAction bool,
) *SubscriptionRequestInfo {
	var op string

	if isSubscribeAction {
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
	// var (
	// 	data    = ""
	// 	method  = "GET"
	// 	expires = strconv.FormatInt(time.Now().Add(time.Hour*24).Unix(), 10)
	// )

	conn, _, err := websocket.DefaultDialer.Dial(b.URL.String(), /*http.Header{
		"Api-Expires":   []string{expires},
		"Api-Key":       []string{os.Getenv("BITMEX_API_KEY")},
		"Api-Signature": []string{generateApiSignature(
			os.Getenv("BITMEX_API_SECRET"),
			method,
			b.URL.Path,
			expires,
			data,
		)},
	}*/ nil)
	if err != nil {
		return fmt.Errorf("error dialing remote server: %s", err)
	}

	// Passing connection to struct in order other handlers
	// to have easy access to websocket connection.
	b.Connection = conn

	return nil
}

// Subscribe to instruments service.
func (b *Bitmex) Subscribe(serverConnection *websocket.Conn) {
	subscriptionHandler(
		b.Connection,
		serverConnection,
		true,
	)
}

// Unsubscribe from instruments service.
func (b *Bitmex) Unsubscribe(serverConnection *websocket.Conn) {
	subscriptionHandler(
		b.Connection,
		serverConnection,
		false,
	)
}
