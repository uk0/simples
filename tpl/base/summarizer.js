(function () {
    const tag = document.currentScript;
    const WORKER_URL = (tag && tag.dataset.worker) || 'http://localhost/embed/summarizer';
    const CSS_URL    = (tag && tag.dataset.css)    || '/summarizer.css';
    const IMAGES_MODE = ((tag && tag.dataset.images) || 'base64').toLowerCase(); // 'base64' | 'url'

    /* 注入 CSS（若未加载） */
    function ensureCSS(href) {
        if ([...document.styleSheets].some(s => s.href && s.href.endsWith(href))) return;
        if ([...document.querySelectorAll('link[rel="stylesheet"]')].some(l => l.href.endsWith(href))) return;
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = href;
        document.head.appendChild(link);
    }
    ensureCSS(CSS_URL);

    /* 注入 Readability（若未加载） */
    function ensureReadability() {
        return new Promise((resolve) => {
            if (window.Readability) return resolve();
            const s = document.createElement('script');
            s.src = 'https://cdn.jsdelivr.net/npm/@mozilla/readability@0.5.0/Readability.min.js';
            s.onload = () => resolve();
            s.onerror = () => resolve();
            document.head.appendChild(s);
        });
    }

    /* DOM：右侧抽屉 + 召唤/关闭 + “开始总结”按钮（父页侧） */
    const drawer = document.createElement('div');
    drawer.id = 'summarizer-drawer';
    drawer.className = 'summarizer-drawer glass';
    drawer.setAttribute('aria-hidden', 'true');
    drawer.innerHTML = `
    <div class="summarizer-drawer__header">
      <div class="actions">
        <button id="summarizer-start" class="btn small" type="button">发送文章</button>
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
    handle.setAttribute('aria-controls', 'summarizer-drawer');
    handle.textContent = '🤔';
    document.body.appendChild(handle);

    const iframe  = drawer.querySelector('#summarizer-iframe');
    const closeBtn= drawer.querySelector('#summarizer-close');
    const startBtn= drawer.querySelector('#summarizer-start');

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
        if (drawer.classList.contains('open') && !drawer.contains(e.target) && e.target !== handle) {
            closeDrawer();
        }
    });

    /* 正文抽取 */
    function extractMainContent() {
        try {
            const cloned = document.cloneNode(true);
            if (window.Readability) {
                const reader = new Readability(cloned);
                const article = reader.parse();
                if (article && article.textContent) {
                    const title = article.title || document.title || '';
                    const text  = `标题：${title}\n\n${article.textContent}`;
                    return (text || '').trim().slice(0, 60000);
                }
            }
        } catch (_) {}
        return (document.body.innerText || '').slice(0, 60000);
    }

    /* 图片（URL 模式） */
    function collectImagesURL(limit = 8) {
        const imgs = Array.from(document.images || []);
        const out = [];
        for (const img of imgs) {
            if (!img.src) continue;
            const w = img.naturalWidth || 0, h = img.naturalHeight || 0;
            if (w < 128 || h < 128) continue;
            const abs = new URL(img.src, location.href).href;
            const rec = { src: abs };
            if (img.alt)   rec.alt   = img.alt;
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

    /* 图片转 Base64（跨源失败即跳过） */
    async function imageToBase64(url) {
        try {
            const res = await fetch(url, { mode: 'cors' });
            if (!res.ok) throw new Error('fetch fail');
            const blob = await res.blob();
            const buf  = await blob.arrayBuffer();
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

    /* 与 iframe postMessage 握手 */
    let iframeReady = false;

    async function buildPayload(extra = '') {
        await ensureReadability();
        const text = extractMainContent();
        let images = [];
        if (IMAGES_MODE === 'base64') images = await collectImagesBase64(8);
        else images = collectImagesURL(8);
        return { text, images, extra };
    }

    async function sendPageToIframe(extra = '', force = false) {
        if (!iframe || !iframe.contentWindow) return;
        if (!iframeReady && !force) return; // 正常路径需 ready；可 force 兜底

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

    // 点击“开始总结”时才发送
    startBtn.addEventListener('click', () => {
        // 优先 ping 一下，避免 ready 丢失
        try { iframe.contentWindow.postMessage({ type: 'qwq-ping' }, new URL(WORKER_URL).origin); } catch {}
        // 稍等 300ms 再发送（无论 ready 与否，force 兜底）
        setTimeout(() => sendPageToIframe('', /*force*/true), 300);
    });

    // 仅记录 ready，不自动打开、不自动发送
    window.addEventListener('message', (ev) => {
        if (!ev || !ev.data) return;
        const okOrigin = new URL(WORKER_URL).origin;
        if (ev.origin !== okOrigin) return;
        if (ev.data.type === 'qwq-ready') {
            console.log('[Summarizer parent] got qwq-ready from iframe');
            iframeReady = true;
        }
    });

    // iframe load 之后只做一次 ping，**不**自动打开或发送
    iframe.addEventListener('load', () => {
        try { iframe.contentWindow.postMessage({ type: 'qwq-ping' }, new URL(WORKER_URL).origin); } catch {}
    });
})();