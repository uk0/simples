(function () {
    const tag = document.currentScript;
    const WORKER_URL = (tag && tag.dataset.worker) || 'http://localhost/embed/summarizer';
    const CSS_URL    = (tag && tag.dataset.css)    || '/summarizer.css';
    const IMAGES_MODE = ((tag && tag.dataset.images) || 'base64').toLowerCase(); // 'base64' | 'url'

    const root = document.documentElement;

    /* ---------- 样式注入 ---------- */
    function ensureCSS(href) {
        if ([...document.styleSheets].some(s => s.href && s.href.endsWith(href))) return;
        if ([...document.querySelectorAll('link[rel="stylesheet"]')].some(l => l.href.endsWith(href))) return;
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = href;
        document.head.appendChild(link);
    }
    ensureCSS(CSS_URL);

    /* ---------- Readability 依赖 ---------- */
    function ensureReadability() {
        return new Promise(resolve => {
            if (window.Readability) return resolve();
            const s = document.createElement('script');
            s.src = 'https://cdn.jsdelivr.net/npm/@mozilla/readability@0.5.0/Readability.min.js';
            s.onload = () => resolve();
            s.onerror = () => resolve();
            document.head.appendChild(s);
        });
    }

    /* ---------- 主题判定：对齐 comment-iframe-loader ---------- */
    function getPreferredTheme() {
        try {
            const attr = root.dataset.theme;
            if (attr === 'dark' || attr === 'light') return attr;
        } catch {}
        const hour = new Date().getHours();
        const isDark = (hour >= 18 || hour < 7);
        return isDark ? 'dark' : 'light';
    }

    /* ---------- DOM 元素 ---------- */
    const drawer = document.createElement('div');
    drawer.id = 'summarizer-drawer';
    drawer.className = 'summarizer-drawer glass';
    drawer.setAttribute('aria-hidden', 'true');

    const theme = getPreferredTheme();
    drawer.dataset.theme = theme;

    drawer.innerHTML = `
    <div class="summarizer-drawer__header">
      <div class="actions">
        <button id="summarizer-start" class="btn small" type="button">总结</button>
        <button id="summarizer-close" class="btn small" type="button">关闭</button>
      </div>
    </div>
    <div class="summarizer-drawer__content">
      <iframe id="summarizer-iframe" class="summarizer-iframe"
              src="${WORKER_URL}" title="Summarizer" loading="lazy"
              allow="autoplay; encrypted-media; fullscreen; clipboard-read; clipboard-write"
              referrerpolicy="no-referrer"
              style="border:0;width:100%;height:98%;border-radius:12px;"></iframe>
    </div>`;
    document.body.appendChild(drawer);

    const handle = document.createElement('button');
    handle.id = 'summarizer-handle';
    handle.className = 'summarizer-handle glass';
    handle.type = 'button';
    handle.textContent = '📝';
    handle.dataset.theme = theme;
    document.body.appendChild(handle);

    const iframe   = drawer.querySelector('#summarizer-iframe');
    const closeBtn = drawer.querySelector('#summarizer-close');
    const startBtn = drawer.querySelector('#summarizer-start');

    /* ---------- 开关逻辑 ---------- */
    function openDrawer() {
        drawer.classList.add('open');
        drawer.setAttribute('aria-hidden', 'false');
    }
    function closeDrawer() {
        drawer.classList.remove('open');
        drawer.setAttribute('aria-hidden', 'true');
    }
    handle.addEventListener('click', openDrawer);
    closeBtn.addEventListener('click', closeDrawer);
    document.addEventListener('keydown', (e) => { if (e.key === 'Escape') closeDrawer(); });
    document.addEventListener('click', (e) => {
        if (drawer.classList.contains('open') && !drawer.contains(e.target) && e.target !== handle)
            closeDrawer();
    });

    /* ---------- 正文抽取 ---------- */
    function extractMainContent() {
        try {
            const cloned = document.cloneNode(true);
            if (window.Readability) {
                const reader = new Readability(cloned);
                const article = reader.parse();
                if (article && article.textContent) {
                    const title = article.title || document.title || '';
                    const text  = `标题：${title}\n\n${article.textContent}`;
                    return text.trim().slice(0, 60000);
                }
            }
        } catch {}
        return (document.body.innerText || '').slice(0, 60000);
    }

    /* ---------- 图片收集 ---------- */
    function collectImagesURL(limit = 8) {
        const imgs = Array.from(document.images || []);
        const out = [];
        for (const img of imgs) {
            if (!img.src) continue;
            const w = img.naturalWidth || 0, h = img.naturalHeight || 0;
            if (w < 128 || h < 128) continue;
            const abs = new URL(img.src, location.href).href;
            const rec = { src: abs };
            if (img.alt) rec.alt = img.alt;
            if (img.title) rec.title = img.title;
            const fig = img.closest('figure');
            if (fig) {
                const cap = fig.querySelector('figcaption');
                if (cap) rec.caption = cap.innerText.trim();
            }
            out.push(rec);
            if (out.length >= limit) break;
        }
        return out;
    }

    async function imageToBase64(url) {
        try {
            const res = await fetch(url, { mode: 'cors' });
            if (!res.ok) throw new Error('fetch fail');
            const blob = await res.blob();
            const buf = await blob.arrayBuffer();
            let bin = '';
            const bytes = new Uint8Array(buf);
            for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i]);
            return btoa(bin);
        } catch {
            return null;
        }
    }

    async function collectImagesBase64(limit = 8) {
        const urls = collectImagesURL(limit).map(x => x.src);
        const out = [];
        for (const u of urls) {
            const b64 = await imageToBase64(u);
            if (b64) out.push(`data:image/png;base64,${b64}`);
            if (out.length >= limit) break;
        }
        return out;
    }

    /* ---------- postMessage 通信 ---------- */
    let iframeReady = false;

    async function buildPayload(extra = '') {
        await ensureReadability();
        const text = extractMainContent();
        const images = (IMAGES_MODE === 'base64') ? await collectImagesBase64(8) : collectImagesURL(8);
        return { text, images, extra };
    }

    async function sendPageToIframe(extra = '', force = false) {
        if (!iframe || !iframe.contentWindow) return;
        if (!iframeReady && !force) return;

        const payload = await buildPayload(extra);
        iframe.contentWindow.postMessage(
            { type: 'qwq-summarize', ...payload },
            new URL(WORKER_URL).origin
        );
        console.log('[Summarizer parent] posted qwq-summarize', {
            textLen: payload.text.length,
            imagesCount: Array.isArray(payload.images) ? payload.images.length : 0,
            mode: IMAGES_MODE
        });
    }

    startBtn.addEventListener('click', () => {
        try { iframe.contentWindow.postMessage({ type: 'qwq-ping' }, new URL(WORKER_URL).origin); } catch {}
        setTimeout(() => sendPageToIframe('', true), 300);
    });

    window.addEventListener('message', (ev) => {
        if (!ev || !ev.data) return;
        const okOrigin = new URL(WORKER_URL).origin;
        if (ev.origin !== okOrigin) return;
        if (ev.data.type === 'qwq-ready') {
            console.log('[Summarizer parent] got qwq-ready from iframe');
            iframeReady = true;
        }
    });

    iframe.addEventListener('load', () => {
        try { iframe.contentWindow.postMessage({ type: 'qwq-ping' }, new URL(WORKER_URL).origin); } catch {}
    });

    /* ---------- 监听全局主题变化 ---------- */
    try {
        const mo = new MutationObserver(muts => {
            for (const m of muts) {
                if (m.attributeName === 'data-theme') {
                    const newTheme = getPreferredTheme();
                    drawer.dataset.theme = newTheme;
                    handle.dataset.theme = newTheme;
                    console.log('[Summarizer] theme updated ->', newTheme);
                }
            }
        });
        mo.observe(root, { attributes: true });
    } catch (e) {
        console.warn('[Summarizer] Theme observer failed:', e);
    }
})();