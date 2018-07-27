package auth

import (
	"testing"
)

func TestMd5sum(t *testing.T) {
	t.Log(md5sum([]byte("testMD5")))
}
