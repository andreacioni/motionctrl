package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"../config"
	"../motion"
)

func Init() {
	glg.Info("Initializing REST API ...")
	var group *gin.RouterGroup
	r := gin.Default()

	if config.Conf.Username != "" && config.Conf.Password != "" {
		glg.Info("Username and password defined, authentication is enabled")
		group = r.Group("/", gin.BasicAuth(gin.Accounts{config.Conf.Username: config.Conf.Password}))
	} else {
		glg.Warn("Username and password not defined, authentication is disabled")
		group = r.Group("/")
	}

	group.GET("/startup", startHandler)
	group.GET("/shutdown", stopHandler)

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}

func startHandler(c *gin.Context) {
	motionDetection := c.Query("startup_motion_detection")
	
	if motionDetection == "true" || motionDetection == "false") {
		motionDetection  := strconv.FormatBool(motionDetection)
	} else {
		motionDetection = false
	}
	
	motion.Startup(motionDetection)
}

func stopHandler(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
