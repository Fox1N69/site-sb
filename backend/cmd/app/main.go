package main

import (
	"site-sb/infra"
	"site-sb/internal/api"
	"site-sb/internal/models"
)

func main() {
	i := infra.New("config/config.json")
	i.SetMode()

	// Выполнить миграцию базы данных
	i.Migrate(&models.User{}, &models.Post{})

	api.NewServer(i).Run()
}
