// One-shot: converts scraped wow-professions data (content.txt) into
// web/src/guides-extra.js — our GUIDES array entries for all professions.
// Run: node scripts/gen-guides.js <content.txt>
const fs = require("fs");

const raw = fs.readFileSync(process.argv[2], "utf8");
// strip the tool's "Result: " wrapper + surrounding quotes, then unescape
let jsonText = raw.trim();
if (jsonText.startsWith("Result: ")) jsonText = jsonText.slice(8);
if (jsonText.startsWith('"')) jsonText = JSON.parse(jsonText); // unescape string
const scraped = JSON.parse(jsonText);

// slug -> display meta (icon texture names from wowhead/zamimg conventions)
const META = {
  alchemy: { name: "Alchemy", icon: "trade_alchemy", blurb: "Cheapest 1-300 route with full shopping lists. Pairs with Herbalism — grow your own herbs and save gold." },
  blacksmithing: { name: "Blacksmithing", icon: "trade_blacksmithing", blurb: "1-300 route covering armor and weapon smithing, plus specializations." },
  enchanting: { name: "Enchanting", icon: "trade_enchanting", blurb: "1-300 route with a disenchanting loop — level Enchanting while leveling other professions." },
  engineering: { name: "Engineering", icon: "trade_engineering", blurb: "1-300 route with the cheapest parts, plus specializations." },
  leatherworking: { name: "Leatherworking", icon: "trade_leatherworking", blurb: "1-300 route for light/medium/heavy leathers, plus specializations. Pairs with Skinning." },
  tailoring: { name: "Tailoring", icon: "trade_tailoring", blurb: "1-300 route — cheap cloth crafts, no gathering needed." },
  mining: { name: "Mining", icon: "trade_mining", blurb: "Zone-by-zone mining route 1-300 with smelting shortcuts and skill bonus items." },
  herbalism: { name: "Herbalism", icon: "trade_herbalism", blurb: "Zone-by-zone herb route 1-300 with skill bonus items." },
  skinning: { name: "Skinning", blurb: "Zone-by-zone skinning route 1-300 with skill requirements per rank.", icon: "trade_skinning" },
  "first-aid": { name: "First Aid", icon: "trade_alchemy", blurb: "1-300 route — bandages from linen to runecloth, cheapest cloth first." },
  fishing: { name: "Fishing", icon: "trade_fishing", blurb: "Zone-by-zone fishing route 1-300 — pair with Cooking for free food." },
};

// Standard rank requirements (universal WoW facts). A note is added when the
// tier CROSSES a rank boundary (75 / 150 / 225) — works for both crafting
// tiers (55-150) and gathering ranges (65-125).
function rankNote(range) {
  const m = String(range).match(/^(\d+)-(\d+)$/);
  if (!m) return "";
  const a = +m[1];
  const b = +m[2];
  // boundary B matters when this tier starts at/below B and continues past it
  const notes = [];
  if (a <= 75 && 75 < b) notes.push("Train Journeyman - requires skill 50 and level 10.");
  if (a <= 150 && 150 < b) notes.push("Train Expert - requires skill 125 and level 20.");
  if (a <= 225 && 225 < b) notes.push("Train Artisan - requires skill 200 and level 35.");
  return notes.join(" ");
}

// "1-55" style range -> {range, note, rows}; drops tiers with no rows
function toSteps(steps) {
  return steps
    .filter((s) => (s.rows || []).length)
    .map((s) => ({
      range: s.range,
      note: (s.notes || []).join(" ") || rankNote(s.range),
      rows: (s.rows || []).map((r) => ({
        skill: r.s,
        qty: r.q ? r.q + "x" : "",
        craft: { name: r.c, id: r.id ? Number(r.id) : null, k: r.k || "spell" },
        reagents: parseReagents(r.rg),
      })),
    }));
}

// "9 Silverleaf, 9 Empty Vial" -> [{qty, name, id:null}]
function parseReagents(str) {
  if (!str) return [];
  return str
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean)
    .map((part) => {
      const m = part.match(/^(\d+)\s+(.+)$/);
      return m ? { qty: m[1], name: m[2], id: null } : { qty: "", name: part, id: null };
    });
}

