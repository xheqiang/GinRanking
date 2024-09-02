package services

import (
	"ginRanking/cache"
	"ginRanking/common"
	"ginRanking/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type VoteService struct {
}

func (v VoteService) Vote(userId, activityId, playerId int) map[string]interface{} {
	player, _ := models.GetPlayerByIDActivityID(playerId, activityId)
	if player.Id == 0 {
		result := map[string]interface{}{
			"status": 303,
			"msg":    "参赛选手信息错误，请重新选择",
			"data":   common.EmptyData,
		}
		return result
	}

	vote, _ := models.GetVoteByUserId(userId, activityId, playerId)
	if vote.Id != 0 {
		result := map[string]interface{}{
			"status": 304,
			"msg":    "您已经投过票了，请勿重复投票",
			"data":   common.EmptyData,
		}
		return result
	}

	_, err := models.AddVote(userId, activityId, playerId)
	if err != nil {
		result := map[string]interface{}{
			"status": 305,
			"msg":    "投票失败，请联系管理员",
			"data":   common.EmptyData,
		}
		return result
	}

	playerInfo, err := models.UpdatePlayerScore(playerId, activityId)
	//fmt.Println(playerInfo)
	if err != nil {
		result := map[string]interface{}{
			"status": 305,
			"msg":    "投票失败，请联系管理员2",
			"data":   common.EmptyData,
		}
		return result
	}

	// 更新排行榜redis
	activityIdStr := strconv.Itoa(activityId)
	rankingKey := "player_ranking_" + activityIdStr
	cache.Redis.ZAdd(cache.Rctx, rankingKey, redis.Z{Score: float64(playerInfo["score"].(int)), Member: playerInfo["player_id"].(int)})
	// 更新过期时间
	cache.Redis.Expire(cache.Rctx, rankingKey, time.Hour*24)

	var data = map[string]interface{}{
		"player_id":   playerInfo["player_id"],
		"player_name": playerInfo["player_name"],
		"score":       playerInfo["score"],
	}

	result := map[string]interface{}{
		"status": 0,
		"msg":    "success",
		"data":   data,
	}
	return result
}
