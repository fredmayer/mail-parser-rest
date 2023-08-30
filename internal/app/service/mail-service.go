package service

import (
	"errors"
	"strconv"

	"github.com/VictorRibeiroLima/converter"
	"github.com/fredmayer/mail-parser-rest/internal/app/models"
	"github.com/fredmayer/mail-parser-rest/internal/modules/mail"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
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

// List of folders(mailboxes)
func (ms *MailService) MailBoxes() ([]string, error) {
	return ms.mail.MailBoxes()
}

// Устанавливаем папку по-умолчанию
func (ms *MailService) SetFolder(folder string) error {
	_, error := ms.mail.SetFolder(folder)
	return error
}

func (ms *MailService) Move(uid int, mailbox string) error {
	list, err := ms.MailBoxes()
	if err != nil {
		return err
	}

	exist := slices.Contains(list, mailbox)
	if !exist {
		return errors.New("Not found " + mailbox + " in mailboxes list")
	}

	return ms.mail.Move(uint32(uid), mailbox)
}

func (ms *MailService) DownloadAttachment(uid int, index int, cxt echo.Context) (*mail.MailAttachmentDto, error) {
	message, err := ms.mail.GetBySid(uid)
	if err != nil {
		return nil, err
	}

	if len(message.Attachments) <= index {
		return nil, errors.New("Not found attachment")
	}
	at := message.Attachments[index]

	reader, err := ms.mail.DownloadAttachment(uid, at.Mime, at.Name)
	if err != nil {
		return nil, err
	}

	at.SetReader(reader)

	return &at, err
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

func (ms *MailService) GetLast(count int, ctx echo.Context) ([]models.MailModel, error) {
	if count > 200 { //max count 200
		count = 200
	}

	res, err := ms.mail.Last(count)
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

	return items, nil
}
