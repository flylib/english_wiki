// 给全局 Message.error 错误弹条追加「复制」按钮:点击复制完整错误详情。
// 在 main.js 引入一次,所有调用点(树/控制台/表单等)自动生效。
import { h } from 'vue';
import { Message } from '@arco-design/web-vue';
import { IconCopy } from '@arco-design/web-vue/es/icon';

const origError = Message.error.bind(Message);

Message.error = (config) => {
  const cfg = typeof config === 'string' ? { content: config } : { ...(config || {}) };
  const text = typeof cfg.content === 'string' ? cfg.content : '';
  if (!text) return origError(config); // 非纯文本内容不处理,原样透传

  if (cfg.duration == null) cfg.duration = 6000; // 错误多留一会儿,方便复制
  cfg.content = () => h(
    'span',
    { style: 'display:inline-flex;align-items:center;gap:8px;max-width:70vw;' },
    [
      h('span', { style: 'word-break:break-all;' }, text),
      h(IconCopy, {
        style: 'cursor:pointer;flex:none;opacity:.75;',
        onClick: (ev) => {
          ev.stopPropagation();
          navigator.clipboard.writeText(text)
            .then(() => Message.success('已复制错误详情'))
            .catch(() => {});
        },
      }),
    ],
  );
  return origError(cfg);
};
