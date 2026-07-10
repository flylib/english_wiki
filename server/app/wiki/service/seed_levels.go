package service

import (
	"gorm.io/gorm"

	"go-admin/app/wiki/models"
)

// CEFR rank: 1=Pre-A1 2=A1 3=A2 4=B1 5=B2 6=C1 7=C2

type levelSeed struct {
	code     string
	name     string
	min, max int
}

// seedLevels 幂等注入三套等级体系;所有档位统一映射到 CEFR rank 区间。
func seedLevels(db *gorm.DB) error {
	var cnt int64
	if err := db.Model(&models.LevelSystem{}).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}

	systems := []struct {
		code, name string
		levels     []levelSeed
	}{
		{"cefr", "英语能力(CEFR)", []levelSeed{
			{"PRE_A1", "Pre-A1 零基础", 1, 1},
			{"A1", "A1 入门", 2, 2},
			{"A2", "A2 初级", 3, 3},
			{"B1", "B1 中级", 4, 4},
			{"B2", "B2 中高级", 5, 5},
			{"C1", "C1 高级", 6, 6},
			{"C2", "C2 精通", 7, 7},
		}},
		{"grade", "教育阶段", []levelSeed{
			{"PRESCHOOL", "学前", 1, 1},
			{"G1", "小学一年级", 1, 1},
			{"G2", "小学二年级", 1, 1},
			{"G3", "小学三年级", 1, 2},
			{"G4", "小学四年级", 1, 2},
			{"G5", "小学五年级", 2, 2},
			{"G6", "小学六年级", 2, 3},
			{"G7", "初一", 2, 3},
			{"G8", "初二", 3, 3},
			{"G9", "初三", 3, 4},
			{"G10", "高一", 3, 4},
			{"G11", "高二", 4, 4},
			{"G12", "高三", 4, 5},
			{"COLLEGE", "大学", 4, 6},
			{"ADULT", "成人", 1, 7},
		}},
		{"exam", "考试体系", []levelSeed{
			{"YLE_STARTERS", "剑桥少儿 Starters", 1, 1},
			{"YLE_MOVERS", "剑桥少儿 Movers", 1, 2},
			{"YLE_FLYERS", "剑桥少儿 Flyers", 2, 2},
			{"KET", "剑桥 KET", 3, 3},
			{"PET", "剑桥 PET", 4, 4},
			{"FCE", "剑桥 FCE", 5, 5},
			{"CET4", "英语四级 CET-4", 4, 5},
			{"CET6", "英语六级 CET-6", 5, 5},
			{"KAOYAN", "考研英语", 5, 6},
			{"TEM4", "专业四级 TEM-4", 5, 5},
			{"TEM8", "专业八级 TEM-8", 6, 7},
			{"IELTS", "雅思 IELTS", 4, 7},
			{"TOEFL", "托福 TOEFL", 5, 7},
		}},
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for si, s := range systems {
			sys := models.LevelSystem{Code: s.code, Name: s.name, Sort: si + 1}
			if err := tx.Create(&sys).Error; err != nil {
				return err
			}
			for li, l := range s.levels {
				lv := models.Level{
					SystemId: sys.Id,
					Code:     l.code,
					Name:     l.name,
					CefrMin:  l.min,
					CefrMax:  l.max,
					Sort:     li + 1,
				}
				if err := tx.Create(&lv).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
