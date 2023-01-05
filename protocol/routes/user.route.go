package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nanpipat/golang-template-hexagonal/internal/handlers"
)

func UserRoutes(api fiber.Router, handler *handlers.UserHandlers) {
	api.Post("/users", handler.Create)
	api.Get("users", handler.Pagination)
	api.Get("/users/:id", handler.Get)
	api.Put("/users/:id", handler.Update)
	api.Delete("/users/:id", handler.Delete)
}
