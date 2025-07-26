const bell = document.getElementById('notif-bell');
const modal = document.getElementById('notification-modal');
const list = document.getElementById('notif-list');
const markAllBtn = document.getElementById('mark-all-read');
const delAllBtn = document.getElementById('delete-all');
const closeBtn = document.getElementById('close-notifications');

async function loadNotifications() {
  try {
    const resp = await fetch('http://localhost:8080/forum/api/user/notifications', { credentials: 'include' });
    if (!resp.ok) throw new Error('failed');
    const data = await resp.json();
    render(data);
    const hasUnread = data.some(n => !n.read_at && n.message);
    bell.classList.toggle('lit', hasUnread);
  } catch(e) { console.error(e); }
}

function render(nots) {
  list.innerHTML = '';
  if (!nots.length) { list.textContent = 'No notifications'; return; }
  nots.forEach(n => {
    const div = document.createElement('div');
    div.className = 'notification-item';
    if (!n.read_at && n.message) div.classList.add('unread');
    div.dataset.id = n.id;

    const icon = document.createElement('span');
    icon.className = 'notification-icon';
    icon.textContent = getIcon(n.type);

    const link = document.createElement('a');
    link.className = 'notif-link';
    link.textContent = n.message || '(deleted)';
    if (n.post_id) {
      let url = `/user/post?id=${encodeURIComponent(n.post_id)}`;
      if (n.comment_id) url += `#${encodeURIComponent(n.comment_id)}`;
      link.href = url;
    } else {
      link.href = '#';
    }

    const time = document.createElement('time');
    time.className = 'notif-time';
    time.textContent = new Date(n.created_at).toLocaleString();

    div.appendChild(icon);
    div.appendChild(link);
    div.appendChild(time);
    list.appendChild(div);
  });
}

function getIcon(type) {
  switch (type) {
    case 'comment':
      return '💬';
    case 'comment_edit':
      return '✏️';
    case 'comment_delete':
      return '❌';
    case 'reaction':
      return '👍';
    default:
      return '🔔';
  }
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
