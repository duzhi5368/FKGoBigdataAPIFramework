package api_simple

import (
	. "api_common"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"

	. "db_engine"
	"fmt"
	"github.com/go-xorm/builder"
	"reflect"
	. "slog"
	"strconv"
)

type mysql_test_parameter struct {
	MysqlHost        string `json:"mysql_host"`
	MysqlPort        int    `json:"mysql_port"`
	MysqlDBName      string `json:"mysql_db_name"`
	MysqlUser        string `json:"mysql_user"`
	MysqlPassword    string `json:"mysql_password"`
	RangeKey         string `json:"range_key"`
	RangeValueList   string `json:"range_value_list"`
	ConditionalKey   string `json:"conditional_key"`
	ConditionalValue string `json:"conditional_value"`
	EchoInfos        string `json:"echo_infos"`
}

type mysql_data_scheme struct {
	TestValue1 string `json:"test_value_1"`
	TestValue2 int    `json:"test_value_2"`
}

func (p mysql_test_parameter) DumpInfo() string {
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

func (p mysql_test_parameter) SafeCheck() error {
	return nil
}

func (p *mysql_test_response) FillBaseByRequest(request *mysql_test_parameter) {
	p.EchoInfos = request.EchoInfos
}

type mysql_test_response struct {
	EchoInfos  string   `json:"echo_infos"`
	TestValue1 []string `json:"test_value_1"`
	TestValue2 []int    `json:"test_value_2"`
}

func OnMySQLTestHandler(req *restful.Request, resp *restful.Response) *ResponceStruct {
	body := req.Request.Body
	defer body.Close()

	content, err := ioutil.ReadAll(body)
	if err != nil {
		return CreateErrResponse("OnMySQLTestHandler", err)
	}

	var reqParameter mysql_test_parameter
	err = json.Unmarshal(content, &reqParameter)
	if err != nil {
		return CreateErrResponse("OnMySQLTestHandler", err)
	}

	err = reqParameter.SafeCheck()
	if err != nil {
		return CreateErrResponse("OnMySQLTestHandler", err)
	}
	Log.Println(reqParameter.DumpInfo())

	var response mysql_test_response
	response.FillBaseByRequest(&reqParameter)

	// BEGIN

	// 构建Mysql引擎
	resMysqlHost := reqParameter.MysqlHost
	resMysqlPort := reqParameter.MysqlPort
	resMysqlDBName := reqParameter.MysqlDBName
	resMysqlUser := reqParameter.MysqlUser
	resMysqlPassword := reqParameter.MysqlPassword
	mysqlHostParameters := MysqlHostParameters{resMysqlHost, resMysqlPort,
		resMysqlDBName, resMysqlUser, resMysqlPassword}
	db, err := MysqlEngine(mysqlHostParameters)
	if err != nil {
		Log.Println("mysql engine error: ", err)
		return CreateErrResponse("OnMySQLTestHandler", err)
	}
	defer db.Close()

	// 构建条件并执行
	result := make([]mysql_data_scheme, 0)
	err = db.Table(new(mysql_data_scheme)).Where(builder.In(reqParameter.RangeKey, reqParameter.RangeValueList).
		And(builder.Eq{reqParameter.ConditionalKey: reqParameter.ConditionalValue})).Find(&result)
	if err != nil {
		Log.Println("mysql engine error: ", err)
		return CreateErrResponse("OnMySQLTestHandler", err)
	}

	// 循环解析数据
	for _, v := range result {
		response.TestValue1 = append(response.TestValue1, v.TestValue1)
		response.TestValue2 = append(response.TestValue2, v.TestValue2)
	}

	// END
	return CreateSuccessResponse(response)
}
