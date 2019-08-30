package server

import (
	"github.com/gin-gonic/gin"
)


func Ping(c *Client, payload interface{})  {

	c.output(gin.H{
		"meta": "pong",
		"data": make(map[string]interface{}),
	})
}
