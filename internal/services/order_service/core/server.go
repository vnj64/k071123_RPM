package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"k071123/internal/services/order_service/domain"
	"k071123/internal/services/order_service/services/config"
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
	app := fiber.New()
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
		AllowCredentials: true, // Убедимся, что можно передавать куки и авторизационные заголовки
		MaxAge:           300,
	})

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
		panic("http server inst start successfully")
	}
}

func (s *HttpServer) App() *fiber.App {
	return s.app
}
