package main

import (
	"github.com/m3owmurrr/DropNote/backend/internal/app"
	"github.com/m3owmurrr/DropNote/backend/internal/utils/config"
)

func main() {
	config.LoadConfig()
	app.RunServer()
}
