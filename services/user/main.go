package main

import (
	"user/handler"
	"user/repo"
	"user/subscriber"

	"context"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/scg130/tools"
	"github.com/scg130/tools/handlers"

	user "user/proto/user"
)

const SRV_NAME = "go.kanshu.service.user"

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
	user.RegisterUserCenterHandler(service.Server(), &handler.UserSrv{repo.User{}})

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.kanshu.service.user", service.Server(), new(subscriber.User))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
