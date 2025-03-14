package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/logger"
)

type AtricleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService,
	l logger.LoggerV1) *AtricleHandler {
	return &AtricleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *AtricleHandler) RegisterArticleRoutes(server *gin.Engine) {
	ag := server.Group("/articles")
	// 新增、修改文章
	ag.POST("/edit", h.Edit)
}

func (h *AtricleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("未发现用户的 session 信息")
	}
	// 检测输入
	id, err := h.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("保存失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}
