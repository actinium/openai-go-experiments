const messages = [];
const input = document.querySelector('#chat-input');
const output = document.querySelector('#chat-messages');
const outputContainer = document.querySelector('.chat-messages-container');

input.focus();

function message(role, msg) {
    return {
        role: role,
        content: msg
    };
}

function messageDiv(msg) {
    const p = document.createElement('p');
    p.innerText = msg.content;

    const div = document.createElement('div');
    div.append(p);
    div.classList.add('msg');
    div.classList.add('msg-' + msg.role);

    return div;
}

async function send(msg) {
    const userMsg = message('user', msg);
    messages.push(userMsg);

    appendMessage(userMsg);
    await sleep(100);
    const msgDiv = appendLoader();

    const resp = await fetch('/chat', {
        method: 'POST',
        body: JSON.stringify(messages)
    });
    msg = await resp.json()
    messages.push(msg);

    replaceLoader(msgDiv, msg);
}

function appendMessage(msg) {
    msgDiv = messageDiv(msg);
    output.append(msgDiv);
    scrollToBottom();
    msgDiv.addEventListener('animationend', () => {
        scrollToBottom();
    });
}

function appendLoader() {
    const img = document.createElement('img');
    img.src = '/assets/loading.svg';
    img.classList.add('loading');

    const div = document.createElement('div');
    div.append(img);
    div.classList.add('msg');
    div.classList.add('msg-assistant');

    output.append(div);
    scrollToBottom();
    div.addEventListener('animationend', () => {
        scrollToBottom();
    });

    return div;
}

function replaceLoader(div, msg) {
    const p = document.createElement('p');
    p.innerText = msg.content;
    div.replaceChildren(p);
    scrollToBottom();
}

function scrollToBottom() {
    outputContainer.scrollTop = outputContainer.scrollHeight;
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

input.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        e.preventDefault();
        const msg = input.value.trim();
        if (msg == '') return;
        input.value = '';
        send(msg);
    }
});
