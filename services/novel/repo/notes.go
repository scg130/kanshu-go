package repo

import (
	"fmt"
	"log"
	"time"
)

type Notes struct {
	Id         int64     `xorm:" pk autoincr INT(11)"`
	UserId     int32     `xorm:"not null default 0 comment('用户id') int(11) index(uid_nid)"`
	NovelId    int32     `xorm:"not null default 0 comment('书籍id') int(11) index(uid_nid)"`
	IsJoin     int32     `xorm:"not null default 0 comment('1 在书架 2 不在书架') tinyint"`
	ChapterNum int32     `xorm:"not null default 0 comment('章节编号') int(11)"`
	IsDelete   int32     `xorm:"not null default 0 comment('是否删除') int(11)"`
	UpdateTime time.Time `xorm:"updated_at not null  DEFAULT CURRENT_TIMESTAMP  timestamp"`
	CreateTime time.Time `xorm:"created_at not null DEFAULT CURRENT_TIMESTAMP timestamp"`
}

func init() {
	notes := new(Notes)
	if isExist, _ := x.IsTableExist(notes); !isExist {
		if err := x.Sync2(notes); err != nil {
			log.Fatal(fmt.Sprintf("sync tables err:%v", err))
		}
	}
}

type SelfData struct {
	Name    string `json:"name"`
	PreNum  int32  `json:"pre_num"`
	NewNum  int32  `json:"new_num"`
	NovelId int32  `json:"novel_id"`
	PreName string `json:"pre_name"`
	NewName string `json:"new_name"`
}

func (self *Notes) RecoveryNote(userId, novelId, chapterNum int32) error {
	m := &Notes{IsDelete: 0}
	num, err := x.Where("user_id=?", userId).And("novel_id=?", novelId).And("chapter_num=?", chapterNum).Cols("is_delete").Update(m)
	fmt.Println(num)
	if err != nil {
		return err
	}
	return err
}

func (n *Notes) DelNote(novelId, userId int32) (err error) {
	m := &Notes{IsDelete: 1}
	_, err = x.Where("user_id=?", userId).And("novel_id=?", novelId).Update(m)
	if err != nil {
		return err
	}
	return err
}

func (n *Notes) Note(novelId, userId int32) (err error) {
	m := &Notes{IsDelete: 1}
	nm, err := x.Where("user_id=?", userId).And("novel_id=?", novelId).Update(m)
	fmt.Println(nm)
	if err != nil {
		return err
	}
	return err
}

func (n *Notes) GetNote(userId, novelId, chapterNum int32) (note Notes, err error) {
	has, err := x.Where("user_id=? and is_delete = 0", userId).Where("novel_id=?", novelId).And("chapter_num=?", chapterNum).Get(&note)
	if !has {
		return note, err
	}
	return note, err
}

func (n *Notes) GetNoteCount(novelId int32) (num int64, err error) {
	note := new(Notes)
	num, err = x.Distinct("user_id").Where("novel_id=? and is_delete = 0", novelId).Count(note)
	return
}

func (n *Notes) CreateNote(userId, novelId, chapterNum, isJoin int32) (err error) {
	note := new(Notes)
	note.UserId = userId
	note.NovelId = novelId
	note.ChapterNum = chapterNum
	note.IsJoin = isJoin
	note.UpdateTime = time.Now()
	note.CreateTime = time.Now()
	note.IsDelete = 0
	_, err = x.Insert(note)
	return
}

func (n *Notes) GetLastNote(userId, novelId int32) (note Notes, err error) {
	has, err := x.Where("user_id=?", userId).And("novel_id=?", novelId).OrderBy("chapter_num desc").Limit(1).Get(&note)
	if !has {
		note.ChapterNum = 1
		return note, err
	}
	return note, err
}

func (n *Notes) GetJoinLastNote(userId, novelId int32) (note Notes, err error) {
	has, err := x.Where("user_id=? and is_join = 1", userId).And("novel_id=?", novelId).OrderBy("id desc").Limit(1).Get(&note)
	if !has {
		note.ChapterNum = 1
		return note, err
	}
	return note, err
}

func (n *Notes) GetNotes(name string, userId, page, size, isEnd int) (data []SelfData, err error) {
	nType := " pre_num < novel.chapter_total"
	if isEnd == 1 {
		nType = " pre_num >= novel.chapter_total"
	}
	sql := fmt.Sprintf("select novel.name as name,novel.chapter_total as total,max(notes.chapter_num) as pre_num,novel.chapter_current as new_num,novel.id as novel_id from notes join novel on notes.novel_id = novel.id where novel.name like ? and notes.user_id = ? and notes.is_delete = 0 group by novel.id having  %s", nType)
	err = x.SQL(sql, "%"+name+"%", userId).Limit(size, size*(page-1)).Find(&data)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func (n *Notes) GetNoteNumByUid(uid int) (int32, error) {
	note := new(Notes)
	num, err := x.Distinct("novel_id").Where("user_id=? and is_delete = 0", uid).Count(note)
	return int32(num), err
}

func (n *Notes) GetNoteJoinNumByUid(uid int) (int32, error) {
	note := new(Notes)
	num, err := x.Distinct("novel_id").Where("user_id=? and is_delete = 0 and is_join = 1", uid).Count(note)
	return int32(num), err
}
