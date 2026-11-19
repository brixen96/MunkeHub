import { GUIDES_EXTRA } from "../web/src/guides-extra.js";

for (const g of GUIDES_EXTRA) {
  console.log("== " + g.name);
  if (g.gathering.length) {
    g.gathering.forEach((x) => {
      console.log("  range " + x.range + " zones:" + x.zones.length);
      x.zones.slice(0, 2).forEach((z) => console.log("    " + JSON.stringify(z)));
      console.log("    note: " + (x.note || "").slice(0, 120));
    });
  }
  if (g.steps.length) {
    const r = g.steps[0].rows && g.steps[0].rows[0];
    console.log("  first step: " + (r ? JSON.stringify(r).slice(0, 200) : "(no rows) " + JSON.stringify(g.steps[0]).slice(0, 150)));
  }
  const sec = g.trainers.sections[0];
  console.log("  trainer[0]: " + JSON.stringify(sec).slice(0, 200));
}