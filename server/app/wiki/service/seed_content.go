package service

import (
	"gorm.io/gorm"

	"go-admin/app/wiki/models"
)

// seedContent 幂等注入初始分类树 + 「英语时态」示例条目(竖切验证用)。
// 注意:仓库自带的 go-admin-db.db 才是内容主库(16 时态 + 从句三类等全量内容,经管理后台/API 维护),
// 本 seed 只是空库兜底的最小样例,不与主库逐条同步。
func seedContent(db *gorm.DB) error {
	var cnt int64
	if err := db.Model(&models.WikiCategory{}).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		mk := func(parent int, name, code, icon string, sort int) (int, error) {
			c := models.WikiCategory{ParentId: parent, Name: name, Code: code, Icon: icon, Sort: sort}
			err := tx.Create(&c).Error
			return c.Id, err
		}

		// 一级:英语时态 / 英语从句 / 英语音标 / 词汇;各时态是时态下的二级分类
		tense, err := mk(0, "英语时态", "tense", "icon-clock-circle", 1)
		if err != nil {
			return err
		}
		intro := models.WikiEntry{
			CategoryId: tense,
			Title:      "英语时态总览",
			Summary:    "时态 = 时间 × 状态,一张图看懂 16 种时态",
			ContentMd:  tenseIntroMd,
			CefrMin:    1,
			CefrMax:    7,
			Status:     "published",
			Sort:       0,
		}
		if err := tx.Create(&intro).Error; err != nil {
			return err
		}
		if _, err = mk(0, "英语从句", "clause", "icon-branch", 2); err != nil {
			return err
		}
		if _, err = mk(0, "句子成分", "sentence-elements", "", 3); err != nil {
			return err
		}
		if _, err = mk(0, "五种基本句型", "basic-sentence-patterns", "", 4); err != nil {
			return err
		}
		if _, err = mk(0, "名词、冠词、介词", "noun-article-preposition", "", 5); err != nil {
			return err
		}
		if _, err = mk(0, "英语音标", "phonetics", "icon-sound", 6); err != nil {
			return err
		}
		if _, err = mk(0, "词汇", "vocabulary", "icon-bookmark", 7); err != nil {
			return err
		}

		for i, e := range tenseEntries {
			cat, err := mk(tense, e.catName, e.catCode, "", i+1)
			if err != nil {
				return err
			}
			row := models.WikiEntry{
				CategoryId: cat,
				Title:      e.title,
				Summary:    e.summary,
				ContentMd:  e.md,
				CefrMin:    e.min,
				CefrMax:    e.max,
				Status:     "published",
				Sort:       i + 1,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type entrySeed struct {
	catName  string // 二级分类名(菜单里显示)
	catCode  string // 二级分类编码
	title    string
	summary  string
	min, max int
	md       string
}

var tenseEntries = []entrySeed{
	{
		"一般现在时", "tense-simple-present",
		"一般现在时 (Simple Present)", "表示习惯、事实与规律,英语中最基础的时态", 2, 2, `
## 结构

- 肯定:主语 + 动词原形(第三人称单数 + **-s/-es**)
- 否定:主语 + **don't / doesn't** + 动词原形
- 疑问:**Do / Does** + 主语 + 动词原形?

## 什么时候用

1. **习惯性动作**:I get up at seven every day. 我每天七点起床。
2. **客观事实/真理**:The sun rises in the east. 太阳从东方升起。
3. **状态**:She likes music. 她喜欢音乐。

## 常见错误

- ❌ He like apples. → ✅ He like**s** apples.(三单要加 -s)
- ❌ Does she likes it? → ✅ Does she **like** it?(疑问句动词还原)

## 标志词

every day / usually / often / always / sometimes / never
`,
	},
	{
		"现在进行时", "tense-present-continuous",
		"现在进行时 (Present Continuous)", "表示此刻正在进行的动作:am/is/are + doing", 2, 3, `
## 结构

主语 + **am / is / are** + 动词 **-ing**

## 什么时候用

1. **此刻正在发生**:I am reading a book now. 我现在正在看书。
2. **现阶段在做的事**:She is learning English this year. 她今年在学英语。
3. **近期确定的计划**(口语常见):We are leaving tomorrow. 我们明天走。

## -ing 拼写规则

| 情况 | 规则 | 例子 |
|---|---|---|
| 一般 | 直接 +ing | play → playing |
| 以不发音 e 结尾 | 去 e +ing | make → making |
| 重读闭音节 | 双写末尾字母 | run → running |

## 常见错误

- ❌ I am study English. → ✅ I am study**ing** English.
- 状态动词一般不用进行时:❌ I am knowing him. → ✅ I **know** him.

## 标志词

now / right now / at the moment / look! / listen!
`,
	},
	{
		"一般过去时", "tense-simple-past",
		"一般过去时 (Simple Past)", "表示过去发生的动作或状态:动词过去式", 3, 3, `
## 结构

- 肯定:主语 + 动词**过去式**(规则动词 +ed;不规则动词需记忆)
- 否定:主语 + **didn't** + 动词原形
- 疑问:**Did** + 主语 + 动词原形?

## 什么时候用

1. **过去某时发生的事**:I visited my grandma yesterday. 我昨天看望了奶奶。
2. **过去的习惯**:He often played football when he was young.

## 常见不规则动词

go → went, see → saw, eat → ate, have → had, do → did, come → came

## 常见错误

- ❌ I didn't went there. → ✅ I didn't **go** there.(didn't 后用原形)
- ❌ Did you saw it? → ✅ Did you **see** it?

## 标志词

yesterday / last week / two days ago / in 2020 / just now
`,
	},
	{
		"一般将来时", "tense-simple-future",
		"一般将来时 (Simple Future)", "表示将要发生的事:will / be going to + 动词原形", 3, 3, `
## 结构

- **will** + 动词原形:临时决定、预测
- **be going to** + 动词原形:事先打算、有迹象的预测

## 两者区别

| | will | be going to |
|---|---|---|
| 决定时机 | 说话时临时决定 | 早有打算 |
| 例子 | I'll help you.(马上决定帮你) | I'm going to study abroad.(早就计划) |
| 预测 | I think it will rain. | Look at the clouds! It's going to rain.(有迹象) |

## 常见错误

- ❌ I will to go. → ✅ I will **go**.(will 后不加 to)
- ❌ He wills go. → ✅ He **will** go.(will 无三单变化)

## 标志词

tomorrow / next week / in the future / soon
`,
	},
	{
		"现在完成时", "tense-present-perfect",
		"现在完成时 (Present Perfect)", "过去的动作对现在有影响:have/has + 过去分词", 4, 4, `
## 结构

主语 + **have / has** + 动词**过去分词**(p.p.)

## 什么时候用

1. **过去发生,影响现在**:I have lost my key.(钥匙丢了,现在进不了门)
2. **到现在为止的经历**:She has been to Beijing twice. 她去过北京两次。
3. **持续到现在**:We have lived here for ten years. 我们在这住了十年。

## 与一般过去时的区别(高频考点)

- I **lost** my key yesterday.(只陈述过去,和现在无关,带过去时间)
- I **have lost** my key.(强调现在没钥匙;不能与 yesterday 等过去时间连用)

## have been to vs have gone to

- has **been to**:去过(已回来) — He has been to Japan.
- has **gone to**:去了(还没回) — He has gone to Japan.

## 标志词

already / yet / just / ever / never / since / for / so far
`,
	},
	{
		"过去进行时", "tense-past-continuous",
		"过去进行时 (Past Continuous)", "过去某一时刻正在进行:was/were + doing", 4, 4, `
## 结构

主语 + **was / were** + 动词 **-ing**

## 什么时候用

1. **过去某时刻正在做**:I was cooking at 6 p.m. yesterday.
2. **长动作被短动作打断**(与 when/while 连用,重点):
   - I **was taking** a shower **when** the phone rang. 我洗澡时电话响了。
   - **While** we **were talking**, the teacher came in.

## when 与 while

- when + 短动作(一般过去时);while + 长动作(过去进行时)
- While I was reading, my brother was playing.(两个长动作并行)

## 常见错误

- ❌ When I was hearing the news… → ✅ When I **heard** the news…(瞬间动词不用进行)
`,
	},
	{
		"过去完成时", "tense-past-perfect",
		"过去完成时 (Past Perfect)", "过去的过去:had + 过去分词", 5, 5, `
## 结构

主语 + **had** + 动词过去分词(所有人称一致)

## 什么时候用

表示**过去某个时间点之前**已经完成的动作——"过去的过去":

- When I arrived, the train **had left**. 我到时火车已经开走了。(开走在到达之前)
- She said she **had seen** the film.(宾语从句中,"看过"早于"说")

## 判断技巧

句中有两个过去动作时,**先发生的用过去完成时**,后发生的用一般过去时。
只有一个过去动作时,通常用一般过去时即可,不要滥用 had done。

## 常见错误

- ❌ After he had finished his homework, he had gone to bed.
- ✅ After he had finished his homework, he **went** to bed.(后一个动作用一般过去时)

## 标志词

by the time… / before / after / when + 过去时间点 / already(过去语境)
`,
	},
	{
		"将来完成时", "tense-future-perfect",
		"将来完成时 (Future Perfect)", "到将来某时为止将已完成:will have + 过去分词", 6, 6, `
## 结构

主语 + **will have** + 动词过去分词

## 什么时候用

表示**到将来某个时间点为止**,动作将已经完成:

- By next month, I **will have worked** here for five years. 到下个月,我在这就工作满五年了。
- By the time you arrive, we **will have finished** dinner. 等你到的时候,我们就已经吃完晚饭了。

## 使用要点

- 几乎总与 **by + 将来时间**(by 2030 / by the time…)连用。
- by the time 引导的从句用一般现在时表将来:By the time he **comes**(不是 will come), …

## 学术/正式写作中的用法

常见于展望、预测类表达:
The project will have been completed before the deadline.(被动式:will have been + p.p.)
`,
	},
}

// tenseIntroMd 「英语时态」一级分类的总览介绍页(含 16 时态矩阵 SVG 图)
var tenseIntroMd = `
## 时态是什么

英语的"时态"其实是两个维度的组合:**时间**(过去 / 现在 / 将来 / 过去将来)× **状态**(一般 / 进行 / 完成 / 完成进行)。4 × 4 组合,就是常说的 **16 种时态**——不用死记硬背,先看懂这张图:

<svg viewBox="0 0 840 566" xmlns="http://www.w3.org/2000/svg" style="max-width:100%;height:auto;font-family:inherit">
<defs><marker id="arr" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto"><path d="M0,0 L6,3 L0,6 Z" fill="var(--color-text-3, #86909c)"/></marker></defs>
<line x1="140" y1="24" x2="648" y2="24" stroke="var(--color-text-3, #86909c)" stroke-width="1.5" marker-end="url(#arr)"/>
<text x="394" y="16" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">时间</text>
<text x="222" y="52" font-size="15" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">过去</text>
<text x="397" y="52" font-size="15" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">现在</text>
<text x="572" y="52" font-size="15" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">将来</text>
<text x="747" y="52" font-size="15" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">过去将来</text>
<text x="70" y="115" font-size="15" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">一般</text>
<text x="70" y="135" font-size="11" fill="var(--color-text-3, #86909c)" text-anchor="middle">客观陈述</text>
<text x="70" y="229" font-size="15" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">进行</text>
<text x="70" y="249" font-size="11" fill="var(--color-text-3, #86909c)" text-anchor="middle">正在发生</text>
<text x="70" y="343" font-size="15" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">完成</text>
<text x="70" y="363" font-size="11" fill="var(--color-text-3, #86909c)" text-anchor="middle">已经完成</text>
<text x="70" y="457" font-size="15" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">完成进行</text>
<text x="70" y="477" font-size="11" fill="var(--color-text-3, #86909c)" text-anchor="middle">一直持续</text>
<rect x="140" y="66" width="155" height="94" rx="10" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="217" y="106" font-size="14" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">一般过去时</text>
<text x="217" y="132" font-size="12" fill="rgb(var(--arcoblue-6, 22,93,255))" text-anchor="middle">did</text>
<rect x="315" y="66" width="155" height="94" rx="10" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="392" y="106" font-size="14" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">一般现在时</text>
<text x="392" y="132" font-size="12" fill="rgb(var(--arcoblue-6, 22,93,255))" text-anchor="middle">do / does</text>
<rect x="490" y="66" width="155" height="94" rx="10" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="567" y="106" font-size="14" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">一般将来时</text>
<text x="567" y="132" font-size="12" fill="rgb(var(--arcoblue-6, 22,93,255))" text-anchor="middle">will do</text>
<rect x="665" y="66" width="155" height="94" rx="10" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="742" y="106" font-size="14" fill="var(--color-text-3, #86909c)" text-anchor="middle">过去将来时</text>
<text x="742" y="132" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">would do</text>
<rect x="140" y="180" width="155" height="94" rx="10" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="217" y="220" font-size="14" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">过去进行时</text>
<text x="217" y="246" font-size="12" fill="rgb(var(--arcoblue-6, 22,93,255))" text-anchor="middle">was/were doing</text>
<rect x="315" y="180" width="155" height="94" rx="10" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="392" y="220" font-size="14" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">现在进行时</text>
<text x="392" y="246" font-size="12" fill="rgb(var(--arcoblue-6, 22,93,255))" text-anchor="middle">am/is/are doing</text>
<rect x="490" y="180" width="155" height="94" rx="10" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="567" y="220" font-size="14" fill="var(--color-text-3, #86909c)" text-anchor="middle">将来进行时</text>
<text x="567" y="246" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">will be doing</text>
<rect x="665" y="180" width="155" height="94" rx="10" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="742" y="220" font-size="14" fill="var(--color-text-3, #86909c)" text-anchor="middle">过去将来进行时</text>
<text x="742" y="246" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">would be doing</text>
<rect x="140" y="294" width="155" height="94" rx="10" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="217" y="334" font-size="14" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">过去完成时</text>
<text x="217" y="360" font-size="12" fill="rgb(var(--arcoblue-6, 22,93,255))" text-anchor="middle">had done</text>
<rect x="315" y="294" width="155" height="94" rx="10" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="392" y="334" font-size="14" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">现在完成时</text>
<text x="392" y="360" font-size="12" fill="rgb(var(--arcoblue-6, 22,93,255))" text-anchor="middle">have/has done</text>
<rect x="490" y="294" width="155" height="94" rx="10" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="567" y="334" font-size="14" font-weight="600" fill="var(--color-text-1, #1d2129)" text-anchor="middle">将来完成时</text>
<text x="567" y="360" font-size="12" fill="rgb(var(--arcoblue-6, 22,93,255))" text-anchor="middle">will have done</text>
<rect x="665" y="294" width="155" height="94" rx="10" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="742" y="334" font-size="14" fill="var(--color-text-3, #86909c)" text-anchor="middle">过去将来完成时</text>
<text x="742" y="360" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">would have done</text>
<rect x="140" y="408" width="155" height="94" rx="10" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="217" y="448" font-size="14" fill="var(--color-text-3, #86909c)" text-anchor="middle">过去完成进行时</text>
<text x="217" y="474" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">had been doing</text>
<rect x="315" y="408" width="155" height="94" rx="10" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="392" y="448" font-size="14" fill="var(--color-text-3, #86909c)" text-anchor="middle">现在完成进行时</text>
<text x="392" y="474" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">have been doing</text>
<rect x="490" y="408" width="155" height="94" rx="10" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="567" y="448" font-size="14" fill="var(--color-text-3, #86909c)" text-anchor="middle">将来完成进行时</text>
<text x="567" y="474" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">will have been doing</text>
<rect x="665" y="408" width="155" height="94" rx="10" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="742" y="448" font-size="14" fill="var(--color-text-3, #86909c)" text-anchor="middle">过去将来完成进行时</text>
<text x="742" y="474" font-size="12" fill="var(--color-text-3, #86909c)" text-anchor="middle">would have been doing</text>
<rect x="140" y="534" width="26" height="16" rx="4" fill="rgba(22, 93, 255, 0.08)" stroke="rgb(var(--arcoblue-6, 22,93,255))" stroke-width="1.5"/>
<text x="174" y="547" font-size="12" fill="var(--color-text-2, #4e5969)">已收录详细讲解(点左侧菜单阅读)</text>
<rect x="420" y="534" width="26" height="16" rx="4" fill="var(--color-fill-1, #f7f8fa)" stroke="var(--color-border-2, #e5e6eb)" stroke-width="1.2" stroke-dasharray="5,4"/>
<text x="454" y="547" font-size="12" fill="var(--color-text-2, #4e5969)">内容陆续补充中</text>
</svg>

## 一个动作,不同时空

以 write(写)为例,同一个动作放进不同的"时间 × 状态":

- I **write** emails every day.(一般现在:习惯,天天写)
- I **am writing** an email now.(现在进行:此刻正在写)
- I **wrote** an email yesterday.(一般过去:昨天写了)
- I **have written** the email.(现在完成:已写完,现在能发)
- By 10 p.m. I **will have written** it.(将来完成:到十点前会写完)

时间词变了,动词形态跟着变——这就是时态的全部秘密。

## 建议学习路线(按水平)

| 阶段 | 建议掌握 |
|---|---|
| A1 入门 | 一般现在时、现在进行时 |
| A2 初级 | 一般过去时、一般将来时 |
| B1 中级 | 现在完成时、过去进行时 |
| B2 中高级 | 过去完成时 |
| C1 高级 | 将来完成时、完成进行类 |

点左侧菜单里的具体时态,直接阅读详细讲解。
`
