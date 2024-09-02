package models

import (
	"errors"
	"ginRanking/cache"
	"ginRanking/common"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Score struct {
	ActivityId int       `json:"activity_id"`
	PlayerId   int       `json:"player_id"`
	Score      int       `json:"score"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Score) TableName() string {
	return "score"
}

func GetPlayerScore(activityId, playerId int) (int, error) {

	rankingKey := "player_ranking_" + strconv.Itoa(activityId)
	playIdStr := strconv.Itoa(playerId)
	playerScore, err := cache.Redis.ZScore(cache.Rctx, rankingKey, playIdStr).Result()

	if err != nil && err != redis.Nil {
		return 0, err
	}

	if err == nil {
		return int(playerScore), nil
	}

	var score Score
	err = DB.Where("activity_id = ? AND player_id = ?", activityId, playerId).First(&score).Error
	/* if err != nil {
		return 0, err
	} */
	// 未找到 也会抛一个错误 当数据库未找到时候 不能抛异常 抛出空数据
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 记录未找到，返回分数为 0 的 Score 结构体
			//return Score{Score: 0}, nil
			cache.Redis.ZAdd(cache.Rctx, rankingKey, redis.Z{Score: float64(0), Member: playIdStr})
			return 0, nil
		}
		// 发生其他错误，需要抛出来
		return 0, err
	}
	cache.Redis.ZAdd(cache.Rctx, rankingKey, redis.Z{Score: float64(score.Score), Member: playIdStr})
	return score.Score, nil
}

func GetPlayerScoreList(activityId int, sort string) ([]Score, error) {

	var scoreList []Score

	err := DB.Where("activity_id = ?", activityId).Order(sort).Find(&scoreList).Error
	if err != nil {
		return nil, err
	}
	return scoreList, nil
}

func UpdatePlayerScore(playerId, activityId int) (map[string]interface{}, error) {

	var score Score

	err := DB.Where("player_id = ? AND activity_id = ?", playerId, activityId).First(&score).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录未找到，则创建一条新的记录
			score = Score{
				PlayerId:   playerId,
				ActivityId: activityId,
				Score:      1, // 初始分数为 1
			}
			err = DB.Create(&score).Error
			if err != nil {
				return common.EmptyData, err
			}
		} else {
			return common.EmptyData, err
		}
	} else {
		// err = DB.Model(&score).Where("player_id = ? AND activity_id = ?", playerId, activityId).Update("score", gorm.Expr("score + ?", 1)).Error
		err = DB.Model(&score).Update("score", gorm.Expr("score + ?", 1)).Error
		if err != nil {
			return common.EmptyData, err
		}
	}

	player, _ := GetPlayerInfo(playerId, activityId)

	data := map[string]interface{}{
		"id":          player.Id,
		"activity_id": player.ActivityId,
		"player_id":   player.PlayerId,
		"player_name": player.PlayerName,
		"score":       score.Score,
		"avatar":      player.Avatar,
		"desc":        player.Desc,
	}

	return data, nil
}
