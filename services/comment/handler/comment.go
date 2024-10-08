package handler

import (
	comment "comment/proto/comment"
	"comment/repo"
	"context"
	"time"
)

type CommentSrv struct {
	CommentRepo repo.Comment
}

func (c *CommentSrv) DianZan(ctx context.Context, in *comment.DianZanRequest, out *comment.CommentResponse) error {
	commentLog, err := c.CommentRepo.GetByID(in.CommentId)
	if err != nil {
		out.Code = -1
		out.Msg = "fail"
		return err
	}
	commentLog.GoodsNum += 1
	err = c.CommentRepo.Update(commentLog)
	if err != nil {
		out.Code = -1
		out.Msg = "fail"
		return err
	}
	out.Code = 0
	out.Msg = "success"
	return nil
}

func (c *CommentSrv) AddComment(ctx context.Context, in *comment.AddCommentRequest, out *comment.CommonResponse) error {
	_, err := c.CommentRepo.Create(repo.Comment{
		UserId:     int64(in.UserId),
		NovelId:    int(in.NovelId),
		Content:    in.Content,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil {
		out.Code = -1
		out.Msg = "fail"
		return err
	}
	out.Code = 0
	out.Msg = "success"
	return nil
}

func (c *CommentSrv) GetComments(ctx context.Context, in *comment.CommentsRequest, out *comment.CommentResponse) error {
	datas, total, err := c.CommentRepo.FindByNovelID(int(in.NovelId), (in.Page-1)*in.Size_, in.Size_)
	if err != nil {
		out.Code = -1
		out.Msg = "fail"
		return err
	}
	comments := make([]*comment.Comment, 0)
	for _, v := range datas {
		comments = append(comments, &comment.Comment{
			Content:   v.Content,
			Username:  v.Username,
			Avatar:    v.Avatar,
			CreatedAt: v.CreateTime.Format("2006/01/02 15:04:05"),
			Id:        int32(v.Id),
			GoodsNum:  v.GoodsNum,
		})
	}
	out.Code = 0
	out.Msg = "success"
	out.Comments = comments
	out.Total = int32(total)
	return nil
}
