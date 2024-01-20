package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/loopassembly/pentathon-backend/controllers"
)

func SetupSheetRoutes(router fiber.Router) {
	
	router.Get("/register", controllers.GetSheet)
	router.Post("/handleGoogleSheetsData",controllers.GetSheet)

}


