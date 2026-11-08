// Prototype chrome: injects the variant switcher, the real site header/footer
// and the start of the next section, so each variant file only holds its section.
(() => {
  const variants = [
    ["00-trenutno.html", "Trenutna verzija"],
    ["01-put.html", "Put s kontrolnom točkom"],
    ["02-ilustracije.html", "Ilustrirani koraci"],
    ["03-prica.html", "Priča u četiri kadra"],
    ["04-razgovor.html", "Razgovor"],
    ["05-glagoli.html", "Četiri glagola"],
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
    <section class="section consulting-section proto-next" id="konzultacije" aria-label="Sljedeća sekcija">
      <div class="wrap">
        <div class="section-head">
          <span class="eyebrow">Konzultacije</span>
          <div><h2>Koristite AI za samostalan razvoj softwarea i trebate pomoć?</h2></div>
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
