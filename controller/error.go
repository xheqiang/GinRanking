package controller

import (
	"ginRanking/util/logger"

	"github.com/gin-gonic/gin"
)

type ErrorController struct {
}

func (u ErrorController) TestErr(ctx *gin.Context) {

	logger.Write("日志信息", "error")

	// 捕获异常 防止前端报错
	/* defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获到异常：", err)
		}
	}() */

	num1 := 1
	num2 := 0
	num3 := num1 / num2

	JsonOutPut(ctx, 404, "success", num3)
}
