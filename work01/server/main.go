package main

import (
	"zg5/work01/server/common/initialize"
	"zg5/work01/server/models"
)

func main() {
	models.Migrate()
	initialize.InitGrpc()
}
