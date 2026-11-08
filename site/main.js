let hint = document.querySelector(".scroll-hint");

let timer = setTimeout(() => {
  hint.classList.remove("hidden");
}, 2000);

window.addEventListener(
  "scroll",
  () => {
    clearTimeout(timer);
    hint.classList.add("hidden");
  },
  { once: true, passive: true },
);

async function loadProjects() {
  let list = document.querySelector(".projects");

  try {
    let res = await fetch("/api/projects");
    if (!res.ok) throw new Error("status " + res.status);
    let projects = await res.json();

    if (projects.length === 0) return;
    list.replaceChildren(...projects.map(makeProject));
  } catch (err) {
    console.error("could not load projects:", err);
  }
}

function makeProject(project) {
  let a = document.createElement("a");
  a.className = "project";
  a.href = project.url;

  let title = document.createElement("h3");
  title.textContent = project.name;
  a.append(title);

  if (project.description) {
    let desc = document.createElement("p");
    desc.textContent = project.description;
    a.append(desc);
  }

  return a;
}

async function loadNotes() {
  let list = document.querySelector(".notes");
  try {
    let res = await fetch("/api/notes");
    if (!res.ok) throw new Error("status " + res.status);
    let notes = await res.json();

    list.replaceChildren(...notes.map(makeNote));
  } catch (err) {
    console.error("could not load notes:", err);
  }
}

function makeNote(note) {
  let item = document.createElement("article");
  item.className = "note";

  let name = document.createElement("span");
  name.className = "note-name";
  name.textContent = note.name;

  let date = document.createElement("time");
  date.className = "note-date";
  date.dateTime = note.created_at;
  date.textContent = new Date(note.created_at).toLocaleDateString();

  let message = document.createElement("p");
  message.textContent = note.message;

  item.append(name, date, message);
  return item;
}

loadProjects();
loadNotes();
setInterval(loadNotes, 15000);
