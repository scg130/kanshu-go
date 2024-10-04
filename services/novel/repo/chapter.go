package repo

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Chapter struct {
	Id      int64  `xorm:" pk autoincr INT(11)"`
	Title   string `xorm:"not null VARCHAR(255)"`
	Content string `xorm:"not null comment('小说简介') text"`
	NovelId int    `xorm:"not null index(idx_nov_id_num) default 0 comment('小说id') int(11)"`
	IsVip   int    `xorm:"not null default 0 comment('是否vip') int(11)"`
	Num     int    `xorm:"not null index(idx_nov_id_num) default 0 comment('章节序号') int(11)"`
	Words   int    `xorm:"not null default 0 comment('总字数') int(11)"`

	UpdateTime time.Time `xorm:"updated_at not null  DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	chapter := new(Chapter)
	if isExist, _ := x.IsTableExist(chapter); !isExist {
		if err := x.Sync2(chapter); err != nil {
			log.Fatal(fmt.Sprintf("sync tables err:%v", err))
		}
	}
}

func (c *Chapter) SetVipChapter(novelId, minChapter, maxChapter, isVip int) error {
	_, err := x.Exec(`update chapter set is_vip=? where novel_id=? and num between ? and ?`, isVip, novelId, minChapter, maxChapter)
	return err
}

func (c *Chapter) GetByNovelId(novelId, page, size int, orderType string) (data []Chapter, err error) {
	if orderType != "asc" && orderType != "desc" && orderType != "" {
		errors.New("orderType is invalid")
		return
	}
	err = x.Where("novel_id=?", novelId).Select(`id, title, novel_id, is_vip, num, words`).OrderBy(fmt.Sprintf("num %s", orderType)).Limit(size, size*(page-1)).Find(&data)
	return data, err
}

func (c *Chapter) GetNewChapter(novelId int) (Chapter, error) {
	var chapter Chapter
	has, err := x.Where("novel_id=?", novelId).OrderBy("num desc").Get(&chapter)
	if !has {
		return chapter, errors.New("no has")
	}
	return chapter, err
}

func (c *Chapter) GetOne(novelId, num int) (Chapter, error) {
	var chapter Chapter
	has, err := x.Where("novel_id=?", novelId).And("num=?", num).Get(&chapter)
	if !has {
		return chapter, errors.New("no has")
	}
	return chapter, err
}
