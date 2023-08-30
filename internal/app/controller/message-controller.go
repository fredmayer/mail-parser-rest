package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/fredmayer/mail-parser-rest/internal/app/service"
	"github.com/fredmayer/mail-parser-rest/internal/app/types"
	"github.com/fredmayer/mail-parser-rest/internal/modules/mail"
	"github.com/fredmayer/mail-parser-rest/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
		return badRequestError(ctx.Path(), err, logrus.Fields{})
	}

	res, err := lc.service.MailService.GetList(rq.Page, ctx)
	if err != nil {
		if errors.Is(err, mail.ErrNotFound) {
			return notFoundError(ctx.Path(), err, logrus.Fields{})
		}
		logging.Log().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	logging.Log().Debug(ctx.Path())
	return ctx.JSON(http.StatusOK, res)
}

func (lc *MessageController) GetLast(ctx echo.Context) error {
	countStr := ctx.QueryParam("count")
	count := 100
	err := errors.New("dummy")
	if countStr != "" {
		if count, err = strconv.Atoi(countStr); err != nil {
			return badRequestError(ctx.Path(), err, logrus.Fields{
				"count": countStr,
			})
		}
	}

	res, err := lc.service.MailService.GetLast(count, ctx)
	if err != nil {
		if errors.Is(err, mail.ErrNotFound) {
			return notFoundError(ctx.Path(), err, logrus.Fields{})
		}
		logging.Log().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	logging.Log().Debug(ctx.Path())
	return ctx.JSON(http.StatusOK, res)
}

func (lc *MessageController) GetView(ctx echo.Context) error {
	uidStr := ctx.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return badRequestError(ctx.Path(), errors.New("uid must be a integer"), logrus.Fields{"uid": uidStr})
	}

	res, err := lc.service.MailService.GetView(uid, ctx)
	if err != nil {

		//Not found
		if errors.Is(err, mail.ErrNotFound) {
			return notFoundError(ctx.Path(), err, logrus.Fields{"uid": uidStr})
		}

		//Custom errors
		logging.Log().WithFields(logrus.Fields{
			"uid":    uidStr,
			"status": http.StatusNotFound,
			"msg":    err.Error(),
		}).Error(ctx.Path())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	logging.Log().WithFields(logrus.Fields{
		"uid": uidStr,
	}).Debug(ctx.Path())
	return ctx.JSON(http.StatusOK, res)
}

func (lc *MessageController) Move(ctx echo.Context) error {
	uidStr := ctx.QueryParam("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return badRequestError(ctx.Path(), errors.New("uid must be a integer"), logrus.Fields{
			"uid": uidStr,
		})
	}

	mailbox := ctx.QueryParam("mailbox")
	if mailbox == "" {
		return badRequestError(ctx.Path(), errors.New("mailbox param is required"), logrus.Fields{
			"uid": uidStr,
		})
	}

	err = lc.service.MailService.Move(uid, mailbox)
	if err != nil {
		return notFoundError(ctx.Path(), err, logrus.Fields{"uid": uidStr, "mailbox": mailbox})
	}

	logging.Log().WithFields(logrus.Fields{
		"uid":     uidStr,
		"mailbox": mailbox,
	}).Debug(ctx.Path())
	return ctx.NoContent(http.StatusNoContent)
}

func (lc *MessageController) DownloadAttachment(ctx echo.Context) error {
	indexStr := ctx.QueryParam("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return badRequestError(ctx.Path(), errors.New("index must be a integer"), logrus.Fields{
			"index": indexStr,
		})
	}

	uidStr := ctx.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return badRequestError(ctx.Path(), errors.New("uid must be a integer"), logrus.Fields{
			"index": indexStr,
			"uid":   uidStr,
		})
	}

	res, err := lc.service.MailService.DownloadAttachment(uid, index, ctx)
	if err != nil {
		return notFoundError(ctx.Path(), err, logrus.Fields{
			"index": indexStr,
			"uid":   uidStr,
		})
	}

	ctx.Response().Writer.Header().Set("Content-Disposition", "attachment; filename="+res.Name)

	logging.Log().WithFields(logrus.Fields{
		"index": indexStr,
		"uid":   uidStr,
	}).Debug(ctx.Path())
	return ctx.Stream(http.StatusOK, res.Mime, res.GetReader())
}
