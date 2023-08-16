package controller

import (
	"context"
	"log"
	"net/http"

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
		log.Printf("error - %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	res, err := lc.service.MailService.GetList(rq.Page)
	if err != nil {
		log.Fatal(err)
	}

	return ctx.JSON(http.StatusOK, res)
}
