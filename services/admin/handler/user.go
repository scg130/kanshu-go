package handler

import (
	"admin/dto"
	"admin/repo"
	"context"
	"errors"
	"time"

	"github.com/scg130/tools"

	"github.com/micro/go-micro/v2/client"

	log "github.com/micro/go-micro/v2/logger"

	admin "admin/proto/admin"
)

type Admin struct {
	Cli  client.Client
	User repo.User
	Menu repo.Menu
	Role repo.Role
}

func (this *Admin) UserDel(ctx context.Context, req *admin.DelRequest, rsp *admin.DelResponse) error {
	log.Info("Received Admin.UserDel request")
	if req.Id == 0 {
		rsp.State = -1
		rsp.Msg = "id is required"
		return errors.New("id is required")
	}
	if err := this.User.Del(req.Id); err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "success"
	return nil
}

func (this *Admin) UserEdit(ctx context.Context, req *admin.EditRequest, rsp *admin.EditResponse) error {
	user := &repo.User{
		Id:         int64(req.Id),
		UserName:   req.Username,
		Email:      req.Email,
		State:      req.State,
		RoleIds:    req.RoleIds,
		UpdateTime: time.Now(),
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	err := this.User.Update(user)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "success"
	return nil
}

func (this *Admin) UserList(ctx context.Context, req *admin.UserListReq, rsp *admin.UserListRep) error {
	log.Info("Received Admin.Call request")
	users, err := this.User.List(req.UserName)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "ok"
	rsp.Users = make([]*admin.User, 0)
	for _, user := range users {
		rsp.Users = append(rsp.Users, &admin.User{
			Id:         int32(user.Id),
			UserName:   user.UserName,
			Email:      user.Email,
			Phone:      user.Phone,
			State:      user.State,
			CreateTime: user.CreateTime.Format(dto.DateTimeFormat),
			RoleIds:    user.RoleIds,
		})
	}
	return err
}

func (this *Admin) Login(ctx context.Context, req *admin.LoginRequest, rsp *admin.LoginResponse) error {
	log.Info("Received Admin.Call request")
	user, err := this.User.FindByName(req.UserName)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}

	if tools.CompareHashAndPasswd(req.Password, user.Password) {
		rsp.State = 0
		rsp.Msg = "ok"
		rsp.Data = &admin.AdminUserInfo{
			Id:      int32(user.Id),
			Name:    user.UserName,
			RoleIds: user.RoleIds,
		}
		return nil
	}
	return errors.New("login failure")
}

func (this *Admin) Reg(ctx context.Context, req *admin.RegRequest, rsp *admin.RegResponse) error {
	log.Info("Received Admin.Call request")
	// micro.NewEvent("test", this.Cli)
	_, err := this.User.Create(repo.User{
		UserName:   req.Username,
		Password:   req.Password,
		Email:      req.Email,
		Phone:      req.Phone,
		State:      req.State,
		UpdateTime: time.Now(),
		CreateTime: time.Now(),
	})
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "ok"
	return err
}
