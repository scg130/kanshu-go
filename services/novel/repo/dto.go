package repo

import (
	"time"
)

type NovelData struct {
	Id             int64  `xorm:" pk autoincr INT(11)" json:"id"`
	Name           string `xorm:"not null unique(name) comment('分类名称') VARCHAR(255)" json:"name"`
	Author         string `xorm:"not null default '' comment('排序') VARCHAR(255)" json:"author"`
	ChapterTotal   int32  `xorm:"not null default 0 comment('总章节数') int(11)" json:"chapter_total"`
	ChapterCurrent int32  `xorm:"not null default 0 comment('最新章节数') int(11)" json:"chapter_current"`
	Img            string `xorm:"not null default '' comment('小说封面图') VARCHAR(300)" json:"img"`
	Intro          string `xorm:"not null default '' comment('小说简介') VARCHAR(500)" json:"intro"`
	CateId         int    `xorm:"not null index(idx_cate_id_sort) default 0 comment('所属分类id') int(11)" json:"cate_id"`
	Words          int32  `xorm:"not null default 0 comment('总字数') int(11)" json:"words"`
	Likes          int32  `xorm:"not null default 0 comment('点赞数') int(11)" json:"likes"`
	UnLikes        int32  `xorm:"not null default 0 comment('不喜欢数') int(11)" json:"un_likes"`
	Sort           int    `xorm:"not null index(idx_cate_id_sort) default 0 comment('排序') int(11)" json:"sort"`

	UpdateTime time.Time `xorm:"updated_at not null  DEFAULT CURRENT_TIMESTAMP  timestamp" json:"update_time"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp" json:"create_time"`

	NewChapter string ` json:"new_chapter"`
	Category   string `json:"category"`
	ViewCounts int    `json:"view_counts"`
	IsCollect  int    `json:"is_collect"`
}
