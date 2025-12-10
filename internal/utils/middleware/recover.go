package middleware

import (
	"github.com/gofiber/fiber/v2"
	"k071123/tools/logger"
	"runtime"
)

func PanicRecovery(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("PANIC: %v", r)

				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				log.Errorf("STACKTRACE: %s", buf[:n])

				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "internal",
					"message": "Internal server error",
				})
			}
		}()

		return c.Next()
	}
}
