package main

import (
	router "ginchat/router"
	utils "ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8080")
}
