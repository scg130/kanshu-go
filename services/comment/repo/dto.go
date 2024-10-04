package repo

import "time"

type CommentData struct {
	Id       int64  `xorm:"id"`
	UserId   int64  `xorm:"user_id"`
	Username string `xorm:"username"`
	Content  string `xorm:"content"`
	Avatar   string `xorm:"avatar"`
	GoodsNum int32  `xorm:"goods_num"`

	CreateTime time.Time `xorm:"create_time"`
}
