package handler

import (
	"charge/repo"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/micro/go-micro/util/log"

	wallet "charge/proto/wallet"
)

type WalletSrv struct {
	WalletRepo    repo.Wallet
	WalletLogRepo repo.WalletLog
	CreateChan    chan *wallet.WalletReq
}

func (self *WalletSrv) Create(ctx context.Context, req *wallet.WalletReq, rsp *wallet.WalletResponse) error {
	flag, err := self.WalletRepo.Create(repo.Wallet{
		UserId:           req.Uid,
		AvailableBalance: req.Amount,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	})
	if err != nil || !flag {
		rsp.State = -1
		log.Log("create order err")
		return nil
	}
	rsp.State = 1
	return nil
}

func (self *WalletSrv) Change(ctx context.Context, req *wallet.WalletReq, rsp *wallet.WalletResponse) error {
	wl, err := self.WalletLogRepo.GetByOrderId(int(req.OrderId))
	if err == nil && wl.Id > 0 {
		rsp.State = -1
		log.Logf("wlog exist order_id %d", wl.Id)
		return errors.New(fmt.Sprintf("wlog exist order_id %d", wl.Id))
	}
	var amount int64
	switch req.Type {
	case wallet.Type_STATE_CHARGE:
		amount = req.Amount
	case wallet.Type_STATE_BUY_VIP:
		amount -= req.Amount
	case wallet.Type_STATE_BUY_CHAPTER:
		amount -= req.Amount
	default:
		rsp.State = -1
		return errors.New("Unknown type")
	}
	w, err := self.WalletRepo.FindByUserId(int(req.Uid))
	if err != nil || w.Id == 0 {
		_, err := self.WalletRepo.Create(repo.Wallet{
			UserId:           req.Uid,
			AvailableBalance: amount,
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
		})
		if err != nil {
			rsp.State = -1
			log.Errorf("create wallet user_id:%d err:%v", req.Uid, err)
			return err
		}
		rsp.State = 1
		return nil
	}
	flag, err := self.WalletRepo.Update(repo.Wallet{
		UserId:           req.Uid,
		AvailableBalance: w.AvailableBalance + amount,
		UpdateTime:       time.Now(),
	})
	if err != nil || !flag {
		rsp.State = -1
		log.Logf("update wallet err:%s", err.Error())
		return errors.New("")
	}
	if req.Type == wallet.Type_STATE_CHARGE {
		self.WalletLogRepo.Create(repo.WalletLog{
			UserId:     req.Uid,
			Amount:     amount,
			Type:       repo.Type_STATE_CHARGE,
			OrderId:    int(req.OrderId),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
	}
	rsp.State = 1
	return nil
}

func (self *WalletSrv) GetOne(ctx context.Context, req *wallet.WalletReq, rsp *wallet.WalletResponse) error {
	w, err := self.WalletRepo.FindByUserId(int(req.Uid))
	if err != nil {
		rsp.State = -1
		log.Log("GetOne wallet err")
		return nil
	}
	rsp.State = 1
	rsp.AvailableBalance = w.AvailableBalance
	return nil
}

func (self *WalletSrv) FindBuyLog(ctx context.Context, req *wallet.LogRequest, rsp *wallet.LogResponse) error {
	rsp.State = 1
	offset := (req.Page - 1) * req.Size_
	limit := req.Size_
	logs, err := self.WalletLogRepo.FindByUserId(int(req.Uid), offset, limit)
	if err != nil {
		rsp.State = -1
		log.Logf("get logs err:%v", err)
		return err
	}
	rsp.Log = make([]*wallet.Log, 0)
	for _, logV := range logs {
		rsp.Log = append(rsp.Log, &wallet.Log{
			NovelName: logV.NovelName,
			Amount:    logV.Amount,
		})
	}

	return nil
}

func (self *WalletSrv) GetChapter(ctx context.Context, req *wallet.BuyChapterRequest, rsp *wallet.LogResponse) error {
	rsp.State = 1
	cLog, err := self.WalletLogRepo.GetChapterByUserIdAndChapterId(int(req.Uid), int(req.ChapterId))
	if err != nil || cLog.Id == 0 {
		rsp.State = -1
		log.Logf("get chapter err:%v", err)
		return err
	}
	rsp.State = 1
	rsp.Log = make([]*wallet.Log, 0)
	rsp.Log = append(rsp.Log, &wallet.Log{
		Amount: cLog.Amount,
	})
	return nil
}

func (self *WalletSrv) BuyChapter(ctx context.Context, req *wallet.BuyChapterRequest, rsp *wallet.WalletResponse) error {
	rsp.State = 1
	flag, err := self.WalletLogRepo.Create(repo.WalletLog{
		UserId:     req.Uid,
		NovelName:  req.NovelName,
		Amount:     req.Amount,
		CreateTime: time.Now(),
		Type:       int(wallet.Type_STATE_BUY_CHAPTER),
		NovelId:    int(req.NovelId),
		UpdateTime: time.Now(),
		ChapterId:  int(req.ChapterId),
	})
	if err != nil || !flag {
		rsp.State = -1
		log.Logf("create log err:%v", err)
		return err
	}
	return nil
}
