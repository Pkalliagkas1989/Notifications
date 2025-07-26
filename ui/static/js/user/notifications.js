const bell = document.getElementById('notif-bell');
const modal = document.getElementById('notification-modal');
const list = document.getElementById('notif-list');
const badge = document.getElementById('notif-count');
const markAllBtn = document.getElementById('mark-all-read');
const delAllBtn = document.getElementById('delete-all');
const closeBtn = document.getElementById('close-notifications');

async function loadNotifications() {
  try {
    const resp = await fetch('http://localhost:8080/forum/api/user/notifications', { credentials: 'include' });
    if (!resp.ok) throw new Error('failed');
    const data = await resp.json();
    render(data);
    const unreadCount = data.filter(n => !n.read_at && n.message).length;
    const hasUnread = unreadCount > 0;
    bell.classList.toggle('lit', hasUnread);
    if (badge) {
      badge.textContent = unreadCount;
      badge.classList.toggle('hidden', !hasUnread);
    }
  } catch(e) { console.error(e); }
}

function render(nots) {
  list.innerHTML = '';
  if (!nots.length) { list.textContent = 'No notifications'; return; }
  nots.forEach(n => {
    const div = document.createElement('div');
    div.className = 'notification-item';
    if (!n.read_at && n.message) div.classList.add('unread');
    div.textContent = n.message || '(deleted)';
    div.dataset.id = n.id;
    list.appendChild(div);
  });
}

bell?.addEventListener('click', () => {
  modal.classList.toggle('hidden');
});

closeBtn?.addEventListener('click', () => { modal.classList.add('hidden'); });

list?.addEventListener('click', async (e) => {
  const id = e.target.dataset.id;
  if (!id) return;
  await fetch(`http://localhost:8080/forum/api/notifications/read/${id}`, { method: 'POST', credentials: 'include' });
  await loadNotifications();
});

markAllBtn?.addEventListener('click', async () => {
  await fetch('http://localhost:8080/forum/api/notifications/read-all', { method:'POST', credentials:'include' });
  await loadNotifications();
});

delAllBtn?.addEventListener('click', async () => {
  await fetch('http://localhost:8080/forum/api/notifications/delete-all', { method:'DELETE', credentials:'include' });
  await loadNotifications();
});

window.addEventListener('DOMContentLoaded', loadNotifications);
