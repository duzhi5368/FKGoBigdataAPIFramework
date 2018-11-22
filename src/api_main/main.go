package api_main

import (
	"api_common"
	"config"
	"fmt"
	"io"
	"net/http"
	. "slog"
	"time"
	"utils"

	"auth_storage"
	"cmd_line"
	"github.com/emicklei/go-restful"
)

type apiStruct struct {
	isAuth     bool
	apiHandler func(request *restful.Request, response *restful.Response) *api_common.ResponceStruct
	desc       string
}

func apiHandler(f func(request *restful.Request, response *restful.Response) *api_common.ResponceStruct, k string) func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		then := time.Now()
		defer func() {
			Log.Infof("%s handler cost %.3fs", k, time.Since(then).Seconds())
		}()
		res := f(request, response)
		io.WriteString(response, res.String())
	}
}

// 注册路由接口
func registerAllRoute(ws *restful.WebService) error {
	ws.Filter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
		Log.Println(request.Request.RemoteAddr, request.Request.RequestURI)
		chain.ProcessFilter(request, response)
	})

	for k, v := range mapAPIHandler {
		Log.Println(utils.GetFunctionName(v.apiHandler), "\t- ", v.desc, "\t["+utils.If(v.isAuth, "需要", "无需").(string)+"签名]")
		apiFunc := apiHandler(v.apiHandler, k)
		if v.isAuth {
			ws.Route(
				ws.POST(k).
					Filter(checkRequestAuth).
					Consumes("application/json").
					To(apiFunc).
					Produces(restful.MIME_JSON))
		} else {
			ws.Route(
				ws.POST(k).
					Consumes("application/json").
					To(apiFunc).
					Produces(restful.MIME_JSON))
		}
	}
	return nil
}

// 启动服务器路由
func StartServe(commandLine cmd_line.SAppCommandLine) error {
	config.InitConfig(commandLine.ConfigFilePath)
	auth_storage.InitKey(commandLine.KeyFilePath)
	
	restful.Filter(OPTIONSFilter)
	restful.Filter(EnableCORSFilter)
	
	ws := new(restful.WebService)
	err := registerAllRoute(ws)
	if err != nil {
		return err
	}

	restful.Add(ws)
	Log.Println("Listen port: ", config.Config.ListenPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.ListenPort), nil)
}

// 开启CORS返回
func EnableCORSFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if origin := req.Request.Header.Get("Origin"); origin != "" {
		resp.AddHeader("Access-Control-Allow-Origin", origin)
	}
	chain.ProcessFilter(req, resp)
}

// 针对CORS的非简单请求做的额外处理
func OPTIONSFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if "OPTIONS" != req.Request.Method {
		chain.ProcessFilter(req, resp)
		return
	}

	archs := req.Request.Header.Get(restful.HEADER_AccessControlRequestHeaders)
	methods := "GET, POST"
	origin := "*"

	resp.AddHeader(restful.HEADER_Allow, methods)
	resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
	resp.AddHeader(restful.HEADER_AccessControlAllowHeaders, archs)
	resp.AddHeader(restful.HEADER_AccessControlAllowMethods, methods)
}
