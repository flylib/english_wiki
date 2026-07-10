import request from '@/utils/request';

const base = '/api/v1/wiki';

// ---- 分类 ----
export function categoryTree() {
  return request({ url: `${base}/category/tree`, method: 'get' });
}
export function addCategory(data) {
  return request({ url: `${base}/category`, method: 'post', data });
}
export function updateCategory(id, data) {
  return request({ url: `${base}/category/${id}`, method: 'put', data });
}
export function removeCategory(id) {
  return request({ url: `${base}/category/${id}`, method: 'delete' });
}

// ---- 条目 ----
export function listEntries(params) {
  return request({ url: `${base}/entry`, method: 'get', params });
}
export function getEntry(id) {
  return request({ url: `${base}/entry/${id}`, method: 'get' });
}
export function addEntry(data) {
  return request({ url: `${base}/entry`, method: 'post', data });
}
export function updateEntry(id, data) {
  return request({ url: `${base}/entry/${id}`, method: 'put', data });
}
export function removeEntry(id) {
  return request({ url: `${base}/entry/${id}`, method: 'delete' });
}

// ---- 等级体系 ----
export function levelSystems() {
  return request({ url: `${base}/level/systems`, method: 'get' });
}

// CEFR rank(1~7)→ 标签
export const CEFR_LABELS = ['', 'Pre-A1', 'A1', 'A2', 'B1', 'B2', 'C1', 'C2'];

export function cefrRangeText(min, max) {
  if (!min && !max) return '';
  if (min === max) return CEFR_LABELS[min] || '';
  return `${CEFR_LABELS[min] || ''}~${CEFR_LABELS[max] || ''}`;
}
