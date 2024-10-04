package repo

import (
	"errors"
	"log"
	"time"
)

type Role struct {
	Id         int64     `xorm:" pk autoincr INT(11)"`
	Name       string    `xorm:"not null default '' unique(name) comment('菜单名称') VARCHAR(128)"`
	MenuIds    []int64   `xorm:"not null comment('菜单ids') blob"`
	UpdateTime time.Time `xorm:"updated_at not null DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	role := new(Role)
	if isExist, _ := x.IsTableExist(role); !isExist {
		if err := x.Sync2(role); err != nil {
			log.Fatal("sync tables err:%v", err)
		}
	}
}

func (u *Role) Del(id int32) error {
	affected, err := x.Where("id = ?", id).Delete(new(Role))
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (u *Role) Update(role *Role) error {
	affected, err := x.Where("id = ?", role.Id).Update(role)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (u *Role) Create(role Role) (bool, error) {
	affected, err := x.Insert(role)
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("create user fail")
	}
	return true, nil
}

func (u *Role) GetById(id int64) (*Role, error) {
	role := new(Role)
	has, err := x.Where("id = ?", id).Get(role)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("not found")
	}
	return role, nil
}

func (u *Role) FindByIds(ids []int64) (roles []*Role, err error) {
	err = x.In("id", ids).Find(&roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u *Role) List(name string) (roles []*Role, err error) {
	query := x.Table("role")
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err = query.Find(&roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
