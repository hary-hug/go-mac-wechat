package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-chat/app/chat/conf"
	"go-chat/app/chat/util/db"
)

// sometimes we need to connect two or more database
// so the Dao struct may include some dao such as redis, mc, sqlite
// or any other datasource's connection
// use these dao to fetch data
type Dao struct {
	ChatDb  *gorm.DB
}

// New init dao connection
func New(cfg *conf.Config) (d *Dao) {

	d = &Dao{
		ChatDb: db.NewMysql(cfg.Mysql.Chat),
	}

	return
}

// close connections
func (d *Dao) Close()  {
	d.ChatDb.Close()
}
