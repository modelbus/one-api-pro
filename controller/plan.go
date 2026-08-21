package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/helper"
	"github.com/modelbus/one-api-pro/model"
)

func GetAllPlans(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	plans, err := model.GetAllPlans(p*config.ItemsPerPage, config.ItemsPerPage)
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
		"data":    plans,
	})
}

func SearchPlans(c *gin.Context) {
	keyword := c.Query("keyword")
	plans, err := model.SearchPlans(keyword)
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
		"data":    plans,
	})
}

func GetPlan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	plan, err := model.GetPlanById(id)
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
		"data":    plan,
	})
}

func AddPlan(c *gin.Context) {
	plan := model.Plan{}
	err := c.ShouldBindJSON(&plan)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if plan.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "套餐名称不能为空",
		})
		return
	}
	if err := plan.ValidateDefaultModel(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	plan.CreatedTime = helper.GetTimestamp()
	plan.UpdatedTime = helper.GetTimestamp()
	err = plan.Insert()
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
		"data":    plan,
	})
}

func UpdatePlan(c *gin.Context) {
	plan := model.Plan{}
	err := c.ShouldBindJSON(&plan)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if plan.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "套餐 ID 不能为空",
		})
		return
	}
	if err := plan.ValidateDefaultModel(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	plan.UpdatedTime = helper.GetTimestamp()
	err = plan.Update()
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
	})
}

func DeletePlan(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := model.DeletePlanById(id)
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
	})
}