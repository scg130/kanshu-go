package main

import (
	"context"
	"novel/handler"
	"novel/repo"
	"novel/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/scg130/tools"
	"github.com/scg130/tools/handlers"

	novel "novel/proto/novel"
)

const SRV_NAME = "go.kanshu.service.novel"

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

	// Register Handler
	novel.RegisterNovelSrvHandler(service.Server(), &handler.NovelSrv{service.Client(), repo.Category{}, repo.Novel{}, repo.Chapter{}, repo.Notes{}})

	// Register Struct as Subscriber
	micro.RegisterSubscriber("novel.read", service.Server(), &subscriber.NovelRead{
		Note:      repo.Notes{},
		Chapter:   repo.Chapter{},
		WalletLog: repo.WalletLog{},
	})

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
