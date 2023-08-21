package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fredmayer/mail-parser-rest/internal/app/controller"
	"github.com/fredmayer/mail-parser-rest/internal/app/service"
	"github.com/fredmayer/mail-parser-rest/internal/configs"
	"github.com/fredmayer/mail-parser-rest/internal/modules/mail"
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
	log.Println("Starting...")
	flag.Parse()

	config := configs.NewConfig(configPath)
	fmt.Println(config)

	_ = mail.Dial()

	// Initialize Echo instance
	e := echo.New()

	// Init service manager
	serviceManager, err := service.NewManager()
	if err != nil {
		log.Fatalln(err)
	}

	//Controllers
	listController := controller.NewListController(ctx, serviceManager)

	//Routes
	e.GET("/mailboxes", listController.GetMailBoxes)
	e.GET("/list", listController.GetList)

	e.GET("/view/:uid", listController.GetView)
	e.POST("/download/:uid", listController.DownloadAttachment)

	// Start server
	s := &http.Server{
		Addr:         config.HTTPAddr,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
	//log.Printf("Server started at %v", config.HTTPAddr)
}
