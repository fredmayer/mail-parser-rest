package models

import "time"

type MailModel struct {
	MessageId   string                `json:"mid"`
	Uid         uint32                `json:"uid"`
	SeqNum      uint32                `json:"sid"`
	From        string                `json:"from"`
	Subject     string                `json:"subject"`
	Date        time.Time             `json:"date"`
	Attachments []MailAttachmentModel `json:"attachments"`
}

type MailAttachmentModel struct {
	Mime string `json:"MIME"`
	Name string `json:"name"`
}
