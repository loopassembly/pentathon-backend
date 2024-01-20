package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/loopassembly/pentathon-backend/initializers"
	"github.com/loopassembly/pentathon-backend/routes"
)

func init() {
	_, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	// initializers.ConnectDB(&config)
}

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	routes.SetupSheetRoutes(app)
 
	app.Listen(":3000")
}
