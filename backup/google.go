package backup

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"

	"github.com/andreacioni/motionctrl/version"
)

type GoogleDriveBackupService struct {
	service *drive.Service
}

func (b GoogleDriveBackupService) Upload(filePath string) error {
	dir, err := b.getRemoteDir()

	if err != nil {
		return fmt.Errorf("Unable to retrieve base directory %v", err)
	}

	file, err := os.Open(filePath)

	if err != nil {
		return fmt.Errorf("Unable open file %s: %v", filePath, err)
	}

	remoteFile := &drive.File{
		Name:     file.Name(),
		MimeType: b.mimeFromExt(filepath.Ext(filePath)),
		Parents:  []string{dir.Id},
	}

	_, err = b.service.Files.Create(remoteFile).Media(file).Do()

	return err
}

func (b GoogleDriveBackupService) mimeFromExt(ext string) string {
	switch ext {
	case ".jpg":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".avi":
		return "video/avi"
	case ".mjpeg":
		return "video/x-motion-jpeg"
	}

	return ""
}

//getRemoteDir check if directory 'motionctrl' exists inside root, if not create it
func (b GoogleDriveBackupService) getRemoteDir() (*drive.File, error) {
	r, err := b.service.Files.List().Q(fmt.Sprintf("'root' in parents and name='%s' and mimeType='application/vnd.google-apps.folder' and trashed = false", version.Name)).PageSize(1).
		Fields("nextPageToken, files(id, name)").Do()

	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve information about Google Drive root directory: %v", err)
	}

	if len(r.Files) == 0 {
		return b.createRemoteDir()
	}

	return r.Files[0], nil

}

func (b GoogleDriveBackupService) createRemoteDir() (*drive.File, error) {
	remoteDir := &drive.File{
		Name:     version.Name,
		Parents:  []string{"root"},
		MimeType: "application/vnd.google-apps.folder",
	}

	remoteDir, err := b.service.Files.Create(remoteDir).Do()

	if err != nil {
		return nil, fmt.Errorf("Unable to create remote directory: %v", err)
	}

	return remoteDir, nil

}

func (b GoogleDriveBackupService) Authenticate() error {
	ctx := context.Background()

	config := b.getConfig()

	client, err := b.getClient(ctx, config)

	if err != nil {
		return fmt.Errorf("Unable to retrieve drive Client %v", err)
	}

	b.service, err = drive.New(client)

	if err != nil {
		return fmt.Errorf("Unable to get service instance %v", err)
	}

	return nil
}

func (b GoogleDriveBackupService) getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "568343557575-pefr1f8otcq5pg5o5gegntbc0f31hm02.apps.googleusercontent.com",
		ClientSecret: "gb-oxnJeOSAbiMq5uymynfOA",
		Scopes:       []string{drive.DriveScope},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func (b GoogleDriveBackupService) getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	cacheFile, err := b.tokenCacheFile()
	if err != nil {
		return nil, fmt.Errorf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := b.tokenFromFile(cacheFile)
	if err != nil {
		tok = b.getTokenFromWeb(config)
		b.saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok), nil
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func (b GoogleDriveBackupService) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func (b GoogleDriveBackupService) tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := usr.HomeDir
	return filepath.Join(tokenCacheDir,
		url.QueryEscape(fmt.Sprintf(".google_drive_%s.json", version.Name))), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func (b GoogleDriveBackupService) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func (b GoogleDriveBackupService) saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
