package repo

import (
	"errors"
	"log"
	"time"
)

const (
	Type_STATE_CHARGE = 1
	Type_STATE_BUY    = 2
)

type WalletLog struct {
	Id        int64  `xorm:" pk autoincr INT(11)"`
	UserId    int64  `xorm:"not null default 0 comment('用户id') int"`
	Amount    int64  `xorm:"not null default 0 comment('读书币') int"`
	Type      int    `xorm:"not null default 0 comment('类型 1 增加 2 减少') int"`
	OrderId   int    `xorm:"not null default 0 comment('订单id') int"`
	NovelId   int    `xorm:"not null default 0 comment('小说id') int"`
	NovelName string `xorm:"not null default 0 comment('小说名称') VARCHAR(255)"`
	ChapterId int    `xorm:"not null default 0 comment('章节id') int"`

	UpdateTime time.Time `xorm:"updated_at not null  DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	wlog := new(WalletLog)
	if isExist, _ := x.IsTableExist(wlog); !isExist {
		if err := x.Sync2(wlog); err != nil {
			log.Fatal("sync tables err:%v", err)
		}
	}
}

func (u *WalletLog) Create(wlog WalletLog) (bool, error) {
	affected, err := x.Insert(wlog)
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("insert wallet log fail")
	}
	return true, nil
}

func (u *WalletLog) FindByUserId(userId int, offset, limit int32) ([]*WalletLog, error) {
	wlogs := make([]*WalletLog, 0)
	err := x.Where("user_id = ? and type = ?", userId, Type_STATE_BUY).
		Select("sum(amount) as amount,min(novel_name) as novel_name,max(id) as id").
		GroupBy("novel_id").
		OrderBy("id desc").
		Limit(int(limit), int(offset)).
		Find(&wlogs)
	if err != nil {
		log.Printf("wallet log find by user_id err:%v", err)
		return wlogs, errors.New("not found")
	}
	return wlogs, nil
}

func (self *WalletLog) GetChapterByUserIdAndChapterId(userId int, ChapterId int) (*WalletLog, error) {
	wlog := new(WalletLog)
	_, err := x.Where("user_id = ? and chapter_id = ?", userId, ChapterId).
		Get(wlog)
	if err != nil {
		log.Printf("wallet log find by user_id err:%v", err)
		return wlog, errors.New("not found")
	}
	return wlog, nil
}

func (u *WalletLog) GetByOrderId(orderId int) (WalletLog, error) {
	var wl WalletLog
	rel, err := x.Where("order_id = ? and type = ?", orderId, Type_STATE_CHARGE).Get(&wl)
	if !rel || err != nil {
		log.Printf("GetByOrderId err:%v", err)
		return wl, errors.New("not found")
	}
	return wl, nil
}
