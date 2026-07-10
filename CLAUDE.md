# CLAUDE.md — english_wiki 项目约定

面向中国用户的**英语学习百科全书**:左侧知识模块树(语法/时态/从句/音标/词汇…),右侧内容按用户水平过滤展示。骨架复制自 db_tool(go-admin 全家桶),数据存后端 SQLite,媒体文件放服务本地目录。

## 架构

```
english_wiki/
├── server/   go-admin 后端(Gin + GORM + SQLite),业务模块在 app/wiki/
└── web/      定制精简版 go-admin-ui(Arco Design Vue3):硬编码路由 + 顶栏壳,不走 sys_menu 动态菜单
```

### 等级体系(核心设计)

**内容只打 CEFR 标**(整数 rank:1=Pre-A1 … 7=C2,存 `wiki_entry.cefr_min/cefr_max` 区间);
学段(小学一年级…大学)与考试(CET4/雅思/KET…)是 `wiki_level` 里映射到 CEFR 区间的档位,**不直接标在内容上**。
用户在阅读页选"我的水平"(任一体系档位)→ 前端换算成 CEFR 区间 → 后端按**区间有交集**过滤条目。
加新体系(如考研、剑桥少儿细分)= 插 `wiki_level_system`/`wiki_level` 数据,内容零改动。

### 表(app/wiki/models/wiki.go,启动时 AutoMigrate + 幂等 seed)

- `wiki_category` 分类树(parent_id;code 唯一)
- `wiki_entry` 条目(Markdown 正文、cefr_min/max、status draft/published)
- `wiki_level_system` / `wiki_level` 等级体系与档位(seed 三套:cefr/grade/exam)
- seed 数据在 `app/wiki/service/seed_levels.go`、`seed_content.go`(仅空库兜底样例);**内容主库是仓库自带的 go-admin-db.db**(16 时态 + 从句三类等全量内容,经管理后台/API 维护)

## 启动(根目录 Makefile)

```bash
make dev        # 一键:后端(:739)+ 前端(:1798),Ctrl+C 一起退出
make backend    # 仅后端(先构建再运行);make frontend 仅前端
make stop       # 停掉占用 739/1798 的进程
make deps       # go mod download + pnpm install(首次)
```

后端 :739 是刻意避开本机 db_tool 的 :8000;构建必须带 `sqlite3` build tag(Makefile 已带)。
自带 `go-admin-db.db` 已 seed;登录 **admin / 123456**(dev 模式验证码随便填)。vite 已代理 `/api/v1` → :739。

## 页面与路由(硬编码在 web/src/router/index.js)

- `/wiki` 阅读页(views/wiki/reader):左分类树 + 顶部水平切换(cascader) + 条目卡片/详情(marked 渲染 Markdown)
  - 分类点击约定:分类**直挂恰好 1 篇**条目 → 直接铺开正文(介绍页/单篇模式,如「英语时态」总览、各具体时态);直挂 0 篇 → 汇总子分类内容出卡片列表;多篇 → 卡片列表。介绍页正文里可嵌 SVG(v-html 直出,配色用 arco CSS 变量适配暗色)。
- `/admin/entry`、`/admin/category` 管理页(入口在右上角头像下拉,仅超管)
- 登录后跳 `/wiki`(views/login/index.vue)

## 关键约定 / 坑

- **arco 树组件**:节点数据里的 `icon` 字段会被当成**渲染函数**,字符串会让整棵树渲染崩溃——所有 a-tree / a-tree-select 的 `field-names` 必须带 `icon: '_icon'` 把它挡掉。
- **arco cascader** 默认只回传**叶子值**(非路径数组),reader 里 cefrRange 已兼容两种形态。
- 新增接口照 `app/wiki/router/wiki.go` 注册;模块经 `cmd/api/wiki.go` 挂进 AppRouters;写操作在 API 内部校验超管(`user.GetRoleName(c) == "admin"`)。
- 端口改动要同步两处:`server/config/settings.sqlite.yml` 与 `web/vite.config.js` 代理。
- 媒体文件(规划中):落 `server/static/uploads/`,DB 只存相对路径(已 gitignore)。
- 后端冒烟:`POST /api/v1/login {username:admin,password:123456,code:0,uuid:x}` 拿 token 打 `/api/v1/wiki/*`。

## 下一步(未做)

- 词库 `word` 表与词汇模块、媒体上传接口、从句/音标内容、阅读端免登录(目前全部接口需 JWT)。
