
package location

import (
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func Routes(app *fiber.App, collections map[string]*mongo.Collection) {
    service := newService(collections)
    handler := Handler{service}

    apiV1 := app.Group("/api/v1")
    locations := apiV1.Group("/locations")

    locations.Post("/", handler.Createlocation)
    locations.Get("/", handler.Getlocations)
    locations.Post("/nearby", handler.GetNearbylocations)
    locations.Get("/:id", handler.Getlocation)
    locations.Patch("/:id", handler.UpdatePartiallocation)
    locations.Delete("/:id", handler.Deletelocation)
}
