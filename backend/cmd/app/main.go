package main

import (
	"site-sb/infra"
	"site-sb/internal/api"
)

func main() {
	i := infra.New("config/config.json")
	i.SetMode()

	api.NewServer(i).Run()
}
