// API access. In production nginx proxies /api -> api container,
// in dev vue-cli proxies /api -> MUNKEHUB_API (see vue.config.js).
// So the base is always just "/api/v1" — no CORS, no hardcoded hosts.

async function get(path) {
  const res = await fetch(`/api/v1${path}`);
  if (!res.ok) throw new Error(`GET ${path} failed: ${res.status}`);
  return res.json();
}

// Bump on every release — shown in the footer so Emil can verify versions.
export const WEB_VERSION = "0.8.7";

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
};