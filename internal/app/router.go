package app

import (
	"github.com/gin-gonic/gin"
	"github.com/reatang/etcdv3_upsync_proxy/internal/handler"
)

const (
	authPrefix     = "/v2/auth"
	keysPrefix     = "/v2/keys"
	machinesPrefix = "/v2/machines"
	membersPrefix  = "/v2/members"
	statsPrefix    = "/v2/stats"
)

func RegisterRouter(gin *gin.Engine) {
	// etcd v2 -> v3
	gin.GET(keysPrefix+"/*key", handler.V2KeysHandle(EtcdClient))

	// grpc health check
	gin.GET("/health-check/*addr", handler.GrpcHealthCheck())
}
