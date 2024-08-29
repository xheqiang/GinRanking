package models

import (
	"encoding/json"
	"fmt"
	"ginRanking/cache"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
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

var (
	playerCache     = make(map[int]Player)
	playerCacheLock sync.RWMutex
)

func (Player) TableName() string {
	return "player"
}

func GetPlayerInfoMom(playerId, activityId int) (Player, error) {

	var player Player

	playerCacheLock.RLock()
	cachedPlayer, found := playerCache[playerId]
	playerCacheLock.RUnlock()

	if found {
		return cachedPlayer, nil
	}

	err := DB.Where("activity_id = ? AND player_id = ?", activityId, playerId).First(&player).Error
	if err != nil {
		return player, err
	}

	playerCacheLock.Lock()
	playerCache[playerId] = player
	playerCacheLock.Unlock()

	return player, nil
}

func GetPlayerInfo(playerId, activityId int) (Player, error) {

	var player Player

	playerInfoKey := fmt.Sprintf("player_info_%d_%d", activityId, playerId)
	cachedPlayer, err := cache.Redis.Get(cache.Rctx, playerInfoKey).Result()

	if err != redis.Nil { // 如果 Redis 错误不是缓存不存在的错误
		return player, err
	}
	if err == nil { // 如果 Redis 命中缓存
		err = json.Unmarshal([]byte(cachedPlayer), &player)
		if err != nil {
			return player, err
		}
		return player, nil
	}

	err = DB.Where("activity_id = ? AND player_id = ?", activityId, playerId).First(&player).Error
	if err != nil {
		return player, err
	}

	// 查询结构写入Redis
	playerJson, err := json.Marshal(player)
	if err == nil {
		cache.Redis.Set(cache.Rctx, playerInfoKey, playerJson, time.Hour*24).Err()
	}

	return player, nil
}

func GetPlayerList(activityId int, sort string) ([]PlayerInfo, error) {

	var playersInfo []PlayerInfo

	var players []Player
	err := DB.Where("activity_id = ?", activityId).Order(sort).Find(&players).Error
	if err != nil {
		return nil, err
	}
	for _, player := range players {
		score, err := GetPlayerScore(activityId, player.PlayerId)
		if err != nil {
			return nil, err
		}
		playerInfo := PlayerInfo{
			Id:         player.Id,
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

func GetPlayerRankingDb(activityId int, sort string) ([]map[string]interface{}, error) {
	// 以score表为准未被投票的无法显示 需要切换维度
	/* scoreList, err := GetPlayerScoreList(activityId, sort)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	for _, scoreInfo := range scoreList {
		playerId := scoreInfo.PlayerId
		score := scoreInfo.Score
		playerInfo, err := GetPlayerInfo(playerId, activityId)
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
	playerList, err := GetPlayerList(activityId, "id asc")
	if err != nil {
		return nil, err
	}
	data := []map[string]interface{}{}
	for _, playerInfo := range playerList {
		playerId := playerInfo.PlayerId
		scoreInfo, err := GetPlayerScore(activityId, playerId)
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

func GetPlayerByIDActivityID(playerId int, activityId int) (Player, error) {

	var player Player

	err := DB.Where("player_id = ? AND activity_id = ?", playerId, activityId).First(&player).Error
	if err != nil {
		return player, err
	}
	return player, nil
}
