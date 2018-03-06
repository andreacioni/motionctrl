package notify

type MockNotifyService struct {
}

func (b *MockNotifyService) Notify(message, file string) error {
	return nil
}

func (b *MockNotifyService) Authenticate() error {
	return nil
}

func (b *MockNotifyService) Stop() error {
	return nil
}
