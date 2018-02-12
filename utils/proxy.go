package utils

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

//ReverseProxy is a courtesy of: https://github.com/gin-gonic/gin/issues/686
func ReverseProxy(target string) gin.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
