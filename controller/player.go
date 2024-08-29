package controller

import (
	"fmt"
	"ginRanking/cache"
	"ginRanking/common"
	"ginRanking/models"
	"ginRanking/util/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type PlayerController struct{}

func (p PlayerController) PlayerList(ctx *gin.Context) {

	activityIdStr := ctx.DefaultPostForm("activity_id", "")
	activityId, _ := strconv.Atoi(activityIdStr)

	playerList, err := models.GetPlayerList(activityId, "id asc")

	if err != nil {
		JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
	}

	JsonOutPut(ctx, 0, "success", playerList)
}

func (p PlayerController) PlayerRankingDb(ctx *gin.Context) {

	activityIdStr := ctx.DefaultPostForm("activity_id", "")
	activityId, _ := strconv.Atoi(activityIdStr)

	rankList, err := models.GetPlayerRankingDb(activityId, "score desc")

	if err != nil {
		JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
		return
	}

	JsonOutPut(ctx, 0, "success", rankList)
}

func (p PlayerController) PlayerRankingRedis(ctx *gin.Context) {

	activityIdStr := ctx.DefaultPostForm("activity_id", "")
	activityId, _ := strconv.Atoi(activityIdStr)

	rankingKey := "player_ranking_" + activityIdStr
	rankList := cache.Redis.ZRevRangeWithScores(cache.Rctx, rankingKey, 0, -1).Val()
	fmt.Println("rankList:", rankList)
	if len(rankList) == 0 {
		scoreList, err := models.GetPlayerScoreList(activityId, "score desc")
		if err != nil {
			JsonOutPut(ctx, 201, "无参赛选手信息", common.EmptyData)
		}
		for _, scoreInfo := range scoreList {
			redisRes := cache.Redis.ZAdd(cache.Rctx, rankingKey, redis.Z{Score: float64(scoreInfo.Score), Member: scoreInfo.PlayerId})
			if redisRes.Err() != nil {
				logger.Error(map[string]interface{}{
					"redis Zadd error": redisRes.Err(),
				})
			}

			rankList = append(rankList, redis.Z{
				Score:  float64(scoreInfo.Score),
				Member: scoreInfo.PlayerId,
			})
		}
		// 更新过期时间
		err = cache.Redis.Expire(cache.Rctx, rankingKey, time.Hour*24).Err()
		if err != nil {
			logger.Error(map[string]interface{}{
				"redis set expire error": err,
			})
		}

	}

	rankInfoList := []map[string]interface{}{}
	for _, rankData := range rankList {
		playerId, _ := strconv.Atoi(rankData.Member.(string))
		score := rankData.Score
		playerInfo, _ := models.GetPlayerInfo(playerId, activityId)
		rankInfo := map[string]interface{}{
			"id":          playerId,
			"activity_id": playerInfo.ActivityId,
			"player_id":   playerInfo.PlayerId,
			"player_name": playerInfo.PlayerName,
			"score":       score,
			"avatar":      playerInfo.Avatar,
			"desc":        playerInfo.Desc,
		}
		rankInfoList = append(rankInfoList, rankInfo)
	}

	JsonOutPut(ctx, 0, "success", rankInfoList)
}
