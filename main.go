package main

import (
	"github.com/borerer/nlib/configs"
	"github.com/borerer/nlib/server"
)

func main() {
	server.Run(configs.GetServerConfig())
}
