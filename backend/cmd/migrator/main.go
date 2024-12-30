package main

import (
	"site-sb/infra"
	"site-sb/internal/models"
)

func main() {
	i := infra.New("config/config.json")
	i.SetMode()
	i.Migrate(
		&models.User{},
		&models.Login{},
		&models.Post{},
	)
}
