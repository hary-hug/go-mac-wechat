package server

import (
	"errors"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat/app/chat/util/helper"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Ws struct {
	Clients  map[*Client]int
	Uids     map[int]*Client
}


func newWs() *Ws {
	return &Ws{
		Clients: make(map[*Client]int),
		Uids: make(map[int]*Client),
	}
}


func (s *Ws) addClient(c *Client, uid int)  {
	s.Clients[c] = uid
	s.Uids[uid] = c
}


func (s *Ws) removeClient(c *Client)  {
	if _, ok := s.Clients[c]; ok {
		delete(s.Clients, c)
		delete(s.Uids, s.Clients[c])
	}
}


func (s *Ws) serve(ctx *gin.Context) (err error) {

	// check token
	token := com.StrTo(ctx.Query("token")).String()

	claims, err := helper.ParseToken(token)

	if err != nil {
		return
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return errors.New("the token has expired")
	}

	// upgrade get to webSocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	// when client reconnect
	// we need to remove invalid current user's connection
	s.removeClient(s.Uids[claims.Uid])
	// add new connection with uid
	s.addClient(client, claims.Uid)


	go client.readLoop()
	go client.writeLoop()

	return nil
}
