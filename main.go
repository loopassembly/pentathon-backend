package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/loopassembly/pentathon-backend/initializers"
	"github.com/loopassembly/pentathon-backend/routes"
	"github.com/loopassembly/pentathon-backend/utils"
)

func init() {
	
	_, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables:", err.Error())
	}

	
	_, err = utils.GetSheetsService()
	if err != nil {
		log.Fatalf("Error initializing Google Sheets service: %v", err)
	}
}


func main() {
	
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	routes.SetupRoutes(app)
 
	app.Listen(":3000")
}
