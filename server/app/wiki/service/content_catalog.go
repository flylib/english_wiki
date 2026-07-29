package service

import (
	"gorm.io/gorm"

	"go-admin/app/wiki/models"
)

// content_catalog.go 是百科的“课程目录”。它把英语拆成可学习、可检索的模块，
// 每个模块至少有一篇总览/方法条目。正文保持 Markdown，后续可以继续在后台细化。

type categorySeed struct {
	parent, name, code, icon string
	sort                     int
}

type catalogEntry struct {
	category, title, summary, md string
	min, max, sort               int
}

var categorySeeds = []categorySeed{
	{"", "英语全景", "english-system", "icon-apps", -1},
	{"", "英语入门", "foundations", "icon-home", 0},
	{"", "词类", "parts-of-speech", "", 8},
	{"", "动词系统", "verb-system", "", 9},
	{"", "语态与非谓语", "voice-nonfinite", "", 10},
	{"", "语法专题", "grammar-topics", "", 11},
	{"", "功能英语", "functional-english", "", 12},
	{"", "英语听力", "listening", "", 13},
	{"", "英语口语", "speaking", "", 14},
	{"", "英语阅读", "reading", "", 15},
	{"", "英语写作", "writing", "", 16},
	{"", "学习方法", "learning-path", "", 17},
	{"", "考试与能力路径", "exam-roadmap", "", 18},
	{"tense", "一般过去时", "tense-simple-past", "", 1}, // 已存在时不会重复创建
	{"tense", "将来进行时", "tense-future-continuous", "", 9},
	{"tense", "过去将来时", "tense-past-future", "", 10},
	{"tense", "过去将来进行时", "tense-past-future-continuous", "", 11},
	{"tense", "过去将来完成时", "tense-past-future-perfect", "", 12},
	{"tense", "过去将来完成进行时", "tense-past-future-perfect-continuous", "", 13},
	{"tense", "现在完成进行时", "tense-present-perfect-continuous", "", 14},
	{"tense", "过去完成进行时", "tense-past-perfect-continuous", "", 15},
	{"tense", "将来完成进行时", "tense-future-perfect-continuous", "", 16},
	{"clause", "名词性从句", "clause-noun", "", 1},
	{"clause", "定语从句", "clause-relative", "", 2},
	{"clause", "状语从句", "clause-adverbial", "", 3},
	{"sentence-elements", "句子成分进阶", "sentence-elements-advanced", "", 1},
	{"foundations", "字母、拼读与课堂英语", "foundation-start", "", 1},
	{"foundations", "第一批高频词", "foundation-high-frequency-words", "", 2},
	{"phonetics", "元音：嘴巴的形状", "phonetics-vowels", "", 1},
	{"phonetics", "辅音：气流和摩擦", "phonetics-consonants", "", 2},
	{"phonetics", "重音、节奏与连读", "phonetics-stress-rhythm", "", 3},
	{"vocabulary", "主题词汇：从身边开始", "vocabulary-themes", "", 1},
	{"vocabulary", "构词法：词根、前缀与后缀", "vocabulary-word-formation", "", 2},
	{"vocabulary", "词块、搭配与短语动词", "vocabulary-collocations", "", 3},
	{"parts-of-speech", "名词与复数", "parts-noun", "", 1},
	{"parts-of-speech", "代词与限定词", "parts-pronoun-determiner", "", 2},
	{"parts-of-speech", "形容词与副词", "parts-adjective-adverb", "", 3},
	{"parts-of-speech", "介词、连词与感叹词", "parts-function-words", "", 4},
	{"verb-system", "动词形式与助动词", "verbs-forms-auxiliaries", "", 1},
	{"verb-system", "情态动词：能力、义务和推测", "verbs-modals", "", 2},
	{"voice-nonfinite", "被动语态", "voice-passive", "", 1},
	{"voice-nonfinite", "不定式", "nonfinite-infinitive", "", 2},
	{"voice-nonfinite", "动名词与分词", "nonfinite-gerund-participle", "", 3},
	{"grammar-topics", "疑问句与否定句", "grammar-questions-negation", "", 1},
	{"grammar-topics", "比较级与最高级", "grammar-comparison", "", 2},
	{"grammar-topics", "条件句", "grammar-conditionals", "", 3},
	{"grammar-topics", "虚拟语气", "grammar-subjunctive", "", 4},
	{"grammar-topics", "直接引语与间接引语", "grammar-reported-speech", "", 5},
	{"grammar-topics", "倒装、强调与省略", "grammar-inversion-emphasis", "", 6},
	{"grammar-topics", "标点与大小写", "grammar-punctuation", "", 7},
	{"functional-english", "问候、介绍与礼貌", "function-greetings", "", 1},
	{"functional-english", "请求、建议与邀请", "function-requests", "", 2},
	{"functional-english", "表达观点、同意与不同意", "function-opinions", "", 3},
	{"listening", "听力入门：先抓大意", "listening-foundation", "", 1},
	{"listening", "听力进阶：预测、定位和复盘", "listening-strategies", "", 2},
	{"speaking", "发音、跟读与流利度", "speaking-fluency", "", 1},
	{"speaking", "日常对话与交际策略", "speaking-conversation", "", 2},
	{"reading", "阅读入门：从词到段落", "reading-foundation", "", 1},
	{"reading", "阅读进阶：结构、推断与批判性阅读", "reading-strategies", "", 2},
	{"writing", "句子、段落与连接", "writing-foundation", "", 1},
	{"writing", "说明文、议论文与正式写作", "writing-academic", "", 2},
	{"learning-path", "如何制定英语学习计划", "learning-plan", "", 1},
	{"learning-path", "间隔复习、输入与输出", "learning-method", "", 2},
	{"exam-roadmap", "考试与 CEFR 路线图", "exam-cefr-roadmap", "", 1},
}

