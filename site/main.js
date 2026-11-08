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
    list.replaceChildren();

    for (let p of projects) {
      list.append(makePanel(p));
    }
  } catch (err) {
    console.error("could not load projects:", err);
  }
}

function makePanel(project) {
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

loadProjects();
