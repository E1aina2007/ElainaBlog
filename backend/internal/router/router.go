// router.go 注册所有 API 路由分组，挂载中间件，统一管理 /api 下的接口
package router

import (
	"ElainaBlog/config"
	"ElainaBlog/config/db"
	"ElainaBlog/internal/article"
	"ElainaBlog/internal/category"
	"ElainaBlog/internal/comment"
	"ElainaBlog/internal/common"
	"ElainaBlog/internal/middleware"
	"ElainaBlog/internal/site"
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
	userController := user.NewController(userService)
	categoryController := category.NewController(userService)
	articleController := article.NewController(userService)
	commentController := comment.NewController(userService)
	uploadStorage := upload.NewLocalStorage(config.GlobalConfig.Upload.Path)
	uploadController := upload.NewController(uploadStorage, config.GlobalConfig.Upload.Size)
	siteController := site.NewController(site.NewService(site.NewRepository(db.DBPool)), userService)

	// 无需鉴权
	r.GET("/health", health)

	r.Static("/uploads", config.GlobalConfig.Upload.Path)

	// api路由组
	apiGroup := r.Group("/api/ui")
	{
		apiGroup.POST("/login", userController.Login)
		apiGroup.POST("/register", userController.Register)
		apiGroup.POST("/refresh", userController.RefreshToken)
		apiGroup.POST("/send-code", userController.SendCode)

		// 站点配置（公开读取，管理员更新）
		apiGroup.GET("/site", siteController.GetList)
		apiGroup.POST("/site", auth.RequireAuth(), siteController.Update)

		// 分类（列表公开）
		apiGroup.GET("/category/list", categoryController.GetList)

		apiGroup.POST("/upload", auth.RequireAuth(), uploadController.Upload)

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

		// 文章管理（需鉴权，仅管理员）
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
