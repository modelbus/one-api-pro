package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Leon-PanPan/one-api-pro/model"
)

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

func AddModelPrice(c *gin.Context) {
	var price model.ModelPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}
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