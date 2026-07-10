<template>
  <div class="wiki-admin">
    <div class="filters">
      <div class="spacer"></div>
      <a-button type="primary" @click="openForm(null, 0)">
        <template #icon><icon-plus /></template>
        新建根分类
      </a-button>
    </div>

    <a-table
      :data="treeData"
      :loading="loading"
      :pagination="false"
      row-key="id"
      :default-expand-all-rows="true"
    >
      <template #columns>
        <a-table-column title="名称" data-index="name" />
        <a-table-column title="编码" data-index="code" :width="160" />
        <a-table-column title="图标" data-index="icon" :width="160" />
        <a-table-column title="条目数" data-index="entryCount" :width="90" />
        <a-table-column title="排序" data-index="sort" :width="80" />
        <a-table-column title="操作" :width="230">
          <template #cell="{ record }">
            <a-button size="mini" type="text" @click="openForm(null, record.id)">加子分类</a-button>
            <a-button size="mini" type="text" @click="openForm(record)">编辑</a-button>
            <a-popconfirm content="确认删除该分类?" type="warning" @ok="onDelete(record)">
              <a-button size="mini" type="text" status="danger">删除</a-button>
            </a-popconfirm>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <a-modal
      v-model:visible="formVisible"
      :title="form.id ? '编辑分类' : '新建分类'"
      :ok-loading="saving"
      @ok="save"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="父分类">
          <a-tree-select
            v-model="form.parentId"
            :data="parentOptions"
            :field-names="{ key: 'id', title: 'name', children: 'children', icon: '_icon' }"
            placeholder="不选 = 根分类"
            allow-clear
          />
        </a-form-item>
        <a-form-item label="名称" required>
          <a-input v-model="form.name" placeholder="如:英语时态" />
        </a-form-item>
        <a-form-item label="编码" required>
          <a-input v-model="form.code" placeholder="唯一英文编码,如 tense" />
        </a-form-item>
        <a-form-item label="图标">
          <a-input v-model="form.icon" placeholder="arco 图标名,如 icon-book(可空)" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model="form.sort" :min="0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { categoryTree, addCategory, updateCategory, removeCategory } from '@/api/wiki';

const treeData = ref([]);
const loading = ref(false);

const load = async () => {
  loading.value = true;
  try {
    const res = await categoryTree();
    treeData.value = res.data || [];
  } finally {
    loading.value = false;
  }
};

// 编辑时父分类候选要排除自己(避免自环;简单处理:过滤掉自己这棵子树)
const pruneSelf = (nodes, selfId) =>
  (nodes || [])
    .filter((n) => n.id !== selfId)
    .map((n) => ({ ...n, children: pruneSelf(n.children, selfId) }));
const parentOptions = computed(() => pruneSelf(treeData.value, form.id));

const formVisible = ref(false);
const saving = ref(false);
const form = reactive({ id: null, parentId: undefined, name: '', code: '', icon: '', sort: 0 });

const openForm = (record, parentId) => {
  if (record) {
    Object.assign(form, {
      id: record.id,
      parentId: record.parentId || undefined,
      name: record.name,
      code: record.code,
      icon: record.icon,
      sort: record.sort,
    });
  } else {
    Object.assign(form, { id: null, parentId: parentId || undefined, name: '', code: '', icon: '', sort: 0 });
  }
  formVisible.value = true;
};

const save = async () => {
  if (!form.name || !form.code) {
    Message.warning('名称与编码必填');
    return false;
  }
  saving.value = true;
  try {
    const data = { ...form, parentId: form.parentId || 0 };
    if (form.id) {
      await updateCategory(form.id, data);
      Message.success('已更新');
    } else {
      await addCategory(data);
      Message.success('已创建');
    }
    formVisible.value = false;
    load();
  } finally {
    saving.value = false;
  }
};

const onDelete = async (record) => {
  try {
    await removeCategory(record.id);
    Message.success('已删除');
    load();
  } catch (e) {
    /* request 拦截器已提示 */
  }
};

onMounted(load);
</script>

<style lang="scss" scoped>
.wiki-admin {
  padding: 16px 24px;
}
.filters {
  display: flex;
  margin-bottom: 16px;
  .spacer {
    flex: 1;
  }
}
</style>
