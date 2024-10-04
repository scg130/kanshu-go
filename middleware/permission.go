package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func Permission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apis, exist := ctx.Get("apis")
		if !exist {
			ctx.Abort()
		}
		authData, exist := ctx.Get("authData")
		if !exist {
			ctx.Abort()
		}
		for _, v := range authData.(map[string]interface{})["role_ids"].([]interface{}) {
			if v.(float64) == 1 {
				ctx.Next()
			}
		}

		uri := ctx.Request.RequestURI
		apisData := apis.([]interface{})
		flag := true
		for _, v := range apisData {
			if strings.Contains(uri, strings.Trim(v.(string), "/*")) {
				flag = false
				break
			}
		}
		if !flag {
			ctx.Abort()
		}
		ctx.Next()
	}
}
