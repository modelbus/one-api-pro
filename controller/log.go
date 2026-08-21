package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/modelbus/one-api-pro/common/config"
	"github.com/modelbus/one-api-pro/common/ctxkey"
	"github.com/modelbus/one-api-pro/model"
)

func enrichLogs(logs []*model.Log) {
	channelIds := make(map[int]bool)
	planIds := make(map[int]bool)
	for _, log := range logs {
		if log.ChannelId > 0 {
			channelIds[log.ChannelId] = true
		}
		if log.PlanId > 0 {
			planIds[log.PlanId] = true
		}
	}

	channelNames := make(map[int]string)
	for id := range channelIds {
		ch, err := model.GetChannelById(id, false)
		if err == nil {
			channelNames[id] = ch.Name
		}
	}

	planNames := make(map[int]string)
	for userPlanId := range planIds {
		up, err := model.GetUserPlanById(userPlanId)
		if err == nil && up != nil {
			p, err := model.GetPlanById(up.PlanId)
			if err == nil {
				planNames[userPlanId] = p.Name
			}
		}
	}

	for _, log := range logs {
		if name, ok := channelNames[log.ChannelId]; ok {
			log.ChannelName = name
		}
		if name, ok := planNames[log.PlanId]; ok {
			log.PlanName = name
		}
	}
}

func GetAllLogs(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	logType, _ := strconv.Atoi(c.Query("type"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	username := c.Query("username")
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	channel, _ := strconv.Atoi(c.Query("channel"))
	logs, err := model.GetAllLogs(logType, startTimestamp, endTimestamp, modelName, username, tokenName, p*config.ItemsPerPage, config.ItemsPerPage, channel)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	enrichLogs(logs)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    logs,
	})
	return
}

func GetUserLogs(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	userId := c.GetInt(ctxkey.Id)
	logType, _ := strconv.Atoi(c.Query("type"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	logs, err := model.GetUserLogs(userId, logType, startTimestamp, endTimestamp, modelName, tokenName, p*config.ItemsPerPage, config.ItemsPerPage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	enrichLogs(logs)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    logs,
	})
	return
}

func SearchAllLogs(c *gin.Context) {
	keyword := c.Query("keyword")
	logs, err := model.SearchAllLogs(keyword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	enrichLogs(logs)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    logs,
	})
	return
}

func SearchUserLogs(c *gin.Context) {
	keyword := c.Query("keyword")
	userId := c.GetInt(ctxkey.Id)
	logs, err := model.SearchUserLogs(userId, keyword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	enrichLogs(logs)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    logs,
	})
	return
}

func GetLogsStat(c *gin.Context) {
	logType, _ := strconv.Atoi(c.Query("type"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	tokenName := c.Query("token_name")
	username := c.Query("username")
	modelName := c.Query("model_name")
	channel, _ := strconv.Atoi(c.Query("channel"))
	quotaNum := model.SumUsedQuota(logType, startTimestamp, endTimestamp, modelName, username, tokenName, channel)
	normalQuota := model.SumUsedQuotaByBillingSource(logType, startTimestamp, endTimestamp, modelName, username, tokenName, channel, 0)
	subscriptionQuota := model.SumUsedQuotaByBillingSource(logType, startTimestamp, endTimestamp, modelName, username, tokenName, channel, 1)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"quota":              quotaNum,
			"normal_quota":       normalQuota,
			"subscription_quota": subscriptionQuota,
		},
	})
	return
}

func GetLogsSelfStat(c *gin.Context) {
	username := c.GetString(ctxkey.Username)
	logType, _ := strconv.Atoi(c.Query("type"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	channel, _ := strconv.Atoi(c.Query("channel"))
	quotaNum := model.SumUsedQuota(logType, startTimestamp, endTimestamp, modelName, username, tokenName, channel)
	normalQuota := model.SumUsedQuotaByBillingSource(logType, startTimestamp, endTimestamp, modelName, username, tokenName, channel, 0)
	subscriptionQuota := model.SumUsedQuotaByBillingSource(logType, startTimestamp, endTimestamp, modelName, username, tokenName, channel, 1)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"quota":              quotaNum,
			"normal_quota":       normalQuota,
			"subscription_quota": subscriptionQuota,
		},
	})
	return
}

func DeleteHistoryLogs(c *gin.Context) {
	targetTimestamp, _ := strconv.ParseInt(c.Query("target_timestamp"), 10, 64)
	if targetTimestamp == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "target timestamp is required",
		})
		return
	}
	count, err := model.DeleteOldLog(targetTimestamp)
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
		"data":    count,
	})
	return
}