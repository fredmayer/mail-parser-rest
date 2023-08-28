package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/fredmayer/mail-parser-rest/internal/app/service"
	"github.com/fredmayer/mail-parser-rest/internal/app/types"
	"github.com/fredmayer/mail-parser-rest/pkg/logging"
	"github.com/labstack/echo/v4"
)

type MessageController struct {
	ctx     context.Context
	service *service.Manager
	//todo add service inject
}

func NewMessageController(ctx context.Context, service *service.Manager) *MessageController {
	return &MessageController{
		ctx:     ctx,
		service: service,
	}
}

func (lc *MessageController) GetList(ctx echo.Context) error {
	var rq types.PageRequest
	err := ctx.Bind(&rq)
	//page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		logging.Log().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	res, err := lc.service.MailService.GetList(rq.Page, ctx)
	if err != nil {
		logging.Log().Error(err.Error())
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	logging.Log().Debug(ctx.Path())
	return ctx.JSON(http.StatusOK, res)
}

func (lc *MessageController) GetView(ctx echo.Context) error {
	uidStr := ctx.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		logging.Log().Warn(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	res, err := lc.service.MailService.GetView(uid, ctx)
	if err != nil {
		log.Printf("error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (lc *MessageController) Move(ctx echo.Context) error {
	uidStr := ctx.QueryParam("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		log.Printf("error - %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	mailbox := ctx.QueryParam("mailbox")
	if mailbox == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("mailbox param is required"))
	}

	err = lc.service.MailService.Move(uid, mailbox)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (lc *MessageController) DownloadAttachment(ctx echo.Context) error {
	indexStr := ctx.QueryParam("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		log.Printf("error - %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	uidStr := ctx.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		log.Printf("error - %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	res, err := lc.service.MailService.DownloadAttachment(uid, index, ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	ctx.Response().Writer.Header().Set("Content-Disposition", "attachment; filename="+res.Name)

	return ctx.Stream(http.StatusOK, res.Mime, res.GetReader())
}
