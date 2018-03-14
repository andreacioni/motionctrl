package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/andreacioni/motionctrl/motion"
	"github.com/andreacioni/motionctrl/utils"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

// isLocalhost middlewares permit requests only from localhost
func isLocalhost(c *gin.Context) {
	ipStr, _, err := net.SplitHostPort(c.Request.RemoteAddr)

	if err == nil {
		fromIp := net.ParseIP(ipStr)

		if b, err := utils.IsLocalIP(fromIp); err != nil || !b {
			glg.Warnf("Rejecting request to %s, IP: %s is not authorized", c.Request.URL.Path, c.Request.RemoteAddr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "call to /internal/* api is allowed only from localhost"})
		} else {
			glg.Debugf("Accepting /internal api call from %s", c.Request.RemoteAddr)
		}

	} else {
		glg.Errorf("Error in parsing request remote address: %s: %v", c.Request.RemoteAddr, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("cannot parse remote address: %s: %v", c.Request.RemoteAddr, err)})
	}

}

// needMotionUp Every request, except for /control* requests, need motion up and running
func needMotionUp(c *gin.Context) {
	if !strings.HasPrefix(fmt.Sprint(c.Request.URL), "/api/control") {

		if motionStarted, err := motion.IsStarted(); err == nil {
			if !motionStarted {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "motion was not started yet"})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": fmt.Errorf("Unable to check if motion is up: %v", err)})
			return
		}

	}
}
