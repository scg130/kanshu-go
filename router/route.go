package router

import (
	"kanshu/endpoint"
	"kanshu/endpoint/novel"
	"kanshu/endpoint/user"
	"kanshu/middleware"
	"net/http"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prometheusF(r *gin.Engine) {
	opsProcessed := promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "The total number of requests to the my_api service.",
		},
	)
	r.Use(func(ctx *gin.Context) {
		if ctx.Request.RequestURI != "/metrics" {
			opsProcessed.Inc()
			ctx.Next()
		}
	})
}

func HttpRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	}, middleware.UserRateLimiter())
	prometheusF(r)
	r.Use(middleware.Tracer())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(middleware.Cors())
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/book/")
	})
	r.Static("/book/", "./resource/")

	novelR(r)
	commentR(r)
	userR(r)
	chargeR(r)
	captchaR(r)

	admin(r)
	return r
}

func commentR(r *gin.Engine) {
	comment := novel.NewCommentSrv()
	cg := r.Group("/comment", middleware.Auth())
	cg.POST("/insert", comment.AddComment)
	cg.GET("/list", comment.CommentList)
	cg.GET("/dianZan", comment.DianZan)
}

func novelR(r *gin.Engine) {
	app := endpoint.NewApp()
	r.GET("/app/banner/selectBannerList", app.AppSelectBannerList)
	r.GET("/banner/selectBannerList", app.SelectBannerList)
	r.GET("/app/course/selectCourse", app.SelectCourse)

	novelSrv := novel.NewNovelSrv()
	cg := r.Group("/novel", middleware.Auth())
	cg.GET("/selectByUserId", novelSrv.SelectByUserId)
	cg.GET("/cate/list", novelSrv.Cates)
	cg.GET("/join-book", novelSrv.JoinBook)
	cg.GET("/detail", novelSrv.NovelDetail)
	cg.GET("/chapters", novelSrv.Chapters)
	cg.GET("/chapter", novelSrv.Chapter)
	cg.GET("/buyChapter", novelSrv.BuyChapter)
}

func chargeR(r *gin.Engine) {
	chargeSrv := endpoint.NewChargeSrv()
	charge := r.Group("/charge", middleware.Auth())
	charge.POST("/create", chargeSrv.CreateOrder)
	charge.GET("/order", chargeSrv.QueryOrder)
	r.POST("/charge/callback", chargeSrv.Callback)
	r.GET("/charge/callback/USD", chargeSrv.USDCallback)
}

func userR(r *gin.Engine) {
	g := r.Group("app/")
	uSrv := user.NewUserSrv()
	lg := g.Group("Login")
	lg.POST("/registerCode", uSrv.RegisterCode)
	lg.POST("/phoneLogin", uSrv.PhoneLogin)
	lg.POST("/forgetPwd", uSrv.ForgetPwd)
	ug := g.Group("user")
	ug.GET("/selectUserById", uSrv.UserInfo)
}

func captchaR(r *gin.Engine) {
	g := r.Group("captcha")
	g.GET("/generate", endpoint.CaptchaGenerate)
	g.GET("/image", endpoint.CaptchaImage)
}
