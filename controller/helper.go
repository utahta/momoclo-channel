package controller

import (
	"bytes"
	"fmt"
	"net/url"
)

func buildURL(u *url.URL, path string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	if len(path) > 0 && path[0] != '/' {
		buf.WriteString("/")
	}
	buf.WriteString(path)

	return buf.String()
}
