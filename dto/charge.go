package dto

type CreateOrderReq struct {
	Amount int64 `json:"amount" form:"amount"`
	//充值渠道 alipay 支付宝 paypal
	Channel string `json:"channel" form:"channel"`
	//支付商品 vip 会员 amount 充值金额 chapter 购买章节
	Subject string `json:"subject" form:"subject"`

	SubjectId int64 `json:"subject_id" form:"subject_id"`
}

type CreateOrderRsp struct {
	Qrcode  []byte `json:"qrcode"`
	OrderId string `json:"order_id"`
	//充值渠道 alipay 支付宝 paypal
	Channel   string `json:"channel"`
	PaypalUrl string `json:"paypal_url"`
}
