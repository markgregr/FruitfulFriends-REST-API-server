package helper

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	tokenHeader = "Authorization"
	tokenPrefix = "Bearer "
)

func ExtractTokenFromHeaders(c *gin.Context) string {
	value := c.GetHeader(tokenHeader)

	if value == "" {
		return ""
	}

	if !strings.HasPrefix(value, tokenPrefix) {
		return ""
	}

	return strings.TrimPrefix(value, tokenPrefix)
}

func GetRemoteAddr(r *http.Request) *string {
	ips, _, _ := net.SplitHostPort(r.RemoteAddr)
	if r.Header.Get("X-FORWARDED-FOR") != "" {
		ips = r.Header.Get("X-FORWARDED-FOR")
	}

	if ips == "" {
		return nil
	}

	ipsSlice := strings.Split(ips, ",")
	if len(ipsSlice) == 0 {
		return nil
	}

	currentIP := ipsSlice[0]
	if currentIP == "" {
		return nil
	}

	return &currentIP
}
