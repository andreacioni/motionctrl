package backup

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/andreacioni/motionctrl/config"
)

func TestBackupCron(t *testing.T) {
	Init(config.Backup{When: "0 0 * * * *", Method: "gdrive"}, ".")
	Shutdown()
}

func TestBackupSize(t *testing.T) {
	Init(config.Backup{When: "10MB", Method: "gdrive"}, ".")
	Shutdown()
}

func TestListFiles(t *testing.T) {
	list, err := listFile(".")
	require.NoError(t, err)
	fmt.Printf("%+v\n", list)
}
