package main

import (
	"fmt"
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
		return
	}
	defer cleanup()

	g := gin.Default()
	app.RegisterRouter(g)

	// 框架错误捕获
	errChan := make(chan error)
	go func() {
		err := endless.ListenAndServe(":2381", g)
		if err != nil {
			errChan <- err
		}
	}()

	// 系统信号捕获
	go func() {
		sign := make(chan os.Signal)
		signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sign)
	}()

	// 输出错误消息
	e := <-errChan
	fmt.Println(`sys interrupt :`, e)
}
