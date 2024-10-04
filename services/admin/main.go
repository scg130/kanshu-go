package main

import (
	"admin/handler"
	"admin/repo"
	"admin/subscriber"
	"context"

	"github.com/micro/go-micro/v2/server"
	"github.com/scg130/tools"
	"github.com/scg130/tools/handlers"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	admin "admin/proto/admin"
)

const (
	SRV_NAME = "go.kanshu.service.admin"
)

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
	admin.RegisterAdminUserHandler(service.Server(), &handler.Admin{
		Cli:  service.Client(),
		User: repo.User{},
		Menu: repo.Menu{},
		Role: repo.Role{},
	})

	// Register Struct as Subscriber
	micro.RegisterSubscriber(SRV_NAME, service.Server(), new(subscriber.Admin))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
