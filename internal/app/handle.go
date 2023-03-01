package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reatang/etcdv3_upsync_proxy/pkg/proxy"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	authPrefix     = "/v2/auth"
	keysPrefix     = "/v2/keys"
	machinesPrefix = "/v2/machines"
	membersPrefix  = "/v2/members"
	statsPrefix    = "/v2/stats"
)

func RegisterRouter(gin *gin.Engine) {
	gin.GET(keysPrefix+"/*key", v2KeysHandle)
}

//////////// Handle ////////////////////

type V2KeysUri struct {
	Key string `uri:"key" binding:"required"`
}

type V2KeysRequest struct {
	Wait      string `form:"wait"`
	WaitIndex int64  `form:"waitIndex"`
	recursive string `form:"recursive"`
}

func v2KeysHandle(ctx *gin.Context) {
	var uri V2KeysUri
	if ctx.ShouldBindUri(&uri) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": 400,
		})
		return
	}
	var query V2KeysRequest
	if ctx.ShouldBindQuery(&query) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorCode": 400,
		})
		return
	}

	response, err := EtcdClientCli.Get(ctx, uri.Key, clientv3.WithPrefix())
	if err != nil {
		return
	}

	event := proxy.Transform(uri.Key, response)

	ctx.JSON(http.StatusOK, event)
}
