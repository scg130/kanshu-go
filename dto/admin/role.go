package admin

type RoleRequest struct {
	Name    string  `json:"name" form:"name" binding:"required"`
	MenuIds []int64 `json:"menu_ids" form:"menu_ids"`
}

type RolesRequest struct {
	Name string `json:"name" form:"name"`
}
