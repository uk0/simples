/* comment-iframe-loader.js
   功能：
   - 懒加载评论 iframe（IntersectionObserver 或可视即加载）
   - 自动在 src 后拼接 theme 参数（优先使用用户偏好 localStorage，否则按时间判断）
   - 监听 <html data-theme> 变化并同步 iframe（会更新 src）
   - 提供全局 API: window.reloadCommentIframe(theme?)
*/

(function () {
    'use strict';

    const CONTAINER_ID = 'comment-widget';
    const PLACEHOLDER_ID = 'comment-iframe-placeholder';
    const IFRAME_ID = 'commentIframe';
    const THEME_KEY = 'site_theme_pref';
    const AUTO_DARK_START = 18; // 18:00
    const AUTO_DARK_END = 7;    // 07:00

    const container = document.getElementById(CONTAINER_ID);
    if (!container) return console.warn('[comment] container not found');

    const placeholder = document.getElementById(PLACEHOLDER_ID);
    if (!placeholder) return console.warn('[comment] placeholder not found');

    const slug = container.getAttribute('data-slug') || '';
    const BASE_EMBED = 'https://cfmmmmmmmmmmmmm.wwwneo.com/embed/area/' + encodeURIComponent(slug);

    // 计算主题：user pref > html[data-theme] > time-based
    function getPreferredTheme() {
        try {
            const pref = localStorage.getItem(THEME_KEY);
            if (pref === 'light' || pref === 'dark') return pref;
        } catch (e) { /* ignore */ }

        // If root dataset has theme set by other scripts, prefer that
        const rootTheme = document.documentElement.dataset.theme;
        if (rootTheme === 'light' || rootTheme === 'dark') return rootTheme;

        // fallback to time-based
        const h = new Date().getHours();
        return (h >= AUTO_DARK_START || h < AUTO_DARK_END) ? 'dark' : 'light';
    }

    function buildUrlWithTheme(base, theme) {
        try {
            const url = new URL(base, window.location.origin);
            // replace or set theme param
            url.searchParams.set('theme', theme);
            return url.toString();
        } catch (e) {
            // fallback: naive append
            let clean = base.replace(/[?&]theme=(light|dark)/i, '');
            clean = clean.replace(/\?$/,'');
            return clean + (clean.indexOf('?') === -1 ? '?' : '&') + 'theme=' + encodeURIComponent(theme);
        }
    }

    // Create iframe and append into placeholder
    function ensureIframe(theme) {
        const existing = document.getElementById(IFRAME_ID);
        const src = buildUrlWithTheme(BASE_EMBED, theme);

        if (existing) {
            // Only update src if different
            if (existing.src !== src) {
                // debounce quick changes
                if (existing._updateTimer) clearTimeout(existing._updateTimer);
                existing._updateTimer = setTimeout(() => { existing.src = src; }, 60);
            }
            return existing;
        }

        const iframe = document.createElement('iframe');
        iframe.id = IFRAME_ID;
        iframe.className = 'comment-iframe';
        iframe.name = '评论区';
        iframe.loading = 'lazy';
        iframe.allow = 'clipboard-read; clipboard-write; encrypted-media; fullscreen';
        iframe.sandbox = 'allow-scripts allow-same-origin allow-forms allow-popups';
        iframe.title = '评论区';
        iframe.src = src;
        iframe.style.border = '0';
        iframe.style.width = '100%';
        iframe.style.height = '100%';

        // replace placeholder content with iframe
        placeholder.innerHTML = '';
        placeholder.appendChild(iframe);

        // attach simple error handling to show message if iframe fails to load
        iframe.addEventListener('error', function () {
            placeholder.innerHTML = '<div class="comment-error">评论区加载失败，请稍后重试</div>';
        });

        return iframe;
    }

    // Lazy load strategies:
    // 1) IntersectionObserver: load when placeholder enters viewport
    // 2) Fallback: load after small delay
    function initLazyLoad() {
        const theme = getPreferredTheme();

        if ('IntersectionObserver' in window) {
            const io = new IntersectionObserver((entries, obs) => {
                entries.forEach(e => {
                    if (e.isIntersecting) {
                        ensureIframe(theme);
                        obs.disconnect();
                    }
                });
            }, { root: null, threshold: 0.12 });
            io.observe(placeholder);

            // Also load after 10s as fallback in case user never scrolls
            const fallbackTimer = setTimeout(() => {
                if (!document.getElementById(IFRAME_ID)) ensureIframe(getPreferredTheme());
            }, 10000);
            // if iframe loaded, clear timer
            const checkLoaded = setInterval(() => {
                if (document.getElementById(IFRAME_ID)) {
                    clearTimeout(fallbackTimer);
                    clearInterval(checkLoaded);
                }
            }, 500);

        } else {
            // legacy: small delay then insert
            setTimeout(() => ensureIframe(theme), 600);
        }
    }

    // Watch for theme changes on root dataset (other scripts may toggle)
    function observeThemeChanges() {
        const root = document.documentElement;
        try {
            const mo = new MutationObserver(muts => {
                muts.forEach(m => {
                    if (m.attributeName === 'data-theme') {
                        const newTheme = root.dataset.theme === 'dark' ? 'dark' : 'light';
                        const iframe = document.getElementById(IFRAME_ID);
                        if (iframe) {
                            // update src to include theme param (reloads iframe)
                            ensureIframe(newTheme);
                        }
                    }
                });
            });
            mo.observe(root, { attributes: true });
        } catch (e) { /* ignore */ }
    }

    // Hook up manual toggle button and persist pref
    function initToggleButton() {
        const btn = document.getElementById('commentThemeToggle');
        if (!btn) return;
        btn.addEventListener('click', () => {
            const cur = getPreferredTheme();
            const next = cur === 'dark' ? 'light' : 'dark';
            try { localStorage.setItem(THEME_KEY, next); } catch (e) {}
            // update html dataset so site-wide scripts can also react
            document.documentElement.dataset.theme = next;
            // update iframe
            ensureIframe(next);
            btn.setAttribute('aria-pressed', String(next === 'dark'));
        });
        // set initial aria state
        const init = getPreferredTheme();
        btn.setAttribute('aria-pressed', String(init === 'dark'));
    }

    // Public API
    window.reloadCommentIframe = function (forceTheme) {
        const t = forceTheme || getPreferredTheme();
        ensureIframe(t);
    };

    window.loadCommentIframeNow = function () {
        ensureIframe(getPreferredTheme());
    };

    // Initialize
    initLazyLoad();
    observeThemeChanges();
    initToggleButton();

    // Also expose debug util
    console.log('[comment] lazy loader initialized');
})();