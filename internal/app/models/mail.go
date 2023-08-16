package models

import "time"

type MailModel struct {
	MessageId string    `json:"id"`
	From      string    `json:"from"`
	Subject   string    `json:"subject"`
	Date      time.Time `json:"date"`
}