var catalogEntries = []catalogEntry{
	{"english-system", "英语知识全景图：从声音到表达", "一张地图看懂英语的六个层次、四项技能和不同阶段的学习重点。", englishSystemMd, 1, 7, 0},
	{"foundations", "英语入门总览：先会用，再学规则", "从声音、词、短句和日常场景建立第一座英语小屋。", foundationsMd, 1, 2, 0},
	{"foundation-start", "字母、拼读与课堂英语", "认识字母名、字母音和最常用的课堂指令。", foundationStartMd, 1, 2, 1},
	{"foundation-high-frequency-words", "第一批高频词：I、you、be 和 have", "用最常见的小词组成自己的第一批句子。", foundationWordsMd, 1, 2, 2},
	{"tense-future-continuous", "将来进行时 (Future Continuous)", "表示将来某一时刻正在进行的动作。", tenseFutureContinuousMd, 5, 6, 9},
	{"tense-past-future", "过去将来时 (Future in the Past)", "从过去的视角描述将来：would do。", tensePastFutureMd, 5, 6, 10},
	{"tense-past-future-continuous", "过去将来进行时", "表示从过去看将来某时正在进行。", tensePastFutureContinuousMd, 6, 7, 11},
	{"tense-past-future-perfect", "过去将来完成时", "表示从过去看将来某时已经完成。", tensePastFuturePerfectMd, 6, 7, 12},
	{"tense-past-future-perfect-continuous", "过去将来完成进行时", "表示从过去看将来某时已持续一段时间。", tensePastFuturePerfectContinuousMd, 7, 7, 13},
	{"tense-present-perfect-continuous", "现在完成进行时", "have/has been doing：动作从过去持续到现在。", tensePresentPerfectContinuousMd, 5, 6, 14},
	{"tense-past-perfect-continuous", "过去完成进行时", "had been doing：过去某时之前一直在做。", tensePastPerfectContinuousMd, 6, 7, 15},
	{"tense-future-perfect-continuous", "将来完成进行时", "到未来某时将已经持续一段时间。", tenseFuturePerfectContinuousMd, 6, 7, 16},
	{"clause-noun", "名词性从句：把一句话当成名词", "that、whether、what、who 引导的主语、宾语和表语从句。", nounClauseMd, 4, 7, 1},
	{"clause-relative", "定语从句：给名词加说明", "who、which、that 和 whose 把两个信息接成一句话。", relativeClauseMd, 4, 6, 2},
	{"clause-adverbial", "状语从句：时间、原因、条件和让步", "because、when、if、although 等连接逻辑关系。", adverbialClauseMd, 3, 7, 3},
	{"sentence-elements-advanced", "长句拆解：主干、修饰和逻辑", "面对长难句先找谓语和主干，再处理插入与从句。", sentenceAdvancedMd, 4, 7, 1},
	{"phonetics-vowels", "元音：嘴巴的形状", "用口型、舌位和长短音区分容易混淆的元音。", vowelsMd, 1, 4, 1},
	{"phonetics-consonants", "辅音：气流和摩擦", "掌握清浊音、摩擦音和汉语里不常见的声音。", consonantsMd, 1, 4, 2},
	{"phonetics-stress-rhythm", "重音、节奏与连读", "英语不是每个词都同样重，节奏会帮助你听懂和说得自然。", stressRhythmMd, 2, 5, 3},
	{"vocabulary-themes", "主题词汇：从身边开始", "家庭、学校、食物、天气、身体、城市、旅行等高频主题。", themesMd, 1, 4, 1},
	{"vocabulary-word-formation", "构词法：词根、前缀与后缀", "从 happy、unhappy、happiness 看懂单词家族。", wordFormationMd, 4, 7, 2},
	{"vocabulary-collocations", "词块、搭配与短语动词", "记住词语常和谁一起出现，表达会更准确自然。", collocationsMd, 3, 7, 3},
	{"parts-noun", "名词与复数：数得清才说得准", "可数/不可数、规则复数、不规则复数和所有格。", nounsMd, 1, 4, 1},
	{"parts-pronoun-determiner", "代词与限定词", "I/me、this/that、some/any、each/every 的用法地图。", pronounsMd, 2, 6, 2},
	{"parts-adjective-adverb", "形容词与副词", "描述事物、动作和程度，理解 -ly 与位置变化。", adjectivesMd, 2, 5, 3},
	{"parts-function-words", "介词、连词与感叹词", "用小词表达空间、时间、原因、转折和情绪。", functionWordsMd, 2, 6, 4},
	{"verbs-forms-auxiliaries", "动词形式与助动词", "原形、三单、过去式、分词，以及 do/be/have 如何帮忙。", verbFormsMd, 2, 6, 1},
	{"verbs-modals", "情态动词：能力、义务和推测", "can、must、should、may、might 等表达态度和可能性。", modalsMd, 3, 7, 2},
	{"voice-passive", "被动语态：关注事情而不是谁做的", "be + 过去分词，以及时态、情态动词和 get-passive。", passiveMd, 4, 7, 1},
	{"nonfinite-infinitive", "不定式：to do 的多种角色", "作目的、主语、宾语、宾补和后置修饰。", infinitiveMd, 4, 7, 2},
	{"nonfinite-gerund-participle", "动名词与分词", "doing 可以是名词，也可以是进行/修饰的一部分。", gerundParticipleMd, 4, 7, 3},
	{"grammar-questions-negation", "疑问句与否定句", "一般疑问、特殊疑问、反意疑问和 not 的位置。", questionsMd, 1, 5, 1},
	{"grammar-comparison", "比较级与最高级", "比较两者、多者以及 as...as、less、the more...the more。", comparisonMd, 2, 6, 2},
	{"grammar-conditionals", "条件句：如果……就……", "零、一、二、三条件句和混合条件句的时间与真实程度。", conditionalsMd, 4, 7, 3},
	{"grammar-subjunctive", "虚拟语气：与事实相反的想象", "wish、if、as if 和正式建议中的虚拟表达。", subjunctiveMd, 5, 7, 4},
	{"grammar-reported-speech", "直接引语与间接引语", "转述别人说的话时，处理时态、人称、指示词和问句。", reportedSpeechMd, 4, 7, 5},
	{"grammar-inversion-emphasis", "倒装、强调与省略", "让句子更正式、更有重点：only、never、it-cleft 等结构。", inversionMd, 6, 7, 6},
	{"grammar-punctuation", "标点与大小写", "逗号、句号、冒号、分号、引号和英文标题大小写。", punctuationMd, 2, 7, 7},
	{"function-greetings", "问候、介绍与礼貌", "从 Hello 到 Nice to meet you，在不同关系中选择合适表达。", greetingsMd, 1, 3, 1},
	{"function-requests", "请求、建议与邀请", "Can you...?、Would you mind...?、Why don't we...? 的语气梯度。", requestsMd, 2, 5, 2},
	{"function-opinions", "表达观点、同意与不同意", "I think、In my view、I see your point but... 让交流更有层次。", opinionsMd, 3, 7, 3},
	{"listening-foundation", "听力入门：先抓大意", "从关键词、语气、场景和图片线索建立听力信心。", listeningFoundationMd, 1, 3, 1},
	{"listening-strategies", "听力进阶：预测、定位和复盘", "听前预测、听中记关键词、听后对照文本，形成闭环。", listeningStrategiesMd, 3, 7, 2},
	{"speaking-fluency", "发音、跟读与流利度", "从单词准确到意群流畅，练习重音、停顿和自我修正。", speakingFluencyMd, 1, 5, 1},
	{"speaking-conversation", "日常对话与交际策略", "开启话题、追问、澄清、接话和结束对话。", conversationMd, 2, 7, 2},
	{"reading-foundation", "阅读入门：从词到段落", "先看标题和图片，再找主旨、细节和上下文线索。", readingFoundationMd, 1, 4, 1},
	{"reading-strategies", "阅读进阶：结构、推断与批判性阅读", "区分事实和观点，理解指代、态度、论证和作者意图。", readingStrategiesMd, 4, 7, 2},
	{"writing-foundation", "句子、段落与连接", "从完整句开始，写出有主题句、支持句和结尾的段落。", writingFoundationMd, 2, 5, 1},
	{"writing-academic", "说明文、议论文与正式写作", "规划观点、组织论据、使用连接词并完成修改清单。", writingAcademicMd, 4, 7, 2},
	{"learning-plan", "如何制定英语学习计划", "按目标、水平、时间和反馈设计可坚持的学习循环。", learningPlanMd, 1, 7, 1},
	{"learning-method", "间隔复习、输入与输出", "把听读输入变成说写输出，使用小步、频繁、可测量的练习。", learningMethodMd, 1, 7, 2},
	{"exam-cefr-roadmap", "考试与 CEFR 路线图", "把 Pre-A1、A1、A2、B1、B2、C1、C2 与学习目标对应起来。", examRoadmapMd, 1, 7, 1},
}

