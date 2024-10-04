package repo

import (
	"errors"
	"log"
	"time"
)

type Menu struct {
	Id         int64     `xorm:" pk autoincr INT(11)"`
	Name       string    `xorm:"not null default '' unique(name) comment('菜单名称') VARCHAR(128)"`
	Path       string    `xorm:"not null default '' comment('url路径') VARCHAR(128)"`
	Api        string    `xorm:"not null default '' comment('api uri') varchar(128)"`
	Icon       string    `xorm:"not null default '' comment('icon') VARCHAR(128)"`
	Pid        int64     `xorm:"not null default 0 comment('父级id') int(11)"`
	State      int32     `xorm:"not null default 2 comment('是否禁用 2 否 1 是') TINYINT(1)"`
	Show       int32     `xorm:"not null default 1 comment('是否在菜单导航栏显示 2 否 1 是') TINYINT(1)"`
	UpdateTime time.Time `xorm:"updated_at not null DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

type MenuData struct {
	Id      int64
	Name    string
	Path    string
	Api     string
	Icon    string
	Pid     int64
	State   int32
	Show    int32
	PidName string
}

func init() {
	menu := new(Menu)
	if isExist, _ := x.IsTableExist(menu); !isExist {
		if err := x.Sync2(menu); err != nil {
			log.Fatal("sync tables err:%v", err)
		}
	}
}

func (u *Menu) Del(id int32) error {
	affected, err := x.Where("id = ?", id).Delete(new(Menu))
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (u *Menu) Update(menu *Menu) error {
	affected, err := x.Where("id = ?", menu.Id).Update(menu)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (u *Menu) Create(menu Menu) (bool, error) {
	affected, err := x.Insert(menu)
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("create user fail")
	}
	return true, nil
}

func (u *Menu) GetById(id int64) (*Menu, error) {
	menu := new(Menu)
	has, err := x.Where("id = ?", id).Get(menu)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("not found")
	}
	return menu, nil
}

func (u *Menu) FindByIds(ids []int64) (menus []*Menu, err error) {
	err = x.In("id", ids).Find(&menus)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (u *Menu) List(name string, offset, limit int64) (menus []*MenuData, total int64, err error) {
	query := x.Table("menu").Join("LEFT", "menu as m1", "m1.id=menu.pid")
	if name != "" {
		query = query.Where("menu.name like ?", "%"+name+"%")
	}
	total, err = query.Select("menu.*, m1.name as pid_name").Limit(int(limit), int(offset)).FindAndCount(&menus)
	if err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

func (u *Menu) AllList() (menus []*Menu, err error) {
	err = x.Where("state = ?", 2).Find(&menus)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (u *Menu) ShowList() (menus []*Menu, err error) {
	err = x.Where("`show` = ? and state = ?", 1, 2).Find(&menus)
	if err != nil {
		return nil, err
	}
	return menus, nil
}
