package admin_srv

import (
	"kanshu/dto/admin"
	go_micro_service_admin "kanshu/proto/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (this *Admin) AddRole(ctx *gin.Context) {
	var req admin.RoleRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	addReq := &go_micro_service_admin.RoleAddReq{
		Name:    req.Name,
		MenuIds: req.MenuIds,
	}
	rep, err := this.AdminCli.RoleAdd(ctx, addReq)
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

func (this *Admin) DelRole(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	rep, err := this.AdminCli.RoleDel(ctx, &go_micro_service_admin.RoleDelReq{
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

func (this *Admin) EditRole(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req admin.RoleRequest
	if err := ctx.Bind(&req); err != nil || id == 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	editReq := &go_micro_service_admin.RoleEditReq{
		Id:      int64(id),
		Name:    req.Name,
		MenuIds: req.MenuIds,
	}

	rep, err := this.AdminCli.RoleEdit(ctx, editReq)
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

func (this *Admin) Roles(ctx *gin.Context) {
	var req admin.RolesRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	rep, _ := this.AdminCli.RoleList(ctx, &go_micro_service_admin.RoleReq{
		Name: req.Name,
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
		Data: rep.Role,
	})
}
