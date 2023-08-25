package mail

import (
	"log"

	"github.com/emersion/go-imap"
)

// Setted current Folder
func (mr *MailReader) SetFolder(folder string) (*imap.MailboxStatus, error) {
	status, err := mr.Cl.Select(folder, false)
	if err != nil {
		return nil, err
	}
	mr.CurrentFolder = status

	return status, nil
}

// Инициализация папки почтового ящика
func (mr *MailReader) initMailbox() *imap.MailboxStatus {
	if mr.CurrentFolder == nil {
		_, err := mr.SetFolder("INBOX")
		if err != nil {
			log.Panic(err)
		}
	}

	return mr.CurrentFolder
}

// List of mailboxes
func (mr *MailReader) MailBoxes() ([]string, error) {
	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- mr.Cl.List("", "*", mailboxes)
	}()

	res := make([]string, 0, 10)
	for m := range mailboxes {
		res = append(res, m.Name)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return res, nil
}
