package service

import (
	"strconv"

	"github.com/fredmayer/mail-parser-rest/internal/app/models"
	"github.com/fredmayer/mail-parser-rest/internal/modules/mail"
	"github.com/labstack/echo/v4"
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

func (ms *MailService) GetList(page int, ctx echo.Context) ([]models.MailModel, error) {
	res, err := ms.mail.List(page)
	if err != nil {
		return nil, err
	}

	var items []models.MailModel
	for _, row := range res.Data {
		items = append(items, models.MailModel{
			MessageId: row.MessageId,
			Uid:       row.Uid,
			SeqNum:    row.SeqNum,
			From:      row.From,
			Subject:   row.Subject,
			Date:      row.Date,
		})
	}

	ctx.Response().Writer.Header().Set("X-Pagination-Page-Count", strconv.Itoa(res.Pages))

	return items, nil
}
