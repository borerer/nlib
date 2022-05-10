package main

import (
	"gitea.home.iloahz.com/iloahz/nlib/app"
	"gitea.home.iloahz.com/iloahz/nlib/configs"
)

func main() {
	app.Run(configs.GetAppConfig())
}
