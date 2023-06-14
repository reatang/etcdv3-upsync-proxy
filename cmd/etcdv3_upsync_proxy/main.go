package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/reatang/etcdv3-upsync-proxy/internal/app"
	"github.com/reatang/etcdv3-upsync-proxy/pkg/conf"
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
	app.RegisterRouter(g)

	srv := endless.NewServer(c.ListenOn, g)

	ech := make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			!strings.Contains(err.Error(), "use of closed network connection") &&
			err != http.ErrServerClosed {
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

	err := <-ech

	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	go func() {
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("server shutdown err:", err)
		}
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("server shutdown timeout:", ctx.Err())
		}
	}

	return err
}
