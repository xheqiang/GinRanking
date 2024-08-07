package controller

import (
	"encoding/json"
	"ginRanking/common"
	"ginRanking/models"
	"ginRanking/util"
	"ginRanking/util/logger"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

// 静态返回 测试专用
func (u UserController) GetStaticUserInfo(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))
	name := ctx.Param("name")

	data := map[string]interface{}{
		"id":   id,
		"name": name,
	}

	JsonOutPut(ctx, 0, "success", data)
}

func (u UserController) GetUserInfoById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	//logger.Debug(map[string]interface{}{"id": id}, "GetUserInfoById")

	user, err := models.GetUserInfoById(id)
	if err != nil {
		JsonOutPut(ctx, 0, "success", "Not Find User")
		logger.Error(map[string]interface{}{"Find User Info Error": err.Error()})
		return
	}
	JsonOutPut(ctx, 0, "success", user)
}

func (u UserController) GetAllUserList(ctx *gin.Context) {
	users, err := models.GetAllUserList()

	if err != nil {
		JsonOutPut(ctx, 0, "success", "Not Find User List")
		logger.Error(map[string]interface{}{"Find User List Error": err.Error()})
		return
	}

	JsonOutPut(ctx, 0, "success", users)
}

func (u UserController) AddUser(ctx *gin.Context) {
	user_name := ctx.DefaultPostForm("user_name", "")
	password := ctx.DefaultPostForm("password", "")
	userId, err := models.AddUser(user_name, util.EncryMd5(password))
	if err != nil {
		JsonOutPut(ctx, 0, "success", "Add User Error")
		logger.Error(map[string]interface{}{"Add User Error": err.Error()})
		return
	}
	resMap := map[string]interface{}{"userId": userId}
	JsonOutPut(ctx, 0, "success", resMap)
}

func (u UserController) UpdateUserName(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	user_name := ctx.DefaultPostForm("user_name", "")
	err := models.UpdateUserName(id, user_name)
	if err != nil {
		JsonOutPut(ctx, 404, "error", "Update User Error")
		logger.Error(map[string]interface{}{"Update User Error": err.Error()})
		return
	}
	JsonOutPut(ctx, 0, "success", common.EmptyData)
}

func (u UserController) DeleteUserById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	userId, err := models.DeleteUserById(id)
	if err != nil {
		JsonOutPut(ctx, 0, "success", "Delete User Error")
		logger.Error(map[string]interface{}{"Delete User Error": err.Error()})
		return
	}
	resMap := map[string]interface{}{"userId": userId}
	JsonOutPut(ctx, 0, "success", resMap)
}

// 开始正式处理业务逻辑
// 注册
func (u UserController) Register(ctx *gin.Context) {
	user_name := ctx.DefaultPostForm("user_name", "")
	password := ctx.DefaultPostForm("password", "")
	confirm_password := ctx.DefaultPostForm("confirm_password", "")

	if user_name == "" || password == "" || confirm_password == "" {
		JsonOutPut(ctx, 103, "请输入正确的信息", common.EmptyData)
		return
	}
	if password != confirm_password {
		JsonOutPut(ctx, 104, "两次密码不一致", common.EmptyData)
		return
	}

	// 查询用户名是否已经存在
	user, _ := models.GetUserInfoByUserName(user_name)
	if user.Id != 0 {
		JsonOutPut(ctx, 105, "用户名已经存在", common.EmptyData)
	}

	userId, err := models.AddUser(user_name, util.EncryMd5(password))
	if err != nil {
		JsonOutPut(ctx, 106, "保存失败", common.EmptyData)
		logger.Error(map[string]interface{}{"AddUser Error": err.Error()})
		return
	}
	data := map[string]interface{}{
		"userId": userId,
	}
	JsonOutPut(ctx, 0, "保存成功", data)
}

func (u UserController) Login(ctx *gin.Context) {
	user_name := ctx.DefaultPostForm("user_name", "")
	password := ctx.DefaultPostForm("password", "")

	if user_name == "" || password == "" {
		JsonOutPut(ctx, 103, "请输入正确的信息", common.EmptyData)
		return
	}

	user, err := models.GetUserByUserName(user_name)

	if err != nil {
		JsonOutPut(ctx, 107, "登录失败，请联系管理员", common.EmptyData)
		return
	}

	if user.Id == 0 {
		JsonOutPut(ctx, 106, "用户名密码不正确", common.EmptyData)
		return
	}

	if user.Password != util.EncryMd5(password) {
		JsonOutPut(ctx, 106, "用户名密码不正确", common.EmptyData)
		return
	}

	session := sessions.Default(ctx)

	// Map 整体放入Redis 无法成功 需要Json序列化
	//var loginInfo = map[string]interface{}{}
	/* loginInfo := make(map[string]interface{})
	loginInfo["UserId"] = user.Id
	loginInfo["UserName"] = user.UserName
	session.Set("LoginInfo", loginInfo) */

	// Map json 序列化后放入Redis 是可以的
	/* loginInfo := make(map[string]interface{})
	loginInfo["UserId"] = user.Id
	loginInfo["UserName"] = user.UserName
	loginInfoJson, _ := json.Marshal(loginInfo)
	session.Set("LoginInfo", string(loginInfoJson)) */

	// 单值存放 正常存取
	/* session.Set("LoginUid", user.Id)
	session.Set("LoginUname", user.Name) */

	// 结构体 整体放入Redis 无法成功 需要Json序列化
	/* loginInfo := common.LoginInfo{
		UserId:   user.Id,
		UserName: user.UserName,
	}
	session.Set("LoginInfo", loginInfo) */

	// 结构体序列化后放入Redis 是可以的
	loginInfo := common.LoginInfo{
		UserId:   user.Id,
		UserName: user.UserName,
	}
	loginInfoJson, _ := json.Marshal(loginInfo)
	session.Set("LoginInfo", string(loginInfoJson))

	session.Save()

	var data = map[string]interface{}{
		"user_id":    user.Id,
		"user_name":  user.UserName,
		"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	JsonOutPut(ctx, 0, "success", data)
}
