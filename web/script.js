// 会话管理前端逻辑
let waitingForAIResponse = false;
let currentSessionId = null;
let sessions = [];
let currentSession = null;

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function () {
    loadSessions();
    createNewSession();
});

// 加载会话列表
async function loadSessions() {
    try {
        const response = await fetch('/api/sessions');
        if (!response.ok) throw new Error('HTTP ' + response.status);

        const data = await response.json();
        sessions = data.sessions || [];
        renderSessionsList();
    } catch (error) {
        console.error('加载会话列表失败:', error);
    }
}

// 渲染会话列表
function renderSessionsList() {
    const sessionsList = document.getElementById('sessions-list');
    sessionsList.innerHTML = '';

    sessions.forEach(session => {
        const sessionItem = document.createElement('div');
        sessionItem.className = 'session-item';
        if (session.id === currentSessionId) {
            sessionItem.classList.add('active');
        }

        const messageCount = session.message_count || 0;
        const updatedAt = new Date(session.updated_at).toLocaleString('zh-CN', {
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });

        sessionItem.innerHTML = `
            <div class="session-name">${session.name}</div>
            <div class="session-info">
                <span>${messageCount} 条消息</span>
                <span>${updatedAt}</span>
            </div>
            <div class="session-actions">
                <button class="action-btn rename" onclick="showRenameModal('${session.id}', '${session.name}')" title="重命名">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="action-btn delete" onclick="showDeleteModal('${session.id}')" title="删除">
                    <i class="fas fa-trash"></i>
                </button>
            </div>
        `;

        sessionItem.onclick = (e) => {
            // 如果点击的是操作按钮，不切换会话
            if (e.target.closest('.session-actions')) {
                return;
            }
            switchToSession(session.id);
        };

        sessionsList.appendChild(sessionItem);
    });
}

// 创建新会话
async function createNewSession() {
    try {
        const sessionName = '新对话 ' + new Date().toLocaleString('zh-CN', {
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });

        const response = await fetch('/api/sessions', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: sessionName })
        });

        if (!response.ok) throw new Error('HTTP ' + response.status);

        const data = await response.json();
        const newSession = data.session;

        // 添加到会话列表
        sessions.unshift(newSession);
        renderSessionsList();

        // 切换到新会话
        switchToSession(newSession.id);

        // 清空聊天框
        clearChatBox();

    } catch (error) {
        console.error('创建新会话失败:', error);
        alert('创建新会话失败: ' + error.message);
    }
}

// 切换到指定会话
async function switchToSession(sessionId) {
    try {
        // 获取会话信息
        const response = await fetch(`/api/sessions/${sessionId}`);
        if (!response.ok) throw new Error('HTTP ' + response.status);

        const data = await response.json();
        currentSession = data.session;
        currentSessionId = sessionId;

        // 更新UI
        renderSessionsList();

        // 加载会话消息
        await loadSessionMessages(sessionId);

        // 更新页面标题
        document.title = `${currentSession.name} - AI 聊天助手`;

    } catch (error) {
        console.error('切换会话失败:', error);
        alert('切换会话失败: ' + error.message);
    }
}

// 加载会话消息
async function loadSessionMessages(sessionId) {
    try {
        const response = await fetch(`/api/sessions/${sessionId}/messages`);
        if (!response.ok) throw new Error('HTTP ' + response.status);

        const data = await response.json();
        const messages = data.messages || [];

        // 清空聊天框
        clearChatBox();

        // 显示消息
        messages.forEach(message => {
            if (message.role === 'user') {
                addMessageToChat('你: ' + message.content, 'user');
            } else if (message.role === 'assistant') {
                addMessageToChat('AI: ' + message.content, 'ai');
            }
        });

        scrollToBottom();

    } catch (error) {
        console.error('加载会话消息失败:', error);
    }
}

// 发送消息
async function sendMessage() {
    if (waitingForAIResponse) return;

    const inputElement = document.getElementById("input");
    const roleSelect = document.getElementById("role-select");
    const message = inputElement.value.trim();
    const role = roleSelect ? roleSelect.value : "general";
    const modeSelect = document.getElementById("mode-select");
    const namespaceInput = document.getElementById("namespace-input");
    const topkInput = document.getElementById("topk-input");
    const debugCheckbox = document.getElementById("debug-checkbox");

    const mode = modeSelect ? modeSelect.value : "rag";
    const namespace = namespaceInput ? namespaceInput.value.trim() : "";
    const topK = topkInput ? parseInt(topkInput.value, 10) || 3 : 3;
    const debug = debugCheckbox ? debugCheckbox.checked : false;

    if (!message) return;
    if (!currentSessionId) {
        alert('请先创建或选择一个会话');
        return;
    }

    inputElement.value = "";

    // 添加用户消息到聊天框
    addMessageToChat("你: " + message, 'user');

    // AI 占位
    const aiEl = addMessageToChat("AI: 正在输入...", 'ai');
    aiEl.classList.add('typing');
    waitingForAIResponse = true;
    scrollToBottom();

    try {
        let url = "/chat";
        let payload = {
            message,
            role,
            session_id: currentSessionId
        };

        // 如果选择 RAG 模式，则走 /rag/chat
        if (mode === "rag") {
            url = "/rag/chat";
            payload = {
                query: message,
                mode: "rag",
                namespace: namespace || undefined,
                top_k: topK,
                debug: debug
            };
        } else if (mode === "normal") {
            // 明确 normal，还是走 /chat，但便于后续扩展
            url = "/chat";
        }

        const response = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        if (!response.ok) throw new Error("HTTP " + response.status);

        const data = await response.json();

        // 移除打字效果
        aiEl.classList.remove('typing');

        // 更新AI回复
        aiEl.textContent = "AI: ";
        const answer = data.reply || data.answer || "出错了，请稍后再试";
        typeText(aiEl, answer);

        // 如果是 RAG 模式且开启 debug，附带命中文档信息
        if (mode === "rag" && debug && data.hit_docs && Array.isArray(data.hit_docs) && data.hit_docs.length > 0) {
            addMessageToChat("命中文档: " + data.hit_docs.join(" | "), "system-message");
        }

        // 显示兜底信息
        if (mode === "rag" && data.fallback) {
            addMessageToChat("提示：知识库无命中，已退化为普通对话。", "system-message");
        }

        // 刷新会话列表
        await loadSessions();

    } catch (error) {
        console.error('发送消息失败:', error);
        aiEl.textContent = "AI: 出错了，请稍后再试";
        aiEl.classList.remove('typing');
        waitingForAIResponse = false;
    }
}

