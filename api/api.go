package api

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/backup"
	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/motion"
	"github.com/andreacioni/motionctrl/utils"
	"github.com/andreacioni/motionctrl/version"
)

// MethodHandler utility struct that contains method and associated handler
type MethodHandler struct {
	method string
	f      func(*gin.Context)
}

var internalHandlersMap = map[string]MethodHandler{
	"/event/start":           MethodHandler{method: http.MethodGet, f: eventStart},
	"/event/end":             MethodHandler{method: http.MethodGet, f: eventEnd},
	"/event/motion/detected": MethodHandler{method: http.MethodGet, f: motionDetected},
	"/event/picture/saved":   MethodHandler{method: http.MethodGet, f: pictureSaved},
}

var handlersMap = map[string]MethodHandler{
	"/control/startup":  MethodHandler{method: http.MethodGet, f: startHandler},
	"/control/shutdown": MethodHandler{method: http.MethodGet, f: stopHandler},
	"/control/status":   MethodHandler{method: http.MethodGet, f: statusHandler},
	"/control/restart":  MethodHandler{method: http.MethodGet, f: restartHandler},
	"/detection/status": MethodHandler{method: http.MethodGet, f: isMotionDetectionEnabled},
	"/detection/start":  MethodHandler{method: http.MethodGet, f: startDetectionHandler},
	"/detection/stop":   MethodHandler{method: http.MethodGet, f: stopDetectionHandler},
	"/camera/stream":    MethodHandler{method: http.MethodGet, f: proxyStream},
	"/camera/snapshot":  MethodHandler{method: http.MethodGet, f: takeSnapshot},
	"/config/list":      MethodHandler{method: http.MethodGet, f: listConfigHandler},
	"/config/set":       MethodHandler{method: http.MethodGet, f: setConfigHandler},
	"/config/get":       MethodHandler{method: http.MethodGet, f: getConfigHandler},
	"/config/write":     MethodHandler{method: http.MethodGet, f: writeConfigHandler},
	"/backup/status":    MethodHandler{method: http.MethodGet, f: backupStatus},
	"/backup/launch":    MethodHandler{method: http.MethodGet, f: backupLaunch},
}

func Init(conf config.Configuration, shutdownHook func()) error {
	glg.Info("Initializing REST API ...")

	if conf.Address == "" || conf.Port <= 0 {
		return fmt.Errorf("Address and/or port not defined in configuration")
	}

	var group *gin.RouterGroup
	router := gin.Default()

	internal := router.Group("/internal", isLocalhost)

	for path, handler := range internalHandlersMap {
		internal.Handle(handler.method, path, handler.f)
	}

	if conf.Username != "" && conf.Password != "" {
		glg.Info("Username and password defined, authentication enabled")
		group = router.Group("/", gin.BasicAuth(gin.Accounts{conf.Username: conf.Password}), needMotionUp)
	} else {
		glg.Warn("Username and password not defined, authentication disabled")
		group = router.Group("/", needMotionUp)
	}

	for path, handler := range handlersMap {
		group.Handle(handler.method, path, handler.f)
	}

	if err := listenAndServe(router, shutdownHook, fmt.Sprintf("%s:%d", conf.Address, conf.Port), conf.Ssl.CertFile, conf.Ssl.KeyFile); err != nil {
		return fmt.Errorf("ListenAndServeTLS fail: %v", err)
	}

	return nil
}

func listenAndServe(router *gin.Engine, shutdownHook func(), addressPort, certFile, keyFile string) error {
	server := endless.NewServer(addressPort, router)

	server.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGINT, shutdownHook)
	server.RegisterSignalHook(endless.PRE_SIGNAL, syscall.SIGTERM, shutdownHook)

	if certFile != "" && keyFile != "" {
		glg.Infof("SSL/TLS enabled using certificate: %s, key: %s", certFile, keyFile)
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			return err
		}
	} else {
		glg.Warn("SSL/TLS NOT enabled")
		if err := server.ListenAndServe(); err != nil {
			return err
		}
	}

	return nil
}

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
	if !strings.HasPrefix(fmt.Sprint(c.Request.URL), "/control") {
		motionStarted := motion.IsStarted()

		if !motionStarted {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "motion was not started yet"})
			return
		}
	}
}

func startHandler(c *gin.Context) {
	motionDetection, err := strconv.ParseBool(c.DefaultQuery("detection", "false"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'detection' parameter must be 'true' or 'false'"})
	} else {
		err = motion.Startup(motionDetection)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "motion started"})
		}
	}

}

func restartHandler(c *gin.Context) {
	err := motion.Restart()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion restarted"})
	}
}

func stopHandler(c *gin.Context) {
	err := motion.Shutdown()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion stopped"})
	}
}

func statusHandler(c *gin.Context) {
	started := motion.IsStarted()

	c.JSON(http.StatusOK, gin.H{"motionStarted": started})
}

func isMotionDetectionEnabled(c *gin.Context) {
	enabled, err := motion.IsMotionDetectionEnabled()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"motionDetectionEnabled": enabled})
	}
}

func startDetectionHandler(c *gin.Context) {
	err := motion.EnableMotionDetection()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion detection started"})
	}
}

func stopDetectionHandler(c *gin.Context) {
	err := motion.DisableMotionDetection()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion detection paused"})
	}
}

//proxyStream is a courtesy of: https://github.com/gin-gonic/gin/issues/686
func proxyStream(c *gin.Context) {
	url, _ := url.Parse(motion.GetStreamBaseURL())
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func takeSnapshot(c *gin.Context) {
	snapFile, err := motion.Snapshot()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.File(snapFile)
	}
}

func listConfigHandler(c *gin.Context) {
	configMap, err := motion.ConfigList()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, configMap)
	}
}

func getConfigHandler(c *gin.Context) {
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'query' parameter not specified"})
	} else {
		config, err := motion.ConfigGet(query)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()}) //TODO improve fail with returned status code from request sent to motion
		} else {
			c.JSON(http.StatusOK, config)
		}
	}
}

func setConfigHandler(c *gin.Context) {
	writeback, _ := strconv.ParseBool(c.DefaultQuery("writeback", "false"))

	nameAndValue := utils.RegexSubmatchTypedMap("/config/set\\?("+motion.KeyValueRegex+"+)=("+motion.KeyValueRegex+"+)", fmt.Sprint(c.Request.URL), motion.ReverseConfigTypeMapper)

	if len(nameAndValue) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'name' and 'value' parameters not specified"})
	} else {
		for k, v := range nameAndValue {
			b := motion.ConfigCanSet(k)
			if b {
				if err := motion.ConfigSet(k, v.(string)); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()}) //TODO improve fail with returned status code from request sent to motion
				} else {

					if writeback {
						err = motion.ConfigWrite()
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
							return
						}
					}
					c.JSON(http.StatusOK, gin.H{k: motion.ConfigTypeMapper(v.(string))})
				}
			} else {
				c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("'%s' parameter cannot be updated through %s", k, version.Name)})
			}
		}
	}
}

func writeConfigHandler(c *gin.Context) {
	err := motion.ConfigWrite()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "configuration written to file"})
	}
}

func backupStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": backup.GetStatus()})
}

func backupLaunch(c *gin.Context) {
	if err := backup.RunNow(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "backup service is running now"})
	}
}
