package db

import (
	"github.com/jinzhu/gorm"
	"log"
)

type MysqlConfig struct {
	Dns         string
	LogMode     bool
	TablePrefix string
	MaxIdle     int
	MaxOpen     int
}


func NewMysql(cfg *MysqlConfig) (db *gorm.DB) {

	var (
		err error
	)

	db, err = gorm.Open("mysql", cfg.Dns)

	if err != nil {
		log.Fatalln(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName;
	}

	db.LogMode(cfg.LogMode)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(cfg.MaxIdle)
	db.DB().SetMaxOpenConns(cfg.MaxOpen)

	return
}
