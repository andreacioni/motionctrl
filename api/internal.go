package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/notify"
)

func eventStart(c *gin.Context) {
	notify.MotionDetectedStart()
}
func eventEnd(c *gin.Context) {
	notify.MotionDetectedStop()
}
func motionDetected(c *gin.Context) { //TODO not sure this is useful by now
	glg.Debugf("Motion detected")
}

func areaDetected(c *gin.Context) {
	glg.Debugf("Area detected")
}

func movieStart(c *gin.Context) {
	glg.Debugf("Movie start")
}

func movieEnd(c *gin.Context) {
	glg.Debugf("Movie end")
}

func cameraLost(c *gin.Context) {
	glg.Debugf("Camera lost")
}

func cameraFound(c *gin.Context) {
	glg.Debugf("Camera found")
}

func pictureSaved(c *gin.Context) {
	picturePath := c.Query("picturepath")

	if picturePath != "" {
		glg.Debugf("Picture saved in: %s", picturePath)
		notify.PhotoSaved(picturePath)
	} else {
		glg.Warnf("'picturepath' not found. Unable to know where picture is.")
	}

}
