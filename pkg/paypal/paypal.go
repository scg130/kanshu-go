package paypal

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/plutov/paypal/v4"
	"github.com/sirupsen/logrus"
)

var (
	CLINET_ID  = ""
	SECRET     = ""
	PAYPAL_ENV = ""
)

type PaypalClient struct {
	Cli *paypal.Client
}

var PaypalCli *PaypalClient

func init() {
	CLINET_ID = os.Getenv("PAYPAL_CLIENT_ID")
	SECRET = os.Getenv("PAYPAL_SECRET")
	PAYPAL_ENV = os.Getenv("PAYPAL_ENV")
}

func New() *PaypalClient {
	if PaypalCli == nil {
		url := paypal.APIBaseLive
		if PAYPAL_ENV == "local" {
			url = paypal.APIBaseSandBox
		}
		c, err := paypal.NewClient(CLINET_ID, SECRET, url)
		if err != nil {
			panic(err)
		}
		ctx := context.Background()
		_, err = c.GetAccessToken(ctx)
		if err != nil {
			panic(err)
		}
		return &PaypalClient{
			Cli: c,
		}
	}
	return PaypalCli
}

func (c *PaypalClient) Create(ctx context.Context, amount string) (orderId string, payUrl string, err error) {
	order, err := c.Cli.CreateOrder(ctx, paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			paypal.PurchaseUnitRequest{
				ReferenceID: "ref-id",
				Amount:      &paypal.PurchaseUnitAmount{Value: amount, Currency: "USD"},
			},
		},
		&paypal.PaymentSource{},
		&paypal.ApplicationContext{ReturnURL: os.Getenv("PAYPAL_RETURN_URL")})
	if err != nil {
		return "", "", err
	}
	return order.ID, order.Links[1].Href, nil
}

func (c *PaypalClient) GetOrder(ctx context.Context, orderId string) (*paypal.Order, error) {
	order, err := c.Cli.GetOrder(ctx, orderId)
	if err != nil {
		return nil, err
	}
	return order, err
}

// 回调(可以利用上面的回调链接实现) orderId 就是返回的orderId := ctx.query("token")
func (c *PaypalClient) PaypalCallback(orderId string) (err error) {
	ctor := paypal.CaptureOrderRequest{}

	order, err := c.Cli.CaptureOrder(context.TODO(), orderId, ctor)
	if err != nil {
		logrus.Info(err, "打款失败")
		return err
	}

	//查看回调完成后订单状态是否支付完成。
	strByte, _ := json.Marshal(order)
	logrus.Info(string(strByte))
	if (*order).Status != "COMPLETED" {
		return errors.New("pay fail")
	}
	return nil
}
