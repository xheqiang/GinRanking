package router

import (
	"ginRanking/controller"
	"ginRanking/util/logger"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gEngine := gin.Default()

	// 引入日志工具
	gEngine.Use(gin.LoggerWithConfig(logger.LoggerToFile()))
	gEngine.Use(logger.Recover)

	gEngine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	/**
	 * 参数获取方法
	 *
	 * 1. url 参数获取
	 	user/info?userId=1
		userId := context.Query("userId")
	 * 2. uri 占位拼接
	 	user/info/1/zhangshan
		ginServer.GET("/user/info/:userId/:userName"
		userId := context.Param("userId")
	 * 3. post 表单获取
	 	user/info
		username := context.PostForm("username")
	 * 4. json方式传递
		rawData, _ := context.GetRawData()
		var m map[string]interface{}
		_ = json.Unmarshal(rawData, &m)
		或者：
		param := make(map[string]interface{})
		err := context.BindJSON(&param)
	 *
	*/

	// 定义Controller快捷方式
	UserController := &controller.UserController{}
	ErrorController := &controller.ErrorController{}

	user := gEngine.Group("/user")
	{
		// http://localhost:8000/user/info/1/zhangshan
		// user.POST("/staticInfo/:id/:name", controller.UserController{}.GetStaticUserInfo)
		user.POST("/staticInfo/:id/:name", UserController.GetStaticUserInfo)

		user.POST("/info", UserController.GetUserInfoById)

		user.POST("/list", UserController.GetAllUserList)

		user.POST("/add", UserController.AddUser)

		user.POST("/update", UserController.UpdateUserName)

		user.POST("/delete", UserController.DeleteUserById)

		user.POST("/register", UserController.Register)
	}

	error := gEngine.Group("/error")
	{
		//error.GET("testErr", controller.ErrorController{}.TestErr)
		error.GET("testErr", ErrorController.TestErr)
	}

	return gEngine
}
