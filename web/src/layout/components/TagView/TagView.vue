<script setup>
import { watch, ref, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useTagViewStore } from '@/store/tagView';
import { onMounted } from 'vue';

const store = useTagViewStore();
const route = useRoute();
const router = useRouter();

const tagListRef = ref(null);
const tagMore = ref(false);

const checkTagListWidth = async () => {
  await nextTick();
  // 获取 tagList 宽度
  const tagListWidth = tagListRef.value.offsetWidth;
  // 获取 tagList 子元素实际宽度
  const tagItemWidth = tagListRef.value.children[0].scrollWidth;
  // 判断 tagList 宽度是否小于 tagItem 实际宽度
  if (tagListWidth < tagItemWidth) {
    // 设置 tagMore 为 true
    tagMore.value = true;
    // 计算向右偏移宽度
    const rightOffset = tagListRef.value.children[0].offsetWidth - tagListRef.value.clientWidth;
    tagListRef.value.children[0].style.right = rightOffset + 'px';
  } else {
    tagMore.value = false;
  }
}

// Tag 标签上一页
const handleTagPrev = () => {
  let leftOffset = tagListRef.value.clientWidth - tagListRef.value.children[0].offsetWidth;

  if (leftOffset < 0) {
    leftOffset = 0;
  }
  tagListRef.value.children[0].style.right = leftOffset + 'px';
}

// Tag 标签下一页
const handleTagNext = () => {
  // 计算向右偏移宽度
  const rightOffset = tagListRef.value.children[0].offsetWidth - tagListRef.value.clientWidth;
  tagListRef.value.children[0].style.right = rightOffset + 'px';
}

// Tag 关闭事件
const handleTagClose = (view, index) => {
  // 判断是否为当前选中
  if (view.path === route.path) {
    // 判断是否为第一个
    if (index !== 0 ) {
      router.push(store.visitedViews[index - 1].path);
    } else {
      router.push(store.visitedViews[index + 1].path)
    }
  }
  store.delVisitedViews(view);
}

const handleDelTag = (v) => {
  switch (v) {
    case 'left':
      store.delLeftVisitedViews(route);
      break;
    case 'right':
      store.delRightVisitedViews(route);
      break;
    case 'all':
      store.delAllVisitedViews(route);
      break;
    default:
      console.log('参数错误')
  }
}

// 监听 store 数据变化
watch(store.visitedViews, () => {
  // 重新计算 tagList 宽度
  checkTagListWidth();
})

// 监听路由变化
watch(() => route.path, () => {
  store.addVisitedViews(route)
}, { immediate: true })

onMounted(() => {
  checkTagListWidth()
})

</script>

<template>
  <div class="tag-view-main">
    <div
      class="tag-card"
      >
      <div class="tag-prev" v-if="tagMore" @click="handleTagPrev">
        <icon-left :size="18"/>
      </div>

      <div class="tag-list" ref="tagListRef">
          <transition-group name="tag-fade">
            <a-tag
              v-for="(view, index) in store.visitedViews"
              :key="view.name"
              :checked="view.meta.title === route.meta.title"
              :color="view.meta.title === route.meta.title ? '#5d87ff' : ''"
              @close="handleTagClose(view, index)"
              :closable="store.visitedViews.length > 1"
              checkable
              :style="{ marginRight: '5px' }"
              >
              <router-link :to="view.path" class="tag-link">
                  {{ view.title }}
              </router-link>
            </a-tag>
          </transition-group>
          <!-- <a-tag
            v-for="name, index in tagListCount"
            :key="index"
            @close="() => tagListCount -= 1"
            closable>
            标题{{ name }}
          </a-tag> -->

      </div>

      <div class="tag-next" v-if="tagMore" @click="handleTagNext">
          <icon-right :size="18"/>
      </div>
    </div>

<!--    <div class="tag-close">-->
<!--      <a-dropdown @select="handleDelTag">-->
<!--        <a-button>-->
<!--          <template #icon>-->
<!--            <icon-down />-->
<!--          </template>-->
<!--        </a-button>-->
<!--        <template #content>-->
<!--          <a-doption value="left">关闭左侧</a-doption>-->
<!--          <a-doption value="right">关闭右侧</a-doption>-->
<!--          <a-doption value="all">关闭全部</a-doption>-->
<!--        </template>-->
<!--      </a-dropdown>-->
<!--    </div>-->
  </div>
</template>

<style lang="scss" scoped>
.tag-view-main {
  padding: 0 15px;
  display: flex;
  justify-content: space-between;
  max-width: 100%;
  min-width: 100%;
  background-color: var(--color-bg-2);
  .tag-card {
    display: flex;
    align-items: center;
    height: 36px;
    overflow: hidden;
    .tag-list {
      position: relative;
      overflow-x: hidden;
      transition: .3s ease-in-out;
    }
    .tag-prev, .tag-next{
      width: 48px;
      cursor: pointer;
      text-align: center;
    }
  }

}

.tag-link {
  color: inherit;
}
.tag-link:hover {
  color: inherit;
}

.tag-fade-enter-active,
.tag-fade-leave-active {
  transition: opacity 0.3s ease-in-out;
}

.tag-fade-enter-from,
.tag-fade-leave-to {
  opacity: 0;
}
</style>