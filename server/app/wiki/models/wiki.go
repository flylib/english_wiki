package models

import "go-admin/common/models"

// CEFR 等级用整数 rank 表示,便于区间过滤:
// 1=Pre-A1 2=A1 3=A2 4=B1 5=B2 6=C1 7=C2

// WikiCategory 知识分类树(左侧菜单:语法/时态/从句/音标/词汇…)
type WikiCategory struct {
	models.Model
	ParentId int    `json:"parentId" gorm:"index;comment:父分类id(0=根)"`
	Name     string `json:"name" gorm:"size:64;comment:分类名称"`
	Code     string `json:"code" gorm:"size:64;uniqueIndex;comment:唯一编码(如 grammar/tense)"`
	Icon     string `json:"icon" gorm:"size:64;comment:图标"`
	Sort     int    `json:"sort" gorm:"comment:排序"`
	models.ModelTime
	models.ControlBy
}

func (WikiCategory) TableName() string { return "wiki_category" }

// WikiEntry 百科条目(一条知识点,内容 Markdown)
type WikiEntry struct {
	models.Model
	CategoryId int    `json:"categoryId" gorm:"index;comment:所属分类"`
	Title      string `json:"title" gorm:"size:255;comment:标题"`
	Summary    string `json:"summary" gorm:"size:512;comment:摘要"`
	ContentMd  string `json:"contentMd" gorm:"type:text;comment:正文(Markdown)"`
	CefrMin    int    `json:"cefrMin" gorm:"index;default:1;comment:适用CEFR下限(1=PreA1..7=C2)"`
	CefrMax    int    `json:"cefrMax" gorm:"index;default:7;comment:适用CEFR上限"`
	Status     string `json:"status" gorm:"size:16;default:published;comment:状态(draft/published)"`
	Sort       int    `json:"sort" gorm:"comment:排序"`
	models.ModelTime
	models.ControlBy
}

func (WikiEntry) TableName() string { return "wiki_entry" }

// LevelSystem 等级体系(cefr=英语能力 / grade=教育阶段 / exam=考试体系)
type LevelSystem struct {
	models.Model
	Code string `json:"code" gorm:"size:32;uniqueIndex;comment:体系编码"`
	Name string `json:"name" gorm:"size:64;comment:体系名称"`
	Sort int    `json:"sort" gorm:"comment:排序"`
	models.ModelTime
}

func (LevelSystem) TableName() string { return "wiki_level_system" }

// Level 体系内的档位,统一映射到 CEFR rank 区间。
// CEFR 体系自身的档位 CefrMin==CefrMax==自己的 rank。
type Level struct {
	models.Model
	SystemId int    `json:"systemId" gorm:"index;comment:所属体系"`
	Code     string `json:"code" gorm:"size:32;comment:档位编码(G1/CET4/B1…)"`
	Name     string `json:"name" gorm:"size:64;comment:档位名称"`
	CefrMin  int    `json:"cefrMin" gorm:"comment:映射CEFR下限rank"`
	CefrMax  int    `json:"cefrMax" gorm:"comment:映射CEFR上限rank"`
	Sort     int    `json:"sort" gorm:"comment:排序"`
	models.ModelTime
}

func (Level) TableName() string { return "wiki_level" }
