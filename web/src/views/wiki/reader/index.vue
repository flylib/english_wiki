<template>
  <div class="wiki-reader">
    <!-- 左侧:知识模块树 -->
    <div class="left" v-show="!treeCollapsed">
      <div class="left-head">知识模块</div>
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
        <div class="markdown-body" v-html="detailHtml"></div>
      </div>

      <!-- 条目列表 -->
      <div v-else class="list">
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
}
.left {
  width: 260px;
  flex-shrink: 0;
  border-right: 1px solid var(--color-border-2);
  padding: 12px 8px;
  overflow: auto;

  .left-head {
    font-weight: 600;
    padding: 4px 8px 12px;
    color: var(--color-text-1);
  }
}
.right {
  flex: 1;
  overflow: auto;
  padding: 16px 24px;
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
