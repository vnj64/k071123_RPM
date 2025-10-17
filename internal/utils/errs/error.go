package errs

import "github.com/gofiber/fiber/v2"

// swagger:models Error
type Error struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type Code int

const (
	InternalServerError Code = 500 // Внутренняя ошибка сервера
	BadRequest          Code = 400 // Некорректный запрос
	Unauthorized        Code = 401 // Неавторизованный доступ
	Forbidden           Code = 403 // Запрещено
	NotFound            Code = 404 // Не найдено
	MethodNotAllowed    Code = 405 // Метод не разрешён
	Conflict            Code = 409 // Конфликт (например, дублирующая запись)
	UnprocessableEntity Code = 422 // Невозможность обработать запрос (например, валидация)
	TooManyRequests     Code = 429 // Слишком много запросов
	ServiceUnavailable  Code = 503 // Сервис временно недоступен
	Timeout             Code = 504 // Таймаут запроса
)

func NewError(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewErrorWithDetails(err *Error, details string) *Error {
	return &Error{
		Code:    err.Code,
		Message: err.Message,
		Details: details,
	}
}

var (
	ErrInternalServerError = NewError(InternalServerError, "Internal server error")
	ErrBadRequest          = NewError(BadRequest, "Bad request")
	ErrUnauthorized        = NewError(Unauthorized, "Unauthorized")
	ErrForbidden           = NewError(Forbidden, "Forbidden")
	ErrNotFound            = NewError(NotFound, "Not found")
	ErrMethodNotAllowed    = NewError(MethodNotAllowed, "Method not allowed")
	ErrConflict            = NewError(Conflict, "Conflict")
	ErrUnprocessableEntity = NewError(UnprocessableEntity, "Unprocessable entity")
	ErrTooManyRequests     = NewError(TooManyRequests, "Too many requests")
	ErrServiceUnavailable  = NewError(ServiceUnavailable, "Service unavailable")
	ErrTimeout             = NewError(Timeout, "Timeout error")
)

func (e *Error) Error() string {
	return e.Message
}

func SendSuccess(ctx *fiber.Ctx, statusCode int, data interface{}) error {
	return ctx.Status(statusCode).JSON(data)
}

func SendError(ctx *fiber.Ctx, err error) error {
	if customErr, ok := err.(*Error); ok {
		return ctx.Status(int(customErr.Code)).JSON(fiber.Map{
			"error":   customErr.Message,
			"details": customErr.Details,
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal server error",
	})
}
