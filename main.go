package main

import "github.com/gin-gonic/gin"

// Initializing global bitmex client.
var bitmex = NewBitmex()

func main() {
	// Using single connection for
	// all users as said in documentation.
	if err := bitmex.Connect(); err != nil {
		panic(err)
	}

	r := gin.Default()

	// Using inline handler to have access to bitmext client
	r.GET("/realtime", handler)

	r.Run()
}
