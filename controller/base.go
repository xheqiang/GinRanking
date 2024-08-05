package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonStruct struct {
	Status int         `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

func JsonOutPut(ctx *gin.Context, status int, msg, data interface{}) {

	json := &JsonStruct{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
	ctx.JSON(http.StatusOK, json)
}
