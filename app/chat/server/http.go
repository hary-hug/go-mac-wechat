package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-chat/app/chat/conf"
	"go-chat/app/chat/dao"
	"log"
	"net/http"
	"time"
)

var (
	DaoInstance  *dao.Dao
	WsInstance   *Ws
)


// init a http server
func New(cfg *conf.Config) (srv *http.Server) {

	// init gin instance
	r := gin.New()

	// set logger
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// cors
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("x-token")
	r.Use(cors.New(config))

	// setup url request router
	setupRoutes(r)

	// set http run mode
	gin.SetMode(cfg.HttpServer.Mode)

	// set http server
	srv = &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.HttpServer.Port),
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.HttpServer.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.HttpServer.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// dao init
	DaoInstance = dao.New(cfg)
	// set websocket instance
	WsInstance  = newWs()

	return

}

// setupRoutes set up http request routes
func setupRoutes(e *gin.Engine) {

	g1 := e.Group("/api")
	{
		g1.POST("/login", Login)
		g1.POST("/register", Register)
	}

	// websocket
	e.GET("/ws", func(c *gin.Context) {

		err := WsInstance.serve(c)
		if err != nil {
			log.Println(err)
		}
	})
}


