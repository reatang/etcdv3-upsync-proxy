package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/reatang/etcdv3_upsync_proxy/internal/app"
)

func main() {
	err, cleanup := app.InitApplication()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Println(`sys interrupt :`, start())
}

func start() error {
	g := gin.Default()
	app.RegisterRouter(g)

	ech := make(chan error)
	go func() {
		err := endless.ListenAndServe(":2381", g)
		if err != nil {
			ech <- err
		}
	}()

	go func() {
		sign := make(chan os.Signal)
		signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
		ech <- fmt.Errorf("%s", <-sign)
	}()

	return <-ech
}
