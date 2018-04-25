package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/backup"
	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/motion"
	"github.com/andreacioni/motionctrl/notify"
	"github.com/andreacioni/motionctrl/utils"
	"github.com/andreacioni/motionctrl/version"
)

// MethodHandler utility struct that contains method and associated handler
type MethodHandler struct {
	method string
	f      func(*gin.Context)
	m      []gin.HandlerFunc
}

var internalHandlersMap = map[string]MethodHandler{
	"/event/start":           {method: http.MethodGet, f: eventStart},
	"/event/end":             {method: http.MethodGet, f: eventEnd},
	"/event/motion/detected": {method: http.MethodGet, f: motionDetected},
	"/event/picture/saved":   {method: http.MethodGet, f: pictureSaved},
}

var handlersMap = map[string]MethodHandler{
	"/control/startup":  {method: http.MethodGet, f: startHandler},
	"/control/shutdown": {method: http.MethodGet, f: stopHandler},
	"/control/status":   {method: http.MethodGet, f: statusHandler},
	"/control/restart":  {method: http.MethodGet, f: restartHandler, m: []gin.HandlerFunc{needMotionUp}},

	"/detection/status": {method: http.MethodGet, f: isMotionDetectionEnabled, m: []gin.HandlerFunc{needMotionUp}},
	"/detection/start":  {method: http.MethodGet, f: startDetectionHandler, m: []gin.HandlerFunc{needMotionUp}},
	"/detection/stop":   {method: http.MethodGet, f: stopDetectionHandler, m: []gin.HandlerFunc{needMotionUp}},

	"/camera/stream":    {method: http.MethodGet, f: proxyStream, m: []gin.HandlerFunc{needMotionUp}},
	"/camera/snapshot":  {method: http.MethodGet, f: takeSnapshot, m: []gin.HandlerFunc{needMotionUp}},
	"/camera/makemovie": {method: http.MethodGet, f: makeMovie, m: []gin.HandlerFunc{needMotionUp}},

	"/config/list":       {method: http.MethodGet, f: listConfigHandler, m: []gin.HandlerFunc{needMotionUp}},
	"/config/set":        {method: http.MethodGet, f: setConfigHandler, m: []gin.HandlerFunc{needMotionUp}},
	"/config/get/:param": {method: http.MethodGet, f: getConfigHandler, m: []gin.HandlerFunc{needMotionUp}},
	"/config/write":      {method: http.MethodGet, f: writeConfigHandler, m: []gin.HandlerFunc{needMotionUp}},

	"/targetdir/list":             {method: http.MethodGet, f: listTargetDir},
	"/targetdir/size":             {method: http.MethodGet, f: sizeTargetDir},
	"/targetdir/get/:filename":    {method: http.MethodGet, f: retrieveFromTargetDir},
	"/targetdir/remove/:filename": {method: http.MethodGet, f: removeFromTargetDir},

	"/backup/status": {method: http.MethodGet, f: backupStatus},
	"/backup/launch": {method: http.MethodGet, f: backupLaunch},

	"/notify/status":     {method: http.MethodGet, f: notifyStatus},
	"/notify/activate":   {method: http.MethodGet, f: notifyActivate},
	"/notify/deactivate": {method: http.MethodGet, f: notifyDeactivate},
}

func Init(conf config.Configuration, shutdownHook func()) error {
	glg.Info("Initializing REST API ...")

	if conf.Address == "" || conf.Port <= 0 {
		return fmt.Errorf("Address and/or port not defined in configuration")
	}

	var group *gin.RouterGroup
	router := gin.Default()

	// /app
	if conf.AppPath != "" {
		glg.Infof("Serving static files from %s to /app", conf.AppPath)
		router.Static("/app", conf.AppPath)
	}

	// /internal
	internal := router.Group("/internal", isLocalhost)

	for path, handler := range internalHandlersMap {
		internal.Handle(handler.method, path, handler.f)
	}

	// /api
	if conf.Username != "" && conf.Password != "" {
		glg.Info("Username and password defined, authentication enabled")
		group = router.Group("/api", gin.BasicAuth(gin.Accounts{conf.Username: conf.Password}))
	} else {
		glg.Warn("Username and password not defined, authentication disabled")
		group = router.Group("/api")
	}

	for path, handler := range handlersMap {
		group.Handle(handler.method, path, append(handler.m, handler.f)...)
	}

	if err := listenAndServe(router, shutdownHook, fmt.Sprintf("%s:%d", conf.Address, conf.Port), conf.Ssl); err != nil {
		return fmt.Errorf("unable to listen & serve: %v", err)
	}

	return nil
}

//From: https://github.com/gin-gonic/gin#graceful-restart-or-stop
func listenAndServe(router *gin.Engine, shutdownHook func(), addressPort string, sslConf config.SSL) error {
	server := &http.Server{
		Addr:    addressPort,
		Handler: router,
	}

	go func() {
		if !sslConf.IsEmpty() {
			glg.Infof("SSL/TLS enabled for API using certificate: %s, key: %s", sslConf.CertFile, sslConf.KeyFile)
			if err := server.ListenAndServeTLS(sslConf.CertFile, sslConf.KeyFile); err != nil {
				glg.Error(err)
			}
		} else {
			glg.Warn("SSL/TLS NOT enabled for API")
			if err := server.ListenAndServe(); err != nil {
				glg.Error(err)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	glg.Debug("Shutdown Server ...")

	shutdownHook()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		glg.Fatal("Server Shutdown:", err)
	}
	glg.Info("Server exiting")

	return nil
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
	if started, err := motion.IsStarted(); err == nil {
		c.JSON(http.StatusOK, gin.H{"motionStarted": started})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Unable to check if motion is up: %v", err)})
	}
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
		glg.Debugf("Snapshot file: %s", snapFile)
		c.File(snapFile)
	}
}

func makeMovie(c *gin.Context) {
	err := motion.MakeMovie()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "movie recording started"})
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
	query := c.Param("param")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'query' parameter not specified"})
	} else {
		config, err := motion.ConfigGet(query)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()}) //TODO improve fail with returned status code from request sent to motion
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{query: config})
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

func listTargetDir(c *gin.Context) {
	if fileList, err := motion.TargetDirListFiles(); err == nil {
		c.JSON(http.StatusOK, fileList)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Unable to list files in target dir: %v", err)})
	}
}

func sizeTargetDir(c *gin.Context) {
	if size, err := motion.TargetDirSize(); err == nil {
		c.JSON(http.StatusOK, gin.H{"size": size})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Unable to list files in target dir: %v", err)})
	}
}

func retrieveFromTargetDir(c *gin.Context) {
	fileName := c.Param("filename")

	if fileName != "" {
		if filePath, err := motion.TargetDirGetFile(fileName); err == nil {
			c.File(filePath)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Unable to get: %s in target dir: %v", fileName, err)})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing 'filename' parameter"})
	}

}

func removeFromTargetDir(c *gin.Context) {
	fileName := c.Param("filename")

	if fileName != "" {
		if err := motion.TargetDirRemoveFile(fileName); err == nil {
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s successfully removed", fileName)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Unable to remove: %s in target dir: %v", fileName, err)})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing 'filename' parameter"})
	}

}

func notifyStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ready":  notify.IsReady(),
		"active": notify.IsActive(),
	})
}

func notifyActivate(c *gin.Context) {
	if err := notify.SetActive(true); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "notify service is ready and active now"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("unable to activate notify: %v", err)})
	}
}

func notifyDeactivate(c *gin.Context) {
	if err := notify.SetActive(false); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "notify service deactivated"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("unable to activate notify: %v", err)})
	}
}
