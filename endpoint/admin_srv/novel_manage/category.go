package novel_manage

import (
	"kanshu/dto"
	"kanshu/dto/admin"
	"kanshu/pkg/spider"
	go_micro_service_novel "kanshu/proto/novel"
	"net/http"
	"strconv"

	selfwrappers "kanshu/wrappers"

	"github.com/gin-gonic/gin"
	"github.com/scg130/tools"
	"github.com/scg130/tools/wrappers"
)

type Novel struct {
	NovelCli go_micro_service_novel.NovelSrvService
}

const NOVEL_SRV_NAME = "go.kanshu.service.novel"

var novelSrv *Novel

func NewNovelSrv() *Novel {
	if novelSrv == nil {
		novelSrv = &Novel{
			NovelCli: go_micro_service_novel.NewNovelSrvService(NOVEL_SRV_NAME, tools.GetMicroClient(
				NOVEL_SRV_NAME,
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
		}
	}
	return novelSrv
}

func (this *Novel) SetVipChapter(c *gin.Context) {
	novelId, _ := strconv.Atoi(c.Param("id"))
	var req admin.SetVipReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	rsp, err := this.NovelCli.SetVipChapter(c, &go_micro_service_novel.SetVipChapterReq{
		NovelId:    int32(novelId),
		MinChapter: int32(req.MinChapter),
		MaxChapter: int32(req.MaxChapter),
		IsVip:      int32(req.IsVip),
	})
	if err != nil || rsp.Code != 0 {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "设置小说vip章节失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
	})
}

func (this *Novel) SpiderNovel(c *gin.Context) {
	var req admin.SpiderNovelReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	go spider.Run(req.PyName, req.ZhName)
	c.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
	})
}

func (this *Novel) EditNovel(c *gin.Context) {
	novelId, _ := strconv.Atoi(c.Param("id"))

	var req admin.EditNovelReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	rsp, err := this.NovelCli.UpdateNovel(c, &go_micro_service_novel.Novel{
		NovelId: int32(novelId),
		Name:    req.Name,
		Author:  req.Author,
		CateId:  int32(req.CateId),
		Img:     req.Img,
		Sort:    int32(req.Sort),
	})
	if err != nil || rsp.Code != 0 {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "编辑小说失败",
		})
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "编辑小说成功",
	})
}

func (this *Novel) GetNovelList(c *gin.Context) {
	var req admin.NovelReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	rsp, err := this.NovelCli.GetNovelList(c, &go_micro_service_novel.NovelListReq{
		Name:     req.Name,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Author:   req.Author,
		CateId:   int32(req.CateId),
	})
	if err != nil || rsp.Code != 0 {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "获取小说列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
		Data: rsp.Novels,
		Pagnation: admin.Pagnation{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    rsp.Pagnation.Total,
		},
	})
}

func (this *Novel) DelCategory(c *gin.Context) {
	cateId, _ := strconv.Atoi(c.Param("id"))
	rsp, err := this.NovelCli.DelCategory(c, &go_micro_service_novel.DelCategoryReq{
		CategoryId: int32(cateId),
	})
	if err != nil || rsp.Code != 0 {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "删除分类失败",
		})
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "删除分类成功",
	})
}

func (this *Novel) EditCategory(c *gin.Context) {
	cateId, _ := strconv.Atoi(c.Param("id"))

	var req admin.CategoryReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	rsp, err := this.NovelCli.UpdateCategory(c, &go_micro_service_novel.Category{
		CateId:  int32(cateId),
		Name:    req.Name,
		Channel: int32(*req.Channel),
		Sort:    int32(req.Sort),
		IsShow:  int32(*req.IsShow),
	})
	if err != nil || rsp.Code != 0 {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "编辑分类失败",
		})
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "编辑分类成功",
	})
}

func (this *Novel) AddCategory(c *gin.Context) {
	var req admin.CategoryReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  err.Error(),
		})

		return
	}

	rsp, err := this.NovelCli.AddCateGory(c, &go_micro_service_novel.AddCateRequest{
		Name:    req.Name,
		Sort:    int32(req.Sort),
		Channel: int32(*req.Channel),
		IsShow:  int32(*req.IsShow),
	})
	if err != nil || rsp.Code != 0 {
		c.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "添加分类失败",
		})
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "添加分类成功",
	})
}

func (this *Novel) GetCategoryList(ctx *gin.Context) {
	var req admin.NovelCategoryReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, admin.Resp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}
	isShow := -1
	if req.IsShow > 0 {
		isShow = req.IsShow
	}
	rep, err := this.NovelCli.GetCateGories(ctx, &go_micro_service_novel.Request{
		Name:   req.Name,
		Page:   int32(req.Page),
		Size_:  int32(req.PageSize),
		IsShow: int32(isShow),
	})
	if err != nil || rep.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	ctx.JSON(http.StatusOK, admin.Resp{
		Code: 0,
		Msg:  "ok",
		Data: rep.Categories,
		Pagnation: admin.Pagnation{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    rep.Pagnation.Total,
		},
	})
}
