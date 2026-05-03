const WS_URL = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/chat`;
const RECONNECT_DELAY = 2000;

const messagesContainer = document.getElementById('messages');
const chatContainer = document.getElementById('chat-container');
const inputForm = document.getElementById('input-form');
const messageInput = document.getElementById('message-input');
const sendBtn = document.getElementById('send-btn');
const typingIndicator = document.getElementById('typing-indicator');
const statusDot = document.getElementById('status-dot');
const newChatBtn = document.getElementById('new-chat-btn');

let ws = null;
let reconnectTimeout = null;
let isWaitingForResponse = false;

function init() {
    connect();
    setupEventListeners();
    scrollToBottom();
}

function setupEventListeners() {
    inputForm.addEventListener('submit', (e) => {
        e.preventDefault();
        sendMessage();
    });

    messageInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });

    newChatBtn.addEventListener('click', () => {
        if (confirm('Start a new chat? Current history will be lost.')) {
            location.reload();
        }
    });

    messageInput.addEventListener('input', () => {
        messageInput.style.height = 'auto';
        messageInput.style.height = Math.min(messageInput.scrollHeight, 120) + 'px';
    });
}

function connect() {
    updateStatus('connecting');
    ws = new WebSocket(WS_URL);

    ws.onopen = () => {
        updateStatus('connected');
        hideWelcome();
    };

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            handleMessage(data);
        } catch (err) {
            console.error('Error parsing message:', err);
        }
    };

    ws.onclose = (event) => {
        updateStatus('disconnected');
        if (!event.wasClean) {
            reconnectTimeout = setTimeout(() => {
                connect();
            }, RECONNECT_DELAY);
        }
    };

    ws.onerror = (error) => {
        updateStatus('error');
    };
}

function disconnect() {
    if (ws) {
        ws.close(1000, 'User disconnected');
    }
    if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
    }
}

function handleMessage(data) {
    const text = data.message;

    hideTyping();
    isWaitingForResponse = false;

    addMessage(text, 'model');
    scrollToBottom();
}

function addMessage(text, sender) {
    hideWelcome();

    const messageEl = document.createElement('div');
    messageEl.className = `message ${sender}`;

    messageEl.textContent = text || '...';

    messagesContainer.appendChild(messageEl);
}

function showTyping() {
    typingIndicator.classList.remove('hidden');
    scrollToBottom();
}

function hideTyping() {
    typingIndicator.classList.add('hidden');
}

function hideWelcome() {
    const welcome = document.querySelector('.welcome-message');
    if (welcome) {
        welcome.style.display = 'none';
    }
}

function scrollToBottom() {
    chatContainer.scrollTop = chatContainer.scrollHeight;
}

function updateStatus(status) {
    statusDot.className = 'status-dot';
    switch (status) {
        case 'connected':
            statusDot.classList.add('connected');
            sendBtn.disabled = false;
            messageInput.disabled = false;
            break;
        case 'connecting':
        case 'disconnected':
        case 'error':
            sendBtn.disabled = true;
            messageInput.disabled = true;
            break;
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function sendMessage() {
    const text = messageInput.value.trim();
    if (!text || !ws || ws.readyState !== WebSocket.OPEN || isWaitingForResponse) {
        return;
    }

    addMessage(text, 'user');
    ws.send(JSON.stringify({ message: text }));
    messageInput.value = '';
    messageInput.style.height = 'auto';
    isWaitingForResponse = true;
    showTyping();
    messageInput.focus();
}

window.addEventListener('beforeunload', disconnect);
document.addEventListener('DOMContentLoaded', init);