package web

import "github.com/gin-gonic/gin"

type AtricleHandler struct {
}

func (h *AtricleHandler) RegisterArticleRoutes(server *gin.Engine) {
	ag := server.Group("/articles")
	// 新增、修改文章
	ag.POST("/edit", h.Edit)
}

func (h *AtricleHandler) Edit(ctx *gin.Context) {

}
