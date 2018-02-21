package backup

import (
	"github.com/kpango/glg"
	"github.com/radovskyb/watcher"

	"../motion"
)

var (
	dirWatcher watcher.Watcher
)

func Init() {

	dirWatcher := watcher.New()

	if err := dirWatcher.AddRecursive(motion.ConfigGet(motion.TargetDir)); err != nil {
		glg.Fatal(err)
	}

}

func setupWatcher() {

}
