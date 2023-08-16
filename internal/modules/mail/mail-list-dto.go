package mail

import "time"

type ListMailDto struct {
	MessageId string
	From      string
	Subject   string
	Date      time.Time
}
