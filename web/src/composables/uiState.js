import { ref } from 'vue';

// 左侧分类树收起状态(Navbar 按钮 与 阅读页共享),持久化到 localStorage
const KEY = 'wiki_tree_collapsed';
export const treeCollapsed = ref(localStorage.getItem(KEY) === '1');

export function toggleTree() {
  treeCollapsed.value = !treeCollapsed.value;
  try { localStorage.setItem(KEY, treeCollapsed.value ? '1' : '0'); } catch (e) { /* 忽略 */ }
}
