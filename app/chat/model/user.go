package model

type User struct {
	Uid        int        `gorm:"column:uid;primary_key" json:"uid"`
	Email      string     `gorm:"column:email" json:"email"`
	Password   string     `gorm:"column:password" json:"password"`
	Salt       string     `gorm:"column:salt" json:"salt"`
	Nickname   string     `gorm:"column:nickname" json:"nickname"`
	Avatar     string     `gorm:"column:avatar" json:"avatar"`
	Signature  string     `gorm:"column:signature" json:"signature"`
	CreateTime int        `gorm:"column:create_time" json:"create_time"`
	CreateIp   string     `gorm:"column:create_ip" json:"create_ip"`
}


type UserFriend struct {
	Uid        int        `gorm:"column:uid" json:"uid"`
	FriendId   int        `gorm:"column:friend_id" json:"friend_id"`
	CreateTime int        `gorm:"column:create_time" json:"create_time"`
}


type UserFriendRequest struct {
	ReqId      int        `gorm:"column:req_id;primary_key" json:"req_id"`
	Uid        int        `gorm:"column:uid" json:"uid"`
	FromUid    int        `gorm:"column:from_uid" json:"from_uid"`
	Status     int        `gorm:"column:status" json:"status"`
	CreateTime int        `gorm:"column:create_time" json:"create_time"`
}


func (User) TableName() string {
	return "user"
}


func (UserFriend) TableName() string {
	return "user_friend"
}


func (UserFriendRequest) TableName() string {
	return "user_friend_request"
}