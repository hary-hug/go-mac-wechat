package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-chat/app/chat/model"
	"go-chat/app/chat/util/helper"
	"strings"
	"time"
)

func GetChats(c *Client, payload string)  {

	var (
		chats      []model.ChatList
		res        []interface{}
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]

	if !ok {
		return
	}

	DaoInstance.ChatDb.Model(model.ChatList{}).
		Where("uid = ?", uid).
		Order("active_time desc").
		Find(&chats)

	if len(chats) <= 0 {
		return
	}


	for i := range chats {

		var (
			chatUser  model.User
			title     string
			avatar    string
			lastMessage    model.ChatMessage
		)

		switch chats[i].ChatType {
		// private chat
		case 1:
			DaoInstance.ChatDb.Model(model.User{}).Where("uid = ?", chats[i].RelateId).Find(&chatUser)
			title = chatUser.Nickname
			avatar = chatUser.Avatar
			break
		// group chat
		case 2:
			break
		}

		// get last message
		DaoInstance.ChatDb.Model(model.ChatMessage{}).Where("chat_id = ?", chats[i].ChatId).Order("msg_id desc").First(&lastMessage)

		item := map[string]interface{} {
			"chat_id": chats[i].ChatId,
			"title":   title,
			"avatar":  avatar,
			"msg":     helper.Substr(lastMessage.Content, 10),
			"new":     chats[i].New,
		}

		res = append(res, item)
	}

	c.output(gin.H{
		"meta": "chats",
		"data": res,
	})

}


// CreateChat create a new chat with friend
func CreateChat(c *Client, payload string)  {

	var (
		chatList  model.ChatList
		friend    model.User
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]

	if !ok {
		return
	}

	params := new(struct {
		FriendId   int  `form:"friend_id" json:"friend_id"`
	})

	json.Unmarshal([]byte(payload), &params)

	// todo check if exists the friend's relationship
	if params.FriendId <= 0 {
		return
	}

	// get friend info
	DaoInstance.ChatDb.Model(model.User{}).Where("uid = ?", params.FriendId).Find(&friend)
	// check if exists
	DaoInstance.ChatDb.Model(model.ChatList{}).Where("uid = ? and relate_id = ?", uid, params.FriendId).Find(&chatList)

	if chatList.ListId <= 0 {
		return
	}

	item := map[string]interface{} {
		"chat_id": chatList.ChatId,
		"title":   friend.Nickname,
		"avatar":  friend.Avatar,
		"msg":     "",
		"new":     0,
	}

	c.output(gin.H{
		"meta": "newChat",
		"data": item,
	})

}


// GetChatMessages return messages of a chat
func GetChatMessages(c *Client, payload string)  {

	var (
		chat         model.ChatList
		chatMessages []model.ChatMessage
		lists        []interface{}
		title        string
		avatar       string
		chatUser     model.User
		res   	     map[string]interface{}
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]

	if !ok {
		return
	}

	params := new(struct {
		ChatId   int  `form:"chat_id" json:"chat_id"`
	})

	json.Unmarshal([]byte(payload), &params)

	DaoInstance.ChatDb.Model(model.ChatList{}).Where("chat_id = ? and uid = ?", params.ChatId, uid).Find(&chat)

	if chat.ChatId <= 0 {
		return
	}

	// get 10 latest message
	DaoInstance.ChatDb.Model(model.ChatMessage{}).
		Where("chat_id = ?", params.ChatId).
		Order("msg_id desc").
		Limit(10).
		Find(&chatMessages)

	if len(chatMessages) <= 0 {
		lists = make([]interface{}, 0)
	}


	DaoInstance.ChatDb.Model(model.ChatUser{})

	// reverse message in order to display in brower
	for i := len(chatMessages)-1 ; i >= 0; i-- {

		var (
			self int
			user  model.User
		)
		if chatMessages[i].Uid == uid {
			self = 1
		}

		DaoInstance.ChatDb.Model(model.User{}).Where("uid = ?", chatMessages[i].Uid).Find(&user)

		item := map[string]interface{} {
			"msg_id":      chatMessages[i].MsgId,
			"msg_type":    chatMessages[i].MsgType,
			"self":        self,
			"nickname":    user.Nickname,
			"avatar":      user.Avatar,
			"content":     strings.Replace(chatMessages[i].Content, "\n", "<br/>", -1),
			"create_time": time.Unix(int64(chatMessages[i].CreateTime),0).Format("2006/01/02 03:04:05"),
		}

		lists = append(lists, item)
	}

	switch chat.ChatType {
	// private chat
	case 1:
		DaoInstance.ChatDb.Model(model.User{}).Where("uid = ?", chat.RelateId).Find(&chatUser)
		title = chatUser.Nickname
		avatar = chatUser.Avatar
		break
		// group chat
	case 2:
		break
	}
	
	res = map[string]interface{}{
		"chat_id": params.ChatId,
		"title":   title,
		"avatar":  avatar,
		"lists":   lists,
	}

	c.output(gin.H{
		"meta": "chatMessages",
		"data": res,
	})
}

// SendChatMessage send a message to friend or group
func SendChatMessage(c *Client, payload string)  {

	var (
		chatMessage  model.ChatMessage
		chatUser     []model.ChatUser
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]

	if !ok {
		return
	}

	params := new(struct {
		ChatId   int  `form:"chat_id" json:"chat_id"`
		Content  string  `form:"content" json:"content"`
	})

	json.Unmarshal([]byte(payload), &params)

	DaoInstance.ChatDb.Model(model.ChatUser{}).Where("chat_id = ?", params.ChatId).Find(&chatUser)

	if len(chatUser) <= 0 {
		return
	}

	// add new message to database
	chatMessage.Uid = uid
	chatMessage.ChatId = params.ChatId
	chatMessage.CreateTime = int(time.Now().Unix())
	chatMessage.Content = params.Content
	chatMessage.MsgType = 1
	DaoInstance.ChatDb.Create(&chatMessage)

	for i := range chatUser {

		var (
			user model.User
			self int
		)

		// dipatch message to client
		sendto := chatUser[i].Uid

		// send to myself
		if uid == sendto {
			self = 1
		}

		// update chat list active time
		go func() {
			DaoInstance.ChatDb.Model(model.ChatList{}).Where("chat_id = ?", params.ChatId).Update("active_time", time.Now().Unix())
		}()

		DaoInstance.ChatDb.Model(model.User{}).Where("uid = ?", chatMessage.Uid).Find(&user)

		item := map[string]interface{} {
			"msg_id":       chatMessage.MsgId,
			"msg_type":     chatMessage.MsgType,
			"self":         self,
			"nickname":     user.Nickname,
			"avatar":       user.Avatar,
			"content":      chatMessage.Content,
			"create_time" : time.Unix(int64(chatMessage.CreateTime),0).Format("2006/01/02 03:04:05"),
		}

		// response result
		res := map[string]interface{} {
			"chat_id": params.ChatId,
			"message": item,
		}

		// get client connection by uid
		client, ok := WsInstance.Uids[sendto]
		if ok {
			// send to online client only
			client.output(gin.H{
				"meta": "newMessage",
				"data": res,
			})
		}

	}


}