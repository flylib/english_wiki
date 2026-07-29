<template>
  <div class="wiki-reader" :class="{ 'kids-mode': isKidsMode }">
    <!-- 左侧:知识模块树 -->
    <div class="left" v-show="!treeCollapsed">
      <div class="left-head">
        <span class="head-emoji">🧭</span>
        <span>知识模块</span>
        <small>Pick a path</small>
      </div>
      <a-tree
        v-if="treeData.length"
        :data="treeData"
        :field-names="{ key: 'id', title: 'name', icon: '_icon' }"
        :default-expand-all="true"
        block-node
        @select="onSelectCategory"
      />
      <a-empty v-else-if="!treeLoading" description="暂无分类" />
    </div>

    <!-- 右侧:内容区 -->
    <div class="right">
      <!-- 顶栏:当前位置 + 身份/等级切换 -->
      <div class="toolbar">
        <div class="crumb">
          <icon-book style="margin-right: 6px" />
          <span>{{ currentCategory ? currentCategory.name : '全部内容' }}</span>
        </div>
        <div class="level-picker">
          <span class="lp-label">我的水平:</span>
          <a-cascader
            v-model="levelValue"
            :options="levelOptions"
            placeholder="不限(看全部)"
            allow-clear
            style="width: 240px"
            @change="onLevelChange"
          />
          <a-tag v-if="cefrRange" color="arcoblue" style="margin-left: 8px">
            适配 {{ cefrRangeLabel }}
          </a-tag>
        </div>
      </div>

      <!-- 学习导航 -->
      <div v-if="!currentCategory && !detail" class="welcome">
        <div class="hero">
          <div class="hero-copy">
            <span class="eyebrow">ENGLISH EXPLORER</span>
            <h1>把英语学成一场小小的冒险</h1>
            <p>从一句 Hello 开始，沿着自己的等级地图，慢慢解锁单词、句子、语法和真实表达。</p>
            <div class="hero-actions">
              <a-button type="primary" @click="selectFirstCategory">开始第一课 <icon-right /></a-button>
              <span class="hero-note">每天 10 分钟，也会有进步 ✨</span>
            </div>
          </div>
          <div class="hero-art" aria-hidden="true">
            <div class="sun">☀️</div>
            <div class="cloud cloud-a">☁️</div>
            <div class="cloud cloud-b">☁️</div>
            <div class="book">📖</div>
            <div class="explorer">🧑‍🚀</div>
            <div class="flag">ABC!</div>
          </div>
        </div>

        <div class="welcome-grid">
          <a-card class="path-card path-kids" hoverable @click="chooseLevel('G1')">
            <div class="path-icon">🎈</div><div><b>小小探险家</b><p>Pre-A1 · A1 · A2</p><span>图画、故事、短句和小游戏</span></div>
          </a-card>
          <a-card class="path-card path-teens" hoverable @click="chooseLevel('G7')">
            <div class="path-icon">🚀</div><div><b>成长挑战者</b><p>B1 · B2</p><span>语法规律、阅读表达与写作</span></div>
          </a-card>
          <a-card class="path-card path-pro" hoverable @click="chooseLevel('CET6')">
            <div class="path-icon">🧠</div><div><b>进阶探索者</b><p>C1 · C2</p><span>精准表达、复杂句与学术英语</span></div>
          </a-card>
        </div>

        <div class="learning-loop">
          <div><span>01</span><b>看懂</b><small>例句 + 图示</small></div><i>→</i>
          <div><span>02</span><b>开口</b><small>跟读 + 替换</small></div><i>→</i>
          <div><span>03</span><b>会用</b><small>练习 + 输出</small></div><i>→</i>
          <div><span>04</span><b>复习</b><small>间隔回看</small></div>
        </div>
        <p class="source-note">内容编排参考 <a href="https://learnenglishkids.britishcouncil.org/" target="_blank">British Council LearnEnglish Kids</a>、<a href="https://www.cambridgeenglish.org/exams-and-tests/qualifications/young-learners/" target="_blank">Cambridge Young Learners</a> 与 CEFR Companion Volume。</p>
      </div>

      <!-- 条目详情 -->
      <div v-if="detail" class="detail">
        <a-button
          v-if="entries.length > 1"
          size="small"
          style="margin-bottom: 12px"
          @click="detail = null"
        >
          <template #icon><icon-left /></template>
          返回列表
        </a-button>
        <h1 class="d-title">{{ detail.title }}</h1>
        <div class="d-meta">
          <a-tag color="green">{{ cefrRangeText(detail.cefrMin, detail.cefrMax) }}</a-tag>
          <span class="d-time">更新于 {{ (detail.updatedAt || '').slice(0, 10) }}</span>
        </div>
        <div class="detail-tip">💡 学习小提示：先读例句，再试着把主语、时间或地点换成自己的内容。</div>
        <div class="markdown-body" v-html="detailHtml"></div>
      </div>

      <!-- 条目列表 -->
      <div v-else-if="currentCategory" class="list">
        <a-spin :loading="listLoading" style="display: block">
          <a-card
            v-for="e in entries"
            :key="e.id"
            class="entry-card"
            hoverable
            @click="openEntry(e.id)"
          >
            <div class="ec-head">
              <span class="ec-title">{{ e.title }}</span>
              <a-tag color="green" size="small">{{ cefrRangeText(e.cefrMin, e.cefrMax) }}</a-tag>
            </div>
            <div class="ec-summary">{{ e.summary }}</div>
          </a-card>
          <a-empty
            v-if="!listLoading && !entries.length"
            description="该模块暂无适配当前水平的内容,试试调整右上角的水平设置"
          />
        </a-spin>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { marked } from 'marked';
