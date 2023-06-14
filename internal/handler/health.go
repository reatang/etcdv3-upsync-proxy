package handler

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/reatang/etcdv3_upsync_proxy/pkg/health"
	"github.com/reatang/etcdv3_upsync_proxy/pkg/xnet"
	"google.golang.org/grpc/credentials/insecure"
)

//////////// Struct ////////////////////

type TargetAddr struct {
	Addr string `uri:"addr" binding:"required"`
}

//////////// Handle ////////////////////

func GrpcHealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req TargetAddr
		if ctx.ShouldBindUri(&req) != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errorCode": 400,
			})
			return
		}
		addr := strings.TrimPrefix(req.Addr, "/")
		var err error

		// 检测是否是私有IP
		host, _, err := net.SplitHostPort(addr)
		if err != nil || !xnet.IsPrivateIP(host) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errorCode": 400,
			})
			return
		}

		// 健康检查
		err = health.Check(ctx, addr, insecure.NewCredentials())
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "NOT_SERVING",
				"msg":    err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": "SERVING",
		})
	}
}
