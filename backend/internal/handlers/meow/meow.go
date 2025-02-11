
package Meow

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service *Service
}

func (h *Handler) CreateMeow(c *fiber.Ctx) error {
	var params CreateMeowParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
		})
	}

	doc := MeowDocument{
		ID:        primitive.NewObjectID(),
		Field1:    params.Field1,
		Field2:    params.Field2,
		Picture:   params.Picture,
		Timestamp: time.Now(),
	}

	_, err := h.service.CreateMeow(&doc); 
    if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create Meow",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(doc)
}

func (h *Handler) GetMeows(c *fiber.Ctx) error {
	Meows, err := h.service.GetAllMeows()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch Meows",
		})
	}

	return c.JSON(Meows)
}

func (h *Handler) GetMeow(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	Meow, err := h.service.GetMeowByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Meow not found",
		})
	}

	return c.JSON(Meow)
}

func (h *Handler) UpdatePartialMeow(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var update UpdateMeowDocument
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.service.UpdatePartialMeow(id, update); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update Meow",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteMeow(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	if err := h.service.DeleteMeow(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Meow",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
