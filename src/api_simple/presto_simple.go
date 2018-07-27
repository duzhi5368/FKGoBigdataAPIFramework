package api_simple

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"reflect"
	"strconv"

	. "api_common"
	. "db_engine"
	"encoding/json"
	"io/ioutil"
	. "slog"
)

type presto_test_parameter struct {
	IsUsePrestoDBTLS bool   `json:"is_use_presto_dbtls"`
	PrestoDBHost     string `json:"presto_db_host"`
	PrestoDBPort     int    `json:"presto_db_port"`
	PrestoDBUer      string `json:"presto_db_uer"`
	PrestoDBPassword string `json:"presto_db_password"`
	PrestoDBCatalog  string `json:"presto_db_catalog"`
	PrestoDBSchema   string `json:"presto_db_schema"`
	PrestoDBTable    string `json:"presto_db_table"`
	SelectRowName    string `json:"select_row_name"`
	RangeKey         string `json:"range_key"`
	RangeMin         int    `json:"range_min"`
	RangeMax         int    `json:"range_max"`
	ConditionalKey   string `json:"conditional_key"`
	ConditionalValue string `json:"conditional_value"`
	SortKey          string `json:"sort_key"`
	SumKey           string `json:"sum_key"`
	AggregationKey   string `json:"aggregation_key"`
	EchoInfos        string `json:"echo_infos"`
}

type presto_date_scheme struct {
	TestValue1 string
	TestValue2 int
}

func (p presto_test_parameter) DumpInfo() string {
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

func (p presto_test_parameter) SafeCheck() error {
	return nil
}

func (p *presto_test_response) FillBaseByRequest(request *presto_test_parameter) {
	p.EchoInfos = request.EchoInfos
}

type presto_test_response struct {
	EchoInfos  string   `json:"echo_infos"`
	TestValue1 []string `json:"test_value_1"`
	TestValue2 []int    `json:"test_value_2"`
}

func OnPrestoTestHandler(req *restful.Request, resp *restful.Response) *ResponceStruct {
	body := req.Request.Body
	defer body.Close()

	content, err := ioutil.ReadAll(body)
	if err != nil {
		return CreateErrResponse("OnPrestoTestHandler", err)
	}

	var reqParameter presto_test_parameter
	err = json.Unmarshal(content, &reqParameter)
	if err != nil {
		return CreateErrResponse("OnPrestoTestHandler", err)
	}

	err = reqParameter.SafeCheck()
	if err != nil {
		return CreateErrResponse("OnPrestoTestHandler", err)
	}
	Log.Println(reqParameter.DumpInfo())

	var response presto_test_response
	response.FillBaseByRequest(&reqParameter)

	// BEGIN

	// 构建Presto引擎
	IsUsePrestoDBTLS := reqParameter.IsUsePrestoDBTLS
	PrestoDBHost := reqParameter.PrestoDBHost
	PrestoDBPort := reqParameter.PrestoDBPort
	PrestoDBUer := reqParameter.PrestoDBUer
	PrestoDBPassword := reqParameter.PrestoDBPassword
	PrestoDBCatalog := reqParameter.PrestoDBCatalog
	PrestoDBSchema := reqParameter.PrestoDBSchema

	prestoHostParameters := PrestoHostParameters{
		IsUsePrestoDBTLS,
		PrestoDBHost,
		PrestoDBPort,
		PrestoDBUer,
		PrestoDBPassword,
		PrestoDBCatalog,
		PrestoDBSchema}
	db, err := PrestoEngine(prestoHostParameters)
	if err != nil {
		Log.Println("Presto error: ", err)
		return CreateErrResponse("OnPrestoTestHandler", err)
	}
	defer db.Close()

	// 组建SQL语句
	var selectSQL = ""

	selectSQL += "SELECT"
	selectSQL += " " + reqParameter.SelectRowName + ","
	selectSQL += " SUM(" + reqParameter.SumKey + ")"
	selectSQL += " FROM " + reqParameter.PrestoDBSchema + "." + reqParameter.PrestoDBTable
	selectSQL += " WHERE " + reqParameter.RangeKey + " BETWEEN " +
		strconv.Itoa(reqParameter.RangeMin) + " AND " + strconv.Itoa(reqParameter.RangeMax)
	if reqParameter.ConditionalKey != "" {
		selectSQL += " AND " + reqParameter.ConditionalKey + " = '" + reqParameter.ConditionalValue + "'"
	}
	selectSQL += " GROUP BY " + reqParameter.AggregationKey

	Log.Println("[Presto SQL] " + selectSQL)

	rows, err := db.Query(selectSQL)
	if err != nil {
		Log.Println(err)
		return CreateErrResponse("OnPrestoTestHandler", err)
	}

	// 循环解析数据
	for rows.Next() {
		line := presto_date_scheme{}
		err = rows.Scan(&line.TestValue1, &line.TestValue2)

		// 填充响应数据
		response.TestValue1 = append(response.TestValue1, line.TestValue1)
		response.TestValue2 = append(response.TestValue2, line.TestValue2)
	}
	err = rows.Err()
	if err != nil {
		Log.Println(err)
		return CreateErrResponse("OnPrestoTestHandler", err)
	}

	// END

	return CreateSuccessResponse(response)
}
