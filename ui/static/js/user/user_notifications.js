const bell = document.getElementById("notification-bell");
const modal = document.getElementById("notification-modal");
const closeBtn = modal.querySelector(".close-btn");
const listEl = document.getElementById("notification-list");
const markAllBtn = document.getElementById("mark-all-read");
const deleteAllBtn = document.getElementById("delete-all");

const sessionURL = "http://localhost:8080/forum/api/session/verify";

let csrfToken = null;

async function loadCSRFToken() {
  try {
    const resp = await fetch(sessionURL, { credentials: "include" });
    if (!resp.ok) return null;
    const data = await resp.json();
    return data.csrf_token || data.CSRFToken;
  } catch (err) {
    console.error("csrf", err);
    return null;
  }
}

async function fetchNotifications() {
  try {
    const resp = await fetch("http://localhost:8080/forum/api/notifications", {
      credentials: "include",
    });
    if (!resp.ok) return [];
    return await resp.json();
  } catch (err) {
    console.error("load notifications", err);
    return [];
  }
}

function renderNotifications(ns) {
  listEl.innerHTML = "";
  if (ns.length === 0) {
    listEl.innerHTML = "<li>No notifications</li>";
    bell.classList.remove("has-unread");
    return;
  }
  ns.forEach((n) => {
    const li = document.createElement("li");
    li.textContent = n.message;
    li.dataset.id = n.id || n.ID;
    li.addEventListener("click", () => markRead(li.dataset.id));
    listEl.appendChild(li);
  });
  bell.classList.add("has-unread");
}

async function markRead(id) {
  await loadCsrfIfNeeded();
  await fetch(`http://localhost:8080/forum/api/notifications/read/${id}`, {
    method: "POST",
    credentials: "include",
    headers: { "X-CSRF-Token": csrfToken },
  });
  loadAndRender();
}

async function markAll() {
  await loadCsrfIfNeeded();
  await fetch("http://localhost:8080/forum/api/notifications/read-all", {
    method: "POST",
    credentials: "include",
    headers: { "X-CSRF-Token": csrfToken },
  });
  loadAndRender();
}

async function deleteAll() {
  await loadCsrfIfNeeded();
  await fetch("http://localhost:8080/forum/api/notifications/delete-all", {
    method: "DELETE",
    credentials: "include",
    headers: { "X-CSRF-Token": csrfToken },
  });
  loadAndRender();
}

async function loadCsrfIfNeeded() {
  if (!csrfToken) {
    csrfToken = await loadCSRFToken();
  }
}

async function loadAndRender() {
  const ns = await fetchNotifications();
  renderNotifications(ns);
}

bell.addEventListener("click", async (e) => {
  e.preventDefault();
  modal.classList.remove("hidden");
  await loadAndRender();
});

closeBtn.addEventListener("click", () => {
  modal.classList.add("hidden");
});

markAllBtn.addEventListener("click", (e) => {
  e.preventDefault();
  markAll();
});
deleteAllBtn.addEventListener("click", (e) => {
  e.preventDefault();
  deleteAll();
});

window.addEventListener("click", (e) => {
  if (e.target === modal) modal.classList.add("hidden");
});
