package api_common

import "encoding/json"

type ResponceStruct struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

func (r *ResponceStruct) String() string {
	ret, _ := json.MarshalIndent(*r, "", " ")
	return string(ret)
}

func CreateErrResponse(module string, err error) *ResponceStruct {
	return &ResponceStruct{Code: -1, Msg: module + ": " + err.Error()}
}

func CreateSuccessResponse(v interface{}) *ResponceStruct {
	return &ResponceStruct{Code: 0, Msg: "success", Data: v}
}
