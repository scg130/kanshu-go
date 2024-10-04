package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/ilylx/gconv"
	"kanshu/dto"
)

func GetUserInfo(ctx *gin.Context) *dto.UserInfo {
	authData, isExist := ctx.Get("authData")
	if !isExist {
		return nil
	}

	userInfo := authData.(map[string]interface{})

	return &dto.UserInfo{
		UserId: gconv.Int64(userInfo["user_id"].(float64)),
		Phone:  gconv.String(userInfo["phone"].(float64)),
	}
}
