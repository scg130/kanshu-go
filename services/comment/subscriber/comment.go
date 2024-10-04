package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	comment "comment/proto/comment"
)

type Comment struct{}

func (e *Comment) Handle(ctx context.Context, msg *comment.CommentsRequest) error {
	log.Log("Handler Received message: ", msg)
	return nil
}

func Handler(ctx context.Context, msg *comment.CommentsRequest) error {
	log.Log("Function Received message: ", msg)
	return nil
}
