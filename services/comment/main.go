package main

import (
	"comment/handler"
	"comment/repo"
	"comment/subscriber"
	"context"

	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	"github.com/scg130/tools"
	"github.com/scg130/tools/handlers"

	comment "comment/proto/comment"
)

const SRV_NAME = "go.kanshu.service.comment"

func main() {

	// New Service
	service := tools.NewService(SRV_NAME, handlers.NewOpentracing(SRV_NAME), func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			h(ctx, req, rsp)
			return nil
		}
	})

	// Initialise service
	service.Init()

	srv := &handler.CommentSrv{
		CommentRepo: repo.Comment{},
	}
	// Register Handler
	comment.RegisterCommentSrvHandler(service.Server(), srv)

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.kanshu.srv.comment", service.Server(), new(subscriber.Comment))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
