// API access. In production nginx proxies /api -> api container,
// in dev vue-cli proxies /api -> MUNKEHUB_API (see vue.config.js).
// So the base is always just "/api/v1" — no CORS, no hardcoded hosts.

async function get(path) {
  const res = await fetch(`/api/v1${path}`);
  if (!res.ok) throw new Error(`GET ${path} failed: ${res.status}`);
  return res.json();
}

// Bump on every release — shown in the footer so Emil can verify versions.
export const WEB_VERSION = "0.8.35";

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
  // My bags: character bag contents valued against cached listings.
  ahBags: (name) => get(`/ah/bags?name=${encodeURIComponent(name)}`),
  // MunkeHub AI advisor — cached 6h server-side; refresh=1 forces a new call.
  ahAdvice: (refresh = false) => get(`/ah/advice${refresh ? "?refresh=1" : ""}`),
  // Price history: batched sparkline points for up to 60 items.
  ahHistory: (ids) => get(`/ah/history?ids=${ids.join(",")}`),
  // Current price + 7-day context for a batch (Guides tab shopping lists).
  ahPrices: (ids) => get(`/ah/prices?ids=${ids.join(",")}`),
  // Merged bag scans from every character (no ?name= → guild-wide).
  ahBagsAll: () => get("/ah/bags"),
  // Weekly AI digest preview (auto-posts on Sundays).
  ahDigest: (refresh = false) => get(`/ah/digest${refresh ? "?refresh=1" : ""}`),
  // Presence: who is in game right now (helper heartbeats).
  presence: () => get("/presence"),
  // Quest log: addon -> helper -> API, with bag progress merged in.
  quests: (name) => get(`/quests?name=${encodeURIComponent(name)}`),
  // AI quest advisor — cached 1h server-side; refresh=1 forces a new call.
  questAdvice: (name, refresh = false) =>
    get(`/quests/advice?name=${encodeURIComponent(name)}${refresh ? "&refresh=1" : ""}`),
  // Which quest IDs other characters share with this one (group-up chips).
  questOverlap: (name) => get(`/quests/overlap?name=${encodeURIComponent(name)}`),
  // Leveling route brackets for a level + race (facts from guides.md sources).
  leveling: (level, race) =>
    get(`/leveling?level=${level}&race=${encodeURIComponent(race || "")}`),
  // Guild Board: crafting orders + events (reads public, writes member token).
  board: () => get("/board"),
  boardOrderAdd: (data, token) =>
    fetch("/api/v1/board/orders", {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
      body: JSON.stringify(data),
    }).then((r) => r.json()),
  boardOrderStatus: (id, data, token) =>
    fetch(`/api/v1/board/orders/${id}`, {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
      body: JSON.stringify(data),
    }).then((r) => r.json()),
  boardOrderDelete: (id, token) =>
    fetch(`/api/v1/board/orders/${id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    }).then((r) => r.json()),
  boardEventAdd: (data, token) =>
    fetch("/api/v1/board/events", {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
      body: JSON.stringify(data),
    }).then((r) => r.json()),
  boardEventDelete: (id, token) =>
    fetch(`/api/v1/board/events/${id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    }).then((r) => r.json()),
  boardRsvp: (id, data, token) =>
    fetch(`/api/v1/board/events/${id}/rsvp`, {
      method: "POST",
      headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
      body: JSON.stringify(data),
    }).then((r) => r.json()),
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