// 知识入库
async function ingestKnowledge() {
    const titleEl = document.getElementById("ingest-title");
    const contentEl = document.getElementById("ingest-content");
    const sourceEl = document.getElementById("ingest-source");
    const nsEl = document.getElementById("ingest-namespace");

    const title = titleEl.value.trim();
    const content = contentEl.value.trim();
    const source = (sourceEl.value || "manual").trim();
    const namespace = (nsEl.value || "default").trim();

    if (!title || !content) {
        alert("请输入标题和内容");
        return;
    }

    try {
        const resp = await fetch("/rag/knowledge", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                title,
                content,
                source,
                namespace
            })
        });

        if (!resp.ok) throw new Error("HTTP " + resp.status);

        const data = await resp.json();
        const chunks = data.chunks || (data.knowledges ? data.knowledges.length : 0);
        alert(`入库成功，片段数：${chunks}`);

        // 清空表单
        titleEl.value = "";
        contentEl.value = "";
        sourceEl.value = "";
        nsEl.value = "";
    } catch (e) {
        console.error("知识入库失败:", e);
        alert("知识入库失败: " + e.message);
    }
}

// 添加消息到聊天框
function addMessageToChat(content, type) {
    const chatBox = document.getElementById("chat-box");
    const messageEl = document.createElement("div");
    messageEl.className = `message ${type}`;
    messageEl.textContent = content;
    chatBox.appendChild(messageEl);
    return messageEl;
}

// 清空聊天框
function clearChatBox() {
    const chatBox = document.getElementById("chat-box");
    chatBox.innerHTML = '';
}

// 打字效果
function typeText(element, text) {
    let i = 0;
    const prefix = "AI: ";
    (function type() {
        if (i < text.length) {
            element.textContent = prefix + text.substring(0, i + 1);
            i++;
            scrollToBottom();
            setTimeout(type, 30);
        } else {
            waitingForAIResponse = false;
        }
    })();
}

// 滚动到底部
function scrollToBottom() {
    const chatBox = document.getElementById("chat-box");
    chatBox.scrollTop = chatBox.scrollHeight;
}

// 显示重命名模态框
function showRenameModal(sessionId, currentName) {
    const modal = document.getElementById('renameModal');
    const input = document.getElementById('renameInput');
    input.value = currentName;
    input.dataset.sessionId = sessionId;
    modal.style.display = 'block';
    input.focus();
}

// 显示删除模态框
function showDeleteModal(sessionId) {
    const modal = document.getElementById('deleteModal');
    modal.dataset.sessionId = sessionId;
    modal.style.display = 'block';
}

// 关闭模态框
function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

// 确认重命名
async function confirmRename() {
    const input = document.getElementById('renameInput');
    const sessionId = input.dataset.sessionId;
    const newName = input.value.trim();

    if (!newName) {
        alert('请输入会话名称');
        return;
    }

    try {
        const response = await fetch(`/api/sessions/${sessionId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: newName })
        });

        if (!response.ok) throw new Error('HTTP ' + response.status);

        // 更新本地数据
        const session = sessions.find(s => s.id === sessionId);
        if (session) {
            session.name = newName;
            session.updated_at = new Date().toISOString();
        }

        // 如果当前会话被重命名，更新页面标题
        if (sessionId === currentSessionId) {
            document.title = `${newName} - AI 聊天助手`;
        }

        // 刷新会话列表
        renderSessionsList();

        // 关闭模态框
        closeModal('renameModal');

    } catch (error) {
        console.error('重命名失败:', error);
        alert('重命名失败: ' + error.message);
    }
}

// 确认删除
async function confirmDelete() {
    const modal = document.getElementById('deleteModal');
    const sessionId = modal.dataset.sessionId;

    try {
        const response = await fetch(`/api/sessions/${sessionId}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('HTTP ' + response.status);

        // 从本地数据中移除
        sessions = sessions.filter(s => s.id !== sessionId);

        // 如果删除的是当前会话，创建新会话
        if (sessionId === currentSessionId) {
            createNewSession();
        } else {
            renderSessionsList();
        }

        // 关闭模态框
        closeModal('deleteModal');

    } catch (error) {
        console.error('删除失败:', error);
        alert('删除失败: ' + error.message);
    }
}

// 切换侧边栏（移动端）
function toggleSidebar() {
    const sidebar = document.getElementById('sidebar');
    sidebar.classList.toggle('open');
}

// 点击模态框外部关闭
window.onclick = function (event) {
    if (event.target.classList.contains('modal')) {
        event.target.style.display = 'none';
    }
}

// 按ESC键关闭模态框
document.addEventListener('keydown', function (event) {
    if (event.key === 'Escape') {
        const modals = document.querySelectorAll('.modal');
        modals.forEach(modal => {
            if (modal.style.display === 'block') {
                modal.style.display = 'none';
            }
        });
    }
});