import {
  categoryTree,
  listEntries,
  getEntry,
  levelSystems,
  cefrRangeText,
} from '@/api/wiki';
import { treeCollapsed } from '@/composables/uiState';

const LEVEL_KEY = 'wiki_my_level';

const treeData = ref([]);
const treeLoading = ref(true);
const currentCategory = ref(null);

const entries = ref([]);
const listLoading = ref(false);

const detail = ref(null);
const detailHtml = computed(() =>
  detail.value ? marked.parse(detail.value.contentMd || '') : ''
);

const isKidsMode = computed(() => {
  const v = levelValue.value;
  const id = Array.isArray(v) ? v[v.length - 1] : v;
  return ['G1', 'G2', 'G3', 'G4', 'G5', 'G6', 'KET'].includes(id);
});

// 等级切换:级联 [体系code, 档位id],映射出 CEFR 区间
const systems = ref([]);
const levelValue = ref(JSON.parse(localStorage.getItem(LEVEL_KEY) || 'null') || undefined);
const levelOptions = computed(() =>
  systems.value.map((s) => ({
    value: s.code,
    label: s.name,
    children: (s.levels || []).map((l) => ({ value: l.id, label: l.name })),
  }))
);
const cefrRange = computed(() => {
  // arco cascader 默认只回传叶子值(档位 id);兼容路径数组形态
  const v = levelValue.value;
  const levelId = Array.isArray(v) ? v[v.length - 1] : v;
  if (!levelId) return null;
  for (const s of systems.value) {
    const hit = (s.levels || []).find((l) => l.id === levelId);
    if (hit) return [hit.cefrMin, hit.cefrMax];
  }
  return null;
});
const cefrRangeLabel = computed(() =>
  cefrRange.value ? cefrRangeText(cefrRange.value[0], cefrRange.value[1]) : ''
);

const loadEntries = async () => {
  listLoading.value = true;
  try {
    const base = { status: 'published', pageSize: 100 };
    if (cefrRange.value) {
      base.cefrMin = cefrRange.value[0];
      base.cefrMax = cefrRange.value[1];
    }
    if (!currentCategory.value) {
      const res = await listEntries(base);
      entries.value = res.data.list || [];
      return;
    }
    // 先取直挂在本分类下的条目(分类自己的介绍页/正文)
    let res = await listEntries({ ...base, categoryId: currentCategory.value.id, includeSub: 0 });
    let list = res.data.list || [];
    if (list.length === 1) {
      // 只有一篇 → 直接铺开正文,免掉多一次点击
      entries.value = list;
      await openEntry(list[0].id);
      return;
    }
    if (!list.length) {
      // 纯目录分类 → 汇总子分类内容成卡片列表
      res = await listEntries({ ...base, categoryId: currentCategory.value.id, includeSub: 1 });
      list = res.data.list || [];
    }
    entries.value = list;
  } finally {
    listLoading.value = false;
  }
};

