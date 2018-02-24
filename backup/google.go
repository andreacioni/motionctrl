package backup

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type GoogleDriveBackupService struct {
	service *drive.Service
}

func (b GoogleDriveBackupService) Upload(file string) error {
	return nil
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

func (b GoogleDriveBackupService) getClient(ctx context.Context, cfg *oauth2.Config) (*http.Client, error) {
	return nil, nil
}
