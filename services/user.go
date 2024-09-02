package services

import (
	"encoding/json"
	"ginRanking/common"
	"ginRanking/models"
	"ginRanking/util"
	"ginRanking/util/logger"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserService struct {
}

type UserInfo struct {
	Id        int    `json:"id"`
	UserName  string `json:"user_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u UserService) GetStaticUserInfo(id int, name string) map[string]interface{} {
	data := map[string]interface{}{
		"id":   id,
		"name": name,
	}
	return data
}

func (u UserService) GetUserInfoById(id int) (UserInfo, error) {
	user, err := models.GetUserInfoById(id)
	if err != nil {
		return UserInfo{}, err
	}

	userInfo := UserInfo{
		Id:        user.Id,
		UserName:  user.UserName,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return userInfo, err
}

func (u UserService) GetAllUserList() ([]UserInfo, error) {
	var userList []UserInfo
	users, err := models.GetAllUserList()
	if err != nil {
		return userList, err
	}

	for _, u := range users {
		userInfo := UserInfo{
			Id:        u.Id,
			UserName:  u.UserName,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		userList = append(userList, userInfo)
	}

	return userList, err
}

func (u UserService) GetUserInfoByUserName(userName string) (UserInfo, error) {
	user, err := models.GetUserInfoByUserName(userName)
	if err != nil {
		return UserInfo{}, err
	}

	userInfo := UserInfo{
		Id:        user.Id,
		UserName:  user.UserName,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return userInfo, err
}

func (u UserService) AddUser(userName, password string) (int, error) {
	password = util.EncryMd5(password)
	userId, err := models.AddUser(userName, password)
	if err != nil {
		return 0, err

	}
	return userId, nil
}

func (u UserService) UpdateUserName(id int, userName string) error {
	err := models.UpdateUserName(id, userName)
	return err
}

func (u UserService) DeleteUserById(id int) (int, error) {
	id, err := models.DeleteUserById(id)
	return id, err
}

func (u UserService) Register(userName, password string) map[string]interface{} {

	// 查询用户名是否已经存在
	user, _ := models.GetUserInfoByUserName(userName)
	if user.Id != 0 {
		result := map[string]interface{}{
			"status": 105,
			"msg":    "用户名已存在",
			"data":   common.EmptyData,
		}
		return result
	}

	// 数据库保存状态
	userId, err := models.AddUser(userName, util.EncryMd5(password))
	if err != nil {
		logger.Error(map[string]interface{}{"AddUser Error": err.Error()})
		result := make(map[string]interface{})
		result["status"] = 106
		result["msg"] = "保存失败"
		result["data"] = common.EmptyData
		return result
	}

	data := map[string]interface{}{
		"userId": userId,
	}
	result := map[string]interface{}{
		"status": 0,
		"msg":    "success",
		"data":   data,
	}

	return result
}

func (u UserService) Login(ctx *gin.Context, userName, password string) map[string]interface{} {
	user, err := models.GetUserByUserName(userName)

	if err != nil {
		result := map[string]interface{}{
			"status": 107,
			"msg":    "登录失败，请联系管理员",
			"data":   common.EmptyData,
		}
		return result
	}

	if user.Id == 0 {
		result := map[string]interface{}{
			"status": 106,
			"msg":    "用户名密码不正确",
			"data":   common.EmptyData,
		}
		return result
	}

	if user.Password != util.EncryMd5(password) {
		result := map[string]interface{}{
			"status": 106,
			"msg":    "用户名密码不正确",
			"data":   common.EmptyData,
		}
		return result
	}

	session := sessions.Default(ctx)

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

	result := map[string]interface{}{
		"status": 0,
		"msg":    "success",
		"data":   data,
	}

	return result
}
