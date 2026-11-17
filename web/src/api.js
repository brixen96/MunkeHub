// API access. In production nginx proxies /api -> api container,
// in dev vue-cli proxies /api -> MUNKEHUB_API (see vue.config.js).
// So the base is always just "/api/v1" — no CORS, no hardcoded hosts.

async function get(path) {
  const res = await fetch(`/api/v1${path}`);
  if (!res.ok) throw new Error(`GET ${path} failed: ${res.status}`);
  return res.json();
}

async function post(path, body) {
  const res = await fetch(`/api/v1${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error(`POST ${path} failed: ${res.status}`);
  return res.json();
}

export const api = {
  guild: () => get("/guild"),
  announcements: () => get("/announcements"),
  createAnnouncement: (data) => post("/announcements", data),
};