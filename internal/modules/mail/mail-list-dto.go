package mail

import "time"

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
	Mime  string
	Name  string
	Index int
}

type ListMailResponse struct {
	Data  []MailDto
	Pages int
}
