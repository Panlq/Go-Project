package main

import (
	"time"

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
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// parameters
	// app.Get("/:value", func(c *fiber.Ctx) error {
	// 	return c.SendString("value: " + c.Params("value"))
	// })

	// options key
	app.Get("/user/:name?", func(c *fiber.Ctx) error {
		if c.Params("name") != "" {
			return c.SendString("Hello " + c.Params("name"))
		}

		return c.SendString("Where is john?")
	})

	// static file
	app.Static("/", "./public")

	// new error
	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(5500, "Custom error message")
	})

	app.Listen(":3080")
}
