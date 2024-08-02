package server

import (
	"strings"
)

const (
	textHtmlContentType        = "text/html"
	textPlainContentType       = "text/plain"
	applicationJSONContentType = "application/json"
	anyContentType             = "*/*"
	defaultContentType         = textPlainContentType
)

func getType(accept string) string {
	l := strings.Split(accept, ",")
	for _, a := range l {
		fa := strings.TrimSpace(a)
		if strings.HasPrefix(textHtmlContentType, fa) {
			return textHtmlContentType
		}
		if strings.HasPrefix(textPlainContentType, fa) {
			return textPlainContentType
		}
		if strings.HasPrefix(applicationJSONContentType, fa) {
			return applicationJSONContentType
		}
	}
	if strings.Contains(accept, anyContentType) {
		return defaultContentType
	}
	return "invalid"
}
