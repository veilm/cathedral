/* ---------- version info ---------- */
const CATHEDRAL_VERSION = 'v0.2.0-recall';

/* ---------- settings menu ---------- */
const settingsBtn = document.getElementById('settings');
const settingsMenu = document.querySelector('.settings-menu');

settingsBtn.addEventListener('click', () => {
  settingsMenu.classList.toggle('open');
});

document.addEventListener('click', (e) => {
  if (!settingsBtn.contains(e.target) && !settingsMenu.contains(e.target)) {
    settingsMenu.classList.remove('open');
  }
});

async function newConversation() {
  try {
    const response = await fetch('/api/new-conversation', {
      method: 'POST'
    });

    if (!response.ok) {
      throw new Error('Failed to create new conversation');
    }

    const data = await response.json();

    // Redirect to the new conversation
    window.location.href = data.url;
  } catch (error) {
    console.error('Failed to create new conversation:', error);
    alert('Failed to create new conversation. Please try again.');
  }
}

async function newConversationContinued() {
  try {
    const response = await fetch('/api/new-conversation-continued', {
      method: 'POST'
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || 'Failed to create continued conversation');
    }

    const data = await response.json();

    // Redirect to the new conversation
    window.location.href = data.url;
  } catch (error) {
    console.error('Failed to create continued conversation:', error);
    alert(`Failed to create continued conversation: ${error.message}`);
  }
}

function sendMessage() {
  // Close the settings menu
  settingsMenu.classList.remove('open');

  // Call the send function that handles message submission
  if (typeof send === 'function') {
    send();
  } else {
    console.error('Send function not available');
  }
}

