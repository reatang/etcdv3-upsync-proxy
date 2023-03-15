package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/reatang/etcdv3_upsync_proxy/internal/app"
	"github.com/reatang/etcdv3_upsync_proxy/internal/ginx"
	"github.com/reatang/etcdv3_upsync_proxy/pkg/conf"
)

var configFile = flag.String("f", "configs/config.yaml", "config file")

func main() {
	flag.Parse()

	var c app.ServerConf
	err := conf.Parse(*configFile, &c)
	if err != nil {
		log.Fatal(err)
	}

	err, cleanup := app.InitApplication(c)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Println(`sys interrupt :`, start(c))
}

func start(c app.ServerConf) error {
	g := gin.Default()
	ginx.RegisterRouter(g)

	srv := endless.NewServer(c.ListenOn, g)

	ech := make(chan error)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			ech <- err
		}
	}()

	go func() {
		_ = srv.RegisterSignalHook(endless.POST_SIGNAL, syscall.SIGINT, func() {
			ech <- fmt.Errorf("%s", syscall.SIGINT)
		})
		_ = srv.RegisterSignalHook(endless.POST_SIGNAL, syscall.SIGTERM, func() {
			ech <- fmt.Errorf("%s", syscall.SIGTERM)
		})
	}()

	return <-ech
}
