package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/notify"
)

func eventStart(c *gin.Context) {
	glg.Debugf("Event start")
}
func eventEnd(c *gin.Context) {
	glg.Debugf("Event end")
}
func motionDetected(c *gin.Context) {
	glg.Debugf("Motion detected")
}

func pictureSaved(c *gin.Context) {
	picturePath := c.Query("picturepath")

	if picturePath != "" {
		glg.Debugf("Picture saved in: %s")
		notify.PhotoSaved(picturePath)
	} else {
		glg.Warnf("'picturepath' not found. Unable to know where picture is.")
	}

}
