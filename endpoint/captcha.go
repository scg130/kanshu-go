package endpoint

import (
	"kanshu/dto"
	"log"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func CaptchaGenerate(ctx *gin.Context) {
	id := captcha.NewLen(4)

	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: gin.H{
			"id": id,
		},
	})
}

func CaptchaImage(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, "id is null")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "image/png")
	if err := captcha.WriteImage(ctx.Writer, id, 200, 52); err != nil {
		log.Println("show captcha error", err)
	}
}
