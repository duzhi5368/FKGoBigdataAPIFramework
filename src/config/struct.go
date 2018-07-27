package config

import (
	"fmt"
	"reflect"
	"strconv"
)

type Configure struct {
	ListenPort                   int    `json:"listen_port"`
	DebugMode                    bool   `json:"debug_mode"`
	PrestoDBTLS                  bool   `json:"presto_db_tls"`
	PrestoDBIP                   string `json:"presto_db_ip"`
	PrestoDBPort                 int    `json:"presto_db_port"`
	PrestoDBCatalog              string `json:"presto_db_catalog"`
	PrestoDBSchema               string `json:"presto_db_schema"`
	PrestoDBUser                 string `json:"presto_db_user"`
	PrestoDBPassword             string `json:"presto_db_password"`
	ESHost                       string `json:"es_host"`
	ESPort                       int    `json:"es_port"`
	ESIndex                      string `json:"es_index"`
	PrestoGameMonitorHourlyTable string `json:"presto_game_monitor_hourly_table"`
	HbaseHost                    string `json:"hbase_host"`
	HbasePort                    int    `json:"hbase_port"`
	MysqlHost                    string `json:"mysql_host"`
	MysqlPort                    int    `json:"mysql_port"`
	MysqlDBName                  string `json:"mysql_db_name"`
	MysqlUser                    string `json:"mysql_user"`
	MysqlPassword                string `json:"mysql_password"`
}

var defaultConfig = []byte(`{
	"listen_port": 5000, 
	"debug_mode": true,
	"presto_db_tls": false,
	"presto_db_ip": "10.180.33.43",
	"presto_db_port": 8086,
	"presto_db_catalog": "hive",
	"presto_db_schema": "test",
	"presto_db_user": "user",
	"presto_db_password": "",
    "es_host": "10.180.33.43",
	 "es_port": 9200,
	"es_index": "test_es_index",
    "presto_game_monitor_hourly_table": "t_monitor_games_dw",
    "hbase_host": "10.180.33.43",
    "hbase_port": 40000,
	"mysql_host": "10.180.33.43",
	"mysql_port": "3306",
	"mysql_db_name": "mldm",
	"mysql_user": "mldm",
	"mysql_password": "mldm123"
	}`)

func (p Configure) DumpInfo() string {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	var result = "==========================================\n"
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
