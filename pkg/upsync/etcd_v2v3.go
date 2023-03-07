package upsync

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/reatang/etcdv3_upsync_proxy/pkg/xetcd/v2store"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var ErrKeyParseError = errors.New("key format error")
var ErrValueParseError = errors.New("value format error")

type (
	UpstreamParams struct {
		Weight      *int64 `json:"weight,omitempty"`
		MaxFails    *int64 `json:"max_fails,omitempty"`
		FailTimeout *int64 `json:"fail_timeout,omitempty"`

		// 上线状态 0,1
		Down *int8 `json:"down,omitempty"`

		// 后备状态 0,1
		Backup *int8 `json:"backup,omitempty"`
	}
)

func Transform(key string, response *clientv3.GetResponse) *v2store.Event {
	nodes := make(v2store.NodeExterns, 0)
	for _, kv := range response.Kvs {

		newKey, newVal, err := formatKey(kv.Key, kv.Value)
		if err != nil {
			log.Println(err, kv)
			continue
		}

		nodes = append(nodes, &v2store.NodeExtern{
			Key:   newKey,
			Value: &newVal,
		})
	}

	event := v2store.NewEvent(v2store.Get, key, uint64(response.Header.Revision), 0)
	event.Node.Nodes = nodes

	return event
}

// kv种类
//
//	  1、upsync 标准
//			key：/upstream/somerpc/<ip>:<port>
//	     val：
//	  2、go-zero风格
//			key：/some_path/somerpc/<lease>
//	     val：<ip>:<port>
//	  3、url 风格
//			key：/some_path/somerpc/<lease>
//	     val：//<ip>:<port>/?weight=1&max_fails=2&fail_timeout=10
func formatKey(key, val []byte) (newKey, newVal string, err error) {
	_key := string(key)
	_v := string(val)

	index := strings.LastIndex(_key, "/")
	keyVal := _key[index+1:]

	// 检测 upsync 标准的
	if strings.Index(keyVal, ":") > 0 {
		newKey = _key
		return
	}

	keyName := _key[:index]

	// 检测是否是url模式
	vUri, err := url.Parse(_v)
	if err == nil {
		// 检测url
		if vUri.Host == "" {
			return "", "", ErrValueParseError
		}

		newKey = fmt.Sprintf("%s/%s", keyName, vUri.Host)

		if vUri.RawQuery != "" {
			j := &UpstreamParams{}
			q := vUri.Query()
			j.Weight = str2pint(q.Get("weight"))
			j.MaxFails = str2pint(q.Get("max_fails"))
			j.FailTimeout = str2pint(q.Get("fail_timeout"))

			vb, _ := json.Marshal(j)
			newVal = string(vb)
		}

		return
	} else {
		// 检测是否是合法端口
		var p string
		if _, p, err = net.SplitHostPort(_v); err != nil {
			return "", "", err
		}
		_, err = strconv.ParseInt(p, 10, 32)
		if err != nil {
			return "", "", err
		}

		//只存在ip:port
		newKey = fmt.Sprintf("%s/%s", keyName, _v)
		err = nil
	}

	return
}

func str2pint(s string) *int64 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil
	}

	return &i
}
