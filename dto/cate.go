package dto

type CatesReq struct {
	Sex int `json:"sex" form:"sex"` //1 男 2 女
}

type JoinReq struct {
	NovelId int32 `json:"novel_id" form:"novel_id"`
}
