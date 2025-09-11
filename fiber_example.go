package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Create new Fiber instance
	app := fiber.New()

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Go Fiber!")
	})

	app.Get("/api/songs", func(c *fiber.Ctx) error {
		// This will be replaced with actual database query
		songs := []map[string]interface{}{
			{"id": 1, "title": "Song 1", "artist": "Artist 1"},
			{"id": 2, "title": "Song 2", "artist": "Artist 2"},
		}
		return c.JSON(songs)
	})

	// Start server
	log.Fatal(app.Listen(":8080"))
}
