const convListEl = document.getElementById("convList");
const chatLogEl = document.getElementById("chatLog");
const chatForm = document.getElementById("chatForm");
const chatInput = document.getElementById("chatInput");
const memoryForm = document.getElementById("memoryForm");
const memoryInput = document.getElementById("memoryInput");
const memoryView = document.getElementById("memoryView");
const newConvBtn = document.getElementById("newConv");
const importConvBtn = document.getElementById("importConv");
const deleteConvBtn = document.getElementById("deleteConv");
const convMenuBtn = document.getElementById("convMenuBtn");
const convMenu = document.getElementById("convMenu");
const openSettingsBtn = document.getElementById("openSettings");
const closeSettingsBtn = document.getElementById("closeSettings");
const settingsModal = document.getElementById("settingsModal");
const themeSelect = document.getElementById("themeSelect");
const settingsBackdrop = settingsModal?.querySelector("[data-close-settings]");
const mainPanels = document.getElementById("mainPanels");
const mobileTabs = document.querySelectorAll(".mobile-tabs .tab");

let currentConvId = null;
let currentMessages = [];
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

function setActivePanel(panel) {
  if (!mainPanels) return;
  mainPanels.dataset.activePanel = panel;
  mobileTabs.forEach((tab) => {
    const isActive = tab.dataset.panel === panel;
    tab.classList.toggle("active", isActive);
    tab.setAttribute("aria-pressed", isActive ? "true" : "false");
  });
}

function isMobile() {
  return window.matchMedia("(max-width: 960px)").matches;
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

function openConvMenu() {
  convMenu.classList.remove("hidden");
  convMenuBtn.setAttribute("aria-expanded", "true");
}

function closeConvMenu() {
  convMenu.classList.add("hidden");
  convMenuBtn.setAttribute("aria-expanded", "false");
}

function toggleConvMenu() {
  if (convMenu.classList.contains("hidden")) openConvMenu();
  else closeConvMenu();
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

const HUMAN_TAG_RE = /^\s*<human\s+timestamp="([^"]+)"\s*>\s*\n?([\s\S]*?)\n?<\/human>\s*$/i;

function normalizeMessage(msg) {
  if (msg.role === "assistant") {
    return { roleLabel: "model", content: msg.content, roleClass: "" };
  }
  if (msg.role !== "user") {
    return { roleLabel: msg.role, content: msg.content, roleClass: "" };
  }

  const content = msg.content || "";
  const humanMatch = content.match(HUMAN_TAG_RE);
  if (humanMatch) {
    const timestamp = humanMatch[1];
    const inner = humanMatch[2] || "";
    return {
      roleLabel: `human · ${timestamp}`,
      content: inner,
      roleClass: "",
    };
  }

  if (content.trim().startsWith("<")) {
    return { roleLabel: "user:system", content, roleClass: "" };
  }

  return {
    roleLabel: "XML NOT DETECTED",
    content,
    roleClass: "role-error",
  };
}

function renderMessages(messages) {
  chatLogEl.innerHTML = "";
  messages.forEach((msg) => {
    const view = normalizeMessage(msg);
    const wrapper = document.createElement("div");
    wrapper.className = msg.loading ? "msg msg-loading" : "msg";
    if (msg.role === "error") wrapper.classList.add("msg-error");

    const role = document.createElement("div");
    role.className = view.roleClass ? `role ${view.roleClass}` : "role";
    role.textContent = view.roleLabel;

    const body = document.createElement("div");
    body.className = "msg-body";
    if (msg.loading) {
      body.innerHTML = '<span class="loading-dots" aria-label="Generating response"><span></span><span></span><span></span></span>';
    } else {
      body.innerHTML = renderMarkdown(view.content || "");
    }

    wrapper.appendChild(role);
    wrapper.appendChild(body);
    chatLogEl.appendChild(wrapper);
  });
  chatLogEl.scrollTop = chatLogEl.scrollHeight;
}

