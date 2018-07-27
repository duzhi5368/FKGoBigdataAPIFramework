package api_simple

import (
	. "api_common"
	. "slog"

	"fmt"

	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"reflect"
	"strconv"
)

type echo_test_parameter struct {
	Msg string `json:"msg"`
}

func (p echo_test_parameter) DumpInfo() string {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	var result = "Struct : " + t.Name() + "\n-------------------- \n"
	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		valueField := v.Field(i)
		var value = ""
		switch valueField.Kind() {
		case reflect.String:
			value = valueField.String()
		case reflect.Int:
			value = strconv.Itoa(int(valueField.Int()))
		}
		result += fmt.Sprintf("%d. %v (%v) = %s \n", i+1, typeField.Name, typeField.Type.Name(), value)
	}
	return result
}

func (p echo_test_parameter) SafeCheck() error {
	if p.Msg == "" {
		return fmt.Errorf("check param failed: YOU MUST SEND SOMETING.")
	}
	return nil
}

type echo_test_response struct {
	Msg string `json:"msg"`
}

func (p *echo_test_response) FillBaseByRequest(request *echo_test_parameter) {
	p.Msg = request.Msg
}

func OnEchoTestHandler(req *restful.Request, resp *restful.Response) *ResponceStruct {
	body := req.Request.Body
	defer body.Close()

	content, err := ioutil.ReadAll(body)
	if err != nil {
		return CreateErrResponse("OnEchoTestHandler", err)
	}

	var reqParameter echo_test_parameter
	err = json.Unmarshal(content, &reqParameter)
	if err != nil {
		return CreateErrResponse("OnEchoTestHandler", err)
	}

	err = reqParameter.SafeCheck()
	if err != nil {
		return CreateErrResponse("OnEchoTestHandler", err)
	}
	Log.Println(reqParameter.DumpInfo())

	var response echo_test_response
	response.FillBaseByRequest(&reqParameter)

	return CreateSuccessResponse(response)
}
