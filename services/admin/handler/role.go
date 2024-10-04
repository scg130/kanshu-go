package handler

import (
	admin "admin/proto/admin"
	"admin/repo"
	"context"
	"errors"
	"time"

	log "github.com/micro/go-micro/v2/logger"
)

func (this *Admin) FindMenuIdsByRoleIds(ctx context.Context, req *admin.IdsReq, rsp *admin.MenuIdsRep) error {
	log.Info("Received Admin.RoleList request")
	roles, err := this.Role.FindByIds(req.Ids)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "ok"
	for _, role := range roles {
		rsp.Ids = append(rsp.Ids, role.MenuIds...)
	}
	return err
}

func (this *Admin) RoleList(ctx context.Context, req *admin.RoleReq, rsp *admin.RoleListRep) error {
	log.Info("Received Admin.RoleList request")
	roles, err := this.Role.List(req.Name)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "ok"
	rsp.Role = make([]*admin.Role, 0)
	for _, role := range roles {
		rsp.Role = append(rsp.Role, &admin.Role{
			Name:    role.Name,
			Id:      role.Id,
			MenuIds: role.MenuIds,
		})
	}
	return err
}

func (this *Admin) RoleAdd(ctx context.Context, req *admin.RoleAddReq, rsp *admin.RoleRep) error {
	log.Info("Received Admin.RoleAdd request")
	_, err := this.Role.Create(repo.Role{
		Name:       req.Name,
		MenuIds:    req.MenuIds,
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
	return nil
}

func (this *Admin) RoleEdit(ctx context.Context, req *admin.RoleEditReq, rsp *admin.RoleRep) error {
	log.Info("Received Admin.RoleEdit request")
	role := &repo.Role{
		Id:         int64(req.Id),
		Name:       req.Name,
		MenuIds:    req.MenuIds,
		UpdateTime: time.Now(),
	}

	err := this.Role.Update(role)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "success"
	return nil
}

func (this *Admin) RoleDel(ctx context.Context, req *admin.RoleDelReq, rsp *admin.RoleRep) error {
	log.Info("Received Admin.RoleDel request")
	if req.Id == 0 {
		rsp.State = -1
		rsp.Msg = "id is required"
		return errors.New("id is required")
	}
	if err := this.Role.Del(int32(req.Id)); err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "success"
	return nil
}
