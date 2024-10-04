package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	wallet "charge/proto/wallet"
)

type Charge struct{}

func (e *Charge) Handle(ctx context.Context, msg *wallet.WalletReq) error {
	log.Log("Handler Received message: ", msg)
	return nil
}

func Handler(ctx context.Context, msg *wallet.WalletReq) error {
	log.Log("Function Received message: ", msg)
	return nil
}
