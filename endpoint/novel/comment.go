package novel

import (
	"github.com/gin-gonic/gin"
	"github.com/ilylx/gconv"
	"github.com/scg130/tools"
	"github.com/scg130/tools/wrappers"
	"kanshu/dto"
	"kanshu/dto/admin"
	"kanshu/endpoint"
	go_micro_service_comment "kanshu/proto/comment"
	selfwrappers "kanshu/wrappers"
	"math"
	"net/http"
)

type Comment struct {
	commentCli go_micro_service_comment.CommentSrvService
}

var commentSrv *Comment

const (
	COMMENT_SRV_NAME = "go.kanshu.service.comment"
)

func NewCommentSrv() *Comment {
	if commentSrv == nil {
		commentSrv = &Comment{
			commentCli: go_micro_service_comment.NewCommentSrvService(COMMENT_SRV_NAME, tools.GetMicroClient(
				COMMENT_SRV_NAME,
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
		}
	}
	return commentSrv
}

func (this *Comment) CommentList(ctx *gin.Context) {
	var req dto.CommentListReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	rep, err := this.commentCli.GetComments(ctx, &go_micro_service_comment.CommentsRequest{
		Page:    int32(req.Page),
		Size_:   int32(req.Size),
		NovelId: int32(req.NovelId),
	})
	if err != nil || rep.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code:  0,
		Msg:   "success",
		Data:  rep.Comments,
		Pages: gconv.Int32(math.Ceil(gconv.Float64(rep.Total) / gconv.Float64(req.Size))),
		Total: rep.Total,
	})
}

func (this *Comment) DianZan(c *gin.Context) {
	var req dto.DianZanReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	uInfo := endpoint.GetUserInfo(c)
	rsp, err := this.commentCli.DianZan(c, &go_micro_service_comment.DianZanRequest{
		UserId:    gconv.Int32(uInfo.UserId),
		CommentId: gconv.Int32(req.CommentId),
	})
	if err != nil || rsp.Code != 0 {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "fail",
		})
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "success",
	})
}

func (this *Comment) AddComment(c *gin.Context) {
	var req dto.AddCommentReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	uInfo := endpoint.GetUserInfo(c)
	rsp, err := this.commentCli.AddComment(c, &go_micro_service_comment.AddCommentRequest{
		Content: req.Content,
		UserId:  gconv.Int32(uInfo.UserId),
		NovelId: req.NovelId,
	})
	if err != nil || rsp.Code != 0 {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "fail",
		})
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "success",
	})
}
