package dto

type Resp struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Total   int32       `json:"total"`
	CurPage int32       `json:"cur_page"`
	CateId  int32       `json:"cate_id"`
	Pages   int32       `json:"pages"`
}
