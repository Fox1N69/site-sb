package mian

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()	
	port := ":3000"

	app.Static("/", "./client/dist")
	
	app.Listen(port)
}