package backup

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/andreacioni/motionctrl/config"
)

func TestBackupCron(t *testing.T) {
	Init(config.Backup{When: "0 0 * * * *", Method: TestMockMethod}, ".")
	Shutdown()
}

func TestBackupSize(t *testing.T) {
	Init(config.Backup{When: "10MB", Method: TestMockMethod}, ".")
	Shutdown()
}

func TestListFiles(t *testing.T) {
	list, err := listFile(".")
	require.NoError(t, err)
	fmt.Printf("%+v\n", list)
}

func TestBackupMethod(t *testing.T) {
	_, err := buildUploadService(GoogleDriveMethod)

	require.NoError(t, err)

	_, err = buildUploadService(TestMockMethod)

	require.NoError(t, err)
}

func TestGetName(t *testing.T) {
	file, err := os.Open("../main.go")

	require.NoError(t, err)
	require.Equal(t, "main.go", file.Name())
}
