package notify

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/andreacioni/motionctrl/config"
)

func TestNotify(t *testing.T) {
	require.NoError(t, Init(config.Notify{Method: "mock"}))
	Shutdown()
}

func TestFailAutheticate(t *testing.T) {
	require.Error(t, Init(config.Notify{Method: "telegram"}))

	MotionDetectedStart()
	MotionDetectedStop()
	PhotoSaved("")

	Shutdown()
}

func TestPhotoLimit(t *testing.T) {
	n := 3
	require.NoError(t, Init(config.Notify{Method: "mock", Photo: n}))

	MotionDetectedStart()

	for i := 0; i < n; i++ {
		PhotoSaved("/tmp") //OK
	}

	PhotoSaved("/tmp") // this will not sent

	MotionDetectedStop() //Reset

	MotionDetectedStart()

	for i := 0; i < n; i++ {
		PhotoSaved("/tmp") //OK
	}

	PhotoSaved("/tmp") // this will not sent

	PhotoSaved("/tmp") // this will not sent

	Shutdown()
}

func TestPhotoLimitWithoutStop(t *testing.T) {
	n := 3
	require.NoError(t, Init(config.Notify{Method: "mock", Photo: n}))

	MotionDetectedStart()

	for i := 0; i < n; i++ {
		PhotoSaved("/tmp") //OK
	}

	PhotoSaved("/tmp") // this will not sent

	//MotionDetectedStop() //NOT RESETTING

	MotionDetectedStart()

	for i := 0; i < n; i++ {
		PhotoSaved("/tmp") //OK
	}

	PhotoSaved("/tmp") // this will not sent

	PhotoSaved("/tmp") // this will not sent

	Shutdown()
}
