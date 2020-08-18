package util

import (
	"net/http"
	"strings"
)

// GetServiceRelativePath retrieves the relative path of a service
func GetServiceRelativePath(r *http.Request, serviceKey string) string {
	return strings.Split(r.URL.String(), serviceKey)[1]
}
