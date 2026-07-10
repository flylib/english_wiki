<template>
  <a-select v-model="dataModel" :options="IconSelect" :trigger-props="{ contentClass: 'iconSelect' }" allow-search allow-clear placeholder="请选择菜单图标">
    <template #label="{ data }">
      <a-space>
        <component :is="data.value" :width="size" :height="size" :size="size.toString()" />
        <span v-if="title">{{data?.value}}</span>
      </a-space>
    </template>
    <template #option="{ data }">
      <component :is="data.value" :width="size" :height="size" :size="size.toString()" />
    </template>
  </a-select>
</template>

<script setup lang="ts">
import { ref, computed, toRefs } from "vue";
import * as ArcoIconModules from "@arco-design/web-vue/es/icon/index.js";
import * as RemixIcon from "@remixicon/vue";

const props = defineProps({
  data: {
    type: String,
  },
  // Akiraka 20240709 是否显示图标名称
  title: {
    type: Boolean,
    default: true
  },
  // Akiraka 20240709 控制图标显示大小
  size: {
    type: [Number, String],
    default: 16
  }
});

const emits = defineEmits(['update:data'])

let dataModel = computed({
  get() {
    return props.data
  },
  set(val) {
    emits('update:data', val)
  }
})

const IconSelect = computed(() => {
  const iconList = [];

  for (let key in RemixIcon) {
    if (RemixIcon[key].name) {
      iconList.push(RemixIcon[key].name);
    }
  }

  for (let key in ArcoIconModules) {
    if (ArcoIconModules[key].name) {
      iconList.push(ArcoIconModules[key].name);
    }
  }
  return iconList;
});

const handleClose = () =>  {
  emits('Change', false);
}

</script>

<style lang="scss">
.iconSelect .arco-select-dropdown-list {
  display: flex;
  flex-wrap: wrap;
}

.iconSelect .arco-select-option {
  width: auto;
}
</style>

