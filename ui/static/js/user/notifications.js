
const bell = document.getElementById('notif-bell');
const modal = document.getElementById('notification-modal');
const list = document.getElementById('notif-list');
const badge = document.getElementById('notif-count');
const markAllBtn = document.getElementById('mark-all-read');
const delAllBtn = document.getElementById('delete-all');
const closeBtn = document.getElementById('close-notifications');

// Hold CSRF token for requests to the API
let csrfToken = null;

// Endpoint used to verify the session and retrieve a CSRF token
const sessionVerifyURL = 'http://localhost:8080/forum/api/session/verify';

async function loadCSRFToken() {
  if (csrfToken) return csrfToken;
  try {
    const resp = await fetch(sessionVerifyURL, { credentials: 'include' });
    if (!resp.ok) throw new Error('session invalid');
    const data = await resp.json();
    csrfToken = data.csrf_token || data.CSRFToken;
    return csrfToken;
  } catch (err) {
    console.warn('Failed to load CSRF token', err);
    return null;
  }
}

async function loadNotifications() {
  try {
    const resp = await fetch(
      `http://localhost:8080/forum/api/user/notifications?limit=${limit}&offset=${offset}`,
      { credentials: "include" },
    );
    if (!resp.ok) throw new Error("failed");
    const data = await resp.json();
    render(data, !reset);

    if (reset) {
      const unreadCount = data.filter((n) => !n.read_at && n.message).length;
      const hasUnread = unreadCount > 0;
      bell.classList.toggle("lit", hasUnread);
      if (badge) {
        badge.textContent = unreadCount;
        badge.classList.toggle("hidden", !hasUnread);
      }
    }

    if (data.length === limit) {
      offset += limit;
      loadMoreBtn.classList.remove("hidden");
    } else {
      loadMoreBtn.classList.add("hidden");
    }
  } catch (e) {
    console.error(e);
  }
}

function render(nots, append = false) {
  if (!append) list.innerHTML = "";
  if (!nots.length) {
    list.textContent = "No notifications";
    return;
  }
  nots.forEach((n) => {
    const div = document.createElement("div");
    div.className = "notification-item";
    if (!n.read_at && n.message) div.classList.add("unread");
    div.dataset.id = n.id;

    const icon = document.createElement("span");
    icon.className = "notification-icon";
    icon.textContent = getIcon(n.type);

    const link = document.createElement("a");
    link.className = "notif-link";
    link.textContent = n.message || "(deleted)";
    if (n.post_id) {
      let url = `/user/post?id=${encodeURIComponent(n.post_id)}`;
      if (n.comment_id) url += `#${encodeURIComponent(n.comment_id)}`;
      link.href = url;
    } else {
      link.href = "#";
    }

    const time = document.createElement("time");
    time.className = "notif-time";
    time.textContent = new Date(n.created_at).toLocaleString();

    div.appendChild(icon);
    div.appendChild(link);
    div.appendChild(time);
    list.appendChild(div);
  });
}

function getIcon(type) {
  switch (type) {
    case "comment":
      return "💬";
    case "comment_edit":
      return "✏️";
    case "comment_delete":
      return "❌";
    case "reaction":
      return "👍";
    default:
      return "🔔";
  }
}

bell?.addEventListener("click", () => {
  modal.classList.toggle("hidden");
});

closeBtn?.addEventListener("click", () => {
  modal.classList.add("hidden");
});


list?.addEventListener('click', async (e) => {
  const item = e.target.closest('.notification-item');
  const id = item?.dataset.id;
  if (!id) return;
  const token = await loadCSRFToken();
  await fetch(`http://localhost:8080/forum/api/notifications/read/${id}`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'X-CSRF-Token': token || ''
    }
  });
  await loadNotifications();
});

markAllBtn?.addEventListener('click', async () => {
  const token = await loadCSRFToken();
  await fetch('http://localhost:8080/forum/api/notifications/read-all', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'X-CSRF-Token': token || ''
    }
  });
  await loadNotifications();
});

delAllBtn?.addEventListener('click', async () => {
  const token = await loadCSRFToken();
  await fetch('http://localhost:8080/forum/api/notifications/delete-all', {
    method: 'DELETE',
    credentials: 'include',
    headers: {
      'X-CSRF-Token': token || ''
    }
  });
  await loadNotifications();
});

loadMoreBtn?.addEventListener("click", () => loadNotifications(false));
window.addEventListener("DOMContentLoaded", () => loadNotifications(true));
