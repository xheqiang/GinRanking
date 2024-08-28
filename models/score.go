package models

import (
	"errors"
	"ginRanking/common"
	"time"

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

func GetPlayerScore(activity_id, player_id int) (Score, error) {
	var score Score

	err := DB.Where("activity_id = ? AND player_id = ?", activity_id, player_id).First(&score).Error
	/* if err != nil {
		return score, err
	} */
	// 未找到 也会抛一个错误 当数据库未找到时候 不能抛异常 抛出空数据
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 记录未找到，返回分数为 0 的 Score 结构体
			return Score{Score: 0}, nil
		}
		// 发生其他错误，需要抛出来
		return score, err
	}
	return score, nil
}

func GetPlayerScoreList(activity_id int, sort string) ([]Score, error) {

	var scoreList []Score

	err := DB.Where("activity_id = ?", activity_id).Order(sort).Find(&scoreList).Error
	if err != nil {
		return nil, err
	}
	return scoreList, nil
}

func UpdatePlayerScore(player_id, activity_id int) (map[string]interface{}, error) {

	var score Score

	err := DB.Where("player_id = ? AND activity_id = ?", player_id, activity_id).First(&score).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录未找到，则创建一条新的记录
			score = Score{
				PlayerId:   player_id,
				ActivityId: activity_id,
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
		// err = DB.Model(&score).Where("player_id = ? AND activity_id = ?", player_id, activity_id).Update("score", gorm.Expr("score + ?", 1)).Error
		err = DB.Model(&score).Update("score", gorm.Expr("score + ?", 1)).Error
		if err != nil {
			return common.EmptyData, err
		}
	}

	player, _ := GetPlayerInfo(player_id, activity_id)

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
