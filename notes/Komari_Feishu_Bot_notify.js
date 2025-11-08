// ======== 配置 ========
const WEBHOOK = "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxxxxxxxxxxxxxx";

// ======== 核心函数 ========

/**
 * 发送文本消息到飞书 Webhook
 * @param {string} message - 消息内容
 * @param {string} title - 消息标题
 * @returns {boolean} 发送结果（同步返回）
 */
function sendMessage(message, title) {
    if (!WEBHOOK || WEBHOOK === 'REPLACE_WITH_WEBHOOK') {
        console.error('飞书 Webhook URL 未设置。消息未发送。');
        return false;
    }

    const url = WEBHOOK;
    const payload = {
        msg_type: "text",
        content: {
            text: title + "\n\n" + message
        }
    };

    try {
        // 使用同步 XMLHttpRequest（适配 goja 等 Go JS 引擎）
        const xhr = new XMLHttpRequest();
        xhr.open('POST', url, false); // false = 同步模式
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.send(JSON.stringify(payload));

        if (xhr.status >= 200 && xhr.status < 300) {
            const result = JSON.parse(xhr.responseText);
            if (result.code !== 0) {
                console.error('Feishu API error:', result.msg);
                return false;
            }
            return true;
        } else {
            console.error('Failed to send message:', xhr.status, xhr.statusText);
            return false;
        }
    } catch (error) {
        console.error('发送消息时出错:', error);
        return false;
    }
}

/**
 * 格式化并发送 Komari 事件通知
 * @param {Object} event - 事件对象
 * @returns {boolean} 发送结果
 */
function sendEvent(event) {
    try {
        // 格式化时间
        const formatTime = function(timeStr) {
            const date = new Date(timeStr);
            return date.toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                timeZone: 'Asia/Shanghai'
            });
        };

        // 格式化文件大小
        const formatBytes = function(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        };

        // 获取事件类型描述
        const getEventTypeDesc = function(eventType) {
            const eventMap = {
                'Offline': '❌ 服务器离线',
                'Online': '✅ 服务器上线',
                'Alert': '⚠️ 监控告警',
                'Renew': '⏰ 服务器已自动续费',
                'Expire': '🚨 服务到期提醒',
                'Test': '🧪 测试通知'
            };
            return eventMap[eventType] || ('📊 ' + eventType);
        };

        // 生成简洁的服务器摘要
        const generateClientSummary = function(client) {
            const parts = [];
            parts.push(client.name || 'Unknown');
            if (client.region) parts.push(client.region);
            return parts.join(' • ');
        };

        const title = getEventTypeDesc(event.event) + ' Komari 通知';
        let message = '';

        message += '时间: ' + formatTime(event.time) + '\n';

        if (event.message && event.message.trim()) {
            message += '说明: ' + event.message + '\n';
        }

        message += '\n';

        if (event.clients && event.clients.length > 0) {
            if (event.clients.length === 1) {
                message += generateClientSummary(event.clients[0]);
            } else {
                message += '影响服务器: ' + event.clients.length + ' 台\n';
                const shown = event.clients.length;
                for (let i = 0; i < shown; i++) {
                    message += (i + 1) + '. ' + generateClientSummary(event.clients[i]) + '\n';
                }
            }
        } else {
            message += '无关联服务器信息';
        }

        // 发送通知
        const success = sendMessage(message, title);
        if (success) {
            console.log('事件通知已发送: ' + event.event);
        } else {
            console.error('事件通知发送失败: ' + event.event);
        }
        return success;

    } catch (error) {
        console.error('发送事件通知时出错:', error);

        // 发送简化的错误通知
        const fallbackMessage = (event.emoji || '') + ' ' + event.event + '\n' + (event.message || '');
        const fallbackTitle = 'Komari 通知';
        try {
            return sendMessage(fallbackMessage, fallbackTitle);
        } catch (fallbackError) {
            console.error('备用通知也失败:', fallbackError);
            return false;
        }
    }
}