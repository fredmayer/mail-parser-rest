package mail

import "time"

type ListMailDto struct {
	MessageId string
	Uid       uint32
	SeqNum    uint32
	From      string
	Subject   string
	Date      time.Time
}

type ListMailResponse struct {
	Data  []ListMailDto
	Pages int
}
