package controller

import (
	"context"
	"net/http"

	"github.com/fredmayer/mail-parser-rest/internal/app/service"
	"github.com/labstack/echo/v4"
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
		echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, res)
}
