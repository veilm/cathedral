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
      div.innerHTML = `<p class="text">${escapeHtml(msg.content)}</p>`;
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

const escapeHtml = (text) => {
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return text.replace(/[&<>"']/g, m => map[m]);
};

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
  div.innerHTML=`<p class="text">${escapeHtml(val)}</p>`;
  chat.appendChild(div);

  // Add loading indicator
  const loadingDiv = document.createElement('div');
  loadingDiv.className = 'message';
  loadingDiv.innerHTML = '<p class="text"><span class="loading"></span></p>';
  chat.appendChild(loadingDiv);
  loadingDiv.scrollIntoView({behavior:'smooth'});

  // Clear input immediately for better UX
  ta.value='';
  ta.disabled = true;

  try {
    // Send to backend API with conversation ID
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

    if (!response.ok) {
      throw new Error(`Server responded with ${response.status}`);
    }

    const data = await response.json();

    // Remove loading indicator
    loadingDiv.remove();

    // Add server response to UI
    const reply=document.createElement('div');
    reply.className='message';
    reply.innerHTML=`<p class="text">${escapeHtml(data.response)}</p>`;
    chat.appendChild(reply);
    reply.scrollIntoView({behavior:'smooth'});

  } catch (error) {
    console.error('Chat error:', error);

    // Remove loading indicator
    loadingDiv.remove();

    // Show error message
    const reply=document.createElement('div');
    reply.className='message';
    reply.innerHTML=`<p class="text" style="color: var(--gold);">
      The stones are silent. The connection to cathedral has been lost.
      Please try again later.
    </p>`;
    chat.appendChild(reply);
    reply.scrollIntoView({behavior:'smooth'});
  } finally {
    ta.disabled = false;
    ta.focus();
  }
};
ta.addEventListener('keydown',e=>{
  if(e.key==='Enter' && !e.shiftKey){e.preventDefault();send();}
});

// Check if we're viewing a specific conversation
window.addEventListener('DOMContentLoaded', () => {
  const path = window.location.pathname;
  if (path.startsWith('/c/')) {
    const conversationId = path.substring(3); // Remove '/c/' prefix
    if (conversationId) {
      currentConversationId = conversationId;
      loadConversation(conversationId);
    }
  }
});