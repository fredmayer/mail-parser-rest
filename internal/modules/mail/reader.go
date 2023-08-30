package mail

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"regexp"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"github.com/fredmayer/mail-parser-rest/internal/configs"
	"github.com/fredmayer/mail-parser-rest/pkg/logging"
	"github.com/maxjust/charmap"
)

type MailReader struct {
	Cl            *client.Client
	CurrentFolder *imap.MailboxStatus
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

func (mr *MailReader) Move(uid uint32, mailbox string) error {
	seqset := new(imap.SeqSet)
	seqset.AddNum(uid)

	return mr.Cl.UidMove(seqset, "LOADED")
}

func (mr *MailReader) DownloadAttachment(uid int, mime string, name string) (io.Reader, error) {
	mbox := mr.initMailbox()
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
	mbox := mr.initMailbox()
	if mbox.Messages == 0 {
		return nil, ErrEmptyMailbox
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

	if msg == nil {
		return nil, ErrNotFound
	}
	from, err := mc.getFirstAddress(msg.Envelope.From)
	if err != nil {
		from = "not setted"
	}
	subject, err := mr.decodeSubject(msg.Envelope.Subject)
	if err != nil {
		subject = msg.Envelope.Subject
	}
	mailDto := MailDto{
		MessageId: msg.Envelope.MessageId,
		Uid:       msg.Uid,
		SeqNum:    msg.SeqNum,
		From:      from,
		Subject:   subject,
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
	mbox := mr.initMailbox()

	r := &ListMailResponse{}
	if mbox.Messages == 0 {
		logging.Log().Warn(ErrEmptyMailbox)
		//log.Println("No messages in mailbox")
		return r, nil
	}
	logging.Log().Debugf("total messages %v", mbox.Messages)

	from := uint32((page-1)*perPage + 1)
	to := uint32(page * perPage)
	if to > mbox.Messages {
		to = mbox.Messages
	}

	messages := make(chan *imap.Message, 10)
	pages := math.Ceil(float64(mbox.Messages) / float64(perPage))

	if page > int(pages) {
		return nil, ErrNotFound
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	done := make(chan error, 1)
	go func() {
		done <- mr.Cl.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, messages)
	}()

	res, err := mr.fillCollection(messages)
	if err != nil {
		logging.Log().Error(err)
	}

	if err := <-done; err != nil {
		logging.Log().Panic(err)
	}

	r.Data = res
	r.Pages = int(pages)

	return r, nil
}

func (mr *MailReader) Last(count int) (*ListMailResponse, error) {
	mbox := mr.initMailbox()

	r := &ListMailResponse{}
	if mbox.Messages == 0 {
		logging.Log().Warn(ErrEmptyMailbox)
		//log.Println("No messages in mailbox")
		return r, nil
	}
	logging.Log().Debugf("total messages %v", mbox.Messages)

	from := mbox.Messages
	to := from - uint32(count)

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	done := make(chan error, 1)
	messages := make(chan *imap.Message, 10)
	go func() {
		done <- mr.Cl.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, messages)
	}()

	res, err := mr.fillCollection(messages)
	if err != nil {
		logging.Log().Error(err)
	}

	if err := <-done; err != nil {
		logging.Log().Panic(err)
	}

	r.Data = res
	r.Pages = 0

	return r, nil
}

/**
* 	Fill collection messages array
 */
func (mr *MailReader) fillCollection(messages chan *imap.Message) ([]MailDto, error) {
	res := make([]MailDto, 0, 10)
	for msg := range messages {
		emailFrom := "not setted"
		if msg.Envelope.From != nil {
			ef, err := mc.getFirstAddress(msg.Envelope.From)
			if err == nil {
				emailFrom = ef
			}
		}

		subject, err := mr.decodeSubject(msg.Envelope.Subject)
		if err != nil {
			subject = msg.Envelope.Subject
		}

		res = append(res, MailDto{
			MessageId: msg.Envelope.MessageId,
			From:      emailFrom,
			Subject:   subject,
			Date:      msg.Envelope.Date,
			Uid:       msg.Uid,
			SeqNum:    msg.SeqNum,
		})
	}

	return res, nil
}

func (mr *MailReader) getFirstAddress(froms []*imap.Address) (string, error) {
	if len(froms) == 0 {
		logging.Log().Warn("not found \"from\" e-mail address")
		return "", errors.New("Not found froms e-mail adresses")
	}
	return froms[0].Address(), nil
}

func (mr *MailReader) decodeSubject(s string) (string, error) {
	// smple regexp = =\?\w*\d?\W?\w?\d{0,4}\?B\?.*=*\?=
	re := regexp.MustCompile(`=\?\w*\d?\W?\w?\d{0,4}\?B\?(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?\?=`)

	matches := re.FindAllString(s, -1)
	if len(matches) > 0 {
		var decodeMatches []string
		for _, ms := range matches {
			//decode str
			res, err := mr.decodeString(ms)
			if err == nil {
				decodeMatches = append(decodeMatches, res)
			}
		}

		if len(decodeMatches) > 0 {
			var rs string
			for i, d := range decodeMatches {
				rs += fmt.Sprintf("${%d}%v", i+1, d)
			}

			res := re.ReplaceAll([]byte(s), []byte(rs))

			return string(res[:]), nil
		}
	}

	if strings.Contains(s, "=?") {
		//fmt.Println(s)
		//s = RFC2047.Decode(s)
		res, err := mr.decodeString(s)
		if err != nil {
			logging.Log().Error(err)
		}
		s = res
	}

	return s, nil
}

func (mr *MailReader) decodeString(s string) (string, error) {

	dec := new(mime.WordDecoder)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "koi8-r":
			content, err := io.ReadAll(input)
			if err != nil {
				return nil, err
			}

			return bytes.NewReader(charmap.KOI8R_to_UTF8(content)), nil

		case "windows-1251":
			content, err := io.ReadAll(input)
			if err != nil {
				return nil, err
			}

			return bytes.NewReader(charmap.CP1251_to_UTF8(content)), nil
		default:
			logging.Log().Errorf("unhandled charset %q", charset)
			return nil, fmt.Errorf("unhandled charset %q", charset)
		}
	}

	res, err := dec.Decode(string(s[:]))
	//res, err := dec.Decode("79PUwdTLyQ")
	if err != nil {
		logging.Log().WithField("string", s).Error(err)
		return s, err
	}
	return res, nil
}
