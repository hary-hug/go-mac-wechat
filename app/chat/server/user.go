package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-chat/app/chat/model"
	"time"
)


func GetFriends(c *Client, payload string)  {

	var (
		users     []model.User
		friends   []model.UserFriend
		friendIds []int
		res       []interface{}
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]
	if !ok {
		return
	}

	DaoInstance.ChatDb.Model(model.UserFriend{}).Where("uid = ?", uid).Find(&friends)

	for _, friend := range friends {
		friendIds = append(friendIds, friend.FriendId)
	}

	if len(friendIds) <= 0 {
		return
	}

	DaoInstance.ChatDb.Model(model.User{}).Where("uid IN (?) ", friendIds).Find(&users)

	for i := range users {

		item := map[string]interface{} {
			"uid":       users[i].Uid,
			"nickname":  users[i].Nickname,
			"avatar":    users[i].Avatar,
			"email":     users[i].Email,
			"signature": users[i].Signature,
		}

		res = append(res, item)
	}

	c.output(gin.H{
		"meta": "friends",
		"data": res,
	})
}


// SendFriendRequest send a friend request
func SendFriendRequest(c *Client, payload string)  {

	var (
		req  model.UserFriendRequest
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]
	if !ok {
		return
	}

	params := new(struct {
		Uid   int  `form:"uid" json:"uid"`
	})

	json.Unmarshal([]byte(payload), &params)

	if uid <= 0 {
		return
	}

	if uid == params.Uid {
		// can not send a request to my self
		return
	}

	req.Uid = params.Uid
	req.FromUid = uid
	req.Status = 0
	req.CreateTime = int(time.Now().Unix())

	DaoInstance.ChatDb.Create(&req)

	if req.ReqId <= 0 {
		return
	}

	item := map[string]interface{} {
		"count": 1,
	}

	client, ok := WsInstance.Uids[params.Uid]
	if ok {
		// send to online client only
		client.output(gin.H{
			"meta": "newContactNo",
			"data": item,
		})
	}

}

// GetReqs return a list of friend request
func GetFriendRequests(c *Client, payload string)  {

	var (
		reqs  []model.UserFriendRequest
		res   []interface{}
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]
	if !ok {
		return
	}

	DaoInstance.ChatDb.Model(model.UserFriendRequest{}).Where("uid = ?", uid).Find(&reqs)

	for i := range reqs {
		var (
			user model.User
		)
		// get request user info
		DaoInstance.ChatDb.Model(model.User{}).Where("uid = ?", reqs[i].FromUid).Find(&user)
		if user.Uid <= 0 {
			continue
		}

		item := map[string]interface{} {
			"req_id":    reqs[i].ReqId,
			"uid":       user.Uid,
			"nickname":  user.Nickname,
			"avatar":    user.Avatar,
			"email":     user.Email,
			"signature": user.Signature,
			"status":    reqs[i].Status,
		}

		res = append(res, item)
	}

	if len(res) <= 0 {
		return
	}

	c.output(gin.H{
		"meta": "friendRequests",
		"data": res,
	})
}

// agree a user friend requst
func AgreeFriendRquest(c *Client, payload string)  {

	var (
		friendRequest  model.UserFriendRequest
		chat  model.Chat
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]
	if !ok {
		return
	}

	params := new(struct {
		ReqId   int  `form:"req_id" json:"req_id"`
	})

	json.Unmarshal([]byte(payload), &params)

	DaoInstance.ChatDb.Model(model.UserFriendRequest{}).
		Where("req_id = ?", params.ReqId).
		Find(&friendRequest)

	// request doest not exist
	if friendRequest.ReqId <= 0 {
		return
	}

	if friendRequest.Status != 0 {
		return
	}


	DaoInstance.AddFriend(uid, friendRequest.FromUid)

	// upadte request status
	go func() {
		DaoInstance.ChatDb.Model(model.UserFriendRequest{}).Where("req_id = ?", params.ReqId).Update("status", 1)
	}()

	uids := [2]int{uid, friendRequest.FromUid}

	for i, id :=range uids {
		var (
			friend model.User
			sendto int
		)
		// get friend info
		DaoInstance.ChatDb.Model(model.User{}).Where("uid = ?", id).Find(&friend)

		item := make(map[string]interface{})

		item["chat_id"] = chat.ChatId
		item["title"] = friend.Nickname
		item["avatar"] = friend.Avatar
		item["msg"] = ""
		item["new"] = 0

		if i == 0 {
			sendto = friendRequest.FromUid
		} else {
			sendto = uid
		}

		// get client connection by uid
		// send to friend
		client, ok := WsInstance.Uids[sendto]
		if ok {
			// send to online client only
			client.output(gin.H{
				"meta": "newChat",
				"data": item,
			})
		}

	}

}


func GetPublicUsers(c *Client, payload string)  {

	var (
		users  []model.User
		res    []interface{}
	)

	DaoInstance.ChatDb.Model(model.User{}).Where("open = ?", 1).Find(&users)

	for i := range users {

		var (
			online int
		)

		_, ok := WsInstance.Uids[(users[i].Uid)]

		if ok {
			online = 1
		}

		item := make(map[string]interface{})
		item["uid"] = users[i].Uid
		item["nickname"] = users[i].Nickname
		item["avatar"] = users[i].Avatar
		item["email"] = users[i].Email
		item["signature"] = users[i].Signature
		item["online"] = online

		res = append(res, item)
	}

	c.output(gin.H{
		"meta": "publicUsers",
		"data": res,
	})
}


func GetMe(c *Client, payload string)  {

	var (
		user model.User
	)

	// get uid by connection
	uid, ok := WsInstance.Clients[c]

	if !ok {
		return
	}

	DaoInstance.ChatDb.Model(model.User{}).Where("uid = ?", uid).Find(&user)

	if user.Uid <= 0 {
		return
	}

	c.output(gin.H{
		"meta": "me",
		"data": user,
	})
}

