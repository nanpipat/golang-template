package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nanpipat/golang-template-hexagonal/helper"
	"github.com/nanpipat/golang-template-hexagonal/internal/core/domain"
	"github.com/nanpipat/golang-template-hexagonal/internal/core/services"
)

type UserHandlers struct {
	service services.IUserServiceInterface
}

func NewUserHandlers(service services.IUserServiceInterface) *UserHandlers {
	return &UserHandlers{
		service: service,
	}
}

func (m *UserHandlers) Create(c *fiber.Ctx) error {
	payload := domain.User{}

	err := c.BodyParser(&payload)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(err.Error())
	}

	result, err := m.service.Create(payload)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (m *UserHandlers) Pagination(c *fiber.Ctx) error {
	result, err := m.service.Pagination(helper.GetPageOptions(c))
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (m *UserHandlers) Get(c *fiber.Ctx) error {
	result, err := m.service.Get(c.Params("id"))
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (m *UserHandlers) Update(c *fiber.Ctx) error {
	payload := domain.User{}

	err := c.BodyParser(&payload)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(err.Error())
	}

	result, err := m.service.Update(c.Params("id"), &payload)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (m *UserHandlers) Delete(c *fiber.Ctx) error {
	err := m.service.Delete(c.Params("id"))
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
