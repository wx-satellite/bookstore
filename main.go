package main

import (
	"context"
	"github.com/wx-satellite/bookstore/server"
	"github.com/wx-satellite/bookstore/store"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s, err := store.New("mem")
	if err != nil {
		panic(err)
	}

	srv := server.New(":8080", s)

	errChan, err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}

	log.Println("web server start ok")
	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-errChan:
		log.Println("web server run failed：", err)
		return
	case <-signalChan:
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		err = srv.Shutdown(ctx)
	}

	if err != nil {
		log.Println("bookstore program exist err：", err)
		return
	}

	log.Println("bookstore program exist ok")
}
