package mail

import (
	"errors"
	"io"
	"log"
	"math"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
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

func (mr *MailReader) MailBoxes() ([]string, error) {
	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- mr.Cl.List("", "*", mailboxes)
	}()

	res := make([]string, 0, 10)
	for m := range mailboxes {
		res = append(res, m.Name)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return res, nil
}

func (mr *MailReader) Move(uid uint32, mailbox string) error {
	seqset := new(imap.SeqSet)
	seqset.AddNum(uid)

	return mr.Cl.UidMove(seqset, "LOADED")
}

func (mr *MailReader) DownloadAttachment(uid int, mime string, name string) (io.Reader, error) {
	//TODO transfer to another function
	mbox, err := mr.Cl.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	if mbox.Messages == 0 {
		log.Println("No messages in mailbox")
		//todo return error
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uint32(uid))

	messages := make(chan *imap.Message, 10)

	done := make(chan error, 1)
	go func() {
		done <- mr.Cl.UidFetch(seqset, []imap.FetchItem{imap.FetchRFC822, imap.FetchEnvelope}, messages)
	}()

	msg := <-messages

	for i, r := range msg.Body {
		log.Println(i)
		entity, err := message.Read(r)
		if err != nil {
			log.Fatal(err)
		}

		multipartReader := entity.MultipartReader()
		for e, err := multipartReader.NextPart(); err != io.EOF; e, err = multipartReader.NextPart() {
			kind, params, cErr := e.Header.ContentType()
			if cErr != nil {
				log.Fatal(cErr)
			}

			nameAt, ok := params["name"]
			if ok {
				if kind == mime && nameAt == name {
					return e.Body, nil
				}
			}
		}
	}

	return nil, errors.New("Not found attachment")

}

// view message by sid
func (mr *MailReader) GetBySid(sid int) (*MailDto, error) {
	mbox, err := mr.Cl.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	if mbox.Messages == 0 {
		log.Println("No messages in mailbox")
		//return nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uint32(sid))

	messages := make(chan *imap.Message, 10)

	done := make(chan error, 1)
	go func() {
		done <- mr.Cl.UidFetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBodyStructure, imap.FetchUid}, messages)
	}()

	msg := <-messages

	from, err := mc.getFirstAddress(msg.Envelope.From)
	if err != nil {
		from = "not setted"
	}
	mailDto := MailDto{
		MessageId: msg.Envelope.MessageId,
		Uid:       msg.Uid,
		SeqNum:    msg.SeqNum,
		From:      from,
		Subject:   msg.Envelope.Subject,
		Date:      msg.Envelope.Date,
	}

	attachments := []MailAttachmentDto{}
	i := 0
	for _, part := range msg.BodyStructure.Parts {
		fname, err := part.Filename()
		if err != nil || len(fname) == 0 {
			continue
		}

		attachments = append(attachments, MailAttachmentDto{
			Mime:  part.MIMEType + "/" + part.MIMESubType,
			Name:  fname,
			Index: i,
		})
		i++
	}

	mailDto.Attachments = attachments

	if err := <-done; err != nil {
		//Todo add normal error 404
		log.Fatal(err)
	}

	return &mailDto, nil
}

// List of emails
func (mr *MailReader) List(page int) (*ListMailResponse, error) {
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
	//log.Printf("Page %d from %d to %d", page, from, to)

	messages := make(chan *imap.Message, 10)
	pages := math.Ceil(float64(mbox.Messages) / float64(perPage))

	if page > int(pages) {
		return nil, errors.New("Not found")
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	done := make(chan error, 1)
	go func() {
		done <- mr.Cl.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, messages)
	}()

	res := make([]MailDto, 0, 10)
	for msg := range messages {
		from, err := mc.getFirstAddress(msg.Envelope.From)
		if err != nil {
			from = "not setted"
		}
		//log.Println(msg.Uid)

		res = append(res, MailDto{
			MessageId: msg.Envelope.MessageId,
			From:      from,
			Subject:   msg.Envelope.Subject,
			Date:      msg.Envelope.Date,
			Uid:       msg.Uid,
			SeqNum:    msg.SeqNum,
		})
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	r := &ListMailResponse{
		Data:  res,
		Pages: int(pages),
	}

	return r, nil
}

func (mr *MailReader) getFirstAddress(froms []*imap.Address) (string, error) {
	if len(froms) == 0 {
		return "", errors.New("Not found froms e-mail adresses")
	}
	return froms[0].Address(), nil
}