function toShopping(lists) {
  return lists.map((s) => ({
    range: s.range,
    items: s.items.map((it) => ({ qty: it.q, id: it.i ? Number(it.i) : null, name: it.n })),
  }));
}

// trainers: [{t:'h'|'ul'|'table'|'p', ...}] -> grouped sections
const HEADER_CELLS = /^(Race|Zone|Trainer|Vendor|Material|Faction|Level)$/i;
function toTrainers(list) {
  const out = { sections: [] };
  let cur = null;
  list.forEach((e) => {
    if (e.t === "h") {
      cur = { title: e.text, lines: [], rows: [], notes: [] };
      out.sections.push(cur);
    } else if (!cur) {
      cur = { title: "", lines: [], rows: [], notes: [] };
      out.sections.push(cur);
    }
    if (e.t === "ul") cur.lines.push(...e.items);
    else if (e.t === "table") {
      // drop header rows (first cell is a known column name)
      cur.rows.push(...e.rows.filter((r) => r[0] && !HEADER_CELLS.test(r[0])));
    } else if (e.t === "p") cur.notes.push(e.text);
  });
  return out;
}

// Fishing zones live in prose/lists the scraper doesn't map — manual facts.
const FISHING_ZONES = {
  "1-75": ["Durotar", "Mulgore", "Tirisfal Glades", "Zephras Isle", "Elwynn Forest", "Dun Morogh", "Teldrassil"],
  "75-150": ["Darnassus", "Stormwind City", "Orgrimmar", "Undercity", "Thunder Bluff"],
  "150-225": ["Dustwallow Marsh", "Alterac Mountains", "Arathi Highlands", "Desolace", "Stranglethorn Vale", "Swamp of Sorrows", "Thousand Needles"],
  "225-300": ["Felwood", "Feralas", "The Hinterlands", "Tanaris", "Un'Goro Crater", "Western Plaguelands"],
};
const FISHING_NOTES = {
  "1-75": "Buy a Fishing Pole and Shiny Baubles from a Fishing Supply vendor.",
  "150-225": "Expert Fishing: buy the book from Old Man Heming in Booty Bay (requires level 20 and Fishing 125).",
  "225-300": "Artisan Fishing is learned via quest.",
};

// gathering: [{range, zones:[names]}] -> steps with zone list + rank note.
// Tiers with neither zones nor a rank note are dropped (prose-only tiers).
function toGathering(list) {
  return list
    .map((g) => {
      const zones = (g.zones || []).length ? g.zones : FISHING_ZONES[g.range] || [];
      const note = [FISHING_NOTES[g.range] || "", rankNote(g.range)].filter(Boolean).join(" ");
      return { range: g.range, note, zones: zones.map((z) => [z]) };
    })
    .filter((g) => g.zones.length || g.note);
}

const guides = [];
for (const [slug, data] of Object.entries(scraped)) {
  const meta = META[slug];
  if (!meta) continue;
  let gathering = data.gathering || [];
  // Fishing's first tier H2 is "How to Fish in WoW Forever (1-75)" — the
  // scraper's range pattern misses it, so prepend it from the manual facts.
  if (slug === "fishing" && !gathering.some((g) => g.range === "1-75")) {
    gathering = [{ range: "1-75", zones: [] }, ...gathering];
  }
  const g = {
    id: slug,
    name: meta.name,
    icon: meta.icon,
    blurb: meta.blurb,
    shopping: toShopping(data.shopping || []),
    steps: toSteps(data.steps || []),
    gathering: toGathering(gathering),
    trainers: toTrainers(data.trainers || []),
  };
  guides.push(g);
}

const out =
  "// GENERATED from scraped guide facts (node scripts/gen-guides.js) — do not edit by hand.\n" +
  "// Facts only: skill ranges, quantities, item IDs, trainer names/zones.\n" +
  "export const GUIDES_EXTRA = " +
  JSON.stringify(guides, null, 2) +
  ";\n";

fs.writeFileSync("web/src/guides-extra.js", out);
console.log("wrote web/src/guides-extra.js:", guides.length, "guides,", out.length, "bytes");
guides.forEach((g) =>
  console.log(" -", g.name, "| shop:", g.shopping.length, "steps:", g.steps.length, "gather:", g.gathering.length, "trainer sections:", g.trainers.sections.length)
);