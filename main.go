package main

import (
	"github.com/borerer/nlib/app"
	"github.com/borerer/nlib/configs"
)

func main() {
	app.Run(configs.GetAppConfig())
}
