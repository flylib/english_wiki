<claude-mem-context>
# Memory Context

# [english_wiki] recent context, 7/31/2026 11:56am GMT+8

Legend: 🎯session 🔴bugfix 🟣feature 🔄refactor ✅change 🔵discovery ⚖️decision 🚨security_alert 🔐security_note
Format: ID TIME TYPE TITLE
Fetch details: get_observations([IDs]) | Search: mem-search skill

Stats: 50 obs (9,374t read) | 1,038,268t work | 99% savings

### Jul 7, 2026
2969 5:03p 🟣 wiki 后端模块编译通过，Task 2 完成
2970 " 🔵 冒烟测试：登录成功但 /wiki/category/tree 返回非法 JSON（python json.tool 报 Extra data）
2971 5:04p 🔵 wiki 路由 404：端口 8000 运行的是旧进程，新 english-wiki 二进制因端口占用启动失败
2972 " 🔵 端口 8000 被 db_tool 的旧 go-admin 进程（PID 14776）占用
2973 " 🟣 wiki 后端端到端冒烟测试通过：分类树 + 种子数据全部正确返回
2974 5:05p 🟣 wiki 后端全功能冒烟测试通过：等级过滤、分类聚合、等级体系、条目详情全部正确
2975 " 🔵 db_tool workbench.vue 布局参考：左面板（可拖拽调宽）+ 右内容区，localStorage 记忆宽度
2976 " 🔵 前端 API 模块规范和 Navbar 结构：wiki API 文件模板及 /workbench 硬编码路径需更新
2977 " 🟣 web/src/api/wiki.js 创建：前端 wiki API 封装层含 CEFR 标签工具函数
2978 5:06p 🟣 前端路由完全替换为 wiki 路由：dbtool 残留路径全部清除
2979 " 🟣 全部三个 Vue 页面创建完成 + dbtool 残留彻底清除
S982 用户询问账号密码 + 待处理分类树重组请求 (Jul 7 at 8:24 PM)
2996 8:26p ⚖️ wiki_category 树结构调整：时态/从句升为一级
### Jul 8, 2026
S981 分类树结构调整：英语时态/英语从句升为一级菜单，各时态降为二级 (Jul 8 at 3:03 PM)
S986 阅读页 UX 重设计需求收到：一级分类显示概览+图表，二级分类内容直接铺开（待实现） (Jul 8 at 3:31 PM)
2997 3:32p ✅ wiki_category 分类树重组：时态/从句升为一级，各时态自建二级分类
2998 " ✅ seed_content.go 所有时态条目补全 catName/catCode 字段并编译通过
2999 " 🟣 实时 API 迁移：分类树重构为两级（时态各自独立二级）
3000 " ✅ seed_content.go 编译验证通过（重复确认）
3001 " 🔵 make stop 在 server/ 子目录调用触发 go-admin 内置 Docker Makefile 而非项目根 Makefile
3002 " 🔵 make stop 必须在项目根目录执行；从 server/ 子目录执行触发 go-admin Docker Makefile
3003 3:34p 🔵 后端直接启动命令（绕过 make）和前端 preview 启动方式
3004 " 🟣 阅读页分类树重组验证通过：时态两级展开 + 水平过滤同时正常
3005 3:35p 🟣 阅读页截图确认：两级分类树 + CET-4 过滤同时正常渲染
3016 " ⚖️ 阅读页 UX 重设计：一级分类显示概览图，二级分类内容直接铺开
S983 分类树重组：英语时态/从句升为一级，各时态建独立二级分类 — 端到端完成 (Jul 8 at 3:35 PM)
S984 阅读页 UX 重设计：一级分类显示介绍+可视化图表，二级分类内容直接铺开不再二次跳转 (Jul 8 at 3:35 PM)
S991 去掉右边的1 — 移除左侧分类树节点上的条目数量角标 (Jul 8 at 3:53 PM)
3017 3:55p 🟣 reader/index.vue loadEntries 智能加载逻辑：单篇自动铺开，空目录聚合子条目
3018 3:56p 🟣 英语时态总览介绍页：SVG 16时态矩阵图 + 一键插入 DB + seed 代码同步
3019 3:57p 🔴 UI 验证通过：英语时态总览页（SVG矩阵图）自动铺开 + 现在完成时内容直接展示
S989 英语时态 UX 重设计：L1分类点击显示带SVG矩阵图的总览介绍页，L2分类点击直接铺开Markdown内容 (Jul 8 at 3:57 PM)
S992 去掉右边的1 — 移除左侧分类树节点上的条目数量角标（entryCount badge） (Jul 8 at 4:00 PM)
S993 去掉右边的1 — 移除左侧分类树节点条目数量角标，已完成并通过截图验证 (Jul 8 at 4:02 PM)
3020 4:11p 🟣 新增三个 L1 菜单分类请求：句子成分、五种基本句型、名词冠词介词
3021 " 🟣 新增三个 L1 语法分类并重排顺序
3022 4:12p 🔵 seed_content.go 未包含新三个分类，需同步补丁
3023 " 🟣 seed_content.go 同步新增三个 L1 分类，音标/词汇排序顺延
3024 " 🟣 seed_content.go 新增三分类补丁编译通过
3025 4:13p 🟣 浏览器刷新验证新增三个 L1 分类菜单显示
3026 " 🟣 截图确认三个新 L1 分类在知识模块树中正常显示
3027 " 🟣 DOM snapshot 确认新三个 L1 分类在左侧导航树中正确渲染
S994 加上三个新菜单：句子成分、五种基本句型、名词冠词介词 — 已完成并验证 (Jul 8 at 4:13 PM)
### Jul 29, 2026
4939 3:16p 🔵 英语百科项目结构与现有内容盘点
4940 3:17p 🔵 seed_content.go 完整时态条目清单：8篇已写入，SVG矩阵图覆盖全16种时态
4941 " ⚖️ 参考剑桥YLE/英国文化协会分级体系设计卡通化内容策略
4942 3:18p 🟣 Reader页新增欢迎引导页和儿童模式(kids-mode)切换
4943 " 🟣 Reader组件完整样式体系：欢迎页英雄区、路径卡片、kids-mode差异化渲染
4944 " 🟣 seed_content.go 为5个空白分类注入初始介绍条目
4945 " 🟣 5个分类的Markdown内容正文写入seed_content.go
4946 3:19p 🔴 trailing whitespace修复：Markdown双空格换行改为&lt;br&gt;标签
4947 3:20p 🔵 Go测试在沙箱中无法访问默认缓存目录，需设置GOCACHE=/tmp
4949 " 🔵 后端API关键行为：CEFR区间交集过滤、子分类BFS遍历、删除保护
4950 " 🔵 seed_content.go最终状态确认：6个Markdown变量已写入，英语从句分类仍无条目
4948 " 🔵 Go wiki包编译通过，无测试文件但seed_content.go变更已验证可编译
4951 3:22p 🔄 seedContent改为增量模式：按code查重后创建，支持已有数据库逐步升级
4952 3:23p 🔄 putEntry辅助函数：按(category_id, title)查重的条目upsert，统一替换所有直接Create调用
4953 3:25p 🔴 content_catalog.go 新建：ensureContentCatalog 框架 + categorySeeds/catalogEntries 数据声明
4954 " 🔴 content_catalog.go Markdown 变量添加状态：第二批成功，第三批失败
4959 3:29p 🔴 新增"英语全景"顶层导航分类和全景地图总览条目

Access 1038k tokens of past work via get_observations([IDs]) or mem-search skill.
</claude-mem-context>