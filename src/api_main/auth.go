package api_main

import (
	"api_common"
	simpleMD5Auth "auth"
	"fmt"
	"io"
	. "slog"

	"github.com/emicklei/go-restful"
)

// 请求签名检查
func checkRequestAuth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	key := req.HeaderParameter("X-API-KEY")
	md5 := req.HeaderParameter("X-DATA-MD5")

	if key == "" || md5 == "" {
		Log.Debugf("key : %s, md5 : %s", key, md5)
		io.WriteString(resp, api_common.CreateErrResponse("auth", fmt.Errorf("authentication failed")).String())
		return
	}
	content := api_common.RawRequest(req)
	mapContent := api_common.MapRequest(req)
	timeStamp := mapContent.Get("timestamp")

	if timeStamp == "" {
		Log.Debugf("timestamp is not set")
		io.WriteString(resp, api_common.CreateErrResponse("auth", fmt.Errorf("authentication failed")).String())
		return
	}

	if err := simpleMD5Auth.NewAuth(key, md5, content).Check(); err != nil {
		Log.Debugf("content is %s", string(content))
		Log.Debugf("key is %s", string(key))
		Log.Debugf("md5 is %s", string(md5))
		Log.Debugf("auth failed %v", err)
		io.WriteString(resp, api_common.CreateErrResponse("auth", fmt.Errorf("authentication failed")).String())
		return
	}
	chain.ProcessFilter(req, resp)
}
