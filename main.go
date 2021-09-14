package main

import "github.com/gin-gonic/gin"

var bitmex = NewBitmex()

func main() {
	r := gin.Default()

	r.GET("/ws", handler)

	r.Run()
}
