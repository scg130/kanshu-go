package admin_srv

import (
	"kanshu/dto/admin"
	go_micro_service_admin "kanshu/proto/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/scg130/tools"
)

func (this *Admin) DelUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	rep, err := this.AdminCli.UserDel(ctx, &go_micro_service_admin.DelRequest{
		Id: int32(id),
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

func (this *Admin) EditUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req admin.RegRequest
	if err := ctx.Bind(&req); err != nil || id == 0 {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
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
	editReq := &go_micro_service_admin.EditRequest{
		Id:       int32(id),
		Username: req.UserName,
		Email:    req.Email,
		Phone:    req.Phone,
		State:    req.State,
		RoleIds:  req.RoleIds,
	}
	if req.Passwd != "" && req.ConfirmPasswd != "" {
		pwd, _ := tools.GeneratePasswd(req.Passwd)
		editReq.Password = pwd
	}

	rep, err := this.AdminCli.UserEdit(ctx, editReq)
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

func (this *Admin) Users(ctx *gin.Context) {
	var req admin.ListInfoRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	rep, err := this.AdminCli.UserList(ctx, &go_micro_service_admin.UserListReq{
		UserName: req.UserName,
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
		Data: rep.Users,
	})
}
