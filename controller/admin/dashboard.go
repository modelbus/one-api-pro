// dashboard.go 管理员仪表盘 HTTP handler
// Admin dashboard HTTP handlers.
// 版本: v0.0.16
// 日期: 2026-09-11
// 作者: opencode
//
// 路由由 router/api.go 在 /api/admin/dashboard 下以 AdminAuth 注册。

package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/modelbus/one-api-pro/model"
)

// GetOverview 处理 GET /api/admin/dashboard/overview。
// 一次性返回 KPI + 7 日趋势，供前端首屏渲染。
// GetOverview handles GET /api/admin/dashboard/overview.
// 版本: v0.0.16
// 日期: 2026-09-11
func GetOverview(c *gin.Context) {
	rawRange := c.Query("range")
	overview, err := model.GetAdminDashboardOverview(rawRange)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    overview,
	})
}

// GetTopUsers 处理 GET /api/admin/dashboard/top-users。
// Query: range=today|7d|30d (default 7d), limit (1-200, default 20)。
// GetTopUsers handles GET /api/admin/dashboard/top-users.
// 版本: v0.0.16
// 日期: 2026-09-11
func GetTopUsers(c *gin.Context) {
	rawRange := c.Query("range")
	limit, _ := strconv.Atoi(c.Query("limit"))

	rows, preset, err := model.GetAdminTopUsers(rawRange, limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"range": preset,
			"items": rows,
		},
	})
}
