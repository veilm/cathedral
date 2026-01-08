const convListEl = document.getElementById("convList");
const chatLogEl = document.getElementById("chatLog");
const chatForm = document.getElementById("chatForm");
const chatInput = document.getElementById("chatInput");
const memoryForm = document.getElementById("memoryForm");
const memoryInput = document.getElementById("memoryInput");
const memoryView = document.getElementById("memoryView");
const newConvBtn = document.getElementById("newConv");

let currentConvId = null;
const urlParams = new URLSearchParams(window.location.search);

function setConvInUrl(convId) {
  const params = new URLSearchParams(window.location.search);
  if (convId) {
    params.set("conv", convId);
  } else {
    params.delete("conv");
  }
  const query = params.toString();
  const nextUrl = query ? `${window.location.pathname}?${query}` : window.location.pathname;
  window.history.replaceState({}, "", nextUrl);
}

function getConvFromUrl() {
  return urlParams.get("conv");
}

async function fetchJSON(url, options = {}) {
  const res = await fetch(url, options);
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || res.statusText);
  }
  return res.json();
}

function renderMessages(messages) {
  chatLogEl.innerHTML = "";
  messages.forEach((msg) => {
    const wrapper = document.createElement("div");
    wrapper.className = "msg";

    const role = document.createElement("div");
    role.className = "role";
    role.textContent = msg.role;

    const body = document.createElement("div");
    body.textContent = msg.content;

    wrapper.appendChild(role);
    wrapper.appendChild(body);
    chatLogEl.appendChild(wrapper);
  });
  chatLogEl.scrollTop = chatLogEl.scrollHeight;
}

async function loadConversations() {
  const conversations = await fetchJSON("/api/conversations");
  convListEl.innerHTML = "";

  conversations.forEach((conv) => {
    const btn = document.createElement("button");
    btn.textContent = conv.id;
    btn.dataset.id = conv.id;
    if (conv.id === currentConvId) btn.classList.add("active");

    btn.addEventListener("click", async () => {
      currentConvId = conv.id;
      setConvInUrl(currentConvId);
      await loadConversation(conv.id);
      await loadConversations();
    });
    convListEl.appendChild(btn);
  });
}

async function loadConversation(id) {
  const convo = await fetchJSON(`/api/conversations/${id}`);
  renderMessages(convo.messages || []);
}

async function sendMessage(message) {
  if (!currentConvId) return;
  await fetchJSON(`/api/conversations/${currentConvId}/message`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ message }),
  });
  await loadConversation(currentConvId);
}

chatForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const message = chatInput.value.trim();
  if (!message) return;
  chatInput.value = "";
  await sendMessage(message);
});

memoryForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const title = memoryInput.value.trim();
  if (!title) return;
  const data = await fetchJSON(`/api/memory/read?title=${encodeURIComponent(title)}`);
  memoryView.textContent = `[${data.path}]\n\n${data.content}`;
});

newConvBtn.addEventListener("click", async () => {
  const conv = await fetchJSON("/api/conversations", { method: "POST" });
  currentConvId = conv.id;
  setConvInUrl(currentConvId);
  await loadConversations();
  await loadConversation(conv.id);
});

(async () => {
  const conversations = await fetchJSON("/api/conversations");
  if (conversations.length > 0) {
    const urlConv = getConvFromUrl();
    const found = urlConv
      ? conversations.find((conv) => conv.id === urlConv)
      : null;
    currentConvId = found ? found.id : conversations[0].id;
    setConvInUrl(currentConvId);
    await loadConversation(currentConvId);
  }
  await loadConversations();
})();
