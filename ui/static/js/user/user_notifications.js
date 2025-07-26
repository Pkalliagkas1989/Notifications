const bell = document.getElementById("notif-bell");
const modal = document.getElementById("notification-modal");
const list = document.getElementById("notification-list");
const readAllBtn = document.getElementById("read-all-btn");
const deleteAllBtn = document.getElementById("delete-all-btn");

async function fetchNotifications() {
  const resp = await fetch("http://localhost:8080/forum/api/notifications", {
    credentials: "include",
  });
  if (!resp.ok) return [];
  return resp.json();
}

function render(notifs) {
  list.innerHTML = "";
  let unread = 0;
  notifs.forEach((n) => {
    const li = document.createElement("li");
    li.textContent = `${n.type} on your post`;
    if (!n.read) unread++;
    const readBtn = document.createElement("button");
    readBtn.textContent = "Read";
    readBtn.addEventListener("click", async () => {
      await fetch(
        `http://localhost:8080/forum/api/notifications/read/${n.id}`,
        { method: "PATCH", credentials: "include" },
      );
      li.classList.add("read");
      bell.classList.remove("has-notif");
    });
    const delBtn = document.createElement("button");
    delBtn.textContent = "Delete";
    delBtn.addEventListener("click", async () => {
      await fetch(
        `http://localhost:8080/forum/api/notifications/delete/${n.id}`,
        { method: "DELETE", credentials: "include" },
      );
      li.remove();
    });
    li.appendChild(readBtn);
    li.appendChild(delBtn);
    list.appendChild(li);
  });
  if (unread > 0) bell.classList.add("has-notif");
}

bell.addEventListener("click", async () => {
  modal.classList.toggle("hidden");
  const notifs = await fetchNotifications();
  render(notifs);
});

readAllBtn.addEventListener("click", async () => {
  await fetch("http://localhost:8080/forum/api/notifications/read-all", {
    method: "PATCH",
    credentials: "include",
  });
  modal.classList.add("hidden");
});

deleteAllBtn.addEventListener("click", async () => {
  await fetch("http://localhost:8080/forum/api/notifications/delete-all", {
    method: "DELETE",
    credentials: "include",
  });
  list.innerHTML = "";
  bell.classList.remove("has-notif");
});
