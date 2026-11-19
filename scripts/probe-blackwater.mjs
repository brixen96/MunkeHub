// Probe script for the blackwater.exchange API — the Auction House data
// source behind /api/v1/ah/*. Run it when the upstream API changes shape:
//
//   node scripts/probe-blackwater.mjs
//
// It prints the status + top-level shape of every endpoint the poller uses.
// The Go poller (api/ah.go) follows the exact same flow:
//
//   1. POST /api/site/session  → sets the bx_market cookie (required, or 401)
//   2. GET  /api/site/markets/eu/15/124/<endpoint> with that cookie
//
// Upstream quirks discovered by probing (do not "fix" these in the poller):
//   - /items accepts ONLY limit=25|50|100 — anything else returns
//     400 {"error":"browse_range"} (the error name is misleading).
//   - /items serves only the first 2,000 matches (page ≤ 20 at limit=100).
//   - /crafting accepts per ≤ 100 (200+ silently falls back to 25).
//   - 429 {"error":"rate_limited","retry_after":N} — hourly budget; the
//     poller paces requests 1.2s apart and aborts the cycle on 429.
//   - Calls need a Referer pointing at the market page.

const BASE = "https://blackwater.exchange";
const MARKET = "/m/eu/15/124";
const API = "/api/site/markets/eu/15/124";
const UA =
  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36";

const cookieJar = new Map();
const cookieHeader = () =>
  [...cookieJar.entries()].map(([k, v]) => `${k}=${v}`).join("; ");

function absorb(res) {
  const list =
    typeof res.headers.getSetCookie === "function"
      ? res.headers.getSetCookie()
      : [];
  for (const c of list) {
    const [pair] = c.split(";");
    const eq = pair.indexOf("=");
    if (eq > 0) cookieJar.set(pair.slice(0, eq).trim(), pair.slice(eq + 1).trim());
  }
}

async function get(path) {
  const res = await fetch(BASE + path, {
    headers: {
      Accept: "application/json",
      "User-Agent": UA,
      Referer: BASE + MARKET,
      ...(cookieJar.size ? { Cookie: cookieHeader() } : {}),
    },
  });
  absorb(res);
  const body = await res.text();
  return { status: res.status, body };
}

const sess = await fetch(BASE + "/api/site/session", {
  method: "POST",
  headers: {
    Accept: "*/*",
    "User-Agent": UA,
    Referer: BASE + MARKET,
    Origin: BASE,
  },
});
absorb(sess);
console.log(`session: ${sess.status} cookies: [${[...cookieJar.keys()]}]`);

const eps = {
  summary: `${API}/summary`,
  deals: `${API}/deals?days=7&discount=20&sold=1&sort=profit&limit=100`,
  "sales(1d)": `${API}/sales?days=1&sort=value&limit=100`,
  "sales(7d)": `${API}/sales?days=7&sort=value&limit=100`,
  items: `${API}/items?sort=sold&dir=desc&limit=100&page=1`,
  crafting: `${API}/crafting?profession=171&per=100&page=1&filter=priced&sort=profit&dir=desc`,
};

for (const [name, path] of Object.entries(eps)) {
  const r = await get(path);
  if (r.status !== 200) {
    console.log(`${name}: ${r.status} ${r.body.slice(0, 100)}`);
    continue;
  }
  const j = JSON.parse(r.body);
  const list = j.items || j.recipes;
  const count = Array.isArray(list) ? list.length : "-";
  const extra =
    j.all !== undefined
      ? ` all=${j.all}`
      : j.total !== undefined
        ? ` total=${j.total}`
        : "";
  console.log(`${name}: 200 rows=${count}${extra} keys=[${Object.keys(j).join(",")}]`);
  if (Array.isArray(list) && list[0]) {
    console.log(`  row keys: [${Object.keys(list[0]).join(",")}]`);
  }
}