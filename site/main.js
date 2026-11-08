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
