package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/wiki/models"
)

type LevelApi struct {
	api.Api
}

type systemWithLevels struct {
	models.LevelSystem
	Levels []models.Level `json:"levels"`
}

// Systems 返回全部等级体系及档位(阅读页身份切换器数据源)
func (e LevelApi) Systems(c *gin.Context) {
	e.MakeContext(c)
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	var systems []models.LevelSystem
	if err := db.Order("sort asc").Find(&systems).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	var levels []models.Level
	if err := db.Order("system_id asc, sort asc").Find(&levels).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	bySys := map[int][]models.Level{}
	for _, l := range levels {
		bySys[l.SystemId] = append(bySys[l.SystemId], l)
	}
	out := make([]systemWithLevels, 0, len(systems))
	for _, s := range systems {
		out = append(out, systemWithLevels{LevelSystem: s, Levels: bySys[s.Id]})
	}
	e.OK(out, "")
}
