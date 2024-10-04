package admin_srv

import (
	"fmt"
	"kanshu/dto/admin"
	"kanshu/env"
	go_micro_service_admin "kanshu/proto/admin"
	selfwrappers "kanshu/wrappers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scg130/tools"
	"github.com/scg130/tools/wrappers"
)

type Admin struct {
	AdminCli go_micro_service_admin.AdminUserService
}

const SRV_NAME = "go.kanshu.service.admin"

var adminSrv *Admin

func NewAdminSrv() *Admin {
	if adminSrv == nil {
		adminSrv = &Admin{
			AdminCli: go_micro_service_admin.NewAdminUserService(SRV_NAME, tools.GetMicroClient(
				SRV_NAME,
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
		}
	}
	return adminSrv
}

func (u *Admin) UserInfo(ctx *gin.Context) {
	uinfo, _ := ctx.Get("authData")
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
		Data: gin.H{
			"username": uinfo.(map[string]interface{})["name"],
			"avatar":   "https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png",
		},
	})
}

// @Summary 登录
// @Description 登录
// @Tags 用户中心
// @Produce json
// @Param body body admin.LoginRequest true "body参数"
// @Success 200 {object}  dto.Resp{data=admin.LoginResp}
// @Router /admin/login [post]
func (u *Admin) Login(ctx *gin.Context) {
	var req admin.LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: 1,
			Msg:  "bad params",
		})
		return
	}

	resp, err := u.AdminCli.Login(ctx, &go_micro_service_admin.LoginRequest{UserName: req.UserName, Password: req.Passwd})
	if resp.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: 1,
			Msg:  resp.Msg,
		})
		return
	}

	mrep, _ := u.AdminCli.FindMenuIdsByRoleIds(ctx, &go_micro_service_admin.IdsReq{
		Ids: resp.Data.RoleIds,
	})
	if mrep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: 1,
			Msg:  resp.Msg,
		})
		return
	}

	mlrep, _ := u.AdminCli.MenuListByIds(ctx, &go_micro_service_admin.IdsReq{
		Ids: mrep.Ids,
	})
	if mlrep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: 1,
			Msg:  resp.Msg,
		})
		return
	}
	paths := make([]string, 0)
	apis := make([]string, 0)
	for _, menu := range mlrep.Menu {
		paths = append(paths, menu.Path)
		apis = append(apis, menu.Api)
	}
	fmt.Println(resp.Data)
	token, err := tools.GenerateToken(env.JwtConf.AdminJwtSecret, map[string]interface{}{
		"udata": resp.Data,
		"paths": paths,
		"apis":  apis,
	}, time.Duration(env.TokenExpire*time.Second))
	if err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: 1,
			Msg:  err.Error(),
		})
		return
	}
	ctx.SetCookie("token", token, env.TokenExpire, "/", ctx.Request.Host, false, false)
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
		Data: admin.LoginResp{Token: token},
	})
}

func (u *Admin) Reg(ctx *gin.Context) {
	var req admin.RegRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: 1,
			Msg:  err.Error(),
		})
		return
	}

	if req.ConfirmPasswd != req.Passwd {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: 1,
			Msg:  "两次密码不一致",
		})
		return
	}
	pwd, _ := tools.GeneratePasswd(req.Passwd)
	state := req.State
	if state == 0 {
		state = 2
	}
	rep, _ := u.AdminCli.Reg(ctx, &go_micro_service_admin.RegRequest{
		Username: req.UserName,
		Password: pwd,
		Email:    req.Email,
		Phone:    req.Phone,
		State:    state,
	})

	if rep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: 1,
			Msg:  rep.Msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
	})
}
