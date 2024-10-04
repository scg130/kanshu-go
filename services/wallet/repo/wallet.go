package repo

import (
	"errors"
	"log"
	"time"
)

type Wallet struct {
	Id               int64 `xorm:" pk autoincr INT(11)"`
	UserId           int64 `xorm:"not null default 0 comment('用户id') int"`
	AvailableBalance int64 `xorm:"not null default 0 comment('支付金额 单位分') int"`

	UpdateTime time.Time `xorm:"updated_at not null  DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	order := new(Wallet)
	if isExist, _ := x.IsTableExist(order); !isExist {
		if err := x.Sync2(order); err != nil {
			log.Fatal("sync tables err:%v", err)
		}
	}
}

func (sefl *Wallet) Update(wallet Wallet) (bool, error) {
	if _, err := x.Where("user_id = ?", wallet.UserId).Update(&wallet); err != nil {
		return false, err
	}
	return true, nil
}

func (u *Wallet) Create(wallet Wallet) (bool, error) {
	affected, err := x.Insert(wallet)
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("insert order fail")
	}
	return true, nil
}

func (u *Wallet) FindByUserId(userId int) (Wallet, error) {
	var wallet Wallet
	rel, err := x.Where("user_id = ?", userId).Get(&wallet)
	if !rel || err != nil {
		log.Printf("wallet find by user_id err:%v", err)
		return wallet, errors.New("not found")
	}
	return wallet, nil
}
