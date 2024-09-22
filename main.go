// main.go
package main

import (
	"auth_go/config"
	"auth_go/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	config.LoadEnv()

	// Create a new Fiber instance
	app := fiber.New()

	// Set up security
	config.SetupSecurity(app)

	// Set up routes
	routes.SetupRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
