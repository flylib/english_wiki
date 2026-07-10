package apis

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/wiki/models"
)

type Category struct {
	api.Api
}

// CategoryNode 树节点
type CategoryNode struct {
	models.WikiCategory
	EntryCount int64           `json:"entryCount"`
	Children   []*CategoryNode `json:"children"`
}

// Tree 返回完整分类树(含各分类已发布条目数)
func (e Category) Tree(c *gin.Context) {
	e.MakeContext(c)
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	var list []models.WikiCategory
	if err := db.Order("sort asc, id asc").Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	// 条目计数(按分类聚合一次查询)
	type cntRow struct {
		CategoryId int
		Cnt        int64
	}
	var cnts []cntRow
	db.Model(&models.WikiEntry{}).Select("category_id, count(*) as cnt").
		Where("status = ?", "published").Group("category_id").Scan(&cnts)
	cntMap := map[int]int64{}
	for _, r := range cnts {
		cntMap[r.CategoryId] = r.Cnt
	}

	nodes := map[int]*CategoryNode{}
	var roots []*CategoryNode
	for _, cat := range list {
		nodes[cat.Id] = &CategoryNode{WikiCategory: cat, EntryCount: cntMap[cat.Id], Children: []*CategoryNode{}}
	}
	for _, cat := range list {
		n := nodes[cat.Id]
		if p, ok := nodes[cat.ParentId]; ok && cat.ParentId != cat.Id {
			p.Children = append(p.Children, n)
		} else {
			roots = append(roots, n)
		}
	}
	e.OK(roots, "")
}

// Insert 新建分类(超管)
func (e Category) Insert(c *gin.Context) {
	e.MakeContext(c)
	if !isSuperAdmin(c) {
		e.Error(403, nil, "仅管理员可管理分类")
		return
	}
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	var m models.WikiCategory
	if err := c.ShouldBindJSON(&m); err != nil {
		e.Error(422, err, "参数错误")
		return
	}
	if m.Name == "" || m.Code == "" {
		e.Error(422, nil, "名称与编码必填")
		return
	}
	if err := db.Create(&m).Error; err != nil {
		e.Error(500, err, "创建失败(编码可能重复)")
		return
	}
	e.OK(m, "创建成功")
}

// Update 更新分类(超管)
func (e Category) Update(c *gin.Context) {
	e.MakeContext(c)
	if !isSuperAdmin(c) {
		e.Error(403, nil, "仅管理员可管理分类")
		return
	}
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var m models.WikiCategory
	if err := db.First(&m, id).Error; err != nil {
		e.Error(404, err, "分类不存在")
		return
	}
	var req models.WikiCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(422, err, "参数错误")
		return
	}
	if req.ParentId == id {
		e.Error(422, nil, "父分类不能是自己")
		return
	}
	updates := map[string]interface{}{
		"parent_id": req.ParentId,
		"name":      req.Name,
		"code":      req.Code,
		"icon":      req.Icon,
		"sort":      req.Sort,
	}
	if err := db.Model(&m).Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(m, "更新成功")
}

// Delete 删除分类(超管;有子分类或条目时拒绝)
func (e Category) Delete(c *gin.Context) {
	e.MakeContext(c)
	if !isSuperAdmin(c) {
		e.Error(403, nil, "仅管理员可管理分类")
		return
	}
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var child, entries int64
	db.Model(&models.WikiCategory{}).Where("parent_id = ?", id).Count(&child)
	db.Model(&models.WikiEntry{}).Where("category_id = ?", id).Count(&entries)
	if child > 0 || entries > 0 {
		e.Error(422, nil, "该分类下还有子分类或条目,不能删除")
		return
	}
	if err := db.Delete(&models.WikiCategory{}, id).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(nil, "删除成功")
}
