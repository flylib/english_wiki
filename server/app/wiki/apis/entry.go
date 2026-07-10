package apis

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"gorm.io/gorm"

	"go-admin/app/wiki/models"
)

type Entry struct {
	api.Api
}

// List 条目列表。
// 过滤:categoryId(默认含子分类)、cefrMin/cefrMax(区间有交集即命中)、keyword、status。
// 阅读端传 status=published;管理端不传 status 看全部。
func (e Entry) List(c *gin.Context) {
	e.MakeContext(c)
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	pageIndex := pageInt(c.Query("pageIndex"), 1)
	pageSize := pageInt(c.Query("pageSize"), 50)

	q := db.Model(&models.WikiEntry{}).
		Select("id, category_id, title, summary, cefr_min, cefr_max, status, sort, created_at, updated_at")

	if v := c.Query("categoryId"); v != "" {
		catId, _ := strconv.Atoi(v)
		ids := []int{catId}
		if c.DefaultQuery("includeSub", "1") == "1" {
			ids = append(ids, descendantIds(db, catId)...)
		}
		q = q.Where("category_id in ?", ids)
	}
	// 等级过滤:条目适用区间 [cefr_min, cefr_max] 与请求区间有交集
	if v := c.Query("cefrMin"); v != "" {
		min, _ := strconv.Atoi(v)
		q = q.Where("cefr_max >= ?", min)
	}
	if v := c.Query("cefrMax"); v != "" {
		max, _ := strconv.Atoi(v)
		q = q.Where("cefr_min <= ?", max)
	}
	if v := c.Query("keyword"); v != "" {
		q = q.Where("title like ? or summary like ?", "%"+v+"%", "%"+v+"%")
	}
	if v := c.Query("status"); v != "" {
		q = q.Where("status = ?", v)
	}

	var count int64
	q.Count(&count)
	var list []models.WikiEntry
	if err := q.Order("sort asc, id asc").Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&list).Error; err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), pageIndex, pageSize, "")
}

// Get 条目详情(含 Markdown 正文)
func (e Entry) Get(c *gin.Context) {
	e.MakeContext(c)
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var m models.WikiEntry
	if err := db.First(&m, id).Error; err != nil {
		e.Error(404, err, "条目不存在")
		return
	}
	e.OK(m, "")
}

type entryReq struct {
	CategoryId int    `json:"categoryId"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	ContentMd  string `json:"contentMd"`
	CefrMin    int    `json:"cefrMin"`
	CefrMax    int    `json:"cefrMax"`
	Status     string `json:"status"`
	Sort       int    `json:"sort"`
}

func (r *entryReq) normalize() (string, bool) {
	if r.Title == "" || r.CategoryId == 0 {
		return "标题与分类必填", false
	}
	if r.CefrMin == 0 {
		r.CefrMin = 1
	}
	if r.CefrMax == 0 {
		r.CefrMax = 7
	}
	if r.CefrMin < 1 || r.CefrMax > 7 || r.CefrMin > r.CefrMax {
		return "CEFR 区间不合法(1~7 且下限≤上限)", false
	}
	if r.Status == "" {
		r.Status = "draft"
	}
	if r.Status != "draft" && r.Status != "published" {
		return "状态只能是 draft/published", false
	}
	return "", true
}

// Insert 新建条目(超管)
func (e Entry) Insert(c *gin.Context) {
	e.MakeContext(c)
	if !isSuperAdmin(c) {
		e.Error(403, nil, "仅管理员可编辑内容")
		return
	}
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	var req entryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(422, err, "参数错误")
		return
	}
	if msg, ok := req.normalize(); !ok {
		e.Error(422, nil, msg)
		return
	}
	m := models.WikiEntry{
		CategoryId: req.CategoryId, Title: req.Title, Summary: req.Summary,
		ContentMd: req.ContentMd, CefrMin: req.CefrMin, CefrMax: req.CefrMax,
		Status: req.Status, Sort: req.Sort,
	}
	if err := db.Create(&m).Error; err != nil {
		e.Error(500, err, "创建失败")
		return
	}
	e.OK(m, "创建成功")
}

// Update 更新条目(超管)
func (e Entry) Update(c *gin.Context) {
	e.MakeContext(c)
	if !isSuperAdmin(c) {
		e.Error(403, nil, "仅管理员可编辑内容")
		return
	}
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var m models.WikiEntry
	if err := db.First(&m, id).Error; err != nil {
		e.Error(404, err, "条目不存在")
		return
	}
	var req entryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(422, err, "参数错误")
		return
	}
	if msg, ok := req.normalize(); !ok {
		e.Error(422, nil, msg)
		return
	}
	updates := map[string]interface{}{
		"category_id": req.CategoryId,
		"title":       req.Title,
		"summary":     req.Summary,
		"content_md":  req.ContentMd,
		"cefr_min":    req.CefrMin,
		"cefr_max":    req.CefrMax,
		"status":      req.Status,
		"sort":        req.Sort,
	}
	if err := db.Model(&m).Updates(updates).Error; err != nil {
		e.Error(500, err, "更新失败")
		return
	}
	e.OK(m, "更新成功")
}

// Delete 删除条目(超管,软删)
func (e Entry) Delete(c *gin.Context) {
	e.MakeContext(c)
	if !isSuperAdmin(c) {
		e.Error(403, nil, "仅管理员可编辑内容")
		return
	}
	db, err := e.GetOrm()
	if err != nil {
		e.Error(500, err, "")
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := db.Delete(&models.WikiEntry{}, id).Error; err != nil {
		e.Error(500, err, "删除失败")
		return
	}
	e.OK(nil, "删除成功")
}

// descendantIds 返回某分类的全部后代分类 id(分类量小,一次载入内存遍历)
func descendantIds(db *gorm.DB, catId int) []int {
	var all []models.WikiCategory
	if err := db.Select("id, parent_id").Find(&all).Error; err != nil {
		return nil
	}
	children := map[int][]int{}
	for _, c := range all {
		children[c.ParentId] = append(children[c.ParentId], c.Id)
	}
	var out []int
	queue := append([]int{}, children[catId]...)
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		out = append(out, id)
		queue = append(queue, children[id]...)
	}
	return out
}
