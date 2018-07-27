package auth

import (
	"auth_storage"
	"fmt"
	"slog"
	"strings"
)

type auth struct {
	key     string
	md5     string
	content []byte
}

func NewAuth(key, md5 string, content []byte) *auth {
	return &auth{
		key,
		md5,
		content,
	}
}
func (a *auth) Check() error {
	keys, err := auth_storage.Query(a.key)
	if err != nil {
		return err
	}
	// post内容在前 secret在后，拼接后求md5
	content := append(a.content, []byte(keys.Secret)...)
	slog.Log.Debug("MD5 content: " + string(content))
	var serverMd5 string = strings.ToLower(md5sum(content))
	var headerMd5 string = strings.ToLower(a.md5)
	if serverMd5 != headerMd5 {
		return fmt.Errorf("auth check failed, [header] %s != %s [server]", headerMd5, serverMd5)
	}
	return nil
}

func GetProductList(key string) []string {
	keys, err := auth_storage.Query(key)
	if err != nil {
		return nil
	}

	return keys.ProductList
}
