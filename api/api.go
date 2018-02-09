package api

import (
	"fmt"
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

	if config.Get().Username != "" && config.Get().Password != "" {
		glg.Info("Username and password defined, authentication enabled")
		group = r.Group("/", gin.BasicAuth(gin.Accounts{config.Get().Username: config.Get().Password}))
	} else {
		glg.Warn("Username and password not defined, authentication disabled")
		group = r.Group("/")
	}

	group.GET("/startup", startHandler)
	group.GET("/shutdown", stopHandler)

	r.Run(fmt.Sprintf("%s:%d", config.Get().Address, config.Get().Port))
}

func startHandler(c *gin.Context) {
	motionDetection, err := strconv.ParseBool(c.Query("detection"))

	if err != nil {
		motionDetection = false
	}

	err = motion.Startup(motionDetection)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
	} else {
		c.JSON(200, gin.H{"message": "motion started"})
	}
}

func stopHandler(c *gin.Context) {
	err := motion.Shutdown()

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
	} else {
		c.JSON(200, gin.H{"message": "motion stopped"})
	}
}
