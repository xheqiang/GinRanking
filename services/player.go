package services

import (
	"fmt"
	"ginRanking/cache"
	"ginRanking/common"
	"ginRanking/models"
	"ginRanking/util/logger"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type PlayerService struct {
}

type PlayerInfo struct {
	Id         int    `json:"id"`
	ActivityId int    `json:"activity_id"`
	PlayerId   int    `json:"player_id"`
	PlayerName string `json:"player_name"`
	Score      int    `json:"score"`
	Avatar     string `json:"avatar"`
	Desc       string `json:"desc"`
}

func (p PlayerService) GetPlayerList(activityId int, sort string) ([]PlayerInfo, error) {
	players, err := models.GetPlayerList(activityId, sort)
	if err != nil {
		return nil, err
	}

	var playersInfo []PlayerInfo
	for _, player := range players {
		score, err := models.GetPlayerScore(activityId, player.PlayerId)
		if err != nil {
			return nil, err
		}
		playerInfo := PlayerInfo{
			Id:         player.Id,
			ActivityId: player.ActivityId,
			PlayerId:   player.PlayerId,
			PlayerName: player.PlayerName,
			Score:      score,
			Avatar:     player.Avatar,
			Desc:       player.Desc,
		}
		playersInfo = append(playersInfo, playerInfo)
	}
	return playersInfo, nil
}

func (p PlayerService) GetPlayerRankingDb(activityId int, sort string) ([]map[string]interface{}, error) {

	// 以score表为准未被投票的无法显示 需要切换维度
	/* scoreList, err := models.GetPlayerScoreList(activityId, sort)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	for _, scoreInfo := range scoreList {
		playerId := scoreInfo.PlayerId
		score := scoreInfo.Score
		playerInfo, err := models.GetPlayerInfo(playerId, activityId)
		if err != nil {
			return nil, err
		}
		dataInfo := map[string]interface{}{
			"id":          playerId,
			"activity_id": playerInfo.ActivityId,
			"player_id":   playerInfo.PlayerId,
			"player_name": playerInfo.PlayerName,
			"score":       score,
			"avatar":      playerInfo.Avatar,
			"desc":        playerInfo.Desc,
		}
		data = append(data, dataInfo)
	} */

	// 以 player 表为准查询 未查询到 分数为0

	playerList, err := models.GetPlayerList(activityId, "id asc")
	if err != nil {
		return nil, err
	}
	data := []map[string]interface{}{}
	for _, playerInfo := range playerList {
		playerId := playerInfo.PlayerId
		score, err := models.GetPlayerScore(activityId, playerId)
		if err != nil {
			return nil, err
		}
		dataInfo := map[string]interface{}{
			"id":          playerId,
			"activity_id": playerInfo.ActivityId,
			"player_id":   playerInfo.PlayerId,
			"player_name": playerInfo.PlayerName,
			"score":       score,
			"avatar":      playerInfo.Avatar,
			"desc":        playerInfo.Desc,
		}
		data = append(data, dataInfo)
	}

	return data, nil
}

func (p PlayerService) PlayerRankingRedis(activityIdStr string) map[string]interface{} {

	activityId, _ := strconv.Atoi(activityIdStr)

	rankingKey := "player_ranking_" + activityIdStr
	rankList := cache.Redis.ZRevRangeWithScores(cache.Rctx, rankingKey, 0, -1).Val()
	fmt.Println("rankList:", rankList)
	if len(rankList) == 0 {
		scoreList, err := models.GetPlayerScoreList(activityId, "score desc")
		if err != nil {
			result := map[string]interface{}{
				"status": 201,
				"msg":    "无参赛选手信息",
				"data":   common.EmptyData,
			}
			return result
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
				Member: strconv.Itoa(scoreInfo.PlayerId), // Redis返回的Member是个字符串 这里保持一致
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

	result := map[string]interface{}{
		"status": 0,
		"msg":    "success",
		"data":   rankInfoList,
	}
	return result
}
