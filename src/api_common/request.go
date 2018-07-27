package api_common

import (
	"net/url"
	"strings"

	"github.com/emicklei/go-restful"
)

func splitContentType(contentType string) string {
	contentType = strings.ToLower(contentType)
	i := strings.Split(contentType, ";")
	return i[0]
}

func getContentType(req *restful.Request) string {
	contentType := req.HeaderParameter("Content-Type")
	return splitContentType(contentType)
}

func RawRequest(req *restful.Request) []byte {
	raw := req.Request.URL.RawQuery
	return []byte(raw)
}

func MapRequest(req *restful.Request) url.Values {
	return req.Request.URL.Query()
}
