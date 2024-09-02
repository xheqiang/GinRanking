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

	if err != nil && err != redis.Nil {
		return player, err
	}

	if err == nil {
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

func GetPlayerList(activityId int, sort string) ([]Player, error) {

	var players []Player
	err := DB.Where("activity_id = ?", activityId).Order(sort).Find(&players).Error
	if err != nil {
		return nil, err
	}

	return players, nil
}

func GetPlayerByIDActivityID(playerId int, activityId int) (Player, error) {

	var player Player

	err := DB.Where("player_id = ? AND activity_id = ?", playerId, activityId).First(&player).Error
	if err != nil {
		return player, err
	}
	return player, nil
}
