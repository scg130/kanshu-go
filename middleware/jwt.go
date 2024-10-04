package middleware

import (
	"kanshu/env"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scg130/tools"
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		token = ctx.Query("token")
		if len(token) == 0 {
			token = ctx.Request.Header.Get("Authorization")
			token = strings.Replace(token, "Bearer ", "", 1)
		}
		if len(token) == 0 {
			token, _ = ctx.Cookie("token")
		}
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, nil)
			ctx.Abort()
			return
		}

		authData, err := tools.ValidToken(token, env.JwtConf.AdminJwtSecret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, nil)
			ctx.Abort()
			return
		}
		data := authData.(map[string]interface{})
		ctx.Set("authData", data["udata"])
		ctx.Set("paths", data["paths"])
		ctx.Set("apis", data["apis"])
		ctx.Set("token", token)
		ctx.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		token = ctx.Query("token")
		if len(token) == 0 {
			token = ctx.Request.Header.Get("Authorization")
			if token == "" {
				token = ctx.Request.Header.Get("token")
			}
			token = strings.Replace(token, "Bearer ", "", 1)
		}
		if len(token) == 0 {
			token, _ = ctx.Cookie("token")
		}
		if len(token) == 0 {
			ctx.JSON(http.StatusUnauthorized, nil)
			ctx.Abort()
			return
		}

		authData, err := tools.ValidToken(token, env.JwtConf.Secret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, nil)
			ctx.Abort()
			return
		}
		data := authData.(map[string]interface{})
		ctx.Set("authData", data)
		ctx.Set("token", token)
		ctx.Next()
	}
}
