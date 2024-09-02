package controller

import (
	"ginRanking/common"
	"ginRanking/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerController struct{}

var PlayerService services.PlayerService

func (p PlayerController) PlayerList(ctx *gin.Context) {

	activityIdStr := ctx.DefaultPostForm("activity_id", "")
	activityId, _ := strconv.Atoi(activityIdStr)

	playerList, err := PlayerService.GetPlayerList(activityId, "id asc")

	if err != nil {
		JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
	}

	JsonOutPut(ctx, 0, "success", playerList)
}

func (p PlayerController) PlayerRankingDb(ctx *gin.Context) {

	activityIdStr := ctx.DefaultPostForm("activity_id", "")
	activityId, _ := strconv.Atoi(activityIdStr)

	rankList, err := PlayerService.GetPlayerRankingDb(activityId, "score desc")

	if err != nil {
		JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
		return
	}

	JsonOutPut(ctx, 0, "success", rankList)
}

func (p PlayerController) PlayerRankingRedis(ctx *gin.Context) {

	activityIdStr := ctx.DefaultPostForm("activity_id", "")
	if activityIdStr == "" {
		JsonOutPut(ctx, 201, "参数错误", common.EmptyData)
		return
	}

	result := PlayerService.PlayerRankingRedis(activityIdStr)

	JsonOutPut(ctx, result["status"].(int), result["msg"], result["data"])
}
