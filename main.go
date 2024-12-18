package main

import (
	"event-booking/config"
	"event-booking/routes"
	"log"

	_ "event-booking/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Fiber app
	app := fiber.New()

	app.Use(cors.New())

	// Connect to the database
	config.ConnectDB()

	// Perform database migration
	config.Migrate()

	// Perform database seed
	config.SeedUsers()
	config.SeedEvents()

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":8080"))
}
