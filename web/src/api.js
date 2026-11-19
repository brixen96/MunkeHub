// API access. In production nginx proxies /api -> api container,
// in dev vue-cli proxies /api -> MUNKEHUB_API (see vue.config.js).
// So the base is always just "/api/v1" — no CORS, no hardcoded hosts.

async function get(path) {
  const res = await fetch(`/api/v1${path}`);
  if (!res.ok) throw new Error(`GET ${path} failed: ${res.status}`);
  return res.json();
}

// Bump on every release — shown in the footer so Emil can verify versions.
export const WEB_VERSION = "0.8.22";

export const api = {
  health: () => fetch("/health").then((r) => r.json()),
  guild: () => get("/guild"),
  announcements: () => get("/announcements"),
  createAnnouncement: (data, token) =>
    fetch("/api/v1/announcements", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    }).then(async (res) => {
      if (!res.ok) throw new Error(`POST failed: ${res.status}`);
      return res.json();
    }),
  characters: () => get("/characters"),
  stats: () => get("/stats"),
  // Character sheet {model, talents, classId} — proxied from the armory's
  // API (their site sends no CORS headers).
  charModel: (slug) => get("/charmodel/" + encodeURIComponent(slug)),
  // Talent tree layout for a class (grid positions, icons, prereqs) —
  // static per class, cached 24h server-side.
  talentLayout: (classId) => get("/talents/" + classId),
  // Auction House — all served from our SQLite cache of blackwater.exchange
  // (poller refreshes hourly), so these are instant and never hit upstream.
  ahSummary: () => get("/ah/summary"),
  ahDeals: () => get("/ah/deals"),
  ahSales: (days = "1") => get(`/ah/sales?days=${days}`),
  ahListings: (q = "", page = 1, limit = 50) =>
    get(
      `/ah/listings?q=${encodeURIComponent(q)}&page=${page}&limit=${limit}`
    ),
  ahCrafting: (profession = "171", sort = "profit") =>
    get(`/ah/crafting?profession=${profession}&sort=${sort}`),
  ahStatus: () => get("/ah/status"),
  // Money makers: flipping opportunities + vendor snipes (pure math, no AI).
  ahFlips: () => get("/ah/flips"),
  ahVendorSnipes: () => get("/ah/vendor-snipes"),
  // Farming table: guild professions × market (ranked by market/day).
  ahFarming: () => get("/ah/farming"),
  // MunkeHub AI advisor — cached 6h server-side; refresh=1 forces a new call.
  ahAdvice: (refresh = false) => get(`/ah/advice${refresh ? "?refresh=1" : ""}`),
  // Price history: batched sparkline points for up to 60 items.
  ahHistory: (ids) => get(`/ah/history?ids=${ids.join(",")}`),
  // Guild watchlist (⭐) — friends-only site, no auth by design.
  ahWatchlist: () => get("/ah/watchlist"),
  ahWatchAdd: (id) =>
    fetch(`/api/v1/ah/watchlist`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ item_id: id }),
    }).then(async (res) => {
      if (!res.ok) throw new Error(`POST failed: ${res.status}`);
      return res.json();
    }),
  ahWatchRemove: (id) =>
    fetch(`/api/v1/ah/watchlist/${id}`, { method: "DELETE" }).then(async (res) => {
      if (!res.ok) throw new Error(`DELETE failed: ${res.status}`);
      return res.json();
    }),
};