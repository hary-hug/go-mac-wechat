package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"reflect"
	"time"
)

const (
	_writeWait  = 10 * time.Second
	_pongWait   = 60 * time.Second
	_pingPeriod = (_pongWait * 9) / 10
	_maxMessageSize = 512
)


type Client struct {
	// current websocket connection
	conn    *websocket.Conn
	// ch transmit data between readLoop and writeLoop
	send    chan []byte
	// request handlers
	handles reflect.Type
}


func (c *Client) close() {
	c.conn.Close()
}


func (c *Client) readLoop() {

	defer func() {
		c.close()
	}()

	c.conn.SetReadLimit(_maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(_pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(_pongWait))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		c.dispatch(raw)
	}

}


func (c *Client) writeLoop() {

	ticker := time.NewTicker(_pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {

		select {
		case message, ok := <- c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// send data to remote client
			c.conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			// ping the remote client
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// send data to client
func (c *Client) output(data map[string]interface{}) {
	result, _ := json.Marshal(data)
	c.send <- result
}


// dispatch 处理客户端发来的数据
func (c *Client) dispatch(raw []byte)  {

	msg := new(struct {
		Event    string   `json:"event"`
		Payload  string   `json:"payload"`
	})

	if err := json.Unmarshal(raw, &msg); err != nil {
		log.Println(err)
		return
	}

	// 注册业务处理事件
	handlers := map[string]interface{}{
		"ping":                   Ping,
		"getChats":               GetChats,
		"createChat":             CreateChat,
		"getChatMessages":        GetChatMessages,
		"sendChatMessage":        SendChatMessage,
		"getFriends":             GetFriends,
		"getMe":                  GetMe,
		"getPublicUsers":         GetPublicUsers,
		"sendFriendRequest":      SendFriendRequest,
		"getFriendRequests":      GetFriendRequests,
		"sendAgreeFriendRequest": AgreeFriendRquest,
	}

	_, err := caller(handlers, msg.Event, c, msg.Payload)

	if err != nil {
		log.Println(err)
	}

}


func caller(handlers map[string]interface{}, event string, params ...interface{}) ([]reflect.Value, error) {

	if _, ok := handlers[event]; !ok {
		return nil, errors.New(fmt.Sprintf("the event: %s does not have the handle function", event))
	}

	f := reflect.ValueOf(handlers[event])

	if len(params) != f.Type().NumIn() {
		return nil, errors.New("the number of input params does not match")
	}

	in := make([]reflect.Value, len(params))

	for k, v := range params {
		in[k] = reflect.ValueOf(v)
	}

	return f.Call(in), nil
}
