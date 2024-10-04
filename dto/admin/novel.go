package admin

type NovelCategoryReq struct {
	Name      string `json:"name" form:"name"`
	Pagnation `json:"pagnation" form:"pagnation"`
	IsShow    int `json:"is_show" form:"is_show"`
}

type NovelReq struct {
	Name      string `json:"name" form:"name"`
	Author    string `json:"author" form:"author"`
	CateId    int    `json:"cate_id" form:"cate_id"`
	Pagnation `json:"pagnation" form:"pagnation"`
}

type CategoryReq struct {
	Name    string `json:"name" form:"name" binding:"required"`
	Sort    int    `json:"sort" form:"sort" binding:"required"`
	IsShow  *int   `json:"is_show" form:"is_show" binding:"required"`
	Channel *int   `json:"channel" form:"channel" binding:"required"`
}

type EditNovelReq struct {
	Name   string `json:"name" form:"name" binding:"required"`
	CateId int    `json:"cate_id" form:"cate_id" binding:"required"`
	Sort   int    `json:"sort" form:"sort" binding:"required"`
	Author string `json:"author" form:"author" binding:"required"`
	Img    string `json:"img" form:"img" binding:"required"`
}

type SpiderNovelReq struct {
	PyName string `json:"py_name" form:"py_name" binding:"required"`
	ZhName string `json:"zh_name" form:"zh_name" binding:"required"`
}

type SetVipReq struct {
	MinChapter int `json:"min_chapter" form:"min_chapter" binding:"required"`
	MaxChapter int `json:"max_chapter" form:"max_chapter" binding:"required"`
	IsVip      int `json:"is_vip" form:"is_vip"`
}
