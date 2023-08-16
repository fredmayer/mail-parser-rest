package service

import (
	"github.com/fredmayer/mail-parser-rest/internal/app/models"
	"github.com/fredmayer/mail-parser-rest/internal/modules/mail"
)

type MailService struct {
	//ctx  context.Context
	mail *mail.MailReader
}

func NewMailService() *MailService {
	return &MailService{
		mail: mail.Get(),
	}
}

func (ms *MailService) GetList(page int) ([]models.MailModel, error) {
	res := ms.mail.List(page)

	var items []models.MailModel
	for _, row := range res {
		items = append(items, models.MailModel{
			MessageId: row.MessageId,
			Uid:       row.Uid,
			SeqNum:    row.SeqNum,
			From:      row.From,
			Subject:   row.Subject,
			Date:      row.Date,
		})
	}

	return items, nil
}
