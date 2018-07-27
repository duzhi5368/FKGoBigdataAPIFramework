package config

import (
	"encoding/json"
	"testing"
)

func Test_GetJson(t *testing.T) {
	b, _ := json.MarshalIndent(&Config, "", " ")
	t.Log(string(b))
}
