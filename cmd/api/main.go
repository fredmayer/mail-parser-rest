package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/fredmayer/mail-parser-rest/internal/app/controller"
	"github.com/fredmayer/mail-parser-rest/internal/app/service"
	"github.com/fredmayer/mail-parser-rest/internal/configs"
	"github.com/fredmayer/mail-parser-rest/internal/modules/mail"
	"github.com/fredmayer/mail-parser-rest/pkg/logging"
	"github.com/labstack/echo/v4"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.toml", "path to config file")
}

func main() {
	ctx := context.Background()
	flag.Parse()

	config := configs.NewConfig(configPath)
	logging.Init(config.LogLevel)
	logging.Log().Info("Starting application")

	_ = mail.Dial()

	// Initialize Echo instance
	e := echo.New()

	// Init service manager
	serviceManager, err := service.NewManager()
	if err != nil {
		log.Fatalln(err)
	}

	//Controllers
	MessageController := controller.NewMessageController(ctx, serviceManager)
	MailBoxController := controller.NewMailBoxController(ctx, serviceManager)

	//Routes
	m := e.Group("/messages")
	m.GET("/list", MessageController.GetList)
	m.GET("/:uid", MessageController.GetView)
	m.GET("/download/:uid", MessageController.DownloadAttachment)
	m.GET("/move", MessageController.Move)

	mls := e.Group("/mails")
	mls.GET("/list", MailBoxController.GetList)
	mls.PUT("/folder", MailBoxController.SetFolder)

	// Start server
	s := &http.Server{
		Addr:         config.HTTPAddr,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
	//log.Printf("Server started at %v", config.HTTPAddr)
}
