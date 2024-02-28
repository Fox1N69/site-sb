package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")

	log.Fatal(app.Listen(port))
}
