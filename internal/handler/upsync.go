package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/reatang/etcdv3_upsync_proxy/internal/app"
	"github.com/reatang/etcdv3_upsync_proxy/pkg/upsync"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//////////// Struct ////////////////////

type V2KeysUri struct {
	Key string `uri:"key" binding:"required"`
}

type V2KeysRequest struct {
	Wait      string `form:"wait"`
	WaitIndex int64  `form:"waitIndex"`
	recursive string `form:"recursive"`
}

//////////// Handle ////////////////////

func V2KeysHandle(ctx *gin.Context) {
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

	key := strings.TrimPrefix(uri.Key, "/")

	response, err := app.EtcdClient.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}

	event := upsync.Transform(uri.Key, response)

	ctx.JSON(http.StatusOK, event)
}
