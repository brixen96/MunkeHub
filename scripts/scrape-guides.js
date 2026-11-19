// Temporary scraper: parses wow-professions.com Forever guide HTML that was
// stashed on window.__guides by the caller. Injected via page.addScriptTag.
// Output: window.__parsed = { slug: guideData }
(function () {
  function parseShopping(doc) {
    const out = [];
    doc.querySelectorAll("h3").forEach((h) => {
      const range = h.textContent.trim();
      if (!/^\d/.test(range)) return;
      let el = h.nextElementSibling;
      while (el && el.tagName !== "H2" && el.tagName !== "H3") {
        if (el.tagName === "TABLE") {
          const items = [];
          el.querySelectorAll("tbody tr").forEach((tr) => {
            const td = tr.querySelectorAll("td")[1];
            if (!td) return;
            const a = td.querySelector("a");
            const m = td.textContent.trim().split("x ");
            if (m.length === 2)
              items.push({ q: m[0], n: m[1], i: a ? (a.href.match(/item=(\d+)/) || [])[1] : null });
          });
          out.push({ range, items });
          break;
        }
        el = el.nextElementSibling;
      }
    });
    return out;
  }

  function parseSteps(doc) {
    const out = [];
    doc.querySelectorAll("h2").forEach((h2) => {
      const title = h2.textContent.trim();
      if (title.indexOf("(") < 0) return;
      const range = title.slice(title.indexOf("(") + 1).replace(")", "");
      const rows = [];
      const notes = [];
      let el = h2.nextElementSibling;
      while (el && el.tagName !== "H2") {
        if (el.classList.contains("route-step")) {
          const r = (el.querySelector(".route-step, .route-range") || {}).textContent ||
            (el.querySelector(".route-range") || {}).textContent || "";
          const p = el.querySelector("p");
          if (p) {
            const a = p.querySelector("a");
            const t = p.textContent.trim();
            const dash = t.indexOf(" - ");
            const head = dash > 0 ? t.slice(0, dash) : t;
            const reag = dash > 0 ? t.slice(dash + 3) : "";
            const sp = head.indexOf("x ");
            const m = a ? a.href.match(/(spell|item)=(\d+)/) : null;
            rows.push({
              s: r,
              q: sp > 0 ? head.slice(0, sp) : "",
              c: sp > 0 ? head.slice(sp + 2) : head,
              id: m ? m[2] : null,
              k: m ? m[1] : null,
              rg: reag,
            });
          }
          // prose notes are the source site's writing — NOT copied (facts only)
        }
        el = el.nextElementSibling;
      }
      out.push({ range, rows, notes: [] });
    });
    return out;
  }

  // Gathering guides: H2 skill ranges. Zone names come from the faction tab
  // buttons (facts). Prose notes are the source site's writing — NOT copied.
  function parseGathering(doc) {
    const out = [];
    doc.querySelectorAll("h2").forEach((h2) => {
      const title = h2.textContent.trim();
      if (!/^\d+-\d+$/.test(title)) return;
      const zones = [];
      let el = h2.nextElementSibling;
      while (el && el.tagName !== "H2") {
        // zone tab buttons inside faction tab containers
        if (el.tagName === "DIV" && /tab-container/.test(el.className || "")) {
          el.querySelectorAll("button[data-tab]").forEach((b) => {
            const z = b.textContent.trim();
            if (z && !/^(Horde|Alliance)$/i.test(z) && zones.indexOf(z) < 0) zones.push(z);
          });
        }
        // zone lists (fishing: capital cities to fish in)
        if (el.tagName === "UL") {
          el.querySelectorAll("li").forEach((li) => {
            const z = li.textContent.trim();
            if (z && z.length < 40 && zones.indexOf(z) < 0) zones.push(z);
          });
        }
        // zone/vendor tables (mining pick vendors etc are NOT zones — skip)
        el = el.nextElementSibling;
      }
      out.push({ range: title, zones });
    });
    return out;
  }

  function parseTrainers(doc) {
    const h2s = Array.from(doc.querySelectorAll("h2"));
    const start = h2s.find((h) => /Trainer Locations/.test(h.textContent));
    if (!start) return [];
    const all = Array.from(doc.querySelectorAll("h2, h3, h4, ul, table, p, details"));
    const si = all.indexOf(start);
    const res = [];
    let faction = "";
    for (let i = si + 1; i < all.length; i++) {
      const el = all[i];
      if (el.tagName === "H2") break;
      if (el.tagName === "H4" || el.tagName === "H3") {
        const t = el.textContent.trim();
        faction = /alliance/i.test(t) ? "Alliance" : /horde/i.test(t) ? "Horde" : "";
        res.push({ t: "h", text: t });
      } else if (el.tagName === "UL") {
        const items = Array.from(el.querySelectorAll("li")).map((li) => li.textContent.trim());
        res.push({ t: "ul", f: faction, items });
      } else if (el.tagName === "TABLE") {
        const rows = Array.from(el.querySelectorAll("tr")).map((tr) =>
          Array.from(tr.querySelectorAll("td,th")).map((c) => c.textContent.trim())
        );
        res.push({ t: "table", f: faction, rows });
      } else if (el.tagName === "P") {
        // prose notes are the source site's writing — NOT copied (facts only)
      }
    }
    return res;
  }

  const parsed = {};
  for (const [slug, html] of Object.entries(window.__guides)) {
    const doc = new DOMParser().parseFromString(html, "text/html");
    const steps = parseSteps(doc);
    parsed[slug] = {
      shopping: parseShopping(doc),
      steps,
      gathering: parseGathering(doc), // always — fishing uses plain range H2s
      trainers: parseTrainers(doc),
    };
  }
  window.__parsed = parsed;
})();