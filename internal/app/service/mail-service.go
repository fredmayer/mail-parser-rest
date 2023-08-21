package service

import (
	"io"
	"strconv"

	"github.com/VictorRibeiroLima/converter"
	"github.com/fredmayer/mail-parser-rest/internal/app/models"
	"github.com/fredmayer/mail-parser-rest/internal/app/types"
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

func (ms *MailService) DownloadAttachment(uid int, params types.AttachmentRequest, cxt echo.Context) (io.Reader, error) {
	reader, err := ms.mail.DownloadAttachment(uid, params.Mime, params.Name)
	if err != nil {
		return nil, err
	}

	return reader, err
}

func (ms *MailService) GetView(uid int, ctx echo.Context) (*models.MailModel, error) {
	message, err := ms.mail.GetBySid(uid)
	if err != nil {
		return nil, err
	}

	at := []models.MailAttachmentModel{}
	converter.Convert(&at, message.Attachments)

	res := models.MailModel{
		MessageId:   message.MessageId,
		Uid:         message.Uid,
		SeqNum:      message.SeqNum,
		From:        message.From,
		Subject:     message.Subject,
		Date:        message.Date,
		Attachments: at,
	}

	return &res, nil
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
