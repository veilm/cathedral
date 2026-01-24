const convListEl = document.getElementById("convList");
const chatLogEl = document.getElementById("chatLog");
const chatForm = document.getElementById("chatForm");
const chatInput = document.getElementById("chatInput");
const memoryForm = document.getElementById("memoryForm");
const memoryInput = document.getElementById("memoryInput");
const memoryView = document.getElementById("memoryView");
const newConvBtn = document.getElementById("newConv");
const openSettingsBtn = document.getElementById("openSettings");
const closeSettingsBtn = document.getElementById("closeSettings");
const settingsModal = document.getElementById("settingsModal");
const themeSelect = document.getElementById("themeSelect");
const settingsBackdrop = settingsModal?.querySelector("[data-close-settings]");

let currentConvId = null;
const urlParams = new URLSearchParams(window.location.search);
const THEME_KEY = "cathedral-theme";

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

function setTheme(theme) {
  document.documentElement.dataset.theme = theme;
  localStorage.setItem(THEME_KEY, theme);
}

function getSavedTheme() {
  return localStorage.getItem(THEME_KEY) || "monotone-light";
}

function openSettings() {
  settingsModal.classList.remove("hidden");
}

function closeSettings() {
  settingsModal.classList.add("hidden");
}

async function fetchJSON(url, options = {}) {
  const res = await fetch(url, options);
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || res.statusText);
  }
  return res.json();
}

function escapeHtml(unsafe) {
  if (unsafe === null || unsafe === undefined) return "";
  return unsafe
    .toString()
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

function parseMarkdownSimple(text) {
  const result = [];
  let i = 0;
  let inBold = false;
  let inItalic = false;
  let boldMarker = null;
  let italicMarker = null;

  while (i < text.length) {
    const char = text[i];
    const nextChar = i + 1 < text.length ? text[i + 1] : null;
    const prevChar = i > 0 ? text[i - 1] : null;

    if (char === "*" || char === "_") {
      let markerCount = 1;
      let j = i + 1;
      while (j < text.length && text[j] === char) {
        markerCount++;
        j++;
      }

      if (markerCount >= 3) {
        if (inBold && inItalic && boldMarker === char && italicMarker === char) {
          result.push("</em></strong>");
          inBold = false;
          inItalic = false;
          boldMarker = null;
          italicMarker = null;
          i += 3;
          continue;
        } else if (!inBold && !inItalic) {
          result.push("<strong><em>");
          inBold = true;
          inItalic = true;
          boldMarker = char;
          italicMarker = char;
          i += 3;
          continue;
        }
      }

      if (markerCount >= 2) {
        if (inBold && boldMarker === char) {
          result.push("</strong>");
          inBold = false;
          boldMarker = null;
          i += 2;
          continue;
        } else if (!inBold) {
          result.push("<strong>");
          inBold = true;
          boldMarker = char;
          i += 2;
          continue;
        }
      }

      if (markerCount === 1) {
        if (char === "_") {
          const atStart = i === 0 || /\s/.test(prevChar);
          const nextIsSpace = !nextChar || /\s/.test(nextChar);

          if (inItalic && italicMarker === "_") {
            const afterIsSpace = i + 1 >= text.length || /\s/.test(text[i + 1]);
            if (afterIsSpace || /[^\w]/.test(text[i + 1])) {
              result.push("</em>");
              inItalic = false;
              italicMarker = null;
              i++;
              continue;
            }
          } else if (!inItalic && atStart && !nextIsSpace) {
            result.push("<em>");
            inItalic = true;
            italicMarker = "_";
            i++;
            continue;
          }
        } else if (char === "*") {
          const beforeSpace = i === 0 || /\s/.test(text[i - 1]);
          const afterSpace = i + 1 >= text.length || /\s/.test(text[i + 1]);

          if (beforeSpace && afterSpace) {
            // Literal asterisk, no formatting.
          } else if (inItalic && italicMarker === "*") {
            if (!beforeSpace) {
              result.push("</em>");
              inItalic = false;
              italicMarker = null;
              i++;
              continue;
            }
          } else if (!inItalic && !afterSpace) {
            result.push("<em>");
            inItalic = true;
            italicMarker = "*";
            i++;
            continue;
          }
        }
      }
    }

    result.push(char);
    i++;
  }

  if (inItalic) result.push("</em>");
  if (inBold) result.push("</strong>");

  return result.join("");
}

function renderMarkdown(text) {
  const recallTags = [];
  const withTokens = text.replace(/<recall>([\s\S]*?)<\/recall>/gi, (match, title) => {
    const token = `%%RECALL_${recallTags.length}%%`;
    recallTags.push(title.trim());
    return token;
  });

  let html = escapeHtml(withTokens);

  html = html.replace(/^######\s+(.+)$/gm, "<h6>$1</h6>");
  html = html.replace(/^#####\s+(.+)$/gm, "<h5>$1</h5>");
  html = html.replace(/^####\s+(.+)$/gm, "<h4>$1</h4>");
  html = html.replace(/^###\s+(.+)$/gm, "<h3>$1</h3>");
  html = html.replace(/^##\s+(.+)$/gm, "<h2>$1</h2>");
  html = html.replace(/^#\s+(.+)$/gm, "<h1>$1</h1>");
  html = html.replace(/^\*\*\*$/gm, "<hr>");

  const codeBlocks = [];
  html = html.replace(/```([\s\S]*?)```/g, (match, code) => {
    const placeholder = `%%CODE_BLOCK_${codeBlocks.length}%%`;
    codeBlocks.push(`<pre><code>${code}</code></pre>`);
    return placeholder;
  });

  const inlineCode = [];
  html = html.replace(/`([^`]+)`/g, (match, code) => {
    const placeholder = `%%INLINE_CODE_${inlineCode.length}%%`;
    inlineCode.push(`<code>${code}</code>`);
    return placeholder;
  });

  html = parseMarkdownSimple(html);

  codeBlocks.forEach((code, i) => {
    html = html.replace(`%%CODE_BLOCK_${i}%%`, code);
  });
  inlineCode.forEach((code, i) => {
    html = html.replace(`%%INLINE_CODE_${i}%%`, code);
  });

  const paragraphs = html.split(/\n\n+/);
  html = paragraphs
    .filter((p) => p.trim())
    .map((p) => `<p>${p}</p>`)
    .join("");

  recallTags.forEach((title, i) => {
    const safeTitle = escapeHtml(title);
    const badge = `<span class="badge badge-recall">&lt;recall&gt;${safeTitle}&lt;/recall&gt;</span>`;
    html = html.replace(`%%RECALL_${i}%%`, badge);
  });

  return html;
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
    body.className = "msg-body";
    body.innerHTML = renderMarkdown(msg.content);

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

openSettingsBtn.addEventListener("click", () => {
  openSettings();
});

closeSettingsBtn.addEventListener("click", () => {
  closeSettings();
});

settingsBackdrop.addEventListener("click", () => {
  closeSettings();
});

themeSelect.addEventListener("change", (event) => {
  setTheme(event.target.value);
});

document.addEventListener("keydown", (event) => {
  if (event.key === "Escape" && !settingsModal.classList.contains("hidden")) {
    closeSettings();
  }
});

(async () => {
  const savedTheme = getSavedTheme();
  setTheme(savedTheme);
  themeSelect.value = savedTheme;

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
