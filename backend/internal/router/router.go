// router.go 注册所有 API 路由分组，挂载中间件，统一管理 /api 下的接口
package router

import (
	"ElainaBlog/config"
	"ElainaBlog/config/db"
	"ElainaBlog/internal/article"
	"ElainaBlog/internal/authorprofile"
	"ElainaBlog/internal/category"
	"ElainaBlog/internal/comment"
	"ElainaBlog/internal/common"
	"ElainaBlog/pkg/util"
	"ElainaBlog/internal/friendlink"
	"ElainaBlog/internal/message"
	"ElainaBlog/internal/middleware"
	"ElainaBlog/internal/notification"
	"ElainaBlog/internal/site"
	"ElainaBlog/internal/siteconfig"
	"ElainaBlog/internal/upload"
	"ElainaBlog/internal/user"
	"ElainaBlog/pkg/rdb"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	// 获取依赖实例
	dbPool := db.DBPool
	redis := rdb.DefaultClient
	tokenMgr := common.JwtAuth

	// 创建仓储层
	userRepo := user.NewRepository(dbPool)
	categoryRepo := category.NewRepository(dbPool)
	commentRepo := comment.NewRepository(dbPool)
	articleRepo := article.NewRepository(dbPool, redis)
	messageRepo := message.NewRepository(dbPool)
	siteRepo := site.NewRepository(dbPool)
	friendLinkRepo := friendlink.NewRepository(dbPool)
	notificationRepo := notification.NewRepository(dbPool)

	// 创建服务层
	userService := user.NewService(userRepo, redis, tokenMgr)
	categoryService := category.NewService(categoryRepo)
	notificationService := notification.NewService(notificationRepo)
	commentService := comment.NewService(commentRepo, articleRepo, notificationService, userService)
	articleService := article.NewService(articleRepo, commentRepo)
	messageService := message.NewService(messageRepo, userRepo, notificationService)
	siteService := site.NewService(siteRepo, redis)
	friendLinkService := friendlink.NewService(friendLinkRepo)

	// 创建中间件
	auth := middleware.NewJwtAuthMiddleware(tokenMgr)
	adminAuth := middleware.NewAdminAuthMiddleware(userService)
	rateLimiter := middleware.NewRateLimitMiddleware(redis)
	uploadLimiter := middleware.NewUploadLimitMiddleware(redis)

	// 创建控制器
	userController := user.NewController(userService, redis)
	categoryController := category.NewController(categoryService)
	articleController := article.NewController(articleService, userService)
	commentController := comment.NewController(commentService, userService)
	uploadStorage := upload.NewLocalStorage(config.GlobalConfig.Upload.Path)
	avatarStorage := upload.NewLocalStorage(config.GlobalConfig.Upload.AvatarPath)
	uploadController := upload.NewController(uploadStorage, config.GlobalConfig.Upload.Size, avatarStorage, config.GlobalConfig.Upload.AvatarSize, userService)
	siteController := site.NewController(siteService, dbPool, redis)
	messageController := message.NewController(messageService, userService)
	siteConfigController := siteconfig.NewController(siteconfig.NewService(siteconfig.NewRepository(dbPool)))
	authorProfileController := authorprofile.NewController(authorprofile.NewService(authorprofile.NewRepository(dbPool)))
	friendLinkController := friendlink.NewController(friendLinkService)
	notificationController := notification.NewController(notificationService)

	// 无需鉴权
	r.GET("/health", health)
	r.Static("/uploads", config.GlobalConfig.Upload.Path)

	// api路由组
	apiGroup := r.Group("/api/ui")
	{
		// 公开接口
		apiGroup.POST("/login", rateLimiter.Limit("login", 10, time.Minute), userController.Login)
		apiGroup.POST("/register", rateLimiter.Limit("register", 5, time.Minute), userController.Register)
		apiGroup.POST("/refresh", rateLimiter.Limit("refresh", 30, time.Minute), userController.RefreshToken)
		apiGroup.POST("/send-code", rateLimiter.Limit("send-code", 5, time.Minute), userController.SendCode)
		apiGroup.POST("/reset-password", rateLimiter.Limit("reset-password", 5, time.Minute), userController.ResetPassword)
		apiGroup.GET("/category/list", categoryController.GetList)
		apiGroup.GET("/article/list", articleController.GetList)
		apiGroup.GET("/article/search", articleController.Search)
		apiGroup.GET("/article/:id", articleController.GetByID)
		apiGroup.GET("/comment/:article_id", commentController.GetList)
		apiGroup.GET("/author/stats", siteController.GetAuthorStats)
		apiGroup.GET("/site-config", siteConfigController.GetPublicConfigs)
		apiGroup.GET("/site-config/quotes", siteConfigController.GetQuotes)
		apiGroup.GET("/author/profile", authorProfileController.Get)
		apiGroup.GET("/message/list", messageController.GetList)
		apiGroup.GET("/friend-link/list", friendLinkController.GetList)
		apiGroup.GET("/favicon", rateLimiter.Limit("favicon", 30, time.Minute), util.FaviconProxy)
		apiGroup.POST("/visit", siteController.RecordVisit)

		// 需要登录的接口
		authGroup := apiGroup.Group("", auth.RequireAuth())
		{
			authGroup.GET("/user/profile", userController.GetProfile)
			authGroup.POST("/user/profile", userController.UpdateProfile)
			authGroup.POST("/user/password", userController.UpdatePassword)
			authGroup.POST("/user/delete", userController.DeleteUser)
			authGroup.POST("/upload", uploadLimiter.Limit(config.GlobalConfig.Upload.RateLimit, time.Duration(config.GlobalConfig.Upload.RateWindow)*time.Second), uploadController.Upload)
			authGroup.POST("/upload/avatar", uploadController.UploadAvatar)
			authGroup.POST("/article/create", articleController.CreateArticle)
			authGroup.POST("/article/update", articleController.UpdateArticle)
			authGroup.POST("/article/delete", articleController.DeleteArticle)
			authGroup.POST("/comment/create", commentController.CreateComment)
			authGroup.POST("/comment/delete", commentController.DeleteComment)
			authGroup.POST("/message/create", messageController.Create)
			authGroup.POST("/message/delete", messageController.Delete)
			authGroup.GET("/notification/list", notificationController.GetList)
			authGroup.GET("/notification/unread-count", notificationController.GetUnreadCount)
			authGroup.POST("/notification/read", notificationController.MarkAsRead)
			authGroup.POST("/notification/read-all", notificationController.MarkAllAsRead)
			authGroup.POST("/notification/delete", notificationController.Delete)
			authGroup.GET("/article/mine", articleController.GetMyList)
			authGroup.GET("/article/mine/:id", articleController.GetMyByID)
			authGroup.POST("/logout", userController.Logout)
		}

		// 需要管理员权限的接口
		adminGroup := apiGroup.Group("", auth.RequireAuth(), adminAuth.RequireAdmin())
		{
			adminGroup.GET("/dashboard/stats", siteController.GetDashboardStats)
			adminGroup.GET("/article/admin-list", articleController.GetAdminList)
			adminGroup.GET("/article/admin/:id", articleController.GetAdminByID)
			adminGroup.GET("/article/:id/uv", articleController.GetArticleUV)
			adminGroup.POST("/category/create", categoryController.Create)
			adminGroup.POST("/category/update", categoryController.Update)
			adminGroup.POST("/category/delete", categoryController.Delete)
			adminGroup.GET("/user/list", userController.GetList)
			adminGroup.GET("/comment/list", commentController.GetAllList)
			adminGroup.GET("/system/status", siteController.GetSystemStatus)
			adminGroup.POST("/cache/clear", siteController.ClearCache)
			adminGroup.GET("/backup/export", siteController.ExportBackup)
			adminGroup.GET("/security/banned-ips", siteController.GetBannedIPs)
			adminGroup.POST("/security/unban", siteController.UnbanIP)
			adminGroup.GET("/site-config/all", siteConfigController.GetAll)
			adminGroup.POST("/site-config/update", siteConfigController.Upsert)
			adminGroup.POST("/site-config/delete", siteConfigController.Delete)
			adminGroup.POST("/author/profile/update", authorProfileController.Update)
			adminGroup.POST("/friend-link/create", friendLinkController.Create)
			adminGroup.POST("/friend-link/update", friendLinkController.Update)
			adminGroup.POST("/friend-link/delete", friendLinkController.Delete)
		}
	}
}

// health 健康检查接口
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "ElainaBlog",
		"msg":     "你好, 我是Elaina!",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}
