package controller

import (
	"encoding/json"
	"fmt"
	"ginRanking/common"
	"ginRanking/models"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	user_id, ok := loginInfo["UserId"].(int)
	if !ok {
		JsonOutPut(ctx, 301, "用户未登录，请登录后操作", common.EmptyData)
		return
	} */

	// map结构 序列化json字符串后 存入session 反序列化可以正常获取
	/* loginInfoStr := session.Get("LoginInfo")
	var loginInfo map[string]interface{}
	json.Unmarshal([]byte(loginInfoStr.(string)), &loginInfo)
	user_id_float := loginInfo["UserId"].(float64)
	user_id := int(user_id_float) */

	// 单值存放 获取
	/* user_id := session.Get("LoginUid").(int) */

	// struct 不处理整体存入无法成功 获取值 为 nil
	/* var loginInfo common.LoginInfo
	if session.Get("LoginInfo") == nil {
		JsonOutPut(ctx, 301, "用户未登录，请登录后操作", common.EmptyData)
		return
	}
	loginInfo = session.Get("LoginInfo").(common.LoginInfo)
	user_id := loginInfo.UserId
	fmt.Println("UserId:", loginInfo.UserId)
	fmt.Println("UserName:", loginInfo.UserName) */

	// struct 序列化json字符串后 存入session 反序列化后可以正常获取
	var loginInfo common.LoginInfo
	loginInfoStr := session.Get("LoginInfo").(string)
	json.Unmarshal([]byte(loginInfoStr), &loginInfo)
	user_id := loginInfo.UserId

	fmt.Println("user_id: ", user_id)

	player_id_str := ctx.DefaultPostForm("player_id", "0")
	player_id, _ := strconv.Atoi(player_id_str)
	activity_id_str := ctx.DefaultPostForm("activity_id", "0")
	activity_id, _ := strconv.Atoi(activity_id_str)

	if player_id == 0 || activity_id == 0 {
		JsonOutPut(ctx, 302, "请输入正确的信息", common.EmptyData)
		return
	}

	player, _ := models.GetPlayerByIDActivityID(player_id, activity_id)
	if player.Id == 0 {
		JsonOutPut(ctx, 303, "参赛选手信息错误，请重新选择", common.EmptyData)
		return
	}

	vote, _ := models.GetVoteByUserId(user_id, activity_id, player_id)
	if vote.Id != 0 {
		JsonOutPut(ctx, 304, "您已经投过票了，请勿重复投票", common.EmptyData)
		return
	}

	_, err := models.AddVote(user_id, activity_id, player_id)
	if err != nil {
		JsonOutPut(ctx, 305, "投票失败，请联系管理员", common.EmptyData)
		return
	}

	playerRes, err := models.UpdatePlayerScore(player_id, activity_id)
	fmt.Println(playerRes)
	if err != nil {
		JsonOutPut(ctx, 305, "投票失败，请联系管理员", common.EmptyData)
		return
	}

	var data = map[string]interface{}{
		"player_id":   playerRes.Id,
		"player_name": playerRes.PlayerName,
		"score":       playerRes.Score,
	}

	JsonOutPut(ctx, 0, "success", data)
}
