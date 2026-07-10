package apis

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

// isSuperAdmin 是否超级管理员(go-admin admin 角色);写操作仅超管可用
func isSuperAdmin(c *gin.Context) bool {
	return user.GetRoleName(c) == "admin"
}

func pageInt(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil && v > 0 {
		return v
	}
	return def
}
