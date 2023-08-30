# What is this repository for? 

This is a REST API microservice of the mail client. To automate the receipt of data through mail attachments. For example: you are sent price lists by email with attached price lists and you need to parse the attachments.

This microservice based on high performance, extensible, minimalist Go web framework [Echo](https://echo.labstack.com/), IMAP client [go-imap](https://github.com/emersion/go-imap). For logging i`m use [Logrus](https://github.com/sirupsen/logrus)

## Config

The entire configuration is in a file `config.toml`

  - `imap`     - IMAP host address
  - `imap_port` - IMAP port
  - `imap_user` - username
  - `imap_password` - password
  - `http_addr` - server bind address and port
  - `log_level` - Level of logs: debug | warn | error

## Build

Before, you need a build programm or use a docker (See the corresponding section).

    make build

The programm has been build `app` application.

## Usage

The following REST API is used:

> [!NOTE]
> The following logic should be implemented in your application:
> 1. We read the messages in order and check for the from field, subject and attachment name
> 2. If the message suits us, then we get an attachment (blob)
> 3. Transfer the letter to one of the folders (LOADED|ERRORS). Well done!
> 4. Repeat

### Manage mailboxes (folders)

1. `GET host/mails/list` - response list of mailboxes (folders).
2. `PUT host/mails/folder` with JSON data: ` {"folder":"FOLDER_NAME"} ` - set current folder of mailbox.

### List of messages
1. `GET host/messages/list?page=1` - list messages with response header info of total pages `X-Pagination-Page-Count`
2. `GET host/messages/last?count=100` - list of last messages. Use query parameter `count`- что бы указать кол-во элементов

### View message (with attachments names)
`GET host/messages/:uid` - where `uid` - id message from list.
Response data of message with array `attachments`. 

### Download attachment
`GET messages/download/:uid?index=0`

- `:uid` - id message from list
- `index` - index attachments from view attachment

**Response:** - blob data

## Feature tasks list

- [ ] Use secure TLS connection IMAP  
- [ ] Migrate go-imap to v2
- [ ] Cover tests with Postman collections

## Show your support

Give a ⭐️ if this project helped you!

**Made with ❤️ :)** Thank you! 