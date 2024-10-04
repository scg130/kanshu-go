package dto

type AddCommentReq struct {
	Content string `json:"content" form:"content" binding:"required"`
	NovelId int32  `json:"novel_id" form:"novel_id" binding:"required"`
}

type CommentListReq struct {
	Page    int   `json:"page" form:"page" binding:"required"`
	Size    int   `json:"size" form:"size" binding:"required"`
	NovelId int32 `json:"novel_id" form:"novel_id" binding:"required"`
}

type DianZanReq struct {
	CommentId int32 `json:"comment_id" form:"comment_id" binding:"required"`
}
