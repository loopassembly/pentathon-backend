// routes/routes.go

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/loopassembly/pentathon-backend/controllers" 
	
)

func SetupRoutes(app *fiber.App) {
	// Solo route (POST)
	app.Post("/solo", controllers.SoloController)

	// Read data route (GET)
	app.Get("/read", controllers.ReadDataHandler)

	// Create data route (POST)
	app.Post("/create", controllers.CreateDataHandler)

	app.Post("/solohandler", controllers.SoloDataHandler)

	app.Post("/teamhandler", controllers.TeamDataHandler)
}


