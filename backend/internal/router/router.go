// router.go 注册所有 API 路由分组，挂载中间件，统一管理 /api 下的接口
package router

import (
	"ElainaBlog/internal/article"
	"ElainaBlog/internal/authorprofile"
	"ElainaBlog/internal/category"
	"ElainaBlog/internal/comment"
	"ElainaBlog/internal/auth"
	"ElainaBlog/internal/config"
	"ElainaBlog/internal/friendlink"
	"ElainaBlog/internal/message"
	"ElainaBlog/internal/middleware/ipban"
	"ElainaBlog/internal/middleware/jwt"
	"ElainaBlog/internal/middleware/ratelimit"
	"ElainaBlog/internal/middleware/uploadlimit"
	"ElainaBlog/internal/notification"
	"ElainaBlog/internal/site"
	"ElainaBlog/internal/siteconfig"
	"ElainaBlog/internal/upload"
	"ElainaBlog/internal/user"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Options 收纳路由装配所需的依赖与配置
type Options struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Jwt    auth.TokenManager
}

// New 创建并返回完整的 gin.Engine（参照 GoFeed 路由层）
func New(opts Options) *gin.Engine {
	// 设置 gin 运行模式：开发模式 Debug，生产模式 Release
	if opts.Config.Dev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 仅信任本机回环代理（生产由宿主机 nginx 回源），
	// 防止客户端伪造 X-Forwarded-For 绕过限流与 IP 封禁
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Printf("设置受信代理失败: %v", err)
	}

	r.Use(gin.Recovery())
	if opts.Config.Dev {
		r.Use(gin.Logger())
	}

	// CORS 配置（前端开发地址固定为 localhost:5173）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// IP 封禁检查：命中即拒绝，作用于全部 /api 路由
	r.Use(ipban.Middleware(opts.Redis))

	// 获取依赖实例
	gormDB := opts.DB
	redis := opts.Redis
	tokenMgr := opts.Jwt

	// 创建仓储层
	userRepo := user.NewRepository(gormDB)
	categoryRepo := category.NewRepository(gormDB)
	commentRepo := comment.NewRepository(gormDB)
	articleRepo := article.NewRepository(gormDB, redis)
	messageRepo := message.NewRepository(gormDB)
	siteRepo := site.NewRepository(gormDB)
	friendLinkRepo := friendlink.NewRepository(gormDB)
	notificationRepo := notification.NewRepository(gormDB)

	// 创建服务层
	userService := user.NewService(userRepo, redis, tokenMgr)
	categoryService := category.NewService(categoryRepo, redis)
	notificationService := notification.NewService(notificationRepo)
	commentService := comment.NewService(commentRepo, articleRepo, notificationService, userService)
	articleService := article.NewService(articleRepo, commentRepo)
	messageService := message.NewService(messageRepo, userRepo, notificationService)
	siteService := site.NewService(siteRepo, redis, gormDB)
	friendLinkService := friendlink.NewService(friendLinkRepo, redis)

	// 创建控制器
	userController := user.NewController(userService, redis)
	categoryController := category.NewController(categoryService)
	articleController := article.NewController(articleService, userService)
	commentController := comment.NewController(commentService, userService)
	uploadStorage := upload.NewLocalStorage(opts.Config.Upload.Path)
	avatarStorage := upload.NewLocalStorage(opts.Config.Upload.AvatarPath)
	uploadController := upload.NewController(uploadStorage, opts.Config.Upload.Size, avatarStorage, opts.Config.Upload.AvatarSize, userService)
	siteController := site.NewController(siteService, gormDB, redis)
	messageController := message.NewController(messageService, userService)
	siteConfigController := siteconfig.NewController(siteconfig.NewService(siteconfig.NewRepository(gormDB), redis))
	authorProfileController := authorprofile.NewController(authorprofile.NewService(authorprofile.NewRepository(gormDB), redis))
	friendLinkController := friendlink.NewController(friendLinkService)
	notificationController := notification.NewController(notificationService)

	// 无需鉴权
	r.GET("/health", health)
	r.Static("/uploads", opts.Config.Upload.Path)

	// api路由组
	apiGroup := r.Group("/api/ui")
	{
		// 公开接口
		apiGroup.POST("/login", ratelimit.Limit(redis, "login", 10, time.Minute), userController.Login)
		apiGroup.POST("/register", ratelimit.Limit(redis, "register", 5, time.Minute), userController.Register)
		apiGroup.POST("/refresh", ratelimit.Limit(redis, "refresh", 30, time.Minute), userController.RefreshToken)
		apiGroup.POST("/send-code", ratelimit.Limit(redis, "send-code", 5, time.Minute), userController.SendCode)
		apiGroup.POST("/reset-password", ratelimit.Limit(redis, "reset-password", 5, time.Minute), userController.ResetPassword)
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
		apiGroup.POST("/visit", siteController.RecordVisit)

		// 需要登录的接口
		authGroup := apiGroup.Group("", jwt.RequireAuth(tokenMgr))
		{
			authGroup.GET("/user/profile", userController.GetProfile)
			authGroup.POST("/user/profile", userController.UpdateProfile)
			authGroup.POST("/user/password", userController.UpdatePassword)
			authGroup.POST("/user/delete", userController.DeleteUser)
			authGroup.POST("/upload", uploadlimit.Limit(redis, opts.Config.Upload.RateLimit, time.Duration(opts.Config.Upload.RateWindow)*time.Second), uploadController.Upload)
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
		adminGroup := apiGroup.Group("", jwt.RequireAuth(tokenMgr), jwt.RequireAdmin(userService))
		{
			adminGroup.GET("/dashboard/stats", siteController.GetDashboardStats)
			adminGroup.GET("/article/admin-list", articleController.GetAdminList)
			adminGroup.GET("/article/admin/:id", articleController.GetAdminByID)
			adminGroup.GET("/article/:id/uv", articleController.GetArticleUV)
			adminGroup.POST("/article/toggle-top", articleController.ToggleTop)
			adminGroup.POST("/category/create", categoryController.Create)
			adminGroup.POST("/category/update", categoryController.Update)
			adminGroup.POST("/category/delete", categoryController.Delete)
			adminGroup.POST("/category/toggle-top", categoryController.ToggleTop)
			adminGroup.GET("/user/list", userController.GetList)
			adminGroup.GET("/comment/list", commentController.GetAllList)
			adminGroup.GET("/system/status", siteController.GetSystemStatus)
			adminGroup.POST("/cache/clear", siteController.ClearCache)
			adminGroup.GET("/backup/export", siteController.ExportBackup)
			adminGroup.GET("/security/banned-ips", siteController.GetBannedIPs)
			adminGroup.POST("/security/ban", siteController.BanIP)
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

	return r
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
