package conf

import (
	"github.com/BurntSushi/toml"
	"go-chat/app/chat/util/db"
	"log"
)

// global config val
var (
	Cfg = &Config{}
)

// The Config structure map to app.toml
type Config struct {
	Common       *Common
	HttpServer   *HttpServer
	Mysql        *Mysql
}

// Config section 'common'
type Common struct {
	JwtSecret  string
	JwtTimeout int
}

// Config section 'server'
type HttpServer struct {
	Mode  string
	Port  string
	ReadTimeout  int
	WriteTimeout int
}

// Config section 'mysql'
type Mysql struct {
	Chat  *db.MysqlConfig
}


// load config file 'app.toml'
func init()  {

	var (
		err error
	)

	_, err = toml.DecodeFile("app.toml", Cfg)

	if err != nil {
		log.Fatalln(err)
	}

}