async function consolidateMemory() {
  // Check if we have a current conversation
  if (!currentConversationId) {
    alert('No active conversation to consolidate. Please start a conversation first.');
    return;
  }

  // Close the settings menu
  settingsMenu.classList.remove('open');

  // Confirm with user
  const confirmed = confirm('This will save the current conversation to your memory store. Continue?');
  if (!confirmed) {
    return;
  }

  // Show loading state
  const originalText = document.querySelector('.settings-menu button[onclick="consolidateMemory()"]').textContent;
  const button = document.querySelector('.settings-menu button[onclick="consolidateMemory()"]');
  button.textContent = 'Consolidating...';
  button.disabled = true;

  try {
    const response = await fetch('/api/consolidate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        conversation_id: currentConversationId,
        compression: 0 // Let backend use its default
      })
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to consolidate: ${errorText}`);
    }

    const data = await response.json();

    // Show success message
    alert(`Memory consolidated successfully!\nEpisode saved at: ${data.episode_path}`);

    // Optionally, create a new conversation after consolidation
    const startNew = confirm('Would you like to start a new conversation?');
    if (startNew) {
      newConversation();
    }

  } catch (error) {
    console.error('Failed to consolidate memory:', error);
    alert(`Failed to consolidate memory: ${error.message}`);
  } finally {
    // Restore button state
    button.textContent = originalText;
    button.disabled = false;
  }
}

async function loadConversation(conversationId) {
  try {
    const response = await fetch(`/api/conversation/${conversationId}`);

    if (!response.ok) {
      throw new Error('Conversation not found');
    }

    const data = await response.json();
    const chat = document.getElementById('chat');

    // Clear current content
    chat.innerHTML = '';

    // Add messages
    for (const msg of data.messages) {
      const div = document.createElement('div');
      div.className = msg.role === 'user' ? 'message user' : 'message';
      div.innerHTML = `<div class="text">${renderMarkdown(msg.content)}</div>`;
      chat.appendChild(div);
    }

    // If no messages, show the default content
    if (data.messages.length === 0) {
      chat.innerHTML = `
        <div class="message">
          <p class="text">
            In this quiet hall, words echo differently.
            Stone remembers what was spoken, and silence holds its own weight.
            Share what moves through your mind—brief meditation or lengthy contemplation.
          </p>
        </div>

        <table class="scripture-table">
          <thead>
            <tr>
              <th>Hour</th>
              <th>Office</th>
              <th>Silence</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Dawn</td>
              <td>Lauds</td>
              <td>Breaking</td>
            </tr>
            <tr>
              <td>Third</td>
              <td>Terce</td>
              <td>Working</td>
            </tr>
            <tr>
              <td>Ninth</td>
              <td>None</td>
              <td>Resting</td>
            </tr>
            <tr>
              <td>Evening</td>
              <td>Vespers</td>
              <td>Gathering</td>
            </tr>
            <tr>
              <td>Night</td>
              <td>Compline</td>
              <td>Complete</td>
            </tr>
          </tbody>
        </table>
      `;
    }
  } catch (error) {
    console.error('Failed to load conversation:', error);
    // Show error message
    const chat = document.getElementById('chat');
    chat.innerHTML = `
      <div class="message">
        <p class="text" style="color: var(--gold);">
          The conversation has been lost to time.
        </p>
      </div>
    `;
  }
}

/* ---------- minimal chat loop ---------- */
const chat=document.getElementById('chat');
const ta=document.querySelector('textarea');
let currentConversationId = null;

// Markdown rendering utilities
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

// Simple character-by-character markdown parser for bold and italic
function parseMarkdownSimple(text) {
	let result = [];
	let i = 0;

	// Current state
	let inBold = false;
	let inItalic = false;
	let boldMarker = null; // Track which marker opened bold ('*' or '_')
	let italicMarker = null; // Track which marker opened italic

	while (i < text.length) {
		const char = text[i];
		const nextChar = i + 1 < text.length ? text[i + 1] : null;
		const prevChar = i > 0 ? text[i - 1] : null;

		if (char === "*" || char === "_") {
			// Count consecutive markers
			let markerCount = 1;
			let j = i + 1;
			while (j < text.length && text[j] === char) {
				markerCount++;
				j++;
			}

			// Handle based on marker count
			if (markerCount >= 3) {
				// Triple marker - could be bold+italic
				if (
					inBold &&
					inItalic &&
					boldMarker === char &&
					italicMarker === char
				) {
					// Close both
					result.push("</em></strong>");
					inBold = false;
					inItalic = false;
					boldMarker = null;
					italicMarker = null;
					i += 3;
					continue;
				} else if (!inBold && !inItalic) {
					// Open both
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
				// Bold marker
				if (inBold && boldMarker === char) {
					// Close bold
					result.push("</strong>");
					inBold = false;
					boldMarker = null;
					i += 2;
					continue;
				} else if (!inBold) {
					// Open bold
					result.push("<strong>");
					inBold = true;
					boldMarker = char;
					i += 2;
					continue;
				}
			}

			// Single marker (italic) - only if we have exactly 1
			if (markerCount === 1) {
				if (char === "_") {
					const atStart = i === 0 || /\s/.test(prevChar);
					const nextIsSpace = !nextChar || /\s/.test(nextChar);

					if (inItalic && italicMarker === "_") {
						// Closing underscore italic - check we're at end of word
						const afterIsSpace = i + 1 >= text.length || /\s/.test(text[i + 1]);
						if (afterIsSpace || /[^\w]/.test(text[i + 1])) {
							result.push("</em>");
							inItalic = false;
							italicMarker = null;
							i++;
							continue;
						}
					} else if (!inItalic && atStart && !nextIsSpace) {
						// Opening underscore italic
						result.push("<em>");
						inItalic = true;
						italicMarker = "_";
						i++;
						continue;
					}
				} else if (char === "*") {
					// Check if this is a literal asterisk (surrounded by spaces)
					const beforeSpace = i === 0 || /\s/.test(text[i - 1]);
					const afterSpace = i + 1 >= text.length || /\s/.test(text[i + 1]);

					if (beforeSpace && afterSpace) {
						// Literal asterisk - skip formatting
					} else if (inItalic && italicMarker === "*") {
						// Check if valid closing (not preceded by space)
						if (!beforeSpace) {
							// Close italic
							result.push("</em>");
							inItalic = false;
							italicMarker = null;
							i++;
							continue;
						}
					} else if (!inItalic && !afterSpace) {
						// Open italic (not followed by space)
						result.push("<em>");
						inItalic = true;
						italicMarker = "*";
						i++;
						continue;
					}
				}
			}
		}

		// Regular character
		result.push(char);
		i++;
	}

	// Close any unclosed tags
	if (inItalic) result.push("</em>");
	if (inBold) result.push("</strong>");

	return result.join("");
}

// Improved markdown renderer
function renderMarkdown(text) {
	// First escape HTML to prevent XSS
	let html = escapeHtml(text);

	// Normalize multiple consecutive newlines to max 2
	html = html.replace(/\n{3,}/g, '\n\n');

	// Headers (h1-h6) - do these first
	html = html.replace(/^######\s+(.+)$/gm, "<h6>$1</h6>");
	html = html.replace(/^#####\s+(.+)$/gm, "<h5>$1</h5>");
	html = html.replace(/^####\s+(.+)$/gm, "<h4>$1</h4>");
	html = html.replace(/^###\s+(.+)$/gm, "<h3>$1</h3>");
	html = html.replace(/^##\s+(.+)$/gm, "<h2>$1</h2>");
	html = html.replace(/^#\s+(.+)$/gm, "<h1>$1</h1>");

	// Horizontal rules (*** alone on a line) - handle before code blocks
	html = html.replace(/^\*\*\*$/gm, "<hr>");

	// Code blocks (```) - protect from further processing
	const codeBlocks = [];
	html = html.replace(/```([\s\S]*?)```/g, function (match, code) {
		const placeholder = `\x00CODE_BLOCK_${codeBlocks.length}\x00`;
		codeBlocks.push("<pre><code>" + code + "</code></pre>");
		return placeholder;
	});

	// Inline code (`) - protect from further processing
	const inlineCode = [];
	html = html.replace(/`([^`]+)`/g, function (match, code) {
		const placeholder = `\x00INLINE_CODE_${inlineCode.length}\x00`;
		inlineCode.push("<code>" + code + "</code>");
		return placeholder;
	});

	// Parse bold and italic
	html = parseMarkdownSimple(html);

	// Restore code blocks and inline code
	codeBlocks.forEach((code, i) => {
		html = html.replace(`\x00CODE_BLOCK_${i}\x00`, code);
	});
	inlineCode.forEach((code, i) => {
		html = html.replace(`\x00INLINE_CODE_${i}\x00`, code);
	});

	// Line breaks
	html = html.replace(/  \n/g, "<br>\n");

	// Split into paragraphs and wrap properly
	const paragraphs = html.split(/\n\n+/);
	html = paragraphs
		.filter((p) => p.trim()) // Remove empty paragraphs
		.map((p) => `<p>${p}</p>`)
		.join("");

	return html;
}

const send=async()=>{
  const val=ta.value.trim();
  if(!val)return;

  // Ensure we have a conversation ID
  if(!currentConversationId) {
    // Create new conversation if none exists
    try {
      const response = await fetch('/api/new-conversation', { method: 'POST' });
      const data = await response.json();
      currentConversationId = data.id;
      window.history.pushState({}, '', `/c/${currentConversationId}`);
    } catch (error) {
      console.error('Failed to create conversation:', error);
      return;
    }
  }

  // Add user message
  const div=document.createElement('div');
  div.className='message user';
  div.innerHTML=`<div class="text">${renderMarkdown(val)}</div>`;
  chat.appendChild(div);

  // Add loading indicator
  const loadingDiv = document.createElement('div');
  loadingDiv.className = 'message';
  loadingDiv.innerHTML = '<p class="text"><span class="loading"></span></p>';
  chat.appendChild(loadingDiv);
  loadingDiv.scrollIntoView({behavior:'smooth'});

  // Clear input immediately for better UX
  ta.value='';
  ta.style.height = 'auto'; // Reset height after sending
  ta.disabled = true;

  // First, submit the message via POST to save it
  try {
    const response = await fetch('/api/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        message: val,
        conversation_id: currentConversationId
      })
    });
  } catch (error) {
    console.error('Chat submission error:', error);
    loadingDiv.innerHTML = '<p class="text" style="color: var(--gold);">Error: Unable to submit message</p>';
    ta.disabled = false;
    ta.focus();
    return;
  }

  // Now use EventSource to receive streaming responses
  let firstResponse = true;
  let recallBadge = null;
  let closedIntentionally = false;

  // Create EventSource connection
  const eventSource = new EventSource(`/api/chat?conversation_id=${encodeURIComponent(currentConversationId)}`);

  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);

      if (data.error) {
        console.error('Server error:', data.error);
        loadingDiv.innerHTML = `<p class="text" style="color: var(--gold);">Error: ${data.error}</p>`;
        closedIntentionally = true;
        eventSource.close();
        ta.disabled = false;
        ta.focus();
        return;
      }

      // Remove loading indicator on first response
      if (firstResponse) {
        loadingDiv.remove();
        firstResponse = false;
      }

      // Add assistant response
      const respDiv = document.createElement('div');
      respDiv.className = 'message';
      respDiv.innerHTML = `<div class="text">${renderMarkdown(data.response)}</div>`;
      chat.appendChild(respDiv);

      // If this response continues (has recall), show indicator
      if (data.continues) {
        // Remove existing recall badge if any
        if (recallBadge) {
          recallBadge.remove();
        }
        recallBadge = document.createElement('div');
        recallBadge.className = 'message';
        recallBadge.innerHTML = '<p class="text" style="opacity: 0.7;"><em>Recalling memory...</em> <span class="loading"></span></p>';
        chat.appendChild(recallBadge);
        recallBadge.scrollIntoView({behavior: 'smooth'});
      } else {
        // Remove recall badge if it exists
        if (recallBadge) {
          recallBadge.remove();
          recallBadge = null;
        }
        // Mark as intentional close before closing
        closedIntentionally = true;
        eventSource.close();
        ta.disabled = false;
        ta.focus();
      }

      respDiv.scrollIntoView({behavior: 'smooth'});
    } catch (error) {
      console.error('Error processing SSE message:', error);
    }
  };

  eventSource.onerror = (error) => {
    console.error('EventSource error:', error);
    eventSource.close();

    // Only show error if we didn't close intentionally
    if (!closedIntentionally) {
      if (firstResponse) {
        loadingDiv.innerHTML = '<p class="text" style="color: var(--gold);">The stones are silent. Connection lost.</p>';
      } else {
        // Show error inline if we were mid-conversation
        const errorDiv = document.createElement('div');
        errorDiv.className = 'message';
        errorDiv.innerHTML = '<p class="text" style="color: var(--gold);"><em>Connection lost during recall.</em></p>';
        chat.appendChild(errorDiv);
      }
    }

    if (recallBadge) {
      recallBadge.remove();
    }
    ta.disabled = false;
    ta.focus();
  }
};
// Auto-resize textarea as user types (max 5 lines)
ta.addEventListener('input', () => {
  ta.style.height = 'auto';
  const maxHeight = parseFloat(getComputedStyle(ta).lineHeight) * 5;
  ta.style.height = Math.min(ta.scrollHeight, maxHeight) + 'px';
});

// Handle keyboard shortcuts - Ctrl+Enter to submit
ta.addEventListener('keydown',e=>{
  if(e.key==='Enter' && (e.ctrlKey || e.metaKey)){
    e.preventDefault();
    send();
  }
});

// Check if we're viewing a specific conversation
window.addEventListener('DOMContentLoaded', () => {
  // Display version
  const versionDisplay = document.getElementById('version-display');
  if (versionDisplay) {
    versionDisplay.textContent = CATHEDRAL_VERSION;
  }

  const path = window.location.pathname;
  if (path.startsWith('/c/')) {
    const conversationId = path.substring(3); // Remove '/c/' prefix
    if (conversationId) {
      currentConversationId = conversationId;
      loadConversation(conversationId);
    }
  }
});