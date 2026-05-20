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
	rateLimiter := middleware.NewRateLimitMiddleware()
	userController := user.NewController(userService)
	categoryController := category.NewController(userService)
	articleController := article.NewController(userService)
	commentController := comment.NewController(userService)
	uploadStorage := upload.NewLocalStorage(config.GlobalConfig.Upload.Path)
	avatarStorage := upload.NewLocalStorage(config.GlobalConfig.Upload.AvatarPath)
	uploadController := upload.NewController(uploadStorage, config.GlobalConfig.Upload.Size, avatarStorage, config.GlobalConfig.Upload.AvatarSize, userService)
	siteController := site.NewController(site.NewService(site.NewRepository(db.DBPool)), userService)
	messageController := message.NewController(userService)
	siteConfigController := siteconfig.NewController(userService)
	authorProfileController := authorprofile.NewController(userService)

	// 无需鉴权
	r.GET("/health", health)

	r.Static("/uploads", config.GlobalConfig.Upload.Path)

	// api路由组
	apiGroup := r.Group("/api/ui")
	{
		apiGroup.POST("/login", rateLimiter.Limit("login", 10, time.Minute), userController.Login)
		apiGroup.POST("/register", rateLimiter.Limit("register", 5, time.Minute), userController.Register)
		apiGroup.POST("/refresh", rateLimiter.Limit("refresh", 30, time.Minute), userController.RefreshToken)
		apiGroup.POST("/send-code", rateLimiter.Limit("send-code", 5, time.Minute), userController.SendCode)

		// 仪表盘统计（管理员）
		apiGroup.GET("/dashboard/stats", auth.RequireAuth(), siteController.GetDashboardStats)

		// 分类（列表公开）
		apiGroup.GET("/category/list", categoryController.GetList)

		// 上传接口（需鉴权，登录用户）
		apiGroup.POST("/upload", auth.RequireAuth(), uploadController.Upload)
		apiGroup.POST("/upload/avatar", auth.RequireAuth(), uploadController.UploadAvatar)

		// 需要鉴权的路由
		userGroup := apiGroup.Group("/user", auth.RequireAuth())
		{
			userGroup.GET("/profile", userController.GetProfile)
			userGroup.GET("/list", userController.GetList)
			userGroup.POST("/profile", userController.UpdateProfile)
			userGroup.POST("/password", userController.UpdatePassword)
			userGroup.POST("/delete", userController.DeleteUser)
		}

		// 分类管理（需鉴权，仅管理员）
		categoryGroup := apiGroup.Group("/category", auth.RequireAuth())
		{
			categoryGroup.POST("/create", categoryController.Create)
			categoryGroup.POST("/update", categoryController.Update)
			categoryGroup.POST("/delete", categoryController.Delete)
		}

		// 文章（列表和详情公开）
		apiGroup.GET("/article/list", articleController.GetList)
		apiGroup.GET("/article/:id", articleController.GetByID)

		// 文章管理（需鉴权）
		articleGroup := apiGroup.Group("/article", auth.RequireAuth())
		{
			articleGroup.POST("/create", articleController.CreateArticle)
			articleGroup.POST("/update", articleController.UpdateArticle)
			articleGroup.POST("/delete", articleController.DeleteArticle)
		}

		// 评论（列表公开）
		apiGroup.GET("/comment/:article_id", commentController.GetList)

		// 评论管理（需鉴权，登录用户）
		commentGroup := apiGroup.Group("/comment", auth.RequireAuth())
		{
			commentGroup.POST("/create", commentController.CreateComment)
			commentGroup.POST("/delete", commentController.DeleteComment)
			commentGroup.GET("/list", commentController.GetAllList)
		}

		// 系统管理（需鉴权，仅管理员）
		systemGroup := apiGroup.Group("/system", auth.RequireAuth())
		{
			systemGroup.GET("/status", siteController.GetSystemStatus)
		}

		// 缓存管理（需鉴权，仅管理员）
		cacheGroup := apiGroup.Group("/cache", auth.RequireAuth())
		{
			cacheGroup.POST("/clear", siteController.ClearCache)
		}

		// 备份管理（需鉴权，仅管理员）
		backupGroup := apiGroup.Group("/backup", auth.RequireAuth())
		{
			backupGroup.GET("/export", siteController.ExportBackup)
		}

		// 安全管理（需鉴权，仅管理员）
		securityGroup := apiGroup.Group("/security", auth.RequireAuth())
		{
			securityGroup.GET("/banned-ips", siteController.GetBannedIPs)
			securityGroup.POST("/unban", siteController.UnbanIP)
		}

		// 作者页统计（公开）
		apiGroup.GET("/author/stats", siteController.GetAuthorStats)

		// 站点配置（公开）
		apiGroup.GET("/site-config", siteConfigController.GetPublicConfigs)
		apiGroup.GET("/site-config/quotes", siteConfigController.GetQuotes)

		// 作者信息（公开）
		apiGroup.GET("/author/profile", authorProfileController.Get)

		// 站点配置管理（需鉴权，仅管理员）
		siteConfigGroup := apiGroup.Group("/site-config", auth.RequireAuth())
		{
			siteConfigGroup.GET("/all", siteConfigController.GetAll)
			siteConfigGroup.POST("/update", siteConfigController.Upsert)
			siteConfigGroup.POST("/delete", siteConfigController.Delete)
		}

		// 作者信息管理（需鉴权，仅管理员）
		authorProfileGroup := apiGroup.Group("/author/profile", auth.RequireAuth())
		{
			authorProfileGroup.POST("/update", authorProfileController.Update)
		}

		// 留言板（列表公开）
		apiGroup.GET("/message/list", messageController.GetList)

		// 留言管理（需鉴权，登录用户）
		messageGroup := apiGroup.Group("/message", auth.RequireAuth())
		{
			messageGroup.POST("/create", messageController.Create)
			messageGroup.POST("/delete", messageController.Delete)
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
