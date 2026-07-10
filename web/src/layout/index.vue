<template>
  <a-layout :style="{ height: '100vh' }">
    <a-layout>
      <a-layout-header>
        <Navbar :collapsed="collapsed" @on-collapse="onCollapse" />
      </a-layout-header>
      <a-layout-content class="layout-content">
        <AppMain />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { ref } from 'vue';
import { AppMain, Navbar } from './components';

// 侧边栏由各页自带(阅读页的分类树),此处不再渲染 sys_menu 菜单。
const collapsed = ref(false);
const onCollapse = () => {
  collapsed.value = !collapsed.value;
};
</script>

<style lang="scss">
@import '../style/index.scss';
@import '../style/transition.scss';
@import '../style/dark-theme.scss';

.layout-content {
  padding: 0;
  background-color: #ffffff;
}
</style>
