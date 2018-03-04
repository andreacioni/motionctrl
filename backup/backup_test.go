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

func TestNoBackupDefined(t *testing.T) {
	_, err := buildUploadService(config.Backup{})
	require.Error(t, err)
}

func TestListFiles(t *testing.T) {
	_, err := os.Create("../testfile")
	require.NoError(t, err)

	_, list, _, err := listFile("../")
	require.NoError(t, err)
	require.EqualValues(t, 5, len(list))
	fmt.Printf("%+v\n", list)

	err = os.Remove("../testfile")
	require.NoError(t, err)

}

func TestBackupMethod(t *testing.T) {
	_, err := buildUploadService(config.Backup{Method: GoogleDriveMethod})

	require.NoError(t, err)

	_, err = buildUploadService(config.Backup{Method: TestMockMethod})

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
	fmt.Println(filepath.Join(os.TempDir(), time.Now().Format("20060102_150405")+".tar.gz"))
}

func TestArchive(t *testing.T) {
	_, fileList, _, err := listFile(".")
	require.NoError(t, err)
	fmt.Println(fileList)

	archive, err := archiveFiles(fileList)
	require.NoError(t, err)
	fmt.Println(archive)

}
