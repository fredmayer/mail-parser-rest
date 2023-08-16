package mail

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/fredmayer/mail-parser-rest/internal/configs"
)

type MailReader struct {
	Cl *client.Client
}

var (
	mc MailReader
)

func Dial() *MailReader {
	cfg := configs.Get()

	c, err := client.Dial(cfg.Imap + ":" + cfg.ImapPort)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Login(cfg.ImapUser, cfg.ImapPassword); err != nil {
		log.Fatal(err)
	}
	//defer c.Logout()

	log.Println("Connected " + cfg.ImapUser)

	mc = MailReader{
		Cl: c,
	}

	return &mc
}

func Get() *MailReader {
	return &mc
}

func (mr *MailReader) List(page int) []ListMailDto {
	perPage := 10
	if page == 0 {
		page = 1
	}
	mbox, err := mr.Cl.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	if mbox.Messages == 0 {
		log.Println("No messages in mailbox")
		//return nil
	}
	log.Printf("Total messages %v \r\n", mbox.Messages)

	from := uint32((page-1)*perPage + 1)
	to := uint32(page * perPage)
	log.Printf("Page %d from %d to %d", page, from, to)

	messages := make(chan *imap.Message, 10)
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	done := make(chan error, 1)
	go func() {
		done <- mr.Cl.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	res := make([]ListMailDto, 0, 10)
	for msg := range messages {
		i := 0
		from := ""
		for f := range msg.Envelope.From {
			//fr[i] = f.MailboxName
			from = msg.Envelope.From[f].Address()
			i++
		}

		res = append(res, ListMailDto{
			MessageId: msg.Envelope.MessageId,
			From:      from,
			Subject:   msg.Envelope.Subject,
			Date:      msg.Envelope.Date,
		})
	}

	return res
}
