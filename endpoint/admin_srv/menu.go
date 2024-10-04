package admin_srv

import (
	"kanshu/dto/admin"
	go_micro_service_admin "kanshu/proto/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// desc: 添加菜单
// tags:	menu
// param:	body	body	admin.MenuRequest	true	"菜单信息"
// produce:	appication/json
// responses: 200 "ok" "admin.Resp{code:0,data:[]*go_micro_service_admin.Menu}"
//
//swagger:route POST /admin/menu/add menu admin addMenu
func (this *Admin) AddMenu(ctx *gin.Context) {
	var req admin.MenuRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	addReq := &go_micro_service_admin.MenuAddReq{
		Name:  req.Name,
		Icon:  req.Icon,
		Path:  req.Path,
		Api:   req.Api,
		State: req.State,
		Show:  req.Show,
		Pid:   req.Pid,
	}
	rep, err := this.AdminCli.MenuAdd(ctx, addReq)
	if err != nil || rep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
	})
}

func (this *Admin) DelMenu(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	rep, err := this.AdminCli.MenuDel(ctx, &go_micro_service_admin.MenuDelReq{
		Id: int64(id),
	})
	if err != nil || rep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
	})
}

func (this *Admin) EditMenu(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req admin.MenuRequest
	if err := ctx.Bind(&req); err != nil || id == 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	editReq := &go_micro_service_admin.MenuEditReq{
		Id:    int64(id),
		Name:  req.Name,
		Icon:  req.Icon,
		Path:  req.Path,
		Api:   req.Api,
		State: req.State,
		Show:  req.Show,
		Pid:   req.Pid,
	}

	rep, err := this.AdminCli.MenuEdit(ctx, editReq)
	if err != nil || rep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
	})
}

func (this *Admin) MenuTree(ctx *gin.Context) {
	rep, _ := this.AdminCli.MenuTree(ctx, &go_micro_service_admin.MenuReq{})
	if rep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: -1,
			Msg:  rep.Msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
		Data: rep.Tree,
	})
}

func (this *Admin) MenuShowTree(ctx *gin.Context) {
	rep, _ := this.AdminCli.MenuShowTree(ctx, &go_micro_service_admin.MenuReq{})
	if rep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: -1,
			Msg:  rep.Msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
		Data: rep.Tree,
	})
}

func (this *Admin) Menus(ctx *gin.Context) {
	var req admin.MenusRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	rep, _ := this.AdminCli.MenuList(ctx, &go_micro_service_admin.MenuReq{
		Name: req.Name,
		Pagnation: &go_micro_service_admin.Pagnation{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	})
	if rep.State != 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: -1,
			Msg:  rep.Msg,
		})
		return
	}
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
		Data: rep.Menu,
		Pagnation: admin.Pagnation{
			Page:     rep.Pagnation.Page,
			PageSize: rep.Pagnation.PageSize,
			Total:    rep.Pagnation.Total,
		},
	})
}
