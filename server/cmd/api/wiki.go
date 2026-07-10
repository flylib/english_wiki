package api

import "go-admin/app/wiki/router"

func init() {
	// 注册 wiki 模块路由
	AppRouters = append(AppRouters, router.InitRouter)
}
