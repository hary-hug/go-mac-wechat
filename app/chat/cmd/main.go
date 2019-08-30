package main

import (
	"context"
	"go-chat/app/chat/conf"
	"go-chat/app/chat/server"
	"log"
	"os"
	"os/signal"
	"time"
)

func main()  {

	srv := server.New(conf.Cfg)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	<-c

	log.Println("Shutdown Server ...")


	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {

		log.Fatal("Server Shutdown:", err)
	}
}
