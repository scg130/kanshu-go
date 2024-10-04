package repo

import (
	"errors"
	"log"
	"time"
)

type Order struct {
	Id         int64  `xorm:" pk autoincr INT(11)"`
	OrderId    string `xorm:"not null default 0 unique(order_id) comment('订单号') varchar(128)"`
	OutTradeNo string `xorm:"not null default '' comment('第三方订单号') VARCHAR(255)"`
	State      int    `xorm:"not null default 0 comment('订单状态 0 创建未支付 1 支付成功 2 支付中 3 已退款') tinyint(1)"`
	Channel    string `xorm:"not null default '' comment('支付渠道 alipay 支付宝 wx 微信') varchar(64)"`
	SubjectId  int64  `xorm:"not null default 0 comment('商品id') int"`
	Subject    string `xorm:"not null default '' comment('购买商品 vip 购买vip amount 充值金额 chapter 购买章节') varchar(64)"`
	Amount     int64  `xorm:"not null default 0 comment('支付金额 单位分') int"`
	UserId     int64  `xorm:"not null default 0 comment('用户id') int"`

	UpdateTime time.Time `xorm:"updated_at not null  DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	order := new(Order)
	if isExist, _ := x.IsTableExist(order); !isExist {
		if err := x.Sync(order); err != nil {
			log.Fatal("sync tables err:%v", err)
		}
	}
}

func (sefl *Order) Update(order Order) (bool,error) {
	if _, err := x.Where(" out_trade_no=?", order.OutTradeNo).Update(&order); err != nil {
		return false,err
	}
	return true, nil
}

func (u *Order) Create(order Order) (bool, error) {
	affected, err := x.Insert(order)
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("insert order fail")
	}
	return true, nil
}

func (u *Order) FindByOrderId(orderId string) (Order, error) {
	var order Order
	rel, err := x.Where("order_id = ?", orderId).Get(&order)
	if !rel || err != nil {
		log.Printf("FindByOrderId err:%v", err)
		return order, errors.New("not found")
	}
	return order, nil
}

func (u *Order) GetByOrderId(outTradeNo string) (Order, error) {
	var order Order
	rel, err := x.Where("out_trade_no = ?", outTradeNo).Get(&order)
	if !rel || err != nil {
		log.Printf("GetByOrderId err:%v", err)
		return order, errors.New("not found")
	}
	return order, nil
}