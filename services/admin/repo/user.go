package repo

import (
	"errors"
	"log"
	"time"
)

type User struct {
	Id         int64     `xorm:" pk autoincr INT(11)"`
	UserName   string    `xorm:"not null default '' unique(user_name) comment('用户名') VARCHAR(128)"`
	Password   string    `xorm:"not null default '' comment('密码') VARCHAR(255)"`
	Email      string    `xorm:"not null default '' unique(email) comment('邮箱') VARCHAR(128)"`
	Phone      int64     `xorm:"not null default 0 unique(phone) comment('手机号') bigint(20)"`
	State      int32     `xorm:"not null default 2 comment('是否锁定 2 否 1 是') TINYINT(1)"`
	RoleIds    []int64   `xorm:"not null comment('role_ids') blob"`
	UpdateTime time.Time `xorm:"updated_at not null DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	user := new(User)
	if isExist, _ := x.IsTableExist(user); !isExist {
		if err := x.Sync2(user); err != nil {
			log.Fatal("sync tables err:%v", err)
		}
	}
}

func (u *User) Del(id int32) error {
	affected, err := x.Where("id = ?", id).Delete(new(User))
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (u *User) Update(user *User) error {
	affected, err := x.Where("id = ?", user.Id).Update(user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (u *User) Create(user User) (bool, error) {
	affected, err := x.Insert(user)
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("create user fail")
	}
	return true, nil
}

func (u *User) FindByName(username string) (*User, error) {
	var user User
	rel, err := x.Where("user_name = ? and state = ?", username, 2).Get(&user)
	if !rel || err != nil {
		log.Printf("user find by name err:%v", err)
		return nil, errors.New("not found")
	}
	return &user, nil
}

func (u *User) List(username string) (users []*User, err error) {
	query := x.Where("1=1")
	if username != "" {
		query = query.Where("user_name like ?", "%"+username+"%")
	}
	err = query.Find(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
