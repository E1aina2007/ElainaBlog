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
	"ElainaBlog/internal/message"
	"ElainaBlog/internal/middleware"
	"ElainaBlog/internal/site"
	"ElainaBlog/internal/siteconfig"
	"ElainaBlog/internal/upload"
	"ElainaBlog/internal/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	// 统一创建 user.Service 实例
	userService := user.NewService(user.NewRepository(db.DBPool))

	auth := middleware.NewJwtAuthMiddleware(common.JwtAuth)
	adminAuth := middleware.NewAdminAuthMiddleware(userService)
	rateLimiter := middleware.NewRateLimitMiddleware()
	userController := user.NewController(userService)
	categoryController := category.NewController()
	articleController := article.NewController(userService)
	commentController := comment.NewController(userService)
	uploadStorage := upload.NewLocalStorage(config.GlobalConfig.Upload.Path)
	avatarStorage := upload.NewLocalStorage(config.GlobalConfig.Upload.AvatarPath)
	uploadController := upload.NewController(uploadStorage, config.GlobalConfig.Upload.Size, avatarStorage, config.GlobalConfig.Upload.AvatarSize, userService)
	siteController := site.NewController(site.NewService(site.NewRepository(db.DBPool)))
	messageController := message.NewController(userService)
	siteConfigController := siteconfig.NewController()
	authorProfileController := authorprofile.NewController()

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
		apiGroup.GET("/category/list", categoryController.GetList)
		apiGroup.GET("/article/list", articleController.GetList)
		apiGroup.GET("/article/:id", articleController.GetByID)
		apiGroup.GET("/comment/:article_id", commentController.GetList)
		apiGroup.GET("/author/stats", siteController.GetAuthorStats)
		apiGroup.GET("/site-config", siteConfigController.GetPublicConfigs)
		apiGroup.GET("/site-config/quotes", siteConfigController.GetQuotes)
		apiGroup.GET("/author/profile", authorProfileController.Get)
		apiGroup.GET("/message/list", messageController.GetList)

		// 需要登录的接口
		authGroup := apiGroup.Group("", auth.RequireAuth())
		{
			authGroup.GET("/user/profile", userController.GetProfile)
			authGroup.POST("/user/profile", userController.UpdateProfile)
			authGroup.POST("/user/password", userController.UpdatePassword)
			authGroup.POST("/user/delete", userController.DeleteUser)
			authGroup.POST("/upload", uploadController.Upload)
			authGroup.POST("/upload/avatar", uploadController.UploadAvatar)
			authGroup.POST("/article/create", articleController.CreateArticle)
			authGroup.POST("/article/update", articleController.UpdateArticle)
			authGroup.POST("/article/delete", articleController.DeleteArticle)
			authGroup.POST("/comment/create", commentController.CreateComment)
			authGroup.POST("/comment/delete", commentController.DeleteComment)
			authGroup.POST("/message/create", messageController.Create)
			authGroup.POST("/message/delete", messageController.Delete)
			authGroup.POST("/logout", userController.Logout)
		}

		// 需要管理员权限的接口
		adminGroup := apiGroup.Group("", auth.RequireAuth(), adminAuth.RequireAdmin())
		{
			adminGroup.GET("/dashboard/stats", siteController.GetDashboardStats)
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
