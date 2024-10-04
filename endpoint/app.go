package endpoint

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilylx/gconv"
	"github.com/scg130/tools"
	"github.com/scg130/tools/wrappers"
	"kanshu/dto"
	"kanshu/env"
	go_micro_service_novel "kanshu/proto/novel"
	selfwrappers "kanshu/wrappers"
	"math/rand"
	"net/http"
	"time"
)

type App struct {
	NovelCli go_micro_service_novel.NovelSrvService
}

var appSrv *App

func NewApp() *App {
	if appSrv == nil {
		appSrv = &App{
			NovelCli: go_micro_service_novel.NewNovelSrvService("go.kanshu.service.novel", tools.GetMicroClient(
				"go.kanshu.service.novel",
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
		}
	}
	return appSrv
}

func (a *App) AppSelectBannerList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: []interface{}{
			gin.H{
				"name":      "免费体验",
				"imageUrl":  fmt.Sprintf("%s/book/imgs/appbg.png", env.AppConf.Domain),
				"state":     1,
				"classify":  5,
				"url":       "/pages/index/index",
				"describes": "快来和我一起看小说吧,好玩的好看的小说都在这里哦",
			},
		},
	})
}

func (a *App) SelectBannerList(ctx *gin.Context) {
	classify := gconv.Int(ctx.Query("classify"))
	data := make([]gin.H, 0)
	if classify == 1 {
		data = append(data, gin.H{"createTime": "2023-10-31 17:42:40",
			"name":      "c",
			"imageUrl":  fmt.Sprintf("%s/book/imgs/banner1.png", env.AppConf.Domain),
			"state":     1,
			"classify":  1,
			"url":       "/package/bookDetails/bookDetails?novel_id=3",
			"sort":      nil,
			"describes": "",
			"course":    nil})
	}
	if classify == 2 {
		data = append(data, gin.H{
			"name":      "推荐",
			"imageUrl":  fmt.Sprintf("%s/book/imgs/tuijian.png", env.AppConf.Domain),
			"state":     1,
			"classify":  2,
			"url":       "/pages/moreBook/moreBook?cate=1",
			"describes": "推荐",
		})
		data = append(data, gin.H{
			"name":      "新书",
			"imageUrl":  fmt.Sprintf("%s/book/imgs/xinshu.png", env.AppConf.Domain),
			"state":     1,
			"classify":  2,
			"url":       "/pages/moreBook/moreBook?cate=2",
			"describes": "新书",
		})
		data = append(data, gin.H{
			"name":      "完结",
			"imageUrl":  fmt.Sprintf("%s/book/imgs/wanjie.png", env.AppConf.Domain),
			"state":     1,
			"classify":  2,
			"url":       "/pages/moreBook/moreBook?cate=3",
			"describes": "完结",
		})

		data = append(data, gin.H{
			"name":      "分类",
			"imageUrl":  fmt.Sprintf("%s/book/imgs/fenlei.png", env.AppConf.Domain),
			"state":     1,
			"classify":  2,
			"url":       "/package/classifty/classifty",
			"describes": "分类",
		})
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func (a *App) SelectCourse(ctx *gin.Context) {
	cate := gconv.Int(ctx.Query("cate"))
	page := gconv.Int(ctx.Query("page"))
	limit := gconv.Int(ctx.Query("limit"))
	userId := gconv.Int32(ctx.Query("user_id"))
	title := ctx.Query("title")

	rsp, err := a.NovelCli.GetNovelsByCateId(ctx, &go_micro_service_novel.Request{
		CateId: gconv.Int32(cate),
		Page:   gconv.Int32(page),
		Size_:  gconv.Int32(limit),
		UserId: userId,
		Name:   title,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	data := make([]gin.H, 0)
	rand.Seed(time.Now().UnixMilli())
	for _, v := range rsp.Novels {
		data = append(data, gin.H{
			"titleImg":   v.Img,
			"source":     float64(rand.Intn(50-40)+40) / 10.0,
			"title":      v.Name,
			"viewCounts": v.ViewCounts,
			"author":     v.Author,
			"novelId":    v.NovelId,
			"isCollect":  v.IsCollect,
			"wordCount":  v.Words,
			"category":   v.CategoryName,
		})
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}
