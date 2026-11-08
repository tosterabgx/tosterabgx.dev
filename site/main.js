let hint = document.querySelector(".scroll-hint");

window.addEventListener("scroll", () => {
  if (window.scrollY > 50) {
    hint.classList.add("hidden")
  }
});
