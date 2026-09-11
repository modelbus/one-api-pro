package router

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/logger"
	"github.com/modelbus/one-api-pro/controller"
	"github.com/modelbus/one-api-pro/middleware"
	"net/http"
	"strings"
)

func SetWebRouter(router *gin.Engine, buildFS embed.FS) {
	indexPath := fmt.Sprintf("web/build/%s/index.html", config.Theme)
	indexPageData, err := buildFS.ReadFile(indexPath)
	if err != nil {
		logger.SysError(fmt.Sprintf("theme %q is not embedded (missing %s): %v; admin UI will return an empty page", config.Theme, indexPath, err))
	}
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())
	router.Use(middleware.Cache())
	router.Use(static.Serve("/", common.EmbedFolder(buildFS, fmt.Sprintf("web/build/%s", config.Theme))))
	router.NoRoute(func(c *gin.Context) {
		// /v1/* 未匹配：返回 OpenAI 兼容错误，供 relay 客户端识别。
		// Unmatched /v1/*: return the OpenAI-compatible error for relay clients.
		if strings.HasPrefix(c.Request.RequestURI, "/v1") {
			controller.RelayNotFound(c)
			return
		}
		// /api/* 未匹配：返回标准 API 404 信封，避免被误认为 relay 路由问题
		// （历史上这里复用 RelayNotFound，导致缺失的管理接口返回
		//  invalid_request_error，排查困难）。
		// Unmatched /api/*: return the standard API 404 envelope. Reusing
		// RelayNotFound here previously produced an OpenAI-style
		// invalid_request_error for missing management routes, which was
		// misleading.
		if strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "接口不存在: " + c.Request.Method + " " + c.Request.URL.Path,
				"data":    nil,
			})
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPageData)
	})
}
