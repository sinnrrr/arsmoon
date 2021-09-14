package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/ws", handler)

	r.Run()
}