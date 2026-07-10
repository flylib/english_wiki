package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/wiki/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerWikiRouter)
}

// 全部接口挂在 JWT 鉴权下;写操作在 API 内校验超管角色。
func registerWikiRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	cat := apis.Category{}
	entry := apis.Entry{}
	level := apis.LevelApi{}

	r := v1.Group("/wiki").Use(authMiddleware.MiddlewareFunc())
	{
		// 分类树
		r.GET("/category/tree", cat.Tree)
		r.POST("/category", cat.Insert)
		r.PUT("/category/:id", cat.Update)
		r.DELETE("/category/:id", cat.Delete)

		// 条目
		r.GET("/entry", entry.List)
		r.GET("/entry/:id", entry.Get)
		r.POST("/entry", entry.Insert)
		r.PUT("/entry/:id", entry.Update)
		r.DELETE("/entry/:id", entry.Delete)

		// 等级体系
		r.GET("/level/systems", level.Systems)
	}
}
