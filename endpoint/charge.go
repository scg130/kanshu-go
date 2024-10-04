package endpoint

import (
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"kanshu/dto"
	"kanshu/env"
	"kanshu/pkg/alipay"
	paypal_pkg "kanshu/pkg/paypal"
	go_micro_service_charge "kanshu/proto/charge"
	go_micro_service_wallet "kanshu/proto/wallet"
	selfwrappers "kanshu/wrappers"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scg130/tools"
	"github.com/scg130/tools/wrappers"
	"github.com/sirupsen/logrus"
)

const OrderPrefixFormat = "20060102150405"

type Charge struct {
	chargeCli go_micro_service_charge.ChargeService
	walletCli go_micro_service_wallet.WalletService
	aliPayCli *alipay.AliPay
	paypalCli *paypal_pkg.PaypalClient
}

var chargeSrv *Charge

func NewChargeSrv() *Charge {
	if chargeSrv == nil {
		chargeSrv = &Charge{
			chargeCli: go_micro_service_charge.NewChargeService("go.kanshu.service.charge", tools.GetMicroClient(
				"go.kanshu.service.charge",
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
			aliPayCli: alipay.NewAliPay(env.AliPayConf.AppId, env.AliPayConf.RsaPrivateKey, env.AliPayConf.RsaPublicKey, env.AliPayConf.NotifyUrl),
			walletCli: go_micro_service_wallet.NewWalletService("go.kanshu.service.wallet", tools.GetMicroClient(
				"go.kanshu.service.wallet",
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
			paypalCli: paypal_pkg.New(),
		}
	}
	return chargeSrv
}

func (self *Charge) generateOrderNo() string {
	return fmt.Sprintf("%s%d", time.Now().Format(OrderPrefixFormat), rand.Intn(999999-100000)+100000)
}

func (self *Charge) QueryOrder(ctx *gin.Context) {
	orderId := ctx.Query("order_id")
	if orderId == "" {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "not found",
		})
		return
	}
	_, isExist := ctx.Get("authData")
	if !isExist {
		ctx.String(http.StatusUnauthorized, "failure")
		return
	}

	rep, err := self.chargeCli.QueryOrder(ctx, &go_micro_service_charge.QueryReq{
		OrderId: orderId,
	})
	if err != nil || rep.State != 1 || rep.Status != int32(go_micro_service_charge.StateType_STATE_PAY_SUCCESS) {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
		})
		return
	}
	if rep.State == 1 && rep.Status == int32(go_micro_service_charge.StateType_STATE_PAY_SUCCESS) {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: 0,
		})
		return
	}
}

// @Summary 创建订单
// @Description 创建订单
// @Tags 订单
// @Produce json
// @Param body body dto.CreateOrderReq true "body参数"
// @Success 200 {object}  dto.Resp{data=dto.CreateOrderRsp}
// @Router /charge/create [post]
func (self *Charge) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderReq
	var err error
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "bad params")
		return
	}
	userInfo := GetUserInfo(ctx)
	subject := fmt.Sprintf("%d-%s", req.Amount, req.Subject)
	tradeOrderId := self.generateOrderNo()
	thirdOrderNo := ""
	paypalUrl := ""
	var imgData []byte
	if req.Channel == "alipay" {
		res, err := self.aliPayCli.CreateOrder(tradeOrderId, fmt.Sprintf("%.2f", float64(req.Amount)/100.00), subject)
		//qrcode.WriteFile(qrcodeStr.(string),qrcode.Medium,256,"./qr.png")
		if err != nil {
			ctx.JSON(http.StatusOK, dto.Resp{
				Code: -1,
				Msg:  err.Error(),
			})
			return
		}
		thirdOrderNo = res.OutTradeOrder
		imgData, _ = qrcode.Encode(res.Qrcode, qrcode.Medium, 256)
	} else if req.Channel == "paypal" {
		thirdOrderNo, paypalUrl, err = self.paypalCli.Create(ctx, fmt.Sprintf("%.2f", float64(req.Amount)/100.00))
		if err != nil {
			ctx.JSON(http.StatusOK, dto.Resp{
				Code: -1,
				Msg:  err.Error(),
			})
			return
		}
	}

	//创建订单
	rep, err := self.chargeCli.Create(ctx, &go_micro_service_charge.ChargeReq{
		Uid:          userInfo.UserId,
		Amount:       req.Amount,
		Channel:      req.Channel,
		Subject:      subject,
		SubjectId:    req.SubjectId,
		State:        go_micro_service_charge.StateType_STATE_NORMAL,
		ThirdOrderNo: thirdOrderNo,
		OrderId:      fmt.Sprintf("ch_%s", tradeOrderId),
	})
	if err != nil {
		log.Println(err)
	}
	if err != nil || rep.State != 1 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "rpc server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Data: dto.CreateOrderRsp{
			Qrcode:    imgData,
			OrderId:   rep.OrderId,
			Channel:   req.Channel,
			PaypalUrl: paypalUrl,
		},
	})
	//ctx.Writer.Write(imgData)
}

