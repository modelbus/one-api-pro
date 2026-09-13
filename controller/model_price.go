package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/model"
)

// ⚠️ Deprecated: 仅 RootAuth 可见，返回 model_price 全量字段（含 enabled=false 与计费详情）。
//   前端下拉等只读场景请改用 ListModelPriceOptions（路径 /api/model_price/options，AdminAuth）。
//   本接口仅保留给 root 角色查看完整模型定价。
//
// Deprecated: prefer ListModelPriceOptions for read-only model lists used in UI dropdowns.
func GetAllModelPrices(c *gin.Context) {
	prices, err := model.GetAllModelPrices()
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
		"data":    prices,
	})
}

// ListModelPriceOptions 返回 model_price 表中所有 enabled=true 模型的 model_name 列表，
// 仅供前端渠道编辑弹窗等场景作为模型候选下拉使用。AdminAuth 即可访问，
// 与 RootAuth 保护的 GetAllModelPrices 隔离。
// 仅返回 model_name 字符串切片，节省带宽，避免把计费字段泄漏到下拉 UI。
// 版本: v0.0.19
// 日期: 2026-09-13
func ListModelPriceOptions(c *gin.Context) {
	var prices []model.ModelPrice
	if err := model.DB.Where("enabled = ?", true).
		Order("model_name asc").
		Find(&prices).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	names := make([]string, 0, len(prices))
	for _, p := range prices {
		if p.ModelName != "" {
			names = append(names, p.ModelName)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    names,
	})
}

func AddModelPrice(c *gin.Context) {
	var price model.ModelPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}
	price.Id = 0
	if price.ModelName == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "模型名称不能为空",
		})
		return
	}
	if price.BillingType == "" {
		price.BillingType = model.BillingTypeToken
	}
	if err := price.Insert(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	model.InitModelPriceCache()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func UpdateModelPrice(c *gin.Context) {
	var price model.ModelPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}
	if price.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "ID不能为空",
		})
		return
	}
	if err := price.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	model.InitModelPriceCache()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func DeleteModelPrice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}
	if err := model.DeleteModelPriceById(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	model.InitModelPriceCache()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func GetAllGroupPrices(c *gin.Context) {
	prices, err := model.GetAllGroupPrices()
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
		"data":    prices,
	})
}

func AddGroupPrice(c *gin.Context) {
	var price model.GroupPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}
	price.Id = 0
	if price.GroupName == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "分组名称不能为空",
		})
		return
	}
	if price.Discount == 0 {
		price.Discount = 1.0
	}
	if err := price.Insert(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	model.InitGroupPriceCache()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func UpdateGroupPrice(c *gin.Context) {
	var price model.GroupPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}
	if price.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "ID不能为空",
		})
		return
	}
	if err := price.Update(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	model.InitGroupPriceCache()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func DeleteGroupPrice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}
	if err := model.DeleteGroupPriceById(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	model.InitGroupPriceCache()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}