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

type SubscriptionRequestInfo struct {
	Op   string `json:"op" validate:"oneof=subscribe unsubscribe"`
	Args string `json:"args" validate:"eq=instrument"`
}

type SubscriptionStatus struct {
	Success   bool                    `json:"success" validate:"required"`
	Subscribe string                  `json:"subscribe" validate:"eq=instrument"`
	Request   SubscriptionRequestInfo `json:"request" validate:"required,dive,required"`
}

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

type Bitmex struct {
	URL        *url.URL
	Connection *websocket.Conn
	Channel    chan InstrumentMessage
}

func NewBitmex() *Bitmex {
	return &Bitmex{
		URL: &url.URL{
			Scheme: "wss",
			Host:   "www.bitmex.com",
			Path:   "realtime",
		},
	}
}

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

// Establish connection to Bitmex websocket server
func (bitmex *Bitmex) Connect() {
	conn, _, err := websocket.DefaultDialer.Dial(bitmex.URL.String(), nil)
	if err != nil {
		log.Fatal("Error dialing remote server: ", err)
	}

	bitmex.Connection = conn
}

func (bitmex *Bitmex) Subscribe() {
	err := bitmex.Connection.WriteJSON(NewSubscriptionRequestInfo(true))
	if err != nil {
		log.Fatal("Error writing subscription message: ", err)
	}

	if err := readMessages(bitmex.Connection, 2); err != nil {
		log.Fatal(err)
	}
}

func (bitmex *Bitmex) Unsubscribe() {
	err := bitmex.Connection.WriteJSON(NewSubscriptionRequestInfo(false))
	if err != nil {
		log.Fatal("Error writing subscription message: ", err)
	}

	if err := readMessages(bitmex.Connection, 0); err != nil {
		log.Fatal(err)
	}
}
