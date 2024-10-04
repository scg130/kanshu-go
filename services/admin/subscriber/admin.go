package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	admin "admin/proto/admin"
)

type Admin struct {
}

func (e *Admin) Handle(ctx context.Context, msg *admin.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *admin.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
