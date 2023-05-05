package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"data": "user data on backend",
		})
	})

	app.Listen(":8080")
}