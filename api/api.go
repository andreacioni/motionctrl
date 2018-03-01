package api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/backup"
	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/motion"
	"github.com/andreacioni/motionctrl/utils"
	"github.com/andreacioni/motionctrl/version"
)

var handlersMap = map[string]func(*gin.Context){
	"/control/startup":  startHandler,
	"/control/shutdown": stopHandler,
	"/control/status":   statusHandler,
	"/control/restart":  restartHandler,
	"/detection/status": isMotionDetectionEnabled,
	"/detection/start":  startDetectionHandler,
	"/detection/stop":   stopDetectionHandler,
	"/camera/stream":    proxyStream,
	"/camera/snapshot":  takeSnapshot,
	"/config/list":      listConfigHandler,
	"/config/set":       setConfigHandler,
	"/config/get":       getConfigHandler,
	"/config/write":     writeConfigHandler,
	"/backup/status":    backupStatus,
}

func Init() {
	glg.Info("Initializing REST API ...")
	var group *gin.RouterGroup
	r := gin.Default()

	if config.GetConfig().Username != "" && config.GetConfig().Password != "" {
		glg.Info("Username and password defined, authentication enabled")
		group = r.Group("/", gin.BasicAuth(gin.Accounts{config.GetConfig().Username: config.GetConfig().Password}), needMotionUp)
	} else {
		glg.Warn("Username and password not defined, authentication disabled")
		group = r.Group("/", needMotionUp)
	}

	for k, v := range handlersMap {
		group.GET(k, v)
	}

	r.Run(fmt.Sprintf("%s:%d", config.GetConfig().Address, config.GetConfig().Port))
}

func needMotionUp(c *gin.Context) {

	/** Every request, except for /control* requests, need motion up and running**/

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
	writeback, err := strconv.ParseBool(c.DefaultQuery("writeback", "false"))

	if err != nil {
		nameAndValue := utils.RegexSubmatchTypedMap("/config/set\\?("+motion.KeyValueRegex+"+)=("+motion.KeyValueRegex+"+)", fmt.Sprint(c.Request.URL), motion.ReverseConfigTypeMapper)

		if len(nameAndValue) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "'name' and 'value' parameters not specified"})
		} else {
			for k, v := range nameAndValue {
				b := motion.ConfigCanSet(k)
				if b {
					err = motion.ConfigSet(k, v.(string))
					if err != nil {
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
					c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("'%s' cannot be updated with %s", k, version.Name)})
				}

			}

		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'writeback' parameter must be true/false"})
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
