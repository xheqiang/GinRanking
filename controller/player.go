package controller

import (
	"ginRanking/common"
	"ginRanking/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerController struct{}

func (p PlayerController) GetPlayerList(ctx *gin.Context) {

	activity_id_str := ctx.DefaultPostForm("activity_id", "")
	activity_id, _ := strconv.Atoi(activity_id_str)

	playerList, err := models.GetPlayerList(activity_id)

	if err != nil {
		JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
	}

	JsonOutPut(ctx, 0, "success", playerList)
}

