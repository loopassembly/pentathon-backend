package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/loopassembly/pentathon-backend/initializers"
	"github.com/loopassembly/pentathon-backend/routes"
	"github.com/loopassembly/pentathon-backend/utils"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
	
	_, err = utils.GetSheetsService()
	if err != nil {
		log.Fatalf("Error initializing Google Sheets service: %v", err)
	}
}


func main() {
	
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:5500", // specify your frontend origin
		AllowMethods: "GET, POST, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	routes.SetupRoutes(app)
 
	app.Listen(":3000")
}
