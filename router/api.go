package router

import (
	"github.com/modelbus/one-api-pro/controller"
	"github.com/modelbus/one-api-pro/controller/auth"
	"github.com/modelbus/one-api-pro/middleware"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))
	apiRouter.Use(middleware.GlobalAPIRateLimit())
	{
		apiRouter.GET("/status", controller.GetStatus)
		apiRouter.GET("/models", middleware.UserAuth(), controller.DashboardListModels)
		apiRouter.GET("/notice", controller.GetNotice)
		apiRouter.GET("/about", controller.GetAbout)
		apiRouter.GET("/home_page_content", controller.GetHomePageContent)
		apiRouter.GET("/verification", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.SendEmailVerification)
		apiRouter.GET("/reset_password", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.SendPasswordResetEmail)
		apiRouter.POST("/user/reset", middleware.CriticalRateLimit(), controller.ResetPassword)
		apiRouter.GET("/oauth/github", middleware.CriticalRateLimit(), auth.GitHubOAuth)
		apiRouter.GET("/oauth/oidc", middleware.CriticalRateLimit(), auth.OidcAuth)
		apiRouter.GET("/oauth/lark", middleware.CriticalRateLimit(), auth.LarkOAuth)
		apiRouter.GET("/oauth/state", middleware.CriticalRateLimit(), auth.GenerateOAuthCode)
		apiRouter.GET("/oauth/wechat", middleware.CriticalRateLimit(), auth.WeChatAuth)
		apiRouter.GET("/oauth/wechat/bind", middleware.CriticalRateLimit(), middleware.UserAuth(), auth.WeChatBind)
		apiRouter.GET("/oauth/email/bind", middleware.CriticalRateLimit(), middleware.UserAuth(), controller.EmailBind)
		apiRouter.POST("/topup", middleware.AdminAuth(), controller.AdminTopUp)

		userRoute := apiRouter.Group("/user")
		{
			userRoute.POST("/register", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.Register)
			userRoute.POST("/login", middleware.CriticalRateLimit(), controller.Login)
			userRoute.GET("/logout", controller.Logout)

			selfRoute := userRoute.Group("/")
			selfRoute.Use(middleware.UserAuth())
			{
				selfRoute.GET("/dashboard", controller.GetUserDashboard)
				selfRoute.GET("/self", controller.GetSelf)
				selfRoute.PUT("/self", controller.UpdateSelf)
				selfRoute.DELETE("/self", controller.DeleteSelf)
				selfRoute.GET("/token", controller.GenerateAccessToken)
				selfRoute.GET("/aff", controller.GetAffCode)
				selfRoute.POST("/topup", controller.TopUp)
				selfRoute.GET("/available_models", controller.GetUserAvailableModels)
			selfRoute.GET("/subscription", controller.GetUserSubscriptionInfo)
			}

			adminRoute := userRoute.Group("/")
			adminRoute.Use(middleware.AdminAuth())
			{
				adminRoute.GET("/", controller.GetAllUsers)
				adminRoute.GET("/search", controller.SearchUsers)
				adminRoute.GET("/:id", controller.GetUser)
				adminRoute.POST("/", controller.CreateUser)
				adminRoute.POST("/manage", controller.ManageUser)
				adminRoute.PUT("/", controller.UpdateUser)
				adminRoute.DELETE("/:id", controller.DeleteUser)
			}
		}
		optionRoute := apiRouter.Group("/option")
		optionRoute.Use(middleware.RootAuth())
		{
			optionRoute.GET("/", controller.GetOptions)
			optionRoute.PUT("/", controller.UpdateOption)
		}
		channelRoute := apiRouter.Group("/channel")
		channelRoute.Use(middleware.AdminAuth())
		{
			channelRoute.GET("/", controller.GetAllChannels)
			channelRoute.GET("/search", controller.SearchChannels)
			channelRoute.GET("/models", controller.ListAllModels)
			channelRoute.GET("/:id", controller.GetChannel)
			channelRoute.GET("/test", controller.TestChannels)
			channelRoute.GET("/test/:id", controller.TestChannel)
			channelRoute.GET("/update_balance", controller.UpdateAllChannelsBalance)
			channelRoute.GET("/update_balance/:id", controller.UpdateChannelBalance)
			channelRoute.POST("/", controller.AddChannel)
			channelRoute.PUT("/", controller.UpdateChannel)
			channelRoute.DELETE("/disabled", controller.DeleteDisabledChannel)
			channelRoute.DELETE("/:id", controller.DeleteChannel)
		}
		tokenRoute := apiRouter.Group("/token")
		tokenRoute.Use(middleware.UserAuth())
		{
			tokenRoute.GET("/", controller.GetAllTokens)
			tokenRoute.GET("/search", controller.SearchTokens)
			tokenRoute.GET("/:id", controller.GetToken)
			tokenRoute.POST("/", controller.AddToken)
			tokenRoute.PUT("/", controller.UpdateToken)
			tokenRoute.DELETE("/:id", controller.DeleteToken)
		}
		redemptionRoute := apiRouter.Group("/redemption")
		redemptionRoute.Use(middleware.AdminAuth())
		{
			redemptionRoute.GET("/", controller.GetAllRedemptions)
			redemptionRoute.GET("/search", controller.SearchRedemptions)
			redemptionRoute.GET("/:id", controller.GetRedemption)
			redemptionRoute.POST("/", controller.AddRedemption)
			redemptionRoute.PUT("/", controller.UpdateRedemption)
			redemptionRoute.DELETE("/:id", controller.DeleteRedemption)
		}
		logRoute := apiRouter.Group("/log")
		logRoute.GET("/", middleware.AdminAuth(), controller.GetAllLogs)
		logRoute.DELETE("/", middleware.AdminAuth(), controller.DeleteHistoryLogs)
		logRoute.GET("/stat", middleware.AdminAuth(), controller.GetLogsStat)
		logRoute.GET("/self/stat", middleware.UserAuth(), controller.GetLogsSelfStat)
		logRoute.GET("/search", middleware.AdminAuth(), controller.SearchAllLogs)
		logRoute.GET("/self", middleware.UserAuth(), controller.GetUserLogs)
		logRoute.GET("/self/search", middleware.UserAuth(), controller.SearchUserLogs)
		groupRoute := apiRouter.Group("/group")
		groupRoute.Use(middleware.AdminAuth())
		{
			groupRoute.GET("/", controller.GetGroups)
		}
		modelPriceRoute := apiRouter.Group("/model_price")
		modelPriceRoute.Use(middleware.RootAuth())
		{
			modelPriceRoute.GET("/", controller.GetAllModelPrices)
			modelPriceRoute.POST("/", controller.AddModelPrice)
			modelPriceRoute.PUT("/", controller.UpdateModelPrice)
			modelPriceRoute.DELETE("/:id", controller.DeleteModelPrice)
		}
		groupPriceRoute := apiRouter.Group("/group_price")
		groupPriceRoute.Use(middleware.RootAuth())
		{
			groupPriceRoute.GET("/", controller.GetAllGroupPrices)
			groupPriceRoute.POST("/", controller.AddGroupPrice)
			groupPriceRoute.PUT("/", controller.UpdateGroupPrice)
			groupPriceRoute.DELETE("/:id", controller.DeleteGroupPrice)
		}
		planRoute := apiRouter.Group("/plan")
		{
			// Public (user-facing) endpoints — no auth.
			planRoute.GET("/list", controller.GetPublicPlans)
			planRoute.GET("/detail/:id", controller.GetPublicPlanDetail)
			// Admin endpoints
			planRoute.GET("/", middleware.AdminAuth(), controller.GetAllPlans)
			planRoute.GET("/search", middleware.AdminAuth(), controller.SearchPlans)
			planRoute.GET("/:id", middleware.AdminAuth(), controller.GetPlan)
			planRoute.POST("/", middleware.RootAuth(), controller.AddPlan)
			planRoute.PUT("/", middleware.RootAuth(), controller.UpdatePlan)
			planRoute.DELETE("/:id", middleware.RootAuth(), controller.DeletePlan)
		}
		// /api/plan/current — authenticated user reads own active plan.
		apiRouter.GET("/plan/current", middleware.UserAuth(), controller.GetCurrentPlan)

		// Orders: /api/order/plan (user self-service) + /api/order/self/:id (user)
		// + /api/order (admin).
		orderSelfRoute := apiRouter.Group("/order/self")
		orderSelfRoute.Use(middleware.UserAuth())
		{
			orderSelfRoute.GET("", controller.GetMyOrders)
			orderSelfRoute.GET("/", controller.GetMyOrders)
			orderSelfRoute.GET("/:id", controller.GetMyOrder)
			orderSelfRoute.POST("/:id/cancel", controller.CancelMyOrder)
		}
		apiRouter.POST("/order/plan", middleware.UserAuth(), controller.CreatePlanOrder)
		orderAdminRoute := apiRouter.Group("/order")
		{
			orderAdminRoute.GET("", middleware.AdminAuth(), controller.GetAllOrders)
			orderAdminRoute.GET("/", middleware.AdminAuth(), controller.GetAllOrders)
			orderAdminRoute.GET("/search", middleware.AdminAuth(), controller.SearchOrders)
			orderAdminRoute.GET("/:id", middleware.AdminAuth(), controller.GetOrder)
			orderAdminRoute.PUT("/:id", middleware.AdminAuth(), controller.MarkOrderPaid)
			orderAdminRoute.DELETE("/:id", middleware.RootAuth(), controller.DeleteOrder)
		}

		// Payment callbacks (no auth; signature-verified) + admin mock.
		apiRouter.POST("/payment/wechat/notify", controller.WechatNotify)
		apiRouter.POST("/payment/alipay/notify", controller.AlipayNotify)
		apiRouter.POST("/payment/mock/notify", middleware.RootAuth(), controller.MockPay)

		// Settings: payment + plan-operations.
		settingRoute := apiRouter.Group("/setting")
		{
			settingRoute.GET("/payment", middleware.RootAuth(), controller.GetPaymentSettings)
			settingRoute.PUT("/payment/:method", middleware.RootAuth(), controller.PutPaymentMethod)
			settingRoute.GET("/plan", middleware.RootAuth(), controller.GetPlanSettings)
			settingRoute.PUT("/plan", middleware.RootAuth(), controller.PutPlanSettings)
		}

		subscriptionRoute := apiRouter.Group("/subscription")
		{
			subscriptionRoute.GET("/self", middleware.UserAuth(), controller.GetUserSubscriptions)
			subscriptionRoute.GET("/", middleware.AdminAuth(), controller.GetAllSubscriptions)
			subscriptionRoute.GET("/search", middleware.AdminAuth(), controller.SearchSubscriptions)
			subscriptionRoute.GET("/:id", middleware.AdminAuth(), controller.GetSubscriptionDetail)
			subscriptionRoute.GET("/:id/usage", middleware.UserAuth(), controller.GetSubscriptionUsage)
			subscriptionRoute.POST("/", middleware.AdminAuth(), controller.AddSubscription)
			subscriptionRoute.PUT("/", middleware.AdminAuth(), controller.UpdateSubscription)
			subscriptionRoute.DELETE("/:id", middleware.AdminAuth(), controller.DeleteSubscription)
		}
		clusterNodeRoute := apiRouter.Group("/cluster_node")
		clusterNodeRoute.Use(middleware.RootAuth())
		{
			clusterNodeRoute.GET("/", controller.GetAllClusterNodes)
			clusterNodeRoute.GET("/:id", controller.GetClusterNode)
			clusterNodeRoute.POST("/", controller.AddClusterNode)
			clusterNodeRoute.PUT("/", controller.UpdateClusterNode)
			clusterNodeRoute.DELETE("/:id", controller.DeleteClusterNode)
			clusterNodeRoute.POST("/:id/enable", controller.EnableClusterNode)
			clusterNodeRoute.GET("/ping/:id", controller.PingClusterNode)
		}
	}
}
