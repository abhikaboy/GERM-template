
package Meow
import (
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

/*
Router maps endpoints to handlers
*/
func Routes(app *fiber.App, collections map[string]*mongo.Collection) {
    service := newService(collections)
    handler := Handler{service}

    // Add a group for API versioning
    apiV1 := app.Group("/api/v1")

    // Add Sample group under API Version 1
    Meows := apiV1.Group("/Meows")

    Meows.Post("/", handler.CreateMeow)
    Meows.Get("/", handler.GetMeows)
    Meows.Get("/:id", handler.GetMeow)
    Meows.Patch("/:id", handler.UpdatePartialMeow)
    Meows.Delete("/:id", handler.DeleteMeow)


}