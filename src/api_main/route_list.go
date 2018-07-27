package api_main

import (
	"api_simple"
)

var (
	mapAPIHandler = map[string]apiStruct{
		"api_test/echo_test": {
			isAuth:     false,
			apiHandler: api_simple.OnEchoTestHandler,
			desc:       "测试通讯Echo样例",
		},
		"api_test/auth_test": {
			isAuth:     true,
			apiHandler: api_simple.OnAuthTestHandler,
			desc:       "测试安全签名样例",
		},
		"api_test/es_test": {
			isAuth:     false,
			apiHandler: api_simple.OnESTestHandler,
			desc:       "测试ES连接查询样例",
		},
		"api_test/hbase_test": {
			isAuth:     false,
			apiHandler: api_simple.OnHBaseTestHandler,
			desc:       "测试HBase查询样例",
		},
		"api_test/sql_test": {
			isAuth:     false,
			apiHandler: api_simple.OnMySQLTestHandler,
			desc:       "测试MySQL查询样例",
		},
		"api_test/presto_test": {
			isAuth:     false,
			apiHandler: api_simple.OnPrestoTestHandler,
			desc:       "测试Presto查询样例",
		},
	}
)
