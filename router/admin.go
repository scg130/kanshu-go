package router

import (
	"kanshu/endpoint/admin_srv"
	"kanshu/endpoint/admin_srv/novel_manage"
	"kanshu/middleware"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func admin(r *gin.Engine) {
	adminSrv := admin_srv.NewAdminSrv()
	r.Use(cors.Default())
	admin := r.Group("/admin")
	admin.POST("/login", adminSrv.Login)
	admin.POST("/reg", adminSrv.Reg)

	g := r.Group("/admin/api")
	g.Use(middleware.AdminAuth(), middleware.Permission())
	g.GET("/info", adminSrv.UserInfo)
	g.POST("/users", adminSrv.Users)
	g.POST("/user/add", adminSrv.Reg)
	g.POST("/user/edit/:id", adminSrv.EditUser)
	g.DELETE("/user/del/:id", adminSrv.DelUser)

	g.POST("/menus", adminSrv.Menus)
	g.GET("/menu/show_tree", adminSrv.MenuShowTree)
	g.GET("/menu/tree", adminSrv.MenuTree)
	g.POST("/menu/add", adminSrv.AddMenu)
	g.POST("/menu/edit/:id", adminSrv.EditMenu)
	g.DELETE("/menu/del/:id", adminSrv.DelMenu)

	g.POST("/roles", adminSrv.Roles)
	g.POST("/role/add", adminSrv.AddRole)
	g.POST("/role/edit/:id", adminSrv.EditRole)
	g.DELETE("/role/del/:id", adminSrv.DelRole)

	novelRouter(r)
}

func novelRouter(r *gin.Engine) {
	adminNovel := novel_manage.NewNovelSrv()
	r.POST("/spider", adminNovel.SpiderNovel)
	gr := r.Group("/novel")
	gr.Use(middleware.AdminAuth(), middleware.Permission())
	gr.POST("/category/list", adminNovel.GetCategoryList)
	gr.POST("/category/add", adminNovel.AddCategory)
	gr.POST("/category/edit/:id", adminNovel.EditCategory)
	gr.DELETE("/category/del/:id", adminNovel.DelCategory)

	gr.POST("/novel/list", adminNovel.GetNovelList)
	gr.POST("/novel/edit/:id", adminNovel.EditNovel)
	gr.POST("/novel/spider", adminNovel.SpiderNovel)
	gr.POST("/novel/vip/:id", adminNovel.SetVipChapter)
}
