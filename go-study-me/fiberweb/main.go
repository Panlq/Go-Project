package main

import (
	"time"

	"github/panlq/gostd/fiberweb/upload"

	"github.com/go-logr/logr"
	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Config struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

type HttpServer struct {
	logger logr.Logger
	app    *fiber.App
	cfg    *Config
}

func NewHttpServer(logger logr.Logger, cfg *Config) *HttpServer {
	app := fiber.New(fiber.Config{
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		Prefork:      false,
	})

	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(fiberLog.New())

	return &HttpServer{
		logger: logger,
		app:    app,
		cfg:    cfg,
	}
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: fiber.DefaultBodyLimit,
		Prefork:   false,
	})

	// 文件上传相关路由
	app.Get("/upload/check", upload.HandleCheck)
	app.Post("/upload/chunk", upload.HandleUploadChunk)
	app.Post("/upload/merge", upload.HandleMerge)

	// 静态文件服务
	app.Static("/", "./public")

	app.Listen(":3080")
}
