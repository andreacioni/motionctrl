package backup

type MockBackupService struct {
}

func (b *MockBackupService) Upload(file string) error {
	return nil
}

func (b *MockBackupService) Authenticate() error {
	return nil
}
