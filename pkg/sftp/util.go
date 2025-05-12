package sftp

import (
	"html"
	"strings"
)

func UrlEncode(s string) string {
	return strings.ReplaceAll(html.EscapeString(s), "+", "%20")
}
