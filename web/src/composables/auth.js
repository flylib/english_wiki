import { computed } from 'vue';
import { useUserStore } from '@/store/userInfo';

// 当前用户是否超级管理员(go-admin admin 角色)
export function useIsSuperAdmin() {
  const store = useUserStore();
  return computed(() => Array.isArray(store.roles) && store.roles.includes('admin'));
}
