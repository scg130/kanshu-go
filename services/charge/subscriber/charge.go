package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	charge "charge/proto/charge"
)

type Charge struct{}

func (e *Charge) Handle(ctx context.Context, msg *charge.ChargeReq) error {
	log.Log("Handler Received message: ", msg)
	return nil
}

func Handler(ctx context.Context, msg *charge.ChargeReq) error {
	log.Log("Function Received message: ", msg)
	return nil
}
