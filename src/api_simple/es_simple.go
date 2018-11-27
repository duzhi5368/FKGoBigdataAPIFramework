package api_simple

import (
	. "api_common"
	. "slog"

	"fmt"

	"config"
	"context"
	. "db_engine"
	"encoding/json"
	"errors"
	"github.com/emicklei/go-restful"
	"gopkg.in/olivere/elastic.v5"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
)

type es_test_parameter struct {
	ESHost           string `json:"es_host"`
	ESPort           int    `json:"es_port"`
	IndexName        string `json:"index_name"`
	RangeKey         string `json:"range_key"`
	RangeMin         int    `json:"range_min"`
	RangeMax         int    `json:"range_max"`
	ConditionalKey   string `json:"conditional_key"`
	ConditionalValue string `json:"conditional_value"`
	SortKey          string `json:"sort_key"`
	EchoInfos        string `json:"echo_infos"`
}

type es_date_scheme struct {
	TestValue string `json:"test_value"`
}

func (p es_test_parameter) DumpInfo() string {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	var result = "\n==========================================\n"
	result += "Struct : 【" + t.Name() + "】\n"
	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		valueField := v.Field(i)
		var value = ""
		switch valueField.Kind() {
		case reflect.Invalid:
			value = "invalid"
		case reflect.String:
			value = "\"" + valueField.String() + "\""
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.Itoa(int(valueField.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			value = strconv.FormatUint(v.Uint(), 10)
		case reflect.Bool:
			value = strconv.FormatBool(bool(valueField.Bool()))
		case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
			value = v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
		default: // reflect.Array, reflect.Struct, reflect.Interface
			value = v.Type().String() + " value"
		}
		result += fmt.Sprintf("%d. %v (%v) = %s \n", i+1, typeField.Name, typeField.Type.Name(), value)
	}
	result += "=========================================="
	return result
}

func (p es_test_parameter) SafeCheck() error {
	return nil
}

func (p *es_test_response) FillBaseByRequest(request *es_test_parameter) {
	p.EchoInfos = request.EchoInfos
}

type es_test_response struct {
	EchoInfos string   `json:"echo_infos"`
	TestValue []string `json:"test_value"`
}

func OnESTestHandler(req *restful.Request, resp *restful.Response) *ResponceStruct {
	body := req.Request.Body
	defer body.Close()

	content, err := ioutil.ReadAll(body)
	if err != nil {
		return CreateErrResponse("OnESTestHandler", err)
	}

	var reqParameter es_test_parameter
	err = json.Unmarshal(content, &reqParameter)
	if err != nil {
		return CreateErrResponse("OnESTestHandler", err)
	}

	err = reqParameter.SafeCheck()
	if err != nil {
		return CreateErrResponse("OnESTestHandler", err)
	}
	Log.Println(reqParameter.DumpInfo())

	var response es_test_response
	response.FillBaseByRequest(&reqParameter)

	// BEGIN

	// 构建ES引擎
	reqESHost := reqParameter.ESHost
	reqESPort := reqParameter.ESPort
	esHostParameters := ESHostParameters{reqESHost, reqESPort}
	client, err := ESEngine(esHostParameters)
	if err != nil {
		Log.Println("es engine error: ", err)
		return CreateErrResponse("OnESTestHandler", err)
	}

	// 检查ESIndex是否存在
	reqESIndex := config.Config.ESIndex
	exists, err := client.IndexExists(reqESIndex).Do(context.Background())
	if err != nil {
		Log.Println("connect es engine error: ", err)
		return CreateErrResponse("OnESTestHandler", err)
	}
	if !exists {
		Log.Println("es index doesn't exist:", reqESIndex)
		return CreateErrResponse("OnESTestHandler",
			errors.New(fmt.Sprintf("index doesn't exist: %s", reqESIndex)))
	}

	// 构建Must条件
	// 区域限制条件
	mustQuery := make([]elastic.Query, 0)
	mustQuery = append(mustQuery, elastic.NewRangeQuery(reqParameter.RangeKey).
		Gte(reqParameter.RangeMin).
		Lt(reqParameter.RangeMax))
	// 单值匹配条件
	if reqParameter.ConditionalKey != "" {
		mustQuery = append(mustQuery, elastic.NewTermQuery(reqParameter.ConditionalKey, reqParameter.ConditionalValue))
	}
	termQuery := elastic.NewBoolQuery().Must(mustQuery...)

	var svc *elastic.ScrollService
	svc = client.Scroll().
		Index(reqESIndex).
		Query(termQuery).
		Pretty(true).
		Sort(reqParameter.SortKey, true).
		Size(5000) // 单次取出5000条
	successedCount := 0
	failedCount := 0

	// 循环解析数据
	for {
		res, err := svc.Do(context.Background())
		if err == io.EOF {
			break
		}
		var esIndex es_date_scheme
		for _, item := range res.Each(reflect.TypeOf(esIndex)) {
			data, ok := item.(es_date_scheme)
			if !ok {
				failedCount++
				Log.Println("expected hit to be serialized as es_date_scheme; got: ", reflect.ValueOf(item))
			} else {
				successedCount++
				response.TestValue = append(response.TestValue, data.TestValue) // 填充返回结果
			}
		}
	}

	Log.Println("ES hit = ", successedCount+failedCount, " successed = ", successedCount, " failed = ", failedCount)
	client.close()
	// END

	return CreateSuccessResponse(response)
}
