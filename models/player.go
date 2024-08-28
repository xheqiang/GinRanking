package models

import (
	"time"
)

type Player struct {
	Id         int       `json:"id"`
	ActivityId int       `json:"activity_id"`
	PlayerId   int       `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Avatar     string    `json:"avatar"`
	Desc       string    `json:"desc"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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

func (Player) TableName() string {
	return "player"
}

func GetPlayerInfo(player_id, activity_id int) (Player, error) {

	var player Player

	err := DB.Where("activity_id = ? AND player_id = ?", activity_id, player_id).First(&player).Error
	if err != nil {
		return player, err
	}
	return player, nil
}

func GetPlayerList(activity_id int, sort string) ([]PlayerInfo, error) {

	var playersInfo []PlayerInfo

	var players []Player
	err := DB.Where("activity_id = ?", activity_id).Order(sort).Find(&players).Error
	if err != nil {
		return nil, err
	}
	for _, player := range players {
		score, err := GetPlayerScore(activity_id, player.PlayerId)
		if err != nil {
			return nil, err
		}
		playerInfo := PlayerInfo{
			Id: player.Id,
			ActivityId: player.ActivityId,
			PlayerId:   player.PlayerId,
			PlayerName: player.PlayerName,
			Score:      score.Score,
			Avatar:     player.Avatar,
			Desc:       player.Desc,
		}
		playersInfo = append(playersInfo, playerInfo)
	}
	return playersInfo, nil
}

func GetPlayerRankingDb(activity_id int, sort string) ([]map[string]interface{}, error) {
	// 以score表为准未被投票的无法显示 需要切换维度
	/* scoreList, err := GetPlayerScoreList(activity_id, sort)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	for _, scoreInfo := range scoreList {
		playerId := scoreInfo.PlayerId
		score := scoreInfo.Score
		playerInfo, err := GetPlayerInfo(playerId, activity_id)
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
	playerList, err := GetPlayerList(activity_id, "id asc")
	if err != nil {
		return nil, err
	}
	data := []map[string]interface{}{}
	for _, playerInfo := range playerList {
		playerId := playerInfo.PlayerId
		scoreInfo, err := GetPlayerScore(playerId, activity_id)
		if err != nil {
			return nil, err
		}
		score := scoreInfo.Score
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

func GetPlayerByIDActivityID(player_id int, activity_id int) (Player, error) {

	var player Player

	err := DB.Where("player_id = ? AND activity_id = ?", player_id, activity_id).First(&player).Error
	if err != nil {
		return player, err
	}
	return player, nil
}