func ensureContentCatalog(db *gorm.DB) error {
	ids := map[string]int{}
	var existingCategories []models.WikiCategory
	if err := db.Select("id, code").Find(&existingCategories).Error; err != nil {
		return err
	}
	for _, category := range existingCategories {
		ids[category.Code] = category.Id
	}
	for _, seed := range categorySeeds {
		var category models.WikiCategory
		err := db.Where("code = ?", seed.code).First(&category).Error
		if err == nil {
			ids[seed.code] = category.Id
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		parentID := 0
		if seed.parent != "" {
			parentID = ids[seed.parent]
			if parentID == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		category = models.WikiCategory{ParentId: parentID, Name: seed.name, Code: seed.code, Icon: seed.icon, Sort: seed.sort}
		if err := db.Create(&category).Error; err != nil {
			return err
		}
		ids[seed.code] = category.Id
	}

	for _, seed := range catalogEntries {
		categoryID := ids[seed.category]
		if categoryID == 0 {
			return gorm.ErrRecordNotFound
		}
		var existing models.WikiEntry
		err := db.Where("category_id = ? AND title = ?", categoryID, seed.title).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		entry := models.WikiEntry{CategoryId: categoryID, Title: seed.title, Summary: seed.summary, ContentMd: seed.md, CefrMin: seed.min, CefrMax: seed.max, Status: "published", Sort: seed.sort}
		if err := db.Create(&entry).Error; err != nil {
			return err
		}
	}
	return nil
}

var englishSystemMd = `
## 英语不是一堆零散规则

把英语看成一座城市：

1. **声音层**：音标、拼读、重音、节奏，让你听得出、说得清。
2. **词汇层**：词义、词性、构词、搭配和短语，让你有材料可用。
3. **句法层**：句子成分、基本句型、时态、语态、从句和非谓语，让表达有骨架。
4. **语篇层**：段落、衔接、主旨、论证和文体，让你读写长内容。
5. **交际层**：问候、请求、澄清、礼貌、观点和文化，让语言适合真实场景。
6. **技能层**：听、说、读、写四项技能相互喂养，不能只学其中一项。

## 按阶段学习什么

| 阶段 | 核心任务 | 重点内容 |
|---|---|---|
| Pre-A1–A1 | 听懂并说出熟悉内容 | 字母音、主题词、be、一般现在时、短对话 |
| A2 | 处理日常生活 | 过去/将来、比较、疑问、常见从句、短文 |
| B1 | 连贯表达经历和观点 | 完成时、被动、条件句、段落写作、复述 |
| B2 | 理解复杂信息并论证 | 长句、非谓语、从句嵌套、正式语体、议论文 |
| C1–C2 | 精准灵活地表达 | 语域、隐含意义、倒装强调、学术/专业语篇 |

## 推荐路线

**声音 → 高频词 → 句子骨架 → 核心时态 → 词法专题 → 从句/非谓语 → 听说读写 → 语域与专业表达**。

任何阶段都要保留真实输入和输出：儿童多用图片、歌曲、故事和游戏；中学生增加规则归纳和短文；成人/进阶学习者增加新闻、书籍、工作场景和观点写作。
`

var foundationsMd = `
## 英语学习的四层楼

1. **声音**：听得出、读得准。
2. **词**：知道意思，也知道常和谁一起用。
3. **句子**：能说清楚谁、做什么、什么时候、在哪里。
4. **篇章**：能听懂故事、读懂文章、表达自己的想法。

小学及零基础阶段先用图片、动作、歌曲和短故事建立理解；进入中学后，再把这些句子归纳成语法规则。规则是地图，不是出发点。

## 第一周小任务

每天 10 分钟：听 3 分钟 → 读 3 分钟 → 说 2 分钟 → 写 2 句。内容只选自己能理解七成以上的材料。
`

var foundationStartMd = `
## 字母名和字母音不是一回事

字母 **B** 的名字读 /biː/，它在 **book** 里的声音是 /b/。先学字母名，再用声音拼读：**s-u-n → sun**。

## 课堂英语小工具箱

**Listen. Look. Repeat. Point to... Open your book. Work in pairs.**

听不懂时可以说：**Could you say that again? / What does ___ mean?**

不要把拼读当成猜谜：看到新词先听示范，再观察字母组合，最后自己读回去。
`

var foundationWordsMd = `
## 第一批“万能小词”

| 人 | be | have | 动作 | 连接 |
|---|---|---|---|---|
| I, you, he, she, we, they | am, is, are | have, has | like, want, go, see, get | and, but, because |

用它们搭句子：**I am happy. She has a cat. We like music because it is fun.**

先会说完整短句，再扩充颜色、数量、时间和地点：**I have a red ball at home.**
`

var tenseFutureContinuousMd = `
## 结构

**will be + doing**：At 8 tomorrow, I **will be studying**.

## 用法

表示将来某个时间点正在进行的动作，常与 **at this time tomorrow / at 8 p.m.** 连用。

它也可以礼貌地询问计划：**Will you be using the car tonight?** 这里不是单纯问“你会不会”，而是确认对方的安排。
`

var tensePastFutureMd = `
## 结构与用法

从过去看将来，常见结构是 **would + 动词原形**：

**She said she would call me.** 她说她会给我打电话。

也可用 **was/were going to + 动词原形** 表示过去的计划：**I was going to leave, but it rained.**

注意：这是“过去视角里的将来”，不是普通将来时的另一种拼写。
`

var tensePastFutureContinuousMd = `
## 结构

**would be doing**：He said that at noon he **would be flying** to London.

## 用法

表示过去某个时刻预想的“将来某时正在发生”。在叙事、新闻转述和间接引语中较常见。

先把时间线画出来：过去的观察点 → 未来的进行动作。
`

var tensePastFuturePerfectMd = `
## 结构

**would have + 过去分词**：She thought she **would have finished** by Friday.

## 用法

从过去看，预计在另一个将来时间点之前已经完成。它与现在完成时的差别不在“完成”，而在观察时间点不同。
`

var tensePastFuturePerfectContinuousMd = `
## 结构

**would have been + doing**：By June, he knew he **would have been working** there for ten years.

这个形式很少用于日常交流，但能精确表达“从过去看，到未来某时已经持续多久”。正式写作中不要为了复杂而滥用。
`

var tensePresentPerfectContinuousMd = `
## 结构

**have/has been + doing**：I **have been learning** English for two years.

## 用法

强调动作从过去开始、一直持续到现在，或刚刚停止但留下明显痕迹：**She is tired because she has been running.**

和现在完成时比较：**I have read three books** 强调数量/结果；**I have been reading for two hours** 强调过程/持续时间。
`

var tensePastPerfectContinuousMd = `
## 结构与用法

**had been + doing** 表示过去某时之前一直持续的动作：**They had been waiting for an hour when the bus arrived.**

重点在“持续多久”，后面的过去事件是参照点。若只强调先后结果，用过去完成时 **had done** 即可。
`

var tenseFuturePerfectContinuousMd = `
## 结构与用法

**will have been + doing**：By next month, I **will have been working** here for a year.

它强调“到将来某时为止持续了多久”。常与 **for + 时长**、**by + 时间点** 连用。日常口语中往往会改用更简单的表达。
`

var nounClauseMd = `
## 名词性从句是什么

一整句话可以像名词一样充当主语、宾语或表语：**I know [that she is busy].**

常见引导词：**that**（陈述内容）、**whether/if**（是否）、**what/who/where/how**（具体信息）。

## 语序

从句内部用陈述语序：**I don't know where he lives.** 不说 *where does he live*。

把复杂句拆成两步：先找到主句动词，再找它后面的“内容盒子”。
`

var relativeClauseMd = `
## 给名词贴一张说明卡

**The girl who is wearing a red hat is my sister.**

先行词是 **the girl**，**who is wearing a red hat** 说明是哪一个女孩。

| 先行词 | 关系词 | 例子 |
|---|---|---|
| 人 | who/that | a teacher who helps me |
| 物 | which/that | a book that I like |
| 所属 | whose | a student whose bag is blue |

限制性从句是辨认对象所必需的信息；非限制性从句用逗号补充信息，通常不用 that。
`

var adverbialClauseMd = `
## 四种常见逻辑

**时间** when/while/before/after；**原因** because/since；**条件** if/unless；**让步** although/even though。

**When I finish my homework, I will play.** 时间从句谈将来时，通常用一般现在时，不说 *when I will finish*。

一个好方法：先说两句短句，再选择它们之间的关系，最后用连接词合并。
`

var sentenceAdvancedMd = `
## 长句三步拆法

1. 圈出所有谓语动词，判断每个动词对应的主语。
2. 删除插入语、介词短语和定语从句，先读主干。
3. 再按 because、although、which、that 等连接词恢复逻辑。

例：**The book that you gave me yesterday is useful.** 主干是 **The book is useful**；中间部分是定语从句。
`

var vowelsMd = `
## 元音从口型开始

长短、开口大小和舌位都会改变元音。对比 **ship/sheep、full/fool、cat/cut**，先听再模仿，不要只看拼写。

练习时用镜子观察口型，录下自己的声音，与示范音逐词对比。准确比快更重要。
`

var consonantsMd = `
## 清浊音和气流

清音声带不振动，浊音会振动：**/p/–/b/, /t/–/d/, /k/–/g/**。用手放在喉咙上能感觉区别。

**think** 的 /θ/ 要让舌尖轻触牙齿；**ship** 的 /ʃ/ 像安静时的“嘘”。每次只练一组对比音，再放入完整句子。
`

var stressRhythmMd = `
## 英语按“意群”跳动

实义词通常更重：**I WANT a NEW BOOK.** 冠词、代词和助动词在句中常弱读。

朗读时用斜线划意群：**When I got home / I called my friend.** 先慢速清楚读，再保持重音、缩短弱读词。
`

var themesMd = `
## 主题词汇的四步

以“食物”为例：看图片认 **apple, rice, hungry** → 听短对话 → 用 **I like... / I don't like...** 表达 → 写购物清单。

每个主题同时收集名词、动词、形容词和一个句型，避免只背孤立名词。儿童阶段可用分类、配对、找不同和角色扮演巩固。
`

var wordFormationMd = `
## 单词会“长家族”

**help → helpful → helpless → unhelpful → helpfulness**。

常见前缀：**un-** 否定、**re-** 再、**mis-** 错误；常见后缀：**-er** 人、**-ness** 名词、**-able** 能够……的、**-ly** 副词。

先判断词性，再根据上下文猜意思，最后查词典确认。词根推测是线索，不是保证。
`

var collocationsMd = `
## 词语要和“伙伴”一起记

说 **make a decision / do homework / take a photo / heavy rain**，比单背 make、do、take 更接近真实语言。

短语动词要整体理解：**look after** 照顾、**find out** 查明、**give up** 放弃。把每个词块放进一个和自己有关的句子。
`

var nounsMd = `
## 名词的三个问题

它能不能数？一个还是多个？需不需要限定？

可数：**a chair / two chairs**；不可数：**some advice / a piece of advice**。复数常见 **-s/-es**，但 **child–children, mouse–mice** 需要单独记。

所有格：**Tom's book**；多人共有：**Tom and Lily's room**。
`

var pronounsMd = `
## 代词替代名词

**I/me/my/mine/myself** 分别常作主语、宾语、限定词、独立所有格和反身代词：**I gave my book to him. The book is mine.**

**some/any** 常用于数量；**each/every** 后接单数名词但意义不同。先看它修饰的是一个整体还是一个个成员。
`

var adjectivesMd = `
## 形容词描述，副词修饰

**a careful student**（描述名词）；**She works carefully**（描述动词）。并非所有副词都以 -ly 结尾：**fast, hard, late** 既可作形容词也可作副词。

顺序常是观点→大小→年龄→颜色→来源→材料→名词：**a lovely small old red box**。真实表达中不要堆太多形容词。
`

var functionWordsMd = `
## 小词决定大关系

**in/on/at** 表示从大范围到具体点；**and/but/so/because** 表示并列、转折、结果和原因；**oh/wow/ouch** 表达即时情绪。

先把两句短句说清楚，再用连接词合并。介词不要只背中文，要连同场景和搭配记忆。
`

var verbFormsMd = `
## 动词会变形

**work–works–worked–working–worked**。第三人称单数、过去式、分词和进行式分别服务于不同结构。

助动词 **do** 帮助疑问和否定，**be** 帮助进行和被动，**have** 帮助完成。助动词出现后，主要动词通常回到原形：**Does he like it?**
`

var modalsMd = `
## 情态动词表达态度

**can** 能力/许可；**must** 强义务或肯定推测；**have to** 外部要求；**should** 建议；**may/might** 可能性。

情态动词后用动词原形，没有三单 **-s**：**She can swim.** 过去或更委婉时，可用 **could/would**，但语气要结合场景理解。
`

var passiveMd = `
## 被动语态

结构是 **be + 过去分词**：**The window was broken.** 重点是窗户和结果，谁打破的并不重要。

时态由 be 承担：**is made / was made / has been made / will be made**。需要说明动作执行者时用 **by**。科技、新闻和说明文里很常见。
`

var infinitiveMd = `
## to do 的常见角色

**To learn takes time.**（主语）  
**I want to learn.**（宾语）  
**I came to help.**（目的）  
**I have a book to read.**（后置修饰）

使役动词 **make/let/have + 人 + 动词原形**；ask/tell/want 常用 **to do**。看到结构先看前面的动词决定后面形式。
`

var gerundParticipleMd = `
## doing 有两张“身份证”

作名词：**Reading is fun. I enjoy reading.**  
作分词：**The boy reading by the window is my brother.**

现在分词还可组成进行时；过去分词可作定语或完成时的一部分。**a boring film / a bored audience** 要分清“令人……的”和“感到……的”。
`

var questionsMd = `
## 先找助动词

**Do you like it? / Is she ready? / Have they arrived? / Can he swim?**

特殊疑问句把疑问词放前面：**Where do you live?** 若疑问词就是主语，不再倒装：**Who called you?**

否定通常把 **not** 放在助动词后：**does not, is not, have not, cannot**。
`

var comparisonMd = `
## 比较的三种距离

短词常用 **-er/-est**：small–smaller–the smallest；长词用 **more/most**：more interesting。

**as...as** 表示一样，**less...than** 表示较不。比较结构后面的对象要同类：**The cost of X is higher than that of Y.**
`

var conditionalsMd = `
## 条件句的时间线

| 类型 | 结构 | 意义 |
|---|---|---|
| 零 | if + 一般现在, 一般现在 | 规律 |
| 一 | if + 一般现在, will do | 真实可能 |
| 二 | if + 一般过去, would do | 现在/将来不太真实 |
| 三 | if + had done, would have done | 过去没有发生 |

主句和 if 从句可以换顺序；if 从句在前时通常用逗号。**unless** = if...not，但不要和 not 重复。
`

var subjunctiveMd = `
## 想象和事实分开

**I wish I had more time.** 表示现在没有那么多时间；**I wish I had studied harder.** 表示对过去的遗憾。

正式建议中可见：**It is important that he be on time.** 日常英语也常说 **is**。学习时先掌握意义，再根据语域选择形式。
`

var reportedSpeechMd = `
## 转述四件事

调整人称、时态、地点时间词和句型：

**“I am busy,” she said. → She said (that) she was busy.**

一般疑问用 **if/whether**；特殊疑问保留疑问词但用陈述语序：**He asked where I lived.** 如果事实仍然成立，时态不一定要机械后移。
`

var inversionMd = `
## 把重点放到句首

正常：**I had never seen such a view.**  
强调：**Never had I seen such a view.**

以 **only, never, rarely, hardly** 等限制性副词开头时，正式英语常倒装助动词和主语。**It was John who called** 是强调句。阅读中能认出即可，写作中要确保结构准确。
`

var punctuationMd = `
## 标点帮助读者呼吸

句首、专有名词和 **I** 要大写；句子之间用句号。逗号可分隔项目、开头状语或非限制性信息，但不能把两个完整句子只用逗号硬连。

冒号引出解释或列表，分号连接关系紧密的完整句。引用时根据英美写作规范统一引号和标点位置。
`

var greetingsMd = `
## 同一句话，不同距离

朋友之间：**Hi! How's it going?**  
第一次见面：**Nice to meet you.**  
正式场合：**Good morning. It's a pleasure to meet you.**

介绍自己可以说：**I'm... / I come from... / I work as...**。结束时说：**It was nice talking to you. See you soon.**
`

var requestsMd = `
## 请求的语气梯度

直接：**Open the window.**  
普通：**Can you open the window?**  
更礼貌：**Could you possibly open the window?**  
征求许可：**Would you mind opening the window?**

建议：**You could... / Why don't we...?**；邀请：**Would you like to...?**。表达感谢并回应：**Thanks. — You're welcome.**
`

var opinionsMd = `
## 观点三件套

提出：**In my view, ...**；给理由：**The main reason is that...**；举例：**For example, ...**。

同意：**I agree. / That's a good point.**  
委婉不同意：**I see what you mean, but... / I'm not sure I agree because...**

先回应对方，再说自己的观点，交流会比直接说 No 更合作。
`

var listeningFoundationMd = `
## 听力不是逐词翻译

听前看标题和图片，猜场景；听中先抓人、地、时间、动作等关键词；听后用一句中文或英文说出大意。

儿童可以从歌曲、押韵、动画短故事开始，重复听同一段比每天换很多陌生材料更有效。听懂七成时就可以跟读。
`

var listeningStrategiesMd = `
## 三遍听法

第一遍只回答“在讲什么”；第二遍定位数字、地点、态度和转折；第三遍打开文本，标出连读、弱读和生词。

听不清时不要立刻暂停每个词：先判断是词汇不认识、声音连在一起，还是注意力错过了。最后用 20 秒复述，确认自己真的听懂。
`

var speakingFluencyMd = `
## 流利度不是说得快

先练意群：**I think / it is a good idea / because it saves time.** 在意群之间停顿，在关键词上重读。

跟读四步：听一句 → 看文本跟读 → 盖住文本复述 → 录音对比。说错时用 **Let me try that again.** 自然修正，不要因为一个错误停掉整句话。
`

var conversationMd = `
## 让对话走下去

打开话题：**How was your weekend?**  
追问：**What happened next? / Why do you think so?**  
澄清：**Do you mean...? / Could you explain that?**  
接话：**Really? That sounds...**

好的对话不是一个人长篇演讲，而是轮流、回应、追问和给对方空间。
`

var readingFoundationMd = `
## 从“猜整体”开始

先看标题、图片和小标题，预测主题；再读第一句和最后一句，找到段落主旨；遇到生词先看前后句，不要每个词都查。

儿童阶段可使用分级读物和重复故事：第一次看图理解，第二次指读，第三次换掉一个词讲自己的版本。
`

var readingStrategiesMd = `
## 读长文的地图

先看结构：观点在哪里？例子支持什么？**however/therefore/for example** 等词会提示逻辑。

推断必须有证据；区分作者说的事实、作者的观点和你的判断。读完用三句话总结：主题、关键证据、作者态度。
`

var writingFoundationMd = `
## 好段落的四块砖

**主题句**告诉读者本段说什么；**支持句**给事实或例子；**解释句**说明为什么；**结尾句**收束或连接下一段。

先写清楚短句，再用 **and, but, because, so, however** 连接。修改顺序：意思 → 句子完整 → 动词时态 → 拼写和标点。
`

var writingAcademicMd = `
## 正式写作的工作流

审题 → 列提纲 → 写段落 → 加证据 → 检查连接 → 修改语言。议论文常见骨架：观点、理由、例子、回应另一种看法、结论。

避免为了“高级”而堆复杂词；优先使用准确的动词和清晰的逻辑。把 **I think**、例子和结论写得具体，文章就会更有力量。
`

var learningPlanMd = `
## 先决定“我要能做什么”

不要只写“提高英语”，改成可观察目标：四周后能用英语介绍自己两分钟；能读懂 A2 短文并说出主旨；能写 120 词邮件。

每周安排：输入 3 次、输出 2 次、复盘 1 次。选择略高于当前水平的材料，完成比收藏更多资源重要。
`

var learningMethodMd = `
## 输入、提取、反馈

看懂是输入，合上书回忆是提取，说写出来是输出，得到纠正才有反馈。一个词至少在不同日子主动回忆四次：当天、第二天、一周后、一个月后。

错题不要只抄答案：记录“我为什么错、下次看到什么信号”。把错误变成下一轮练习的题目。
`

var examRoadmapMd = `
## CEFR 是能力地图

| 阶段 | 能做什么 |
|---|---|
| Pre-A1 | 认识熟悉词、跟随极简指令 |
| A1 | 介绍自己，理解非常熟悉的短句 |
| A2 | 处理日常简单交流，读懂短文本 |
| B1 | 处理旅行、学习和工作中的常见情况 |
| B2 | 理解较复杂文本，清晰表达观点 |
| C1 | 灵活、准确地处理长篇和专业内容 |
| C2 | 几乎无障碍地理解并精确表达细微意义 |

考试是阶段性测量，不是英语的全部。先按听说读写的真实任务学习，再选择适合的考试材料。
`
