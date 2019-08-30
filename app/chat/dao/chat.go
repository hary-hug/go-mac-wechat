package dao

import (
	"go-chat/app/chat/model"
	"time"
)

const (
	// 私聊
	_privateType = 1
	// 群聊
	_groupType   = 2
)

// AddFriend make a new friend data and chat record
func (d *Dao) AddFriend(uid int, friendId int) (chatId int, err error) {

	var (
		chat model.Chat
	)

	// begin transaction
	ts := d.ChatDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			ts.Rollback()
		}
	}()

	now := int(time.Now().Unix())
	// create friendship
	if err = ts.Create(&model.UserFriend{
		Uid: uid,
		FriendId: friendId,
		CreateTime: now,
	}).Error; err != nil {
		return
	}
	if err = ts.Create(&model.UserFriend{
		Uid: friendId,
		FriendId: uid,
		CreateTime: now,
	}).Error; err != nil {
		return
	}

	// create a new chat
	chat.ChatType = _privateType
	chat.CreateTime = now
	if err = ts.Create(&chat).Error; err != nil {
		return
	}

	// create chat list of me
	if err = ts.Create(&model.ChatList{
		ChatId:     chat.ChatId,
		Uid:        uid,
		ChatType:   _privateType,
		RelateId:   friendId,
		CreateTime: now,
		ActiveTime: now,
	}).Error; err != nil {
		return
	}

	// create chat list of my friend
	if err = ts.Create(&model.ChatList{
		ChatId:     chat.ChatId,
		Uid:        friendId,
		ChatType:   _privateType,
		RelateId:   uid,
		CreateTime: now,
		ActiveTime: now,
	}).Error; err != nil {
		return
	}

	// create chat user
	if err = ts.Create(&model.ChatUser{
		ChatId: chat.ChatId,
		Uid:    uid,
	}).Error; err != nil {
		return
	}

	if err = ts.Create(&model.ChatUser{
		ChatId: chat.ChatId,
		Uid:    friendId,
	}).Error; err != nil {
		return
	}

	if err = ts.Commit().Error; err != nil {
		return
	}

	return chat.ChatId, nil
}