package backup

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	list, err := listFile("../")
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

	require.Equal(t, "image/jpeg", mime.TypeByExtension(".jpeg"))
	require.Equal(t, "image/jpeg", mime.TypeByExtension(".jpg"))
}

func TestJoinArchive(t *testing.T) {
	fmt.Println(filepath.Join(os.TempDir(), time.Now().Format("20060102_150405")+"tar.gz"))
}

func TestArchive(t *testing.T) {
	fileList, err := listFile(".")
	require.NoError(t, err)
	fmt.Println(fileList)

	archive, err := archiveFiles(fileList)
	require.NoError(t, err)
	fmt.Println(archive)

}
