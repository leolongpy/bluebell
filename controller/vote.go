package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PostVoteController(c *gin.Context) {
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, err)
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
