package auth

import (
	"crypto/md5"
	"fmt"
)

func md5sum(content []byte) string {
	return fmt.Sprintf("%x", md5.Sum(content))
}
