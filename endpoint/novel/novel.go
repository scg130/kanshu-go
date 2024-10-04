package novel

import (
	"github.com/gin-gonic/gin"
	"github.com/ilylx/gconv"
	"github.com/scg130/tools"
	"github.com/scg130/tools/wrappers"
	"kanshu/dto"
	"kanshu/endpoint"
	go_micro_service_novel "kanshu/proto/novel"
	go_micro_service_user "kanshu/proto/user"
	go_micro_service_wallet "kanshu/proto/wallet"
	selfwrappers "kanshu/wrappers"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type Novel struct {
	novelCli  go_micro_service_novel.NovelSrvService
	walletCli go_micro_service_wallet.WalletService
	userCli   go_micro_service_user.UserCenterService
}

var novelSrv *Novel

const (
	NOVEL_SRV_NAME  = "go.kanshu.service.novel"
	WALLET_SRV_NAME = "go.kanshu.service.wallet"
	USER_SRV_NAME   = "go.kanshu.service.user"
)

func NewNovelSrv() *Novel {
	if novelSrv == nil {
		novelSrv = &Novel{
			novelCli: go_micro_service_novel.NewNovelSrvService(NOVEL_SRV_NAME, tools.GetMicroClient(
				NOVEL_SRV_NAME,
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
			walletCli: go_micro_service_wallet.NewWalletService(WALLET_SRV_NAME, tools.GetMicroClient(
				WALLET_SRV_NAME,
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
			userCli: go_micro_service_user.NewUserCenterService(USER_SRV_NAME, tools.GetMicroClient(
				USER_SRV_NAME,
				wrappers.NewTracerWrapper(),
				wrappers.NewRateLimitClientWrapper(100),
				selfwrappers.NewHystrixWrapper(),
			)),
		}
	}
	return novelSrv
}

// @Summary 获取小说详情
// @Description 获取小说详情
// @Tags novel
// @Produce json
// @Param novel_id query int true "query参数"
// @Success 200 {object}  dto.Resp{data=go_micro_service_novel.Novel}
// @Router /novel/detail [get]
func (n *Novel) NovelDetail(ctx *gin.Context) {
	var req dto.NovelRequest
	if err := ctx.BindQuery(&req); err != nil || req.NovelId == 0 {
		ctx.String(http.StatusBadRequest, "bad params")
		return
	}

	//span := opentracing.StartSpan("novel")
	//defer span.Finish()
	//c := opentracing.ContextWithSpan(ctx, span)
	userInfo := endpoint.GetUserInfo(ctx)
	resp, err := n.novelCli.GetNovelById(ctx, &go_micro_service_novel.Request{NovelId: req.NovelId, UserId: int32(userInfo.UserId)})
	if err != nil || resp.Code != 0 || resp.Novel == nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: gin.H{
			"novel":         resp.Novel,
			"chapter_total": resp.Novel.ChapterTotal,
			"chapter_num":   resp.Novel.ChapterNum,
		},
	})
}

// @Summary 章节列表
// @Description 章节列表
// @Tags novel
// @Produce json
// @Param query query dto.NovelRequest true "query参数"
// @Success 200 {object}  dto.Resp{data=[]go_micro_service_novel.Chapter}
// @Router /novel/chapters [get]
func (n *Novel) Chapters(ctx *gin.Context) {
	var req dto.NovelRequest
	if err := ctx.BindQuery(&req); err != nil || req.Page == 0 || req.Limit == 0 {
		ctx.String(http.StatusBadRequest, "bad params")
		return
	}

	resp, err := n.novelCli.GetChaptersByNovelId(ctx, &go_micro_service_novel.Request{NovelId: req.NovelId, Page: req.Page, Size_: req.Limit})
	if err != nil || resp.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}

	rsp, err := n.novelCli.GetNovelById(ctx, &go_micro_service_novel.Request{
		NovelId: req.NovelId,
	})
	if err != nil || resp.Code != 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Resp{
		Code:  0,
		Msg:   "ok",
		Data:  resp.Chapters,
		Pages: gconv.Int32(math.Ceil(gconv.Float64(rsp.Novel.ChapterTotal) / gconv.Float64(req.Limit))),
	})
}

// @Summary 获取章节
// @Description 获取章节
// @Tags novel
// @Produce json
// @Param chapter_id query int true "query参数"
// @Success 200 {object}  dto.Resp{data=go_micro_service_novel.Chapter}
// @Router /novel/chapter [get]
func (n *Novel) Chapter(ctx *gin.Context) {
	var req dto.NovelRequest
	if err := ctx.BindQuery(&req); err != nil || req.NovelId == 0 || req.ChapterNum == 0 {
		ctx.String(http.StatusBadRequest, "bad params")
		return
	}

	userInfo := endpoint.GetUserInfo(ctx)

	resp, err := n.novelCli.GetChapterById(ctx, &go_micro_service_novel.Request{NovelId: req.NovelId, Num: req.ChapterNum, UserId: int32(userInfo.UserId)})
	if err != nil || resp.Code != 0 || resp.Chapter == nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "failure",
		})
		return
	}
	if resp.Chapter.IsVip == go_micro_service_novel.VipType_IS_VIP {
		rsp, err := n.walletCli.GetChapter(ctx, &go_micro_service_wallet.BuyChapterRequest{
			Uid:       userInfo.UserId,
			ChapterId: int64(resp.Chapter.ChapterId),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
		if rsp.State != 1 {
			ctx.JSON(http.StatusOK, dto.Resp{
				Code: 1,
				Msg:  "该章节属性vip章节,请先购买或开通vip。",
				Data: gin.H{
					"chapter_id":  resp.Chapter.ChapterId,
					"chapter_num": resp.Chapter.Num,
					"is_have":     0,
				},
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "ok",
		Data: gin.H{
			"chapter": resp.Chapter,
			"is_have": 1,
		},
	})
}

// @Summary 购买章节
// @Description 购买章节
// @Tags novel
// @Produce json
// @Param chapter_id query int true "query参数"
// @Success 200 {object}  dto.Resp{}
// @Router /novel/buy_chapter [get]
func (self *Novel) BuyChapter(ctx *gin.Context) {
	var req dto.BuyRequest
	if err := ctx.BindQuery(&req); err != nil {
		ctx.String(http.StatusBadRequest, "bad params")
		return
	}

	userInfo := endpoint.GetUserInfo(ctx)

	crep, err := self.walletCli.GetChapter(ctx, &go_micro_service_wallet.BuyChapterRequest{
		Uid:       userInfo.UserId,
		ChapterId: req.ChapterId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "购买失败,chapterId 错误！",
		})
		return
	}
	if len(crep.Log) > 0 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "已经购买过！",
		})
		return
	}
	chapter, err := self.novelCli.GetChapterById(ctx, &go_micro_service_novel.Request{
		NovelId: req.NovelId,
		Num:     req.Num,
	})
	if err != nil || chapter.Code != 0 || chapter.Chapter == nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "购买失败",
		})
		return
	}
	novel, err := self.novelCli.GetNovelById(ctx, &go_micro_service_novel.Request{
		NovelId: chapter.Chapter.NovelId,
	})
	if err != nil || novel.Code != 0 || novel.Novel == nil {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "购买失败",
		})
		return
	}
	wRsp, err := self.walletCli.GetOne(ctx, &go_micro_service_wallet.WalletReq{
		Uid: userInfo.UserId,
	})
	if err != nil || wRsp.AvailableBalance < 100 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -3,
			Msg:  "余额不足,请先充值！",
		})
		return
	}
	resp, err := self.walletCli.BuyChapter(ctx, &go_micro_service_wallet.BuyChapterRequest{
		Uid:       userInfo.UserId,
		ChapterId: req.ChapterId,
		NovelId:   int64(novel.Novel.NovelId),
		NovelName: novel.Novel.Name,
		Amount:    100,
	})
	if err != nil || resp.State != 1 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "购买失败",
		})
		return
	}
	changeRep, err := self.walletCli.Change(ctx, &go_micro_service_wallet.WalletReq{
		Uid:    userInfo.UserId,
		Amount: 100,
		Type:   go_micro_service_wallet.Type_STATE_BUY_CHAPTER,
	})
	if err != nil || changeRep.State != 1 {
		ctx.JSON(http.StatusOK, dto.Resp{
			Code: -1,
			Msg:  "购买失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "购买成功",
	})
	return
}