func (c *Charge) USDCallback(ctx *gin.Context) {
	orderId := ctx.Query("token")
	defer ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/package/wallet/recharge", env.AppConf.FrontUrl))

	rep, err := c.chargeCli.QueryOrderByThirdOrderId(ctx, &go_micro_service_charge.ChargeReq{
		ThirdOrderNo: orderId,
	})
	if err != nil || rep.State != 1 {
		logrus.Error(errors.New(fmt.Sprintf("get order err:%s", orderId)))
		return
	}

	if rep.Status == int32(go_micro_service_charge.StateType_STATE_PAY_SUCCESS) {
		logrus.Info(fmt.Sprintf("order_id:%s  callbak retry", orderId))
		return
	}

	err = c.paypalCli.PaypalCallback(orderId)
	if err != nil {
		logrus.Error(err)
		return
	}
	order, err := c.paypalCli.GetOrder(ctx, orderId)
	if err != nil {
		logrus.Error(err)
		return
	}
	amount := order.PurchaseUnits[0].Amount.Value
	chargeRsp, err := c.chargeCli.ChargeSuccess(ctx, &go_micro_service_charge.ChargeReq{
		ThirdOrderNo: orderId,
		State:        go_micro_service_charge.StateType_STATE_PAY_SUCCESS,
		Subject:      fmt.Sprintf("subject%s", amount),
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	if chargeRsp.State != 1 {
		log.Println("success", chargeRsp)
		return
	}
	famount, _ := strconv.ParseFloat(amount, 64)
	walletRsp, err := c.walletCli.Change(ctx, &go_micro_service_wallet.WalletReq{
		Uid:     chargeRsp.UserId,
		Amount:  int64(famount * 600),
		Type:    go_micro_service_wallet.Type_STATE_CHARGE,
		OrderId: chargeRsp.OrderIdInt,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	if err != nil || walletRsp.State != 1 {
		logrus.Println(err, walletRsp)
		return
	}
	return
}

func (c *Charge) Callback(ctx *gin.Context) {
	rsp := c.aliPayCli.Callback(ctx.Request, ctx.Writer)
	if !rsp.Success {
		logrus.Println("callback error", rsp)
		return
	}

	chargeRsp, err := c.chargeCli.ChargeSuccess(ctx, &go_micro_service_charge.ChargeReq{
		ThirdOrderNo: rsp.OutTradeNo,
		State:        go_micro_service_charge.StateType_STATE_PAY_SUCCESS,
		Subject:      rsp.Subject,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	if chargeRsp.State != 1 {
		log.Println("success", chargeRsp)
		return
	}

	amount, _ := strconv.ParseFloat(rsp.TotalAmount, 64)
	walletRsp, err := c.walletCli.Change(ctx, &go_micro_service_wallet.WalletReq{
		Uid:     chargeRsp.UserId,
		Amount:  int64(amount * 100),
		Type:    go_micro_service_wallet.Type_STATE_CHARGE,
		OrderId: chargeRsp.OrderIdInt,
	})
	if err != nil {
		logrus.Error(err)
	}
	if err != nil || walletRsp.State != 1 {
		logrus.Println(err, walletRsp)
	}
	return
}
