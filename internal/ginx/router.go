package ginx

import "github.com/gin-gonic/gin"

func RegisterRouter(gin *gin.Engine) {
	gin.GET(keysPrefix+"/*key", v2KeysHandle)
}
