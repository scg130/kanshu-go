package dto

type RegUserRequest struct {
	Phone         int64  `json:"phone" form:"phone" binding:"required"`
	Passwd        string `json:"passwd" form:"passwd" binding:"required"`
	PasswdConfirm string `json:"passwd_confirm" form:"passwd_confirm" binding:"required"`
	Captcha
}

type ForgetRequest struct {
	Phone  int64  `json:"phone" form:"phone" binding:"required"`
	Passwd string `json:"passwd" form:"passwd" binding:"required"`
	Code   string `json:"code" form:"code" binding:"required"`
}

type UserRequest struct {
	Phone  int64  `json:"phone" form:"phone" binding:"required"`
	Passwd string `json:"passwd" form:"passwd" binding:"required"`
}

type LoginResp struct {
	Token          string `json:"token"`
	UserId         int64  `json:"user_id"`
	Username       string `json:"username"`
	Avatar         string `json:"avatar"`
	Phone          string `json:"phone"`
	InvitationCode string `json:"invitation_code"`
	Sex            int    `json:"sex"` // 1 男 2 女 3 其它
}

type Captcha struct {
	Id   string `json:"id" form:"id" binding:"required"`
	Code string `json:"code" form:"code" binding:"required"`
}

type UserInfo struct {
	UserId         int64   `json:"user_id"`
	Username       string  `json:"username"`
	Phone          string  `json:"phone"`
	Avatar         string  `json:"avatar"`
	InvitationCode string  `json:"invitation_code"`
	Sex            int     `json:"sex"` // 1 男 2 女 3 其它
	Money          float64 `json:"money"`
	InviteMoney    float64 `json:"invite_money"`
	IsVip          int     `json:"is_vip"` //1 是 0 否
	VipEndTime     string  `json:"vip_end_time"`
	IntegralNum    int     `json:"integral_num"` //积分
	LoveNum        int     `json:"love_num"`     //书架
	NoteNum        int     `json:"note_num"`     //阅读数
}
