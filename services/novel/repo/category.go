package repo

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Category struct {
	Id         int64     `xorm:" pk autoincr INT(11)"`
	Name       string    `xorm:"not null unique(name) comment('分类名称') VARCHAR(255)"`
	Sort       int32     `xorm:"not null default 0 comment('排序') int(11)"`
	Channel    int32     `xorm:"not null default 0 comment('频道 1男生 2女生') tinyint"`
	IsShow     int32     `xorm:"not null default 0 comment('是否显示 1显示 2不显示 ') tinyint"`
	UpdateTime time.Time `xorm:"updated_at not null  DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	category := new(Category)
	if isExist, _ := x.IsTableExist(category); !isExist {
		if err := x.Sync2(category); err != nil {
			log.Fatal(fmt.Sprintf("sync tables err:%v", err))
		}
	}
}

func (u *Category) Del(id int32) error {
	affected, err := x.Where("id = ?", id).Delete(new(Category))
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (u *Category) Update(cate *Category) error {
	_, err := x.Where("id = ?", cate.Id).Update(cate)
	if err != nil {
		return err
	}

	return nil
}

func (u *Category) Create(cate Category) (bool, error) {
	affected, err := x.Insert(cate)
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("create category fail")
	}
	return true, nil
}

func (c *Category) Get(page, size, isShow int, name string, sex int32) ([]Category, int64, error) {
	data := make([]Category, 0)
	query := x.Where("1=1")
	if isShow != -1 {
		query = query.Where("is_show = ?", isShow)
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if sex != 0 {
		query = query.Where("channel = ?", sex)
	}
	total, err := query.Limit(size, size*(page-1)).FindAndCount(&data)

	return data, total, err
}

func (c *Category) GetOne(cateId int) (Category, error) {
	var cate Category
	has, err := x.Where("id=?", cateId).Get(&cate)
	if !has {
		return cate, errors.New("no has")
	}
	return cate, err
}
