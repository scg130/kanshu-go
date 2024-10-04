package handler

import (
	"charge/repo"
	"context"
	"errors"
	"time"

	"github.com/micro/go-micro/util/log"

	charge "charge/proto/charge"
)

type ChargeSrv struct {
	OrderRepo  repo.Order
	CreateChan chan *charge.ChargeReq
}

func (self *ChargeSrv) Create(ctx context.Context, r *charge.ChargeReq, rsp *charge.ChargeResponse) error {
	self.CreateChan <- r

	req := <-self.CreateChan
	flag, err := self.OrderRepo.Create(repo.Order{
		UserId:     req.Uid,
		OrderId:    req.OrderId,
		OutTradeNo: req.ThirdOrderNo,
		Amount:     req.Amount,
		State:      int(charge.StateType_STATE_NORMAL),
		Channel:    req.Channel,
		SubjectId:  req.SubjectId,
		Subject:    req.Subject,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil || !flag {
		rsp.State = -1
		log.Logf("create order err:%s", err.Error())
		return nil
	}
	rsp.State = 1
	rsp.OrderId = req.OrderId
	return nil
}

func (self *ChargeSrv) ChargeSuccess(ctx context.Context, req *charge.ChargeReq, rsp *charge.ChargeResponse) error {
	log.Log("Received Charge.Call request")
	order, err := self.OrderRepo.GetByOrderId(req.ThirdOrderNo)
	if err != nil {
		rsp.State = -1
		log.Logf("update order err:%s", err.Error())
		return errors.New("")
	}
	flag, err := self.OrderRepo.Update(repo.Order{
		OutTradeNo: req.ThirdOrderNo,
		State:      int(charge.StateType_STATE_PAY_SUCCESS),
		UpdateTime: time.Now(),
	})
	if err != nil || !flag {
		rsp.State = -1
		log.Logf("update order err:%s", err.Error())
		return errors.New("")
	}
	rsp.State = 1
	rsp.OrderId = order.OrderId
	rsp.UserId = order.UserId
	rsp.OrderIdInt = order.Id
	return nil
}

func (self *ChargeSrv) QueryOrder(ctx context.Context, req *charge.QueryReq, rsp *charge.QueryRsp) error {
	order, err := self.OrderRepo.FindByOrderId(req.OrderId)
	if err != nil || order.Id == 0 {
		rsp.State = -1
		log.Logf("QueryOrder err:%s", err.Error())
		return nil
	}
	rsp.State = 1
	rsp.OrderId = order.OrderId
	rsp.Status = int32(order.State)
	return nil
}

func (self *ChargeSrv) QueryOrderByThirdOrderId(ctx context.Context, req *charge.ChargeReq, rsp *charge.ChargeResponse) error {
	order, err := self.OrderRepo.GetByOrderId(req.ThirdOrderNo)
	if err != nil || order.Id == 0 {
		rsp.State = -1
		log.Logf("QueryOrder err:%s", err.Error())
		return nil
	}
	rsp.State = 1
	rsp.OrderId = order.OrderId
	rsp.Status = int32(order.State)
	return nil
}
