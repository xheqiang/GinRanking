package controller

import (
	"ginRanking/cache"
	"ginRanking/common"
	"ginRanking/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type PlayerController struct{}

func (p PlayerController) PlayerList(ctx *gin.Context) {

	activity_id_str := ctx.DefaultPostForm("activity_id", "")
	activity_id, _ := strconv.Atoi(activity_id_str)

	playerList, err := models.GetPlayerList(activity_id, "id asc")

	if err != nil {
		JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
	}

	JsonOutPut(ctx, 0, "success", playerList)
}

func (p PlayerController) PlayerRanking(ctx *gin.Context) {

	activity_id_str := ctx.DefaultPostForm("activity_id", "")
	activity_id, _ := strconv.Atoi(activity_id_str)

	playerList, err := models.GetPlayerList(activity_id, "score desc")

	if err != nil {
		JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
		return
	}

	JsonOutPut(ctx, 0, "success", playerList)
}

func (p PlayerController) PlayerRankingRedis(ctx *gin.Context) {

	activity_id_str := ctx.DefaultPostForm("activity_id", "")
	activity_id, _ := strconv.Atoi(activity_id_str)

	rankingKey := "player_ranking_" + activity_id_str
	rankList := cache.Redis.ZRevRangeWithScores(cache.Rctx, rankingKey, 0, -1).Val()
	if len(rankList) == 0 {
		playerList, err := models.GetPlayerList(activity_id, "score desc")
		if err != nil {
			JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
		}
		for _, playerData := range playerList {
			cache.Redis.ZAdd(cache.Rctx, rankingKey, redis.Z{Score: float64(playerData.Score), Member: playerData.PlayerId})

			rankList = append(rankList, redis.Z{
				Score:  float64(playerData.Score),
				Member: playerData.PlayerId,
			})
		}
	}

	rankInfo := map[string]interface{}{}
	for _, rankData := range rankList {
		rankInfo["playerId"] = rankData.Member
		rankInfo["score"] = rankData.Score
	}

	JsonOutPut(ctx, 0, "success", rankInfo)
}
