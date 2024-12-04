package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SuccessJsonStruct struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type ErrorJsonStruct struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ReturnSuccess(ctx *gin.Context, code int, data interface{}, msg string) {
	json := &SuccessJsonStruct{Code: code, Data: data, Msg: msg}
	ctx.JSON(http.StatusOK, json)
}

func ReturnError(ctx *gin.Context, code int, msg string) {
	json := &ErrorJsonStruct{Code: code, Msg: msg}
	ctx.JSON(http.StatusInternalServerError, json)
}
