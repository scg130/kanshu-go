package admin

type LoginRequest struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Passwd   string `json:"passwd" form:"passwd" binding:"required"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type RegRequest struct {
	UserName      string  `json:"user_name" form:"user_name" binding:"required"`
	Passwd        string  `json:"passwd" form:"passwd"`
	ConfirmPasswd string  `json:"confirm_passwd" form:"confirm_passwd"`
	Email         string  `json:"email" form:"email" binding:"required"`
	Phone         int64   `json:"phone" form:"phone" binding:"required"`
	State         int32   `json:"state" form:"state" `
	RoleIds       []int64 `json:"role_ids" form:"role_ids"`
}

type ListInfoRequest struct {
	UserName string `json:"user_name" form:"user_name"`
}
