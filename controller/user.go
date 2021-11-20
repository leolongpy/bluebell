package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err)
		return
	}
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err)
		return
	}
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"atoken":    user.ATOKEN,
		"rtoken":    user.RTOKEN,
	})
	return
}
