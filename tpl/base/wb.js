const root = document.documentElement;
// 切换主题函数
window.setTheme = function (isDark) {
    root.dataset.theme = isDark ? 'dark' : 'light';
}

// 根据时间判断
window.checkTimeAndSetTheme = function() {
    const now = new Date();
    const hour = now.getHours();
    const isDark = (hour >= 18 || hour < 7);
    setTheme(isDark);
}

// 初始化与定时检测
window.checkTimeAndSetTheme();
setInterval(window.checkTimeAndSetTheme, 60 * 1000); // 每分钟检查一次