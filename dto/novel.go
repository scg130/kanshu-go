package dto

type NovelRequest struct {
	NovelId    int32 `json:"novel_id" form:"novel_id"`
	Page       int32 `json:"page" form:"page"`
	Limit      int32 `json:"limit" form:"limit"`
	ChapterNum int32 `json:"chapter_num" form:"chapter_num"`
}

type BuyRequest struct {
	Page      int32 `json:"page" form:"page"`
	Size      int32 `json:"size" form:"size"`
	ChapterId int64 `json:"chapter_id" form:"chapter_id"`
	NovelId   int32 `json:"novel_id" form:"novel_id"`
	Num       int32 `json:"num" form:"num"`
}
