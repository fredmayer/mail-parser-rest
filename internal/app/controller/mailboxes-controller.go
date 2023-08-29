package controller

import (
	"context"
	"net/http"

	"github.com/fredmayer/mail-parser-rest/internal/app/service"
	"github.com/fredmayer/mail-parser-rest/internal/app/types"
	"github.com/fredmayer/mail-parser-rest/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type MailBoxController struct {
	ctx     context.Context
	service *service.Manager
	//todo add service inject
}

func NewMailBoxController(ctx context.Context, service *service.Manager) *MailBoxController {
	return &MailBoxController{
		ctx:     ctx,
		service: service,
	}
}

func (mb *MailBoxController) GetList(ctx echo.Context) error {
	res, err := mb.service.MailService.MailBoxes()
	if err != nil {
		return notFoundError(ctx.Path(), err, logrus.Fields{})
	}

	logging.Log().Debug(ctx.Path())
	return ctx.JSON(http.StatusOK, res)
}

func (mb *MailBoxController) SetFolder(ctx echo.Context) error {
	var req types.FolderRequest
	err := ctx.Bind(&req)
	if err != nil {
		return badRequestError(ctx.Path(), err, logrus.Fields{})
	}

	err = mb.service.MailService.SetFolder(req.Folder)
	if err != nil {
		return notFoundError(ctx.Path(), err, logrus.Fields{})
	}

	logging.Log().Debug(ctx.Path())
	return ctx.NoContent(http.StatusCreated)
}
