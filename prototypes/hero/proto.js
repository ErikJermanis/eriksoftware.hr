// Prototype chrome: injects the variant switcher, the real site header/footer
// and the start of the next section, so each variant file only holds its hero.
(() => {
  const variants = [
    ["00-trenutno.html", "Trenutna verzija"],
    ["01-precrtano.html", "Precrtano"],
    ["02-ilustracije.html", "Ilustracije"],
    ["03-po-mjeri.html", "Po mjeri"],
    ["04-spojeno.html", "Sve spojeno"],
    ["05-dok-radite.html", "Dok Vi radite"],
  ];
  const file = location.pathname.split("/").pop();
  const i = Math.max(0, variants.findIndex(([f]) => f === file));
  const prev = variants[(i - 1 + variants.length) % variants.length];
  const next = variants[(i + 1) % variants.length];

  const bar = `
    <div class="proto-bar">
      <a href="index.html">← Sve</a>
      <span class="proto-bar-title">${variants[i][0].slice(0, 2)} · ${variants[i][1]}</span>
      <nav aria-label="Prototipovi">
        <a href="${prev[0]}" aria-label="Prethodna varijanta">‹</a>
        <a href="${next[0]}" aria-label="Sljedeća varijanta">›</a>
      </nav>
    </div>`;

  const header = `
    <header class="site-header">
      <div class="wrap header-inner">
        <a class="brand" href="index.html" aria-label="Erik Software — naslovnica">
          <img src="../../public/eriksoftware_logo.svg" width="42" height="42" alt="">
          <span>Erik Software</span>
        </a>
        <button class="menu-toggle" type="button" aria-controls="site-nav" aria-expanded="false">
          <span class="menu-icon" aria-hidden="true"></span>
          <span>Izbornik</span>
        </button>
        <nav id="site-nav" aria-label="Glavna navigacija">
          <a href="#kako-radim">Kako radim</a>
          <a href="#konzultacije">Konzultacije</a>
          <a href="#o-meni">O meni</a>
          <a href="#">Blog</a>
        </nav>
      </div>
    </header>`;

  const teaser = `
    <section class="section process-section proto-next" id="kako-radim" aria-label="Sljedeća sekcija">
      <div class="wrap">
        <div class="section-head">
          <span class="eyebrow">Proces</span>
          <div>
            <h2>Kako radim?</h2>
            <p>Četiri koraka, od prvog razgovora do alata koji koristite svaki dan.</p>
          </div>
        </div>
      </div>
    </section>`;

  const footer = `
    <footer class="site-footer">
      <div class="wrap footer-inner">
        <div class="footer-identity">
          <span>Erik Jermaniš / software i automatizacija</span>
          <span>© ${new Date().getFullYear()} ERIK SOFTWARE</span>
        </div>
        <div class="footer-links">
          <a href="mailto:erik@eriksoftware.hr">erik@eriksoftware.hr</a>
          <a href="#">RSS</a>
          <a href="#main">Natrag na vrh ↑</a>
        </div>
      </div>
    </footer>`;

  const main = document.querySelector("main");
  main.insertAdjacentHTML("beforebegin", bar + header);
  main.insertAdjacentHTML("beforeend", teaser);
  main.insertAdjacentHTML("afterend", footer);

  const siteHeader = document.querySelector(".site-header");
  const toggle = siteHeader.querySelector(".menu-toggle");
  const closeMenu = () => {
    siteHeader.classList.remove("menu-open");
    toggle.setAttribute("aria-expanded", "false");
  };
  toggle.addEventListener("click", () => {
    const open = siteHeader.classList.toggle("menu-open");
    toggle.setAttribute("aria-expanded", String(open));
  });
  siteHeader.querySelector("nav").addEventListener("click", (e) => {
    if (e.target.closest("a")) closeMenu();
  });
})();

// Shared by the variants: [data-in] elements get .is-in once they are first
// seen (one-off moves) and .is-live while on screen (loops pause off-screen).
// With reduced motion everything is shown in its final state straight away.
(() => {
  const els = [...document.querySelectorAll("[data-in]")];
  const calm = matchMedia("(prefers-reduced-motion: reduce)").matches;
  if (calm || !("IntersectionObserver" in window)) {
    els.forEach((el) => el.classList.add("is-in"));
    return;
  }
  const io = new IntersectionObserver(
    (entries) =>
      entries.forEach((entry) => {
        if (entry.isIntersecting) entry.target.classList.add("is-in");
        entry.target.classList.toggle("is-live", entry.isIntersecting);
        entry.target.dispatchEvent(new CustomEvent(entry.isIntersecting ? "hx:show" : "hx:hide"));
      }),
    { threshold: 0.35 }
  );
  // Let fonts and layout settle first so the hero does not animate mid-load.
  const start = () => setTimeout(() => els.forEach((el) => io.observe(el)), 250);
  document.fonts ? document.fonts.ready.then(start) : start();
})();
