package routes

import (
	"server/internal/controller"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", controller.HomePage)
}
