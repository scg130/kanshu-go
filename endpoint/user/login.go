package user

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/ilylx/gconv"
	"github.com/scg130/tools"
	"github.com/scg130/tools/wrappers"
	"github.com/sirupsen/logrus"
	"kanshu/dto"
	"kanshu/env"
	go_micro_service_novel "kanshu/proto/novel"
	go_micro_service_user "kanshu/proto/user"
	go_micro_service_wallet "kanshu/proto/wallet"
	selfwrappers "kanshu/wrappers"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	UserCli   go_micro_service_user.UserCenterService
	WalletCli go_micro_service_wallet.WalletService
	NovelCli  go_micro_service_novel.NovelSrvService
}

var userSrv *User

func NewUserSrv() *User {
	if userSrv == nil {
		userSrv = &User{
			UserCli: go_micro_service_user.NewUserCenterService("go.kanshu.service.user", tools.GetMicroClient(
				"go.kanshu.service.user",
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
			WalletCli: go_micro_service_wallet.NewWalletService("go.kanshu.service.wallet", tools.GetMicroClient(
				"go.kanshu.service.wallet",
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
			NovelCli: go_micro_service_novel.NewNovelSrvService("go.kanshu.service.novel", tools.GetMicroClient(
				"go.kanshu.service.novel",
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
		}
	}
	return userSrv
}

// @Summary 注册
// @Description 注册
// @Tags 用户中心
// @Produce json
// @Param body body dto.UserRequest true "body参数"
// @Success 200 {object}  dto.Resp{}
// @Failure 500 {string} string "服务异常"
// @Router /app/Login/registerCode [post]
func (u *User) RegisterCode(ctx *gin.Context) {
	var req dto.RegUserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	if req.Phone < 10000000000 || req.Phone > 99999999999 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "phone is valid",
		})
		return
	}

	if len(req.Passwd) < 6 || len(req.Passwd) > 20 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "password length is between 6 and 20",
		})
		return
	}
	if req.Passwd != req.PasswdConfirm {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "confirm_passwd isn't same of password",
		})
		return
	}
	if !captcha.VerifyString(req.Id, req.Code) {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "code is invalid",
		})
		return
	}

	rsp, err := u.UserCli.Register(ctx, &go_micro_service_user.Request{
		Phone:  req.Phone,
		Avatar: fmt.Sprintf("%s/book/imgs/avatar%d.jpg", env.AppConf.Domain, rand.Intn(5)+1),
		Passwd: req.Passwd})
	if err != nil || rsp == nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "register failure",
		})
		return
	}
	if rsp.Code == -2 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "phone 已经被注册",
		})
		return
	}
	if rsp.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "register failure",
		})
		return
	}
	token, err := tools.GenerateToken(env.JwtConf.Secret, rsp.Data, time.Duration(time.Second*env.TokenExpire))
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusOK, "login fail")
		return
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: map[string]interface{}{
			"token":   token,
			"user_id": rsp.Data.UserId,
			"phone":   rsp.Data.Phone,
		},
	})
}

// @Summary 注册
// @Description 注册
// @Tags 用户中心
// @Produce json
// @Param body body dto.UserRequest true "body参数"
// @Success 200 {object}  dto.Resp{}
// @Failure 500 {string} string "服务异常"
// @Router /app/Login/ForgetPwd [post]
func (u *User) ForgetPwd(ctx *gin.Context) {
	var req dto.ForgetRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	if req.Phone < 10000000000 || req.Phone > 99999999999 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "phone is valid",
		})
		return
	}

	if len(req.Passwd) < 6 || len(req.Passwd) > 20 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "password length is between 6 and 20",
		})
		return
	}
	if req.Code != "123456" {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "验证码错误",
		})
		return
	}
	uRsp, err := u.UserCli.Find(ctx, &go_micro_service_user.Request{Phone: req.Phone})
	fmt.Println(uRsp, err)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	if uRsp.Code == -2 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  uRsp.Msg,
		})
		return
	}
	rsp, err := u.UserCli.ResetPasswd(ctx, &go_micro_service_user.Request{UserId: uRsp.Data.UserId, Passwd: req.Passwd})
	if err != nil || rsp == nil || rsp.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: nil,
	})
}

// @Summary 登录
// @Description 登录
// @Tags 用户中心
// @Produce json
// @Param body body dto.UserRequest true "body参数"
// @Success 200 {object}  dto.Resp{data=dto.LoginResp}
// @Router /app/Login/phoneLogin [post]
func (u *User) PhoneLogin(ctx *gin.Context) {
	var req dto.UserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	resp, err := u.UserCli.Login(ctx, &go_micro_service_user.Request{Phone: req.Phone, Passwd: req.Passwd})
	if err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	if resp.Code == -2 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  resp.Msg,
		})
		return
	}
	if resp.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	token, err := tools.GenerateToken(env.JwtConf.Secret, resp.Data, time.Duration(time.Second*env.TokenExpire))
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusOK, "login fail")
		return
	}
	ctx.SetCookie("token", token, env.TokenExpire, "/", "/", false, false)
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: dto.LoginResp{
			Token:          token,
			UserId:         resp.Data.UserId,
			Username:       gconv.String(resp.Data.Phone),
			Phone:          gconv.String(resp.Data.Phone),
			InvitationCode: "test_code",
			Sex:            1,
		},
	})
}

// @Summary 查找用户
// @Description 通过手机号查找用户
// @Tags 用户中心
// @Param phone path int true "手机号"
// @Success 200 {object}  dto.Resp{data=dto.UserInfo}
// @Failure 500 {string} string "服务异常"
// @Router /app/UserVip/selectUserVip?user_id= [get]
func (u *User) UserInfo(ctx *gin.Context) {
	userId := ctx.Query("user_id")

	rsp, err := u.UserCli.GetById(ctx, &go_micro_service_user.Request{UserId: gconv.Int64(userId)})
	if err != nil || rsp.Code != 0 || rsp.Data == nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "not found userInfo by userId " + gconv.String(userId),
		})
		return
	}

	wRsp, err := u.WalletCli.GetOne(ctx, &go_micro_service_wallet.WalletReq{Uid: gconv.Int64(userId)})
	if err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "not found wallet by userId " + gconv.String(userId),
		})
		return
	}

	nRsp, err := u.NovelCli.GetNoteNum(ctx, &go_micro_service_novel.NoteNumReq{
		UserId: gconv.Int32(userId),
	})
	if err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "fail find notes num by userId " + gconv.String(userId),
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: dto.UserInfo{
			UserId:         rsp.Data.UserId,
			Username:       gconv.String(rsp.Data.Phone),
			Phone:          gconv.String(rsp.Data.Phone),
			InvitationCode: "test_code",
			Money:          gconv.Float64(float64(wRsp.AvailableBalance) / 100.00),
			InviteMoney:    1.23,
			IntegralNum:    0,
			NoteNum:        int(nRsp.Num),
			LoveNum:        int(nRsp.JoinNum),
			Avatar:         rsp.Data.Avatar,
		},
	})
}
