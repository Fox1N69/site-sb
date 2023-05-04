package server 

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	port := ":8080"

	app.Static("/", "/client/index.html")

	app.Listen(port)

}
