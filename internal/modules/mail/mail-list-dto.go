package mail

import (
	"io"
	"time"
)

type MailDto struct {
	MessageId   string
	Uid         uint32
	SeqNum      uint32
	From        string
	Subject     string
	Date        time.Time
	Attachments []MailAttachmentDto
}

type MailAttachmentDto struct {
	Mime   string
	Name   string
	Index  int
	reader io.Reader
}

type ListMailResponse struct {
	Data  []MailDto
	Pages int
}

func (ma *MailAttachmentDto) SetReader(r io.Reader) {
	ma.reader = r
}

func (ma *MailAttachmentDto) GetReader() io.Reader {
	return ma.reader
}
