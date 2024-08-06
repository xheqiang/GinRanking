package controller

import (
	"ginRanking/common"
	"ginRanking/models"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type VoteController struct{}

func (v VoteController) Vote(ctx *gin.Context) {

	session := sessions.Default(ctx)
	loginInfo, ok := session.Get("LoginInfo").(map[string]interface{})
	if !ok {
		JsonOutPut(ctx, 301, "用户未登录，请登录后操作", common.EmptyData)
		return
	}
	user_id, ok := loginInfo["UserId"].(int)
	if !ok {
		JsonOutPut(ctx, 301, "用户未登录，请登录后操作", common.EmptyData)
		return
	}

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
