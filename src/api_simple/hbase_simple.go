package api_simple

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"reflect"
	"strconv"

	. "api_common"
	. "db_engine"
	. "slog"
	"strings"
)

type hbase_test_parameter struct {
	HBaseHost      string   `json:"hbase_host"`
	HBasePort      int      `json:"hbase_port"`
	HBaseTable     string   `json:"hbase_table"`
	IndexValueList []string `json:"index_value_list"` // 查询条件
	ColumnName1    string   `json:"column_name_1"`
	ColumnName2    string   `json:"column_name_2"`
}

func (p hbase_test_parameter) DumpInfo() string {
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

func (p hbase_test_parameter) SafeCheck() error {
	return nil
}

func (p *hbase_test_response) FillBaseByRequest(request *hbase_test_parameter) {
	p.ColumnName1 = request.ColumnName1
	p.ColumnName2 = request.ColumnName2
}

type hbase_test_response struct {
	ColumnName1 string   `json:"column_name_1"`
	ValueList1  []string `json:"value_list_1"`
	ColumnName2 string   `json:"column_name_2"`
	ValueList2  []string `json:"value_list_2"`
}

type hbase_date_scheme struct {
	ColumnValue1 string `json:"column_value_1"`
	ColumnValue2 string `json:"column_value_2"`
}

func OnHBaseTestHandler(req *restful.Request, resp *restful.Response) *ResponceStruct {
	body := req.Request.Body
	defer body.Close()

	content, err := ioutil.ReadAll(body)
	if err != nil {
		return CreateErrResponse("OnHBaseTestHandler", err)
	}

	var reqParameter hbase_test_parameter
	err = json.Unmarshal(content, &reqParameter)
	if err != nil {
		return CreateErrResponse("OnHBaseTestHandler", err)
	}

	err = reqParameter.SafeCheck()
	if err != nil {
		return CreateErrResponse("OnHBaseTestHandler", err)
	}
	Log.Println(reqParameter.DumpInfo())

	var response hbase_test_response
	response.FillBaseByRequest(&reqParameter)

	// BEGIN

	// 构建ES引擎
	reqHBaseHost := reqParameter.HBaseHost
	reqHBasePort := reqParameter.HBasePort
	hbaseHostParameters := HBaseHostParameters{reqHBaseHost, reqHBasePort}
	client, err := HBaseEngine(hbaseHostParameters)
	if err != nil {
		Log.Println("HBase engine error: ", err)
		return CreateErrResponse("OnHBaseTestHandler", err)
	}

	if err = client.Open(); err != nil {
		Log.Println("HBase engine error: ", err)
		return CreateErrResponse("OnHBaseTestHandler", err)
	}
	defer client.Close()

	// 检查HBaseTable
	reqHBaseTableName := reqParameter.HBaseTable
	TableNameArr, err := client.GetTableNames()
	if err != nil {
		Log.Println("HBase engine error: ", err)
		return CreateErrResponse("OnHBaseTestHandler", err)
	}

	var findTable = false
	for _, tableName := range TableNameArr {
		if tableName == reqHBaseTableName {
			findTable = true
			break
		}
	}
	if !findTable {
		Log.Println("HBase engine error: ", err)
		return CreateErrResponse("OnHBaseTestHandler",
			errors.New(fmt.Sprintf("HBase table doesn't exist : %s", reqHBaseTableName)))
	}

	// 构建条件
	byteHBaseRowkeyArray := make([][]byte, 0)
	for _, row := range reqParameter.IndexValueList {
		byteHBaseRowkeyArray = append(byteHBaseRowkeyArray, []byte(row))
	}
	columusKey := make([]string, 0)
	columusKey = append(columusKey, reqParameter.ColumnName1)
	columusKey = append(columusKey, reqParameter.ColumnName2)
	HBaseSearchResult, err := client.GetRowsWithColumns(reqHBaseTableName,
		byteHBaseRowkeyArray, columusKey, nil)
	if err != nil {
		Log.Println("HBase engine error: ", err)
		return CreateErrResponse("OnHBaseTestHandler", err)
	}

	// 循环解析数据
	jsonResult := make([]map[string]interface{}, 0)
	for _, rowResult := range HBaseSearchResult {
		mapRow := make(map[string]interface{})
		for index, value := range rowResult.Columns {
			value := value.Value
			colFamily := strings.Split(index, ":")
			if len(colFamily) == 2 {
				mapRow[colFamily[1]] = string(value)
			} else {
				mapRow[index] = string(value)
			}
		}

		jsonResult = append(jsonResult, mapRow)
	}

	b, err := json.MarshalIndent(jsonResult, "", " ")
	if err != nil {
		Log.Println("HBase engine error: ", err)
		return CreateErrResponse("OnHBaseTestHandler", err)
	}

	// 填充返回结果
	var HBaseResult []hbase_date_scheme
	err = json.Unmarshal(b, &HBaseResult)
	if err != nil {
		Log.Println("HBase engine error: ", err)
		return CreateErrResponse("OnHBaseTestHandler", err)
	}

	for _, v := range HBaseResult {
		response.ValueList1 = append(response.ValueList1, v.ColumnValue1)
		response.ValueList2 = append(response.ValueList2, v.ColumnValue2)
	}

	// END

	return CreateSuccessResponse(response)
}
