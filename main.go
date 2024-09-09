package main

import (
	"ginRanking/router"
	"ginRanking/util/logger"
)

func main() {
	ginServer := router.Router()

	ginServer.Run(":8002")
	logger.Debug(map[string]interface{}{}, "gin server start")
}
