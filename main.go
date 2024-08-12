package main

import (
	"ginRanking/router"
)

func main() {
	ginServer := router.Router()

	ginServer.Run(":8002")
}




