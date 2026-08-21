package router

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/cluster"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/logger"
	"net/http"
	"os"
	"strings"
)

func SetRouter(router *gin.Engine, buildFS embed.FS) {
	SetApiRouter(router)
	SetDashboardRouter(router)
	SetRelayRouter(router)
	cluster.RegisterRoutes(router)
	frontendBaseUrl := os.Getenv("FRONTEND_BASE_URL")
	if config.IsMasterNode && frontendBaseUrl != "" {
		frontendBaseUrl = ""
		logger.SysLog("FRONTEND_BASE_URL is ignored on master node")
	}
	if frontendBaseUrl == "" {
		SetWebRouter(router, buildFS)
	} else {
		frontendBaseUrl = strings.TrimSuffix(frontendBaseUrl, "/")
		router.NoRoute(func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%s", frontendBaseUrl, c.Request.RequestURI))
		})
	}
}
