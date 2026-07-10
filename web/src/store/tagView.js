import { defineStore } from 'pinia';
import { ref } from 'vue';
import { sessionStorage } from '@/utils/storage';

export const useTagViewStore = defineStore('tagView', () => {
    const visitedViews = ref(sessionStorage.getItem('tagList') || []);
    const cachedViews = ref([]);

    function addVisitedViews(view) {
        console.log(view.meta.title)

        // Akiraka 20240228 限制标签数量
        if (visitedViews.value.length >= 20) {
          visitedViews.value.splice(0,1)
        }

        if (visitedViews.value.some(v => v.path === view.path || v.meta.title === view.meta.title)) return;

        visitedViews.value.push(
            Object.assign({}, view, {
                title: view.meta.title || 'no-name'
            })
        )
        sessionStorage.setItem('tagList', visitedViews.value);
    }

    function delVisitedViews(view) {
        for (const [i, v] of visitedViews.value.entries()) {
            if (v.path === view.path) {
                visitedViews.value.splice(i, 1);
                break;
            }
        }
        sessionStorage.setItem('tagList', visitedViews.value);
    }

    function delAllVisitedViews(view) {
        visitedViews.value = visitedViews.value.filter(v => v.path === view.path);
        sessionStorage.setItem('tagList', visitedViews.value);
    }

    function delLeftVisitedViews(view) {
        for (const [i, v] of visitedViews.value.entries()) {
            if (v.path === view.path) {
                visitedViews.value = visitedViews.value.slice(i);
                break;
            }
        }
        sessionStorage.setItem('tagList', visitedViews.value);
    }
    function delRightVisitedViews(view) {
        for (const [i, v] of visitedViews.value.entries()) {
            if (v.path === view.path) {
                visitedViews.value = visitedViews.value.slice(0, i + 1);
                break;
            }
        }
        sessionStorage.setItem('tagList', visitedViews.value);
    }

    function addCachedViews(view) {
        if (cachedViews.value.includes(view.name)) return;
        if (!view.meta.noCache) cachedViews.value.push(view.name);
    }

    function delCachedViews(view) {
        visitedViews.value = visitedViews.value.filter(v => v !== view.name);
    }

    function delCachedVisitedViews() {
        cachedViews.value = [];
    }

    return {
        visitedViews,
        cachedViews,
        addVisitedViews,
        delVisitedViews,
        delAllVisitedViews,
        delLeftVisitedViews,
        delRightVisitedViews
    }
})