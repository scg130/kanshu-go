package repo

import (
	"errors"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type Comment struct {
	Id         int64     `xorm:" pk autoincr INT(11)"`
	UserId     int64     `xorm:"not null default 0 comment('用户id') int"`
	NovelId    int       `xorm:"not null default 0 comment('小说id') int"`
	Content    string    `xorm:"not null default '' comment('内容') varchar(5000)"`
	GoodsNum   int       `xorm:"not null default 0 comment('点赞') int"`
	UpdateTime time.Time `xorm:"updated_at not null  DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	wlog := new(Comment)
	if isExist, _ := x.IsTableExist(wlog); !isExist {
		if err := x.Sync2(wlog); err != nil {
			log.Fatal("sync tables err:%v", err)
		}
	}
}

func (u *Comment) GetByID(id int32) (comment Comment, err error) {
	_, err = x.ID(id).Get(&comment)
	if err != nil {
		return
	}
	return
}

func (u *Comment) Update(comment Comment) error {
	_, err := x.ID(comment.Id).Update(&comment)
	if err != nil {
		return err
	}
	return nil
}

func (u *Comment) Create(wlog Comment) (bool, error) {
	affected, err := x.Insert(wlog)
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("insert comment fail")
	}
	return true, nil
}

func (u *Comment) FindByNovelID(novelID int, offset, limit int32) ([]*CommentData, int64, error) {
	wlogs := make([]*CommentData, 0)
	query := x.Table("comment").
		Join("left", "user", "comment.user_id = user.id")
	total, err := query.Where("comment.novel_id = ?", novelID).
		Select(`comment.goods_num as goods_num,comment.id as id,comment.content as content,comment.created_at as create_time,user.phone as username,user.avatar as avatar`).
		OrderBy("comment.created_at desc").
		Limit(int(limit), int(offset)).
		FindAndCount(&wlogs)
	if err != nil {
		logrus.Errorf("comment find by novelID:%d err:%v", novelID, err)
		return wlogs, 0, errors.New("not found")
	}
	return wlogs, total, nil
}