func (self *Novel) SelectByUserId(ctx *gin.Context) {
	page := gconv.Int32(ctx.Query("page"))
	limit := gconv.Int32(ctx.Query("limit"))
	//1书架 2点赞 3历史记录
	classify := gconv.Int32(ctx.Query("classify"))
	userInfo := endpoint.GetUserInfo(ctx)
	rsp, err := self.novelCli.GetNovelsByUserId(ctx, &go_micro_service_novel.RequestByUserId{
		Page:     gconv.Int32(page),
		Size_:    gconv.Int32(limit),
		UserId:   gconv.Int32(userInfo.UserId),
		Classify: gconv.Int32(classify),
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
			"titleImg":             v.Img,
			"source":               float64(rand.Intn(50-40)+40) / 10.0,
			"title":                v.Name,
			"viewCounts":           v.ViewCounts,
			"author":               v.Author,
			"novelId":              v.NovelId,
			"isCollect":            v.IsCollect,
			"wordCount":            v.Words,
			"courseClassification": v.CategoryName,
			"chapterTotal":         v.ChapterTotal,
			"coursePage":           v.ChapterNum,
		})
	}
	ctx.JSON(http.StatusOK, dto.Resp{
		Code: 0,
		Msg:  "success",
		Data: gin.H{
			"records": data,
			"pages":   rsp.Total,
		},
	})
}
