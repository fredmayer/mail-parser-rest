package service

type Manager struct {
	MailService *MailService
}

func NewManager() (*Manager, error) {
	return &Manager{
		MailService: NewMailService(),
	}, nil
}
