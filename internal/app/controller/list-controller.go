package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/fredmayer/mail-parser-rest/internal/app/service"
	"github.com/fredmayer/mail-parser-rest/internal/app/types"
	"github.com/labstack/echo/v4"
)

type ListController struct {
	ctx     context.Context
	service *service.Manager
	//todo add service inject
}

func NewListController(ctx context.Context, service *service.Manager) *ListController {
	return &ListController{
		ctx:     ctx,
		service: service,
	}
}

func (lc *ListController) GetList(ctx echo.Context) error {
	var rq types.PageRequest
	err := ctx.Bind(&rq)
	//page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		log.Printf("error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	res, err := lc.service.MailService.GetList(rq.Page, ctx)
	if err != nil {
		log.Fatal(err)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (lc *ListController) GetView(ctx echo.Context) error {
	uidStr := ctx.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		log.Printf("error - %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	res, err := lc.service.MailService.GetView(uid, ctx)
	if err != nil {
		log.Printf("error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (lc *ListController) DownloadAttachment(ctx echo.Context) error {
	var params types.AttachmentRequest
	err := ctx.Bind(&params)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	uidStr := ctx.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		log.Printf("error - %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	reader, err := lc.service.MailService.DownloadAttachment(uid, params, ctx)
	if err != nil {
		echo.NewHTTPError(http.StatusNotFound, err)
	}

	ctx.Response().Writer.Header().Set("Content-Disposition", "attachment; filename="+params.Name)
	//w.Header().Set("Content-Disposition", "attachment; filename=WHATEVER_YOU_WANT")
	//w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	return ctx.Stream(http.StatusOK, params.Mime, reader)
}
