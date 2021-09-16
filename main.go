package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Initializing global bitmex client.
var bitmex = NewBitmex()

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Using single connection for
	// all users as said in documentation.
	if err := bitmex.Connect(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/realtime", handler)

	r.Run()
}
