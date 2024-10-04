package novel

import (
	"github.com/gin-gonic/gin"
	"kanshu/dto"
	"kanshu/endpoint"
	go_micro_service_novel "kanshu/proto/novel"
	"net/http"
)

// @Summary 分类列表
// @Description 分类列表
// @Tags novel
// @Produce json
// @Param query query dto.NovelRequest true "query参数"
// @Success 200 {object}  dto.Resp{data=[]go_micro_service_novel.Category}
// @Router /novel/cate/list [get]
func (n *Novel) Cates(ctx *gin.Context) {
	var req dto.CatesReq
	if err := ctx.BindQuery(&req); err != nil {
		ctx.String(http.StatusBadRequest, "bad params")
		return
	}

	resp, err := n.novelCli.GetCateGories(ctx, &go_micro_service_novel.Request{Page: 1, Size_: 9999, Sex: int32(req.Sex), IsShow: int32(1)})
	if err != nil || resp.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: resp.Categories,
	})
}

// @Summary 加入书架
// @Description 我的书架
// @Tags novel
// @Produce json
// @Param query query dto.NovelRequest true "query参数"
// @Success 200 {object}  dto.Resp{data=[]go_micro_service_novel.Note}
// @Router /novel/join-book [get]
func (n *Novel) JoinBook(ctx *gin.Context) {
	var req dto.JoinReq
	if err := ctx.BindQuery(&req); err != nil {
		ctx.String(http.StatusBadRequest, "bad params")
		return
	}

	userInfo := endpoint.GetUserInfo(ctx)

	resp, err := n.novelCli.JoinNote(ctx, &go_micro_service_novel.Request{UserId: int32(userInfo.UserId), NovelId: req.NovelId})
	if err != nil || resp.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: nil,
	})
}
