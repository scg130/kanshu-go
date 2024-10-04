package alipay

type CreateResp struct {
	Qrcode        string `json:"qrcode"`
	OutTradeOrder string `json:"out_trade_order"`
}

type SuccessResp struct {
	Success     bool   `json:"success"`
	OutTradeNo  string `json:"out_trade_no"`
	TotalAmount string `json:"total_amount"`
	PayAmount   string `json:"pay_amount"`
	BuyerId     int    `json:"buyer_id"`
	Subject     string `json:"subject"`
	TradeNo     string `json:"trade_no"`
}
