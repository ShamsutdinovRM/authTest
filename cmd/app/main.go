package main

import "authTest/pkg/app"

const ConfigPath = "configs"

func main() {
	app.Run(ConfigPath)
}
