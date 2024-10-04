package admin

type MenuRequest struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Path  string `json:"path" form:"path" binding:"required"`
	Api   string `json:"api" form:"api" binding:"required"`
	Icon  string `json:"icon" form:"icon" binding:"required"`
	Show  int32  `json:"show" form:"show" binding:"required"`
	Pid   int64  `json:"pid" form:"pid"`
	State int32  `json:"state" form:"state" binding:"required"`
}

type MenusRequest struct {
	Name      string `json:"name" form:"name"`
	Pagnation `json:"pagnation" form:"pagnation"`
}

type Pagnation struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"page_size" form:"page_size"`
	Total    int64 `json:"total" form:"total"`
}
