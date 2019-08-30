package model

type Chat struct {
	ChatId     int    `gorm:"column:chat_id;primary_key" json:"chat_id"`
	ChatType   int    `gorm:"column:chat_type" json:"chat_type"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
}


type ChatList struct {
	ListId     int    `gorm:"column:list_id;primary_key" json:"list_id"`
	ChatId     int    `gorm:"column:chat_id" json:"chat_id"`
	ChatType   int    `gorm:"column:chat_type" json:"chat_type"`
	Uid        int    `gorm:"column:uid" json:"uid"`
	RelateId   int    `gorm:"column:relate_id" json:"relate_id"`
	New        int    `gorm:"column:new" json:"new"`
	ActiveTime int    `gorm:"column:active_time" json:"active_time"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
}


type ChatMessage struct {
	MsgId      int    `gorm:"column:msg_id;primary_key" json:"msg_id"`
	ChatId     int    `gorm:"column:chat_id" json:"chat_id"`
	Uid        int    `gorm:"column:uid" json:"uid"`
	MsgType    int    `gorm:"column:msg_type" json:"msg_type"`
	Content    string `gorm:"column:content" json:"content"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
}


type ChatUser struct {
	ChatId     int    `gorm:"column:chat_id" json:"chat_id"`
	Uid        int    `gorm:"column:uid" json:"uid"`
}


func (Chat) TableName() string  {
	return "chat"
}


func (ChatList) TableName() string  {
	return "chat_list"
}


func (ChatMessage) TableName() string  {
	return "chat_message"
}


func (ChatUser) TableName() string  {
	return "chat_user"
}