function appendMessage(msg) {
  currentMessages.push(msg);
  renderMessages(currentMessages);
}

function replaceLoadingMessage(msg) {
  const index = currentMessages.findIndex((item) => item.loading);
  if (index === -1) {
    appendMessage(msg);
    return;
  }
  currentMessages[index] = msg;
  renderMessages(currentMessages);
}

function removeLoadingMessage() {
  const index = currentMessages.findIndex((item) => item.loading);
  if (index === -1) return;
  currentMessages.splice(index, 1);
  renderMessages(currentMessages);
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
      if (isMobile()) setActivePanel("chat");
    });
    convListEl.appendChild(btn);
  });
}

async function loadConversation(id) {
  const convo = await fetchJSON(`/api/conversations/${id}`);
  currentMessages = convo.messages || [];
  renderMessages(currentMessages);
}

async function sendMessage(message) {
  if (!currentConvId) return;
  try {
    const appendRes = await fetchJSON(`/api/conversations/${currentConvId}/message`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message }),
    });
    appendMessage(appendRes.message);
  } catch (err) {
    appendMessage({
      role: "error",
      content: `<cathedral>\nFailed to write message: ${err.message}\n</cathedral>`,
    });
    return;
  }

  appendMessage({ role: "assistant", content: "", loading: true });

  try {
    const genRes = await fetchJSON(`/api/conversations/${currentConvId}/generate`, {
      method: "POST",
    });
    replaceLoadingMessage(genRes.message);
  } catch (err) {
    removeLoadingMessage();
    appendMessage({
      role: "error",
      content: `<cathedral>\nGeneration failed: ${err.message}\n</cathedral>`,
    });
  }
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
  closeConvMenu();
  const conv = await fetchJSON("/api/conversations", { method: "POST" });
  currentConvId = conv.id;
  setConvInUrl(currentConvId);
  await loadConversations();
  await loadConversation(conv.id);
  if (isMobile()) setActivePanel("chat");
});

importConvBtn.addEventListener("click", async () => {
  closeConvMenu();
  const path = window.prompt("Conversation path to import:");
  if (!path) return;
  const conv = await fetchJSON("/api/conversations/import", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ path }),
  });
  currentConvId = conv.id;
  setConvInUrl(currentConvId);
  await loadConversations();
  await loadConversation(conv.id);
  if (isMobile()) setActivePanel("chat");
});

deleteConvBtn.addEventListener("click", async () => {
  closeConvMenu();
  if (!currentConvId) return;
  const ok = window.confirm(`Unlink conversation ${currentConvId}?`);
  if (!ok) return;

  await fetchJSON(`/api/conversations/${currentConvId}`, { method: "DELETE" });
  const conversations = await fetchJSON("/api/conversations");
  if (conversations.length === 0) {
    currentConvId = null;
    currentMessages = [];
    setConvInUrl(null);
    renderMessages([]);
    await loadConversations();
    return;
  }
  currentConvId = conversations[0].id;
  setConvInUrl(currentConvId);
  await loadConversations();
  await loadConversation(currentConvId);
  if (isMobile()) setActivePanel("chat");
});

mobileTabs.forEach((tab) => {
  tab.addEventListener("click", () => {
    setActivePanel(tab.dataset.panel);
  });
});

convMenuBtn.addEventListener("click", (event) => {
  event.stopPropagation();
  toggleConvMenu();
});

document.addEventListener("click", (event) => {
  if (!convMenu.contains(event.target) && event.target !== convMenuBtn) {
    closeConvMenu();
  }
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
  if (event.key === "Escape" && !convMenu.classList.contains("hidden")) {
    closeConvMenu();
    return;
  }
  if (event.key === "Escape" && !settingsModal.classList.contains("hidden")) {
    closeSettings();
  }
});

(async () => {
  const savedTheme = getSavedTheme();
  setTheme(savedTheme);
  themeSelect.value = savedTheme;
  setActivePanel("chat");

  const conversations = await fetchJSON("/api/conversations");
  closeConvMenu();
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
