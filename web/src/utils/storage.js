// LocalStorage

export const storage = {
  getKeys() {
    const keys = [];
    for (let i = 0; i < window.localStorage.length; i++) {
      keys.push(window.localStorage.key(i));
    }

    return keys;
  },
  setItem(key, val) {
    if (typeof key !== 'string') {
      key = key.toString();
    }

    if (key === undefined || key.trim().length === 0) throw new Error('key 参数不能为空或者undefined');

    window.localStorage.setItem(key, JSON.stringify(val));
  },
  getItem(key) {
    const val = window.localStorage.getItem(key);
    return JSON.parse(val);
  },
  clearAllKeys() {
    window.localStorage.clear();
  },
  removeItem(key) {
    window.localStorage.removeItem(key);
  }
}

export const sessionStorage = {
  getKeys() {
    const keys = [];
    for (let i = 0; i < window.sessionStorage.length; i++) {
      keys.push(window.sessionStorage.key(i));
    }

    return keys;
  },
  setItem(key, val) {
    if (typeof key !== 'string') {
      key = key.toString();
    }

    if (key === undefined || key.trim().length === 0) throw new Error('key 参数不能为空或者undefined');

    window.sessionStorage.setItem(key, JSON.stringify(val));
  },
  getItem(key) {
    const val = window.sessionStorage.getItem(key);
    return JSON.parse(val);
  },
  clearAllKeys() {
    window.sessionStorage.clear();
  },
  removeItem(key) {
    window.sessionStorage.removeItem(key);
  }
}