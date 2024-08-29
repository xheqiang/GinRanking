package controller

import (
	"encoding/json"
	"fmt"
	"ginRanking/cache"
	"ginRanking/common"
	"ginRanking/models"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type VoteController struct{}

func (v VoteController) Vote(ctx *gin.Context) {

	session := sessions.Default(ctx)

	// map结构 非序列化 存入session不成功
	/* loginInfo, ok := session.Get("LoginInfo").(map[string]interface{})
	fmt.Println("LoginInfo:", loginInfo)
	if !ok {
		JsonOutPut(ctx, 301, "用户未登录，请登录后操作", common.EmptyData)
		return
	}
	userId, ok := loginInfo["UserId"].(int)
	if !ok {
		JsonOutPut(ctx, 301, "用户未登录，请登录后操作", common.EmptyData)
		return
	} */

	// map结构 序列化json字符串后 存入session 反序列化可以正常获取
	/* loginInfoStr := session.Get("LoginInfo")
	var loginInfo map[string]interface{}
	json.Unmarshal([]byte(loginInfoStr.(string)), &loginInfo)
	userIdFloat := loginInfo["UserId"].(float64)
	userId := int(user_id_float) */

	// 单值存放 获取
	/* userId := session.Get("LoginUid").(int) */

	// struct 不处理整体存入无法成功 获取值 为 nil
	/* var loginInfo common.LoginInfo
	if session.Get("LoginInfo") == nil {
		JsonOutPut(ctx, 301, "用户未登录，请登录后操作", common.EmptyData)
		return
	}
	loginInfo = session.Get("LoginInfo").(common.LoginInfo)
	userId := loginInfo.UserId
	fmt.Println("UserId:", loginInfo.UserId)
	fmt.Println("UserName:", loginInfo.UserName) */

	// struct 序列化json字符串后 存入session 反序列化后可以正常获取
	var loginInfo common.LoginInfo
	loginInfoStr := session.Get("LoginInfo").(string)
	json.Unmarshal([]byte(loginInfoStr), &loginInfo)
	userId := loginInfo.UserId

	fmt.Println("user_id: ", userId)

	playerIdStr := ctx.DefaultPostForm("player_id", "0")
	playerId, _ := strconv.Atoi(playerIdStr)
	activityIdStr := ctx.DefaultPostForm("activity_id", "0")
	activityId, _ := strconv.Atoi(activityIdStr)

	if playerId == 0 || activityId == 0 {
		JsonOutPut(ctx, 302, "请输入正确的信息", common.EmptyData)
		return
	}

	player, _ := models.GetPlayerByIDActivityID(playerId, activityId)
	if player.Id == 0 {
		JsonOutPut(ctx, 303, "参赛选手信息错误，请重新选择", common.EmptyData)
		return
	}

	vote, _ := models.GetVoteByUserId(userId, activityId, playerId)
	if vote.Id != 0 {
		JsonOutPut(ctx, 304, "您已经投过票了，请勿重复投票", common.EmptyData)
		return
	}

	_, err := models.AddVote(userId, activityId, playerId)
	if err != nil {
		JsonOutPut(ctx, 305, "投票失败，请联系管理员", common.EmptyData)
		return
	}

	playerInfo, err := models.UpdatePlayerScore(playerId, activityId)
	//fmt.Println(playerInfo)
	if err != nil {
		JsonOutPut(ctx, 305, "投票失败，请联系管理员", common.EmptyData)
		return
	}

	// 更新排行榜redis
	rankingKey := "player_ranking_" + activityIdStr
	cache.Redis.ZAdd(cache.Rctx, rankingKey, redis.Z{Score: float64(playerInfo["score"].(int)), Member: playerInfo["player_id"].(int)})
	// 更新过期时间
	cache.Redis.Expire(cache.Rctx, rankingKey, time.Hour*24)

	var data = map[string]interface{}{
		"player_id":   playerInfo["player_id"],
		"player_name": playerInfo["player_name"],
		"score":       playerInfo["score"],
	}

	JsonOutPut(ctx, 0, "success", data)
}
