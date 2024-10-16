package main

import (
	minigin "github.com/rosas99/streaming/examples/mini-gin/gin"
	"log"
	"time"
)

func Logger() minigin.HandlerFunc {
	return func(c *minigin.Context) {
		t := time.Now()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
		c.Next()
	}
}
