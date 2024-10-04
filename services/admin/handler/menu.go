package handler

import (
	admin "admin/proto/admin"
	"admin/repo"
	"context"
	"errors"
	"time"

	log "github.com/micro/go-micro/v2/logger"
)

func (this *Admin) MenuListByIds(ctx context.Context, req *admin.IdsReq, rsp *admin.MenuListRep) error {
	log.Info("Received Admin.MenuListByIds request")
	menus, err := this.Menu.FindByIds(req.Ids)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "ok"
	for _, menu := range menus {
		rsp.Menu = append(rsp.Menu, &admin.Menu{
			Name:  menu.Name,
			Path:  menu.Path,
			Api:   menu.Api,
			Icon:  menu.Icon,
			Show:  menu.Show,
			State: menu.State,
			Pid:   menu.Pid,
			Id:    menu.Id,
		})
	}

	return nil
}

func makeTree(pid int64, menus []*repo.Menu) []*admin.Tree {
	var result []*admin.Tree
	for _, menu := range menus {
		if menu.Pid == pid {
			result = append(result, &admin.Tree{
				Id:       menu.Id,
				Label:    menu.Name,
				Name:     menu.Name,
				Children: makeTree(menu.Id, menus),
			})
		}
	}
	return result
}

func (this *Admin) MenuTree(ctx context.Context, req *admin.MenuReq, rsp *admin.MenuTreeRep) error {
	log.Info("Received Admin.MenuTree request")
	menus, err := this.Menu.AllList()
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "ok"
	rsp.Tree = makeTree(0, menus)
	return nil
}

func (this *Admin) MenuShowTree(ctx context.Context, req *admin.MenuReq, rsp *admin.MenuTreeRep) error {
	log.Info("Received Admin.MenuShowTree request")
	menus, err := this.Menu.ShowList()
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "ok"
	rsp.Tree = makeTree(0, menus)
	return nil
}

func (this *Admin) MenuList(ctx context.Context, req *admin.MenuReq, rsp *admin.MenuListRep) error {
	log.Info("Received Admin.Menu request")
	menus, total, err := this.Menu.List(req.Name, (req.Pagnation.Page-1)*req.Pagnation.PageSize, req.Pagnation.PageSize)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "ok"
	rsp.Pagnation = &admin.Pagnation{
		Page:     req.Pagnation.Page,
		PageSize: req.Pagnation.PageSize,
		Total:    total,
	}
	rsp.Menu = make([]*admin.Menu, 0)
	for _, menu := range menus {
		rsp.Menu = append(rsp.Menu, &admin.Menu{
			Name:    menu.Name,
			Path:    menu.Path,
			Api:     menu.Api,
			Icon:    menu.Icon,
			Show:    menu.Show,
			State:   menu.State,
			Pid:     menu.Pid,
			Id:      menu.Id,
			PidName: menu.PidName,
		})
	}
	return err
}

func (this *Admin) MenuAdd(ctx context.Context, req *admin.MenuAddReq, rsp *admin.MenuRep) error {
	log.Info("Received Admin.MenuAdd request")
	_, err := this.Menu.Create(repo.Menu{
		Name:       req.Name,
		Path:       req.Path,
		Api:        req.Api,
		Icon:       req.Icon,
		Show:       req.Show,
		State:      req.State,
		Pid:        req.Pid,
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

func (this *Admin) MenuEdit(ctx context.Context, req *admin.MenuEditReq, rsp *admin.MenuRep) error {
	log.Info("Received Admin.MenuEdit request")
	menu := &repo.Menu{
		Id:         int64(req.Id),
		Name:       req.Name,
		Path:       req.Path,
		Api:        req.Api,
		Icon:       req.Icon,
		Show:       req.Show,
		State:      req.State,
		Pid:        req.Pid,
		UpdateTime: time.Now(),
	}

	err := this.Menu.Update(menu)
	if err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "success"
	return nil
}

func (this *Admin) MenuDel(ctx context.Context, req *admin.MenuDelReq, rsp *admin.MenuRep) error {
	log.Info("Received Admin.MenuDel request")
	if req.Id == 0 {
		rsp.State = -1
		rsp.Msg = "id is required"
		return errors.New("id is required")
	}
	if err := this.Menu.Del(int32(req.Id)); err != nil {
		rsp.State = -1
		rsp.Msg = err.Error()
		return err
	}
	rsp.State = 0
	rsp.Msg = "success"
	return nil
}