const onSelectCategory = (keys, { node }) => {
  currentCategory.value = node;
  detail.value = null;
  loadEntries();
};

const onLevelChange = () => {
  localStorage.setItem(LEVEL_KEY, JSON.stringify(levelValue.value || null));
  detail.value = null;
  loadEntries();
};

const chooseLevel = (code) => {
  for (const system of systems.value) {
    const level = (system.levels || []).find((item) => item.code === code);
    if (level) {
      levelValue.value = level.id;
      onLevelChange();
      return;
    }
  }
};

const selectFirstCategory = () => {
  const first = treeData.value[0];
  if (first) onSelectCategory([first.id], { node: first });
};

const openEntry = async (id) => {
  const res = await getEntry(id);
  detail.value = res.data;
};

onMounted(async () => {
  try {
    const [treeRes, sysRes] = await Promise.all([categoryTree(), levelSystems()]);
    treeData.value = treeRes.data || [];
    systems.value = sysRes.data || [];
  } finally {
    treeLoading.value = false;
  }
  loadEntries();
});
</script>

<style lang="scss" scoped>
.wiki-reader {
  display: flex;
  height: calc(100vh - 60px);
  overflow: hidden;
  background: var(--color-fill-1);
}
.left {
  width: 250px;
  flex-shrink: 0;
  border-right: 1px solid var(--color-border-2);
  padding: 12px 8px;
  overflow: auto;

  .left-head {
    display: flex;
    align-items: center;
    gap: 7px;
    font-weight: 700;
    padding: 8px 10px 14px;
    color: var(--color-text-1);
    .head-emoji { font-size: 20px; }
    small { margin-left: auto; color: var(--color-text-3); font-size: 10px; font-weight: 500; }
  }
}
.right {
  flex: 1;
  overflow: auto;
  padding: 20px clamp(18px, 4vw, 54px) 48px;
}
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  .crumb {
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text-1);
    display: flex;
    align-items: center;
  }
  .level-picker {
    display: flex;
    align-items: center;
    .lp-label {
      margin-right: 8px;
      color: var(--color-text-2);
    }
  }
}
.entry-card {
  margin-bottom: 12px;
  cursor: pointer;

  .ec-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .ec-title {
      font-size: 15px;
      font-weight: 600;
      color: var(--color-text-1);
    }
  }
  .ec-summary {
    margin-top: 6px;
    color: var(--color-text-3);
    font-size: 13px;
  }
}
.welcome { max-width: 1080px; margin: 0 auto; }
.hero {
  min-height: 270px;
  display: flex;
  overflow: hidden;
  border-radius: 26px;
  background: linear-gradient(120deg, #fff7e8 0%, #fff1d6 54%, #dff7f2 100%);
  box-shadow: 0 14px 34px rgba(246, 171, 77, .12);
}
.hero-copy { flex: 1; padding: 38px 34px; position: relative; z-index: 1; }
.eyebrow { color: #e46b38; font-size: 11px; letter-spacing: 2px; font-weight: 800; }
.hero h1 { margin: 10px 0 12px; max-width: 450px; color: #3f352d; font-size: clamp(28px, 4vw, 42px); line-height: 1.15; }
.hero p { max-width: 520px; color: #6f625a; font-size: 15px; line-height: 1.8; }
.hero-actions { display: flex; align-items: center; gap: 14px; margin-top: 24px; }
.hero-note { color: #9b806b; font-size: 12px; }
.hero-art { width: 40%; min-width: 240px; position: relative; overflow: hidden; }
.sun { position: absolute; top: 28px; right: 54px; font-size: 42px; }
.cloud { position: absolute; opacity: .8; font-size: 42px; }
.cloud-a { top: 54px; left: 18px; }.cloud-b { top: 112px; right: 12px; font-size: 30px; }
.book { position: absolute; bottom: 28px; left: 20%; font-size: 88px; transform: rotate(-8deg); }
.explorer { position: absolute; bottom: 24px; left: 52%; font-size: 94px; transform: rotate(5deg); }
.flag { position: absolute; right: 16px; bottom: 76px; padding: 7px 10px; color: #fff; font-weight: 800; background: #ef8961; border-radius: 12px 12px 12px 2px; transform: rotate(8deg); }
.welcome-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; margin: 22px 0; }
.path-card { border: 0; border-radius: 18px; cursor: pointer; }
.path-card :deep(.arco-card-body) { display: flex; align-items: center; gap: 14px; padding: 18px; }
.path-icon { font-size: 32px; }.path-card b { color: var(--color-text-1); font-size: 16px; }.path-card p { margin: 3px 0; color: var(--color-text-2); font-size: 12px; }.path-card span { color: var(--color-text-3); font-size: 12px; }
.path-kids { background: #fff8df; }.path-teens { background: #eaf7ff; }.path-pro { background: #f1edff; }
.learning-loop { display: flex; align-items: center; justify-content: space-around; padding: 18px 10px; border-radius: 18px; background: var(--color-bg-2); }
.learning-loop div { display: grid; grid-template-columns: 32px auto; column-gap: 8px; align-items: center; }.learning-loop span { grid-row: span 2; display: grid; place-items: center; width: 30px; height: 30px; border-radius: 50%; color: #fff; background: #6d9eff; font-size: 11px; }.learning-loop b { font-size: 14px; }.learning-loop small { color: var(--color-text-3); font-size: 11px; }.learning-loop i { color: var(--color-text-3); font-size: 20px; font-style: normal; }
.source-note { margin: 18px 0; color: var(--color-text-3); font-size: 12px; text-align: center; }.source-note a { color: rgb(var(--arcoblue-6)); }
.detail-tip { margin: 0 0 20px; padding: 11px 14px; border-radius: 12px; color: #8b642f; background: #fff8e5; font-size: 13px; }
.kids-mode .right { background: linear-gradient(180deg, #fffdf4 0, var(--color-fill-1) 280px); }
.kids-mode .left { background: #fffdf4; }
.kids-mode .left :deep(.arco-tree-node) { margin: 3px 0; border-radius: 12px; }
.kids-mode .left :deep(.arco-tree-node-title) { font-size: 14px; }
.kids-mode .entry-card { border-radius: 18px; border: 2px solid #ffe7a9; background: #fffdf7; }
.kids-mode .entry-card .ec-title { color: #d36a42; }
.detail {
  max-width: 860px;

  .d-title {
    font-size: 24px;
    margin: 8px 0;
    color: var(--color-text-1);
  }
  .d-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    .d-time {
      color: var(--color-text-3);
      font-size: 13px;
    }
  }
}
@media (max-width: 800px) {
  .toolbar { align-items: flex-start; flex-direction: column; gap: 10px; }
  .welcome-grid { grid-template-columns: 1fr; }
  .hero-art { width: 34%; min-width: 130px; }.explorer { left: 30%; font-size: 74px; }.book { left: 0; font-size: 64px; }.sun { right: 15px; font-size: 30px; }
  .hero-copy { padding: 28px 22px; }.hero h1 { font-size: 28px; }.hero-note { display: none; }
  .learning-loop { flex-wrap: wrap; gap: 12px; }.learning-loop i { display: none; }
}
</style>

<style lang="scss">
/* Markdown 渲染样式(全局,不 scoped) */
.markdown-body {
  color: var(--color-text-1);
  line-height: 1.8;
  font-size: 14px;

  h2 {
    font-size: 18px;
    margin: 20px 0 10px;
    padding-bottom: 6px;
    border-bottom: 1px solid var(--color-border-2);
  }
  h3 {
    font-size: 16px;
    margin: 16px 0 8px;
  }
  p {
    margin: 8px 0;
  }
  ul,
  ol {
    padding-left: 22px;
    margin: 8px 0;
  }
  li {
    margin: 4px 0;
  }
  code {
    background: var(--color-fill-2);
    padding: 1px 5px;
    border-radius: 3px;
    font-size: 13px;
  }
  blockquote {
    border-left: 3px solid rgb(var(--arcoblue-5));
    background: var(--color-fill-1);
    margin: 10px 0;
    padding: 6px 12px;
    color: var(--color-text-2);
  }
  table {
    border-collapse: collapse;
    margin: 12px 0;
    width: 100%;
    max-width: 640px;

    th,
    td {
      border: 1px solid var(--color-border-2);
      padding: 6px 12px;
      text-align: left;
    }
    th {
      background: var(--color-fill-2);
    }
  }
  strong {
    color: rgb(var(--arcoblue-6));
  }
}
</style>
