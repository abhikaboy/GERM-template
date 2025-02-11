
package location

import (
    "errors"
    "log/slog"
    "time"

    "github.com/abhikaboy/GERM-template/internal/xerr"
    "github.com/abhikaboy/GERM-template/internal/xvalidator"
    go_json "github.com/goccy/go-json"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
    service *Service
}

func (h *Handler) Createlocation(c *fiber.Ctx) error {
    ctx := c.Context()
    var location locationDocument
    var params CreatelocationParams

    slog.LogAttrs(ctx, slog.LevelInfo, "Inserting location")
    err := c.BodyParser(&params)
    if err != nil {
        return err
    }
    errs := xvalidator.Validator.Validate(params)
    if len(errs) > 0 {
        return c.Status(fiber.StatusBadRequest).JSON(errs)
    }

    location = locationDocument{
        Field1:    params.Field1,
        Field2:    params.Field2,
        Location:  params.Location,
        Picture:   params.Picture,
        Timestamp: time.Now(),
        ID:        primitive.NewObjectID(),
    }

    result, err := h.service.Insertlocation(location)
    if err != nil {
        sErr := err.(mongo.WriteException)
        if sErr.HasErrorCode(121) {
            return xerr.WriteException(c, sErr)
        }
    }

    return c.JSON(result)
}

func (h *Handler) Getlocations(c *fiber.Ctx) error {
    locations, err := h.service.GetAlllocations()
    if err != nil {
        return err
    }
    return c.JSON(locations)
}

func (h *Handler) Getlocation(c *fiber.Ctx) error {
    id, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(xerr.BadRequest(err))
    }

    location, err := h.service.GetlocationByID(id)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.Status(fiber.StatusNotFound).JSON(xerr.NotFound("location", "id", id.Hex()))
        }
        return err
    }
    return c.JSON(location)
}

func (h *Handler) UpdatePartiallocation(c *fiber.Ctx) error {
    id, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(xerr.BadRequest(err))
    }

    var partialUpdate UpdatelocationDocument
    if err := go_json.Unmarshal(c.Body(), &partialUpdate); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(xerr.InvalidJSON())
    }

    err = h.service.UpdatePartiallocation(id, partialUpdate)
    if err != nil {
        return err
    }

    return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Deletelocation(c *fiber.Ctx) error {
    id, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(xerr.BadRequest(err))
    }

    if err := h.service.Deletelocation(id); err != nil {
        return err
    }

    return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetNearbylocations(c *fiber.Ctx) error {
    var params GetNearbylocationsParams

    err := c.BodyParser(&params)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(xerr.InvalidJSON())
    }

    locations, err := h.service.GetNearbylocations(params.Location, params.Radius)
    if err != nil {
        return err
    }

    return c.JSON(locations)
}
