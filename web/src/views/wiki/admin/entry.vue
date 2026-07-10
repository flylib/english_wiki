<template>
  <div class="wiki-admin">
    <div class="filters">
      <a-tree-select
        v-model="filter.categoryId"
        :data="treeData"
        :field-names="{ key: 'id', title: 'name', children: 'children', icon: '_icon' }"
        placeholder="按分类过滤"
        allow-clear
        style="width: 220px"
        @change="load"
      />
      <a-select v-model="filter.status" placeholder="状态" allow-clear style="width: 130px" @change="load">
        <a-option value="published">已发布</a-option>
        <a-option value="draft">草稿</a-option>
      </a-select>
      <a-input-search
        v-model="filter.keyword"
        placeholder="搜索标题/摘要"
        style="width: 220px"
        allow-clear
        @search="load"
        @press-enter="load"
        @clear="load"
      />
      <div class="spacer"></div>
      <a-button type="primary" @click="openForm()">
        <template #icon><icon-plus /></template>
        新建条目
      </a-button>
    </div>

    <a-table
      :data="rows"
      :loading="loading"
      :pagination="{ total, current: page, pageSize, showTotal: true }"
      row-key="id"
      @page-change="(p) => { page = p; load(); }"
    >
      <template #columns>
        <a-table-column title="ID" data-index="id" :width="70" />
        <a-table-column title="标题" data-index="title" />
        <a-table-column title="分类" :width="140">
          <template #cell="{ record }">{{ catName(record.categoryId) }}</template>
        </a-table-column>
        <a-table-column title="适用等级" :width="120">
          <template #cell="{ record }">
            <a-tag color="green" size="small">{{ cefrRangeText(record.cefrMin, record.cefrMax) }}</a-tag>
          </template>
        </a-table-column>
        <a-table-column title="状态" :width="100">
          <template #cell="{ record }">
            <a-tag :color="record.status === 'published' ? 'arcoblue' : 'gray'" size="small">
              {{ record.status === 'published' ? '已发布' : '草稿' }}
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="排序" data-index="sort" :width="80" />
        <a-table-column title="更新时间" :width="120">
          <template #cell="{ record }">{{ (record.updatedAt || '').slice(0, 10) }}</template>
        </a-table-column>
        <a-table-column title="操作" :width="160">
          <template #cell="{ record }">
            <a-button size="mini" type="text" @click="openForm(record)">编辑</a-button>
            <a-popconfirm content="确认删除该条目?" type="warning" @ok="onDelete(record)">
              <a-button size="mini" type="text" status="danger">删除</a-button>
            </a-popconfirm>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <!-- 编辑弹窗 -->
    <a-modal
      v-model:visible="formVisible"
      :title="form.id ? '编辑条目' : '新建条目'"
      :width="920"
      :ok-loading="saving"
      @ok="save"
    >
      <a-form :model="form" layout="vertical">
        <div class="form-row">
          <a-form-item label="标题" required style="flex: 2">
            <a-input v-model="form.title" placeholder="如:一般现在时 (Simple Present)" />
          </a-form-item>
          <a-form-item label="所属分类" required style="flex: 1">
            <a-tree-select
              v-model="form.categoryId"
              :data="treeData"
              :field-names="{ key: 'id', title: 'name', children: 'children', icon: '_icon' }"
              placeholder="选择分类"
            />
          </a-form-item>
        </div>
        <a-form-item label="摘要">
          <a-input v-model="form.summary" placeholder="列表页展示的一句话简介" />
        </a-form-item>
        <div class="form-row">
          <a-form-item label="适用等级下限(CEFR)" style="flex: 1">
            <a-select v-model="form.cefrMin">
              <a-option v-for="(l, i) in CEFR_LABELS.slice(1)" :key="i + 1" :value="i + 1">{{ l }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="适用等级上限(CEFR)" style="flex: 1">
            <a-select v-model="form.cefrMax">
              <a-option v-for="(l, i) in CEFR_LABELS.slice(1)" :key="i + 1" :value="i + 1">{{ l }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="状态" style="flex: 1">
            <a-radio-group v-model="form.status" type="button">
              <a-radio value="draft">草稿</a-radio>
              <a-radio value="published">发布</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="排序" style="flex: 0 0 100px">
            <a-input-number v-model="form.sort" :min="0" />
          </a-form-item>
        </div>
        <a-form-item label="正文(Markdown)">
          <a-textarea
            v-model="form.contentMd"
            :auto-size="{ minRows: 14, maxRows: 24 }"
            placeholder="支持 Markdown:## 标题、**加粗**、表格、列表…"
            class="md-input"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import {
  categoryTree,
  listEntries,
  getEntry,
  addEntry,
  updateEntry,
  removeEntry,
  cefrRangeText,
  CEFR_LABELS,
} from '@/api/wiki';

const treeData = ref([]);
const catMap = ref({});
const rows = ref([]);
const loading = ref(false);
const total = ref(0);
const page = ref(1);
const pageSize = 20;

const filter = reactive({ categoryId: undefined, status: undefined, keyword: '' });

const flattenCats = (nodes) => {
  for (const n of nodes || []) {
    catMap.value[n.id] = n.name;
    flattenCats(n.children);
  }
};
const catName = (id) => catMap.value[id] || id;

const load = async () => {
  loading.value = true;
  try {
    const params = { pageIndex: page.value, pageSize };
    if (filter.categoryId) params.categoryId = filter.categoryId;
    if (filter.status) params.status = filter.status;
    if (filter.keyword) params.keyword = filter.keyword;
    const res = await listEntries(params);
    rows.value = res.data.list || [];
    total.value = res.data.count || 0;
  } finally {
    loading.value = false;
  }
};

// 表单
const formVisible = ref(false);
const saving = ref(false);
const emptyForm = () => ({
  id: null, categoryId: undefined, title: '', summary: '',
  contentMd: '', cefrMin: 1, cefrMax: 7, status: 'draft', sort: 0,
});
const form = reactive(emptyForm());

const openForm = async (record) => {
  Object.assign(form, emptyForm());
  if (record) {
    // 拉详情拿正文
    const res = await getEntry(record.id);
    Object.assign(form, res.data);
  }
  formVisible.value = true;
};

const save = async () => {
  if (!form.title || !form.categoryId) {
    Message.warning('标题与分类必填');
    return false;
  }
  saving.value = true;
  try {
    if (form.id) {
      await updateEntry(form.id, form);
      Message.success('已更新');
    } else {
      await addEntry(form);
      Message.success('已创建');
    }
    formVisible.value = false;
    load();
  } finally {
    saving.value = false;
  }
};

const onDelete = async (record) => {
  await removeEntry(record.id);
  Message.success('已删除');
  load();
};

onMounted(async () => {
  const res = await categoryTree();
  treeData.value = res.data || [];
  flattenCats(treeData.value);
  load();
});
</script>

<style lang="scss" scoped>
.wiki-admin {
  padding: 16px 24px;
}
.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
  .spacer {
    flex: 1;
  }
}
.form-row {
  display: flex;
  gap: 16px;
}
.md-input :deep(textarea) {
  font-family: 'SF Mono', Menlo, Consolas, monospace;
  font-size: 13px;
}
</style>
