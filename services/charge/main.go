package main

import (
	"charge/handler"
	"charge/repo"
	"charge/subscriber"
	"context"

	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	"github.com/scg130/tools"
	"github.com/scg130/tools/handlers"

	charge "charge/proto/charge"
)

const SRV_NAME = "go.kanshu.service.charge"

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

	srv := &handler.ChargeSrv{
		OrderRepo:  repo.Order{},
		CreateChan: make(chan *charge.ChargeReq, 10),
	}
	// Register Handler
	charge.RegisterChargeHandler(service.Server(), srv)

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.kanshu.srv.charge", service.Server(), new(subscriber.Charge))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
