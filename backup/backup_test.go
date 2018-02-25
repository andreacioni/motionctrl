package backup

import (
	"fmt"
	"path/filepath"
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
	require.Equal(t, "main.go", filepath.Base("../main.go"))
	require.Equal(t, "file", filepath.Base("/abc/def/file"))
}

func TestExtToMime(t *testing.T) {

	require.Equal(t, ".jpeg", filepath.Ext("/abc/test.jpeg"))
	require.Equal(t, ".jpg", filepath.Ext("/abc/test.jpg"))

	require.Equal(t, "image/jpeg", GoogleDriveBackupService{}.mimeFromExt(".jpeg"))
	require.Equal(t, "image/jpeg", GoogleDriveBackupService{}.mimeFromExt(".jpg"))
}
