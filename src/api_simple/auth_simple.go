package api_simple

import (
	. "api_common"
	"github.com/emicklei/go-restful"
	"io/ioutil"
)

type auth_test_response struct {
}

func OnAuthTestHandler(req *restful.Request, resp *restful.Response) *ResponceStruct {
	body := req.Request.Body
	defer body.Close()

	_, err := ioutil.ReadAll(body)
	if err != nil {
		return CreateErrResponse("OnAuthTestHandler", err)
	}

	var response auth_test_response
	return CreateSuccessResponse(response)
}
