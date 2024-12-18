package routes

import (
	"event-booking/controllers"
	"event-booking/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/login", controllers.Login)

	secured := app.Group("/api", middleware.JWTMiddleware)
	secured.Get("/events", controllers.GetEvents)
	secured.Post("/events/:id/approve", controllers.ApproveEvent)
	secured.Post("/events/:id/reject", controllers.RejectEvent)
}
