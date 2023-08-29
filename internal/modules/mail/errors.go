package mail

import "errors"

var ErrEmptyMailbox = errors.New("mailbox is empty")
var ErrNotFound = errors.New("not found messgaes")
