package router

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"

	"go-admin/app/wiki/service"
	common "go-admin/common/middleware"
)

// InitRouter 由 cmd/api/wiki.go 注册进 AppRouters,在 DB 初始化之后执行
func InitRouter() {
	h := sdk.Runtime.GetEngine()
	r, ok := h.(*gin.Engine)
	if !ok {
		log.Fatal("not support other engine")
		os.Exit(-1)
	}

	// 自动建表 + 幂等注入种子数据
	for _, db := range sdk.Runtime.GetDb() {
		if err := service.Migrate(db); err != nil {
			log.Errorf("wiki migrate error: %s", err.Error())
		}
		break
	}

	authMiddleware, err := common.AuthInit()
	if err != nil {
		log.Fatalf("JWT Init Error, %s", err.Error())
	}
	initRouter(r, authMiddleware)
}
