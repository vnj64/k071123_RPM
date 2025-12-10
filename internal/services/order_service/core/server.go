package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"k071123/internal/services/order_service/domain"
	"k071123/internal/services/order_service/services/config"
	"k071123/internal/utils/middleware"
	"runtime"
	"strings"
)

type HttpServer struct {
	app *fiber.App
	ctx domain.Context
}

type Server interface {
	Start()
	App() *fiber.App
}

func NewHttpServer() Server {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
	})

	ctx := InitCtx()

	var methods = []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete}
	var headers = []string{
		fiber.HeaderAccept,
		fiber.HeaderAuthorization,
		fiber.HeaderContentType,
		fiber.HeaderContentLength,
		fiber.HeaderAcceptEncoding,
	}

	corsConfig := cors.New(cors.Config{
		AllowOrigins: strings.Join([]string{
			"http://localhost:7803/", "http://localhost:8129/",
		}, ", "),
		AllowMethods:     strings.Join(methods, ", "),
		AllowHeaders:     strings.Join(headers, ", "),
		AllowCredentials: true,
		MaxAge:           300,
	})

	app.Use(middleware.LoggerMiddleware(ctx.Services().Logger()))
	app.Use(middleware.PanicRecovery(ctx.Services().Logger()))

	app.Use(corsConfig)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("context", ctx)
		return c.Next()
	})

	return &HttpServer{
		app: app,
		ctx: ctx,
	}
}

func (s *HttpServer) Start() {
	cfg := config.Make()

	runtime.GOMAXPROCS(runtime.NumCPU())
	err := s.app.Listen(":" + cfg.HttpPort())
	if err != nil {
		panic("http server isn't started successfully")
	}
}

func (s *HttpServer) App() *fiber.App {
	return s.app
}
