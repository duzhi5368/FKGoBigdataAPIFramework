package auth_storage

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Log(loadKeyFromFile())
	t.Log(Query("1247fbc0e373637ff9fb08f0c81722f3"))
}

//------------------------------------------------------------
// 生成一个UUID
func newSecretKey() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

//------------------------------------------------------------
// 生成一个UUID
func newAccessKey() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func Test_GenerateAllProductConfig(t *testing.T) {
	var auth Keys
	productList := []string{
		"A01", "A02", "A03", "A04", "A05", "A06", "B01", "B05", "B06", "B07", "C01", "C02", "C07",
	}
	auth.Version = "1.0.0"

	for _, pid := range productList {
		key := Key{}
		key.ProductList = []string{pid}
		key.Access = newAccessKey()
		key.Secret = newSecretKey()
		auth.Keys = append(auth.Keys, key)
	}
	key := Key{}
	key.ProductList = productList
	key.Access = newAccessKey()
	key.Secret = newSecretKey()
	auth.Keys = append(auth.Keys, key)
	b, _ := json.MarshalIndent(auth, "", " ")
	fmt.Println(string(b))
}
