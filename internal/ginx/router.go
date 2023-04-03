package ginx

import "github.com/gin-gonic/gin"

func RegisterRouter(gin *gin.Engine) {
	// etcd v2 -> v3
	gin.GET(keysPrefix+"/*key", v2KeysHandle)

	// grpc health check
	gin.GET("/health-check/*addr", grpcHealthCheck)
}
