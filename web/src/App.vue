<template>
  <div class="page">
    <!-- ================= HEADER ================= -->
    <header class="banner">
      <div class="banner-inner">
        <div class="crest">🐒</div>
        <div class="titles">
          <h1>MunkeHub</h1>
          <p class="guild-line">
            <span class="guild-name">Munkeklubben</span>
            <span class="sep">•</span>
            <span v-if="guild">{{ guild.server }}</span>
            <span v-else-if="guildError" class="error">API unreachable</span>
            <span v-else class="muted">…</span>
          </p>
        </div>
        <a class="download-btn" href="/api/v1/download" download>
          ⬇ Download Addon
        </a>
      </div>
      <nav class="tabs">
        <button :class="{ active: tab === 'roster' }" @click="tab = 'roster'">
          Guild Roster
        </button>
        <button :class="{ active: tab === 'news' }" @click="tab = 'news'">
          Announcements
        </button>
      </nav>
    </header>

    <main>
      <!-- ================= ROSTER ================= -->
      <section v-show="tab === 'roster'">
        <div class="panel">
          <h2 class="panel-title">Characters</h2>
          <p v-if="charError" class="error">Could not load the roster.</p>
          <p v-else-if="!characters" class="muted">Loading…</p>
          <p v-else-if="characters.length === 0" class="muted">
            No characters synced yet — run the MunkeHub helper in-game.
          </p>
          <table v-else class="roster">
            <thead>
              <tr>
                <th>Character</th>
                <th>Race</th>
                <th>Class</th>
                <th>Lvl</th>
                <th>iLvl</th>
                <th>Owner</th>
                <th>Updated</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="c in characters" :key="c.id">
                <td data-label="Character">
                  <span class="char-name" :class="classSlug(c.class)">
                    {{ c.name }}
                  </span>
                  <span class="realm">{{ c.realm }}</span>
                </td>
                <td data-label="Race">{{ c.race }}</td>
                <td data-label="Class" :class="classSlug(c.class)">
                  {{ c.class }}
                </td>
                <td data-label="Lvl" class="lvl">{{ c.level }}</td>
                <td data-label="iLvl">{{ c.ilvl }}</td>
                <td data-label="Owner" class="owner">{{ c.owner }}</td>
                <td data-label="Updated" class="when">
                  {{ formatDate(c.updated_at) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- ================= ANNOUNCEMENTS ================= -->
      <section v-show="tab === 'news'">
        <div class="panel">
          <h2 class="panel-title">Guild News</h2>
          <p v-if="annError" class="error">Could not load announcements.</p>
          <p v-else-if="!announcements" class="muted">Loading…</p>
          <p v-else-if="announcements.length === 0" class="muted">
            The notice board is empty.
          </p>
          <article v-for="a in announcements" :key="a.id" class="notice">
            <div class="notice-head">
              <h3>{{ a.title }}</h3>
              <time>{{ formatDate(a.created_at) }}</time>
            </div>
            <p>{{ a.body }}</p>
          </article>

          <details class="post-box">
            <summary>Post an announcement (guild members)</summary>
            <form @submit.prevent="postAnnouncement">
              <input
                v-model="newTitle"
                type="text"
                placeholder="Title"
                required
                maxlength="120"
              />
              <textarea
                v-model="newBody"
                placeholder="Message for Munkeklubben…"
                rows="3"
              ></textarea>
              <div class="post-row">
                <input
                  v-model="postToken"
                  type="password"
                  placeholder="Your member token"
                  required
                />
                <button type="submit" :disabled="posting">
                  {{ posting ? "Posting…" : "Post" }}
                </button>
              </div>
              <p v-if="postError" class="error">{{ postError }}</p>
              <p v-if="postOk" class="ok">Posted ✓</p>
            </form>
          </details>
        </div>
      </section>
    </main>

    <footer>
      <span>MunkeHub — private hall of Munkeklubben • WoW: Forever</span>
    </footer>
  </div>
</template>

<script>
import { api } from "./api";

// WoW class colours — iconic, instantly readable for players
const CLASS_COLORS = {
  warrior: "#C79C6E",
  paladin: "#F58CBA",
  hunter: "#AAD372",
  rogue: "#FFF569",
  priest: "#FFFFFF",
  deathknight: "#C41E3A",
  shaman: "#0070DE",
  mage: "#69CCF0",
  warlock: "#9482C9",
  monk: "#00FF98",
  druid: "#FF7D0A",
  demonhunter: "#A330C9",
  evoker: "#33937F",
};

export default {
  name: "App",
  data() {
    return {
      tab: "roster",
      guild: null,
      guildError: false,
      announcements: null,
      annError: false,
      characters: null,
      charError: false,
      newTitle: "",
      newBody: "",
      postToken: localStorage.getItem("munkehub_token") || "",
      posting: false,
      postError: "",
      postOk: false,
    };
  },
  async created() {
    try {
      this.guild = await api.guild();
    } catch {
      this.guildError = true;
    }
    try {
      this.announcements = await api.announcements();
    } catch {
      this.annError = true;
    }
    try {
      this.characters = await api.characters();
    } catch {
      this.charError = true;
    }
  },
  methods: {
    formatDate(iso) {
      if (!iso) return "";
      return new Date(iso).toLocaleString();
    },
    classSlug(cls) {
      return "c-" + (cls || "").toLowerCase().replace(/[^a-z]/g, "");
    },
    async postAnnouncement() {
      this.posting = true;
      this.postError = "";
      this.postOk = false;
      try {
        await api.createAnnouncement(
          { title: this.newTitle, body: this.newBody },
          this.postToken
        );
        localStorage.setItem("munkehub_token", this.postToken);
        this.postOk = true;
        this.newTitle = "";
        this.newBody = "";
        this.announcements = await api.announcements();
      } catch {
        this.postError = "Post failed — check your token.";
      } finally {
        this.posting = false;
      }
    },
  },
};
</script>

<style>
/* ================================================================
   MunkeHub — WoW-flavoured guild hub
   Dark stone panels, gold trim, class colours, parchment notices
   ================================================================ */
:root {
  color-scheme: dark;
  --gold: #c8a44d;
  --gold-bright: #ffd166;
  --stone: #1c1b1f;
  --stone-light: #26252b;
  --edge: #3a3840;
  --ink: #dcd7cc;
  --muted: #8f8a7e;
}
* {
  box-sizing: border-box;
}
body {
  margin: 0;
  font-family: "Segoe UI", system-ui, -apple-system, sans-serif;
  background: radial-gradient(ellipse at top, #23212a 0%, #131217 70%) fixed;
  color: var(--ink);
}
.page {
  max-width: 860px;
  margin: 0 auto;
  padding: 1.5rem 1rem 3rem;
}

/* ---------------- banner ---------------- */
.banner {
  border: 1px solid var(--gold);
  border-radius: 10px;
  background: linear-gradient(180deg, #2b2933 0%, #1a1920 100%);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}
.banner-inner {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.4rem 1.5rem 1rem;
}
.crest {
  font-size: 2.6rem;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.6));
}
.download-btn {
  margin-left: auto;
  align-self: center;
  background: linear-gradient(180deg, #d9b45c, #a8842f);
  border: 1px solid var(--gold);
  border-radius: 6px;
  color: #1a1508;
  font-weight: 700;
  font-size: 0.9rem;
  padding: 0.55rem 1.1rem;
  text-decoration: none;
  white-space: nowrap;
  transition: filter 0.15s;
}
.download-btn:hover {
  filter: brightness(1.12);
}
.titles h1 {
  margin: 0;
  font-family: Georgia, "Times New Roman", serif;
  font-size: 2rem;
  letter-spacing: 0.06em;
  color: var(--gold-bright);
  text-shadow: 0 2px 6px rgba(0, 0, 0, 0.8);
}
.guild-line {
  margin: 0.2rem 0 0;
  color: var(--muted);
  font-size: 0.95rem;
}
.guild-name {
  color: var(--gold);
  font-weight: 600;
}
.sep {
  margin: 0 0.4rem;
  opacity: 0.5;
}

/* ---------------- tabs ---------------- */
.tabs {
  display: flex;
  border-top: 1px solid var(--edge);
  background: rgba(0, 0, 0, 0.25);
}
.tabs button {
  flex: 1;
  padding: 0.7rem 1rem;
  background: transparent;
  border: none;
  border-right: 1px solid var(--edge);
  color: var(--muted);
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
}
.tabs button:last-child {
  border-right: none;
}
.tabs button:hover {
  color: var(--ink);
}
.tabs button.active {
  color: var(--gold-bright);
  background: rgba(200, 164, 77, 0.1);
  box-shadow: inset 0 -2px 0 var(--gold);
}

/* ---------------- panels ---------------- */
main {
  margin-top: 1.4rem;
}
.panel {
  background: var(--stone);
  border: 1px solid var(--edge);
  border-radius: 10px;
  padding: 1.2rem 1.3rem;
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.4);
}
.panel-title {
  margin: 0 0 1rem;
  font-family: Georgia, serif;
  font-size: 1.15rem;
  color: var(--gold);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  border-bottom: 1px solid var(--edge);
  padding-bottom: 0.6rem;
}

/* ---------------- roster ---------------- */
.roster {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.92rem;
}
.roster th,
.roster td {
  text-align: left;
  padding: 0.55rem 0.6rem;
  border-bottom: 1px solid var(--edge);
}
.roster th {
  color: var(--muted);
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.roster tbody tr:hover {
  background: rgba(255, 255, 255, 0.03);
}
.char-name {
  font-weight: 700;
  font-size: 1rem;
}
.realm {
  display: block;
  color: var(--muted);
  font-size: 0.75rem;
  font-weight: 400;
}
.lvl {
  color: var(--gold-bright);
  font-weight: 600;
}
.owner {
  color: var(--muted);
}
.when {
  color: var(--muted);
  font-size: 0.8rem;
}

/* class colours (via .c-<class> on name + class cells) */
.c-warrior { color: #c79c6e; }
.c-paladin { color: #f58cba; }
.c-hunter { color: #aad372; }
.c-rogue { color: #fff569; }
.c-priest { color: #ffffff; }
.c-deathknight { color: #c41e3a; }
.c-shaman { color: #0070de; }
.c-mage { color: #69ccf0; }
.c-warlock { color: #9482c9; }
.c-monk { color: #00ff98; }
.c-druid { color: #ff7d0a; }
.c-demonhunter { color: #a330c9; }
.c-evoker { color: #33937f; }

/* ---------------- notices ---------------- */
.notice {
  background: var(--stone-light);
  border: 1px solid var(--edge);
  border-left: 3px solid var(--gold);
  border-radius: 6px;
  padding: 0.9rem 1.1rem;
  margin: 0.8rem 0;
}
.notice-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 1rem;
}
.notice h3 {
  margin: 0;
  color: var(--gold-bright);
  font-size: 1.02rem;
}
.notice time {
  color: var(--muted);
  font-size: 0.78rem;
  white-space: nowrap;
}
.notice p {
  margin: 0.5rem 0 0;
  line-height: 1.5;
}

/* ---------------- post form ---------------- */
.post-box {
  margin-top: 1.5rem;
  background: var(--stone-light);
  border: 1px dashed var(--edge);
  border-radius: 8px;
  padding: 0.75rem 1rem;
}
.post-box summary {
  cursor: pointer;
  color: var(--muted);
  font-size: 0.9rem;
}
.post-box form {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: 0.75rem;
}
.post-box input,
.post-box textarea {
  background: #131217;
  border: 1px solid var(--edge);
  border-radius: 6px;
  color: var(--ink);
  padding: 0.55rem 0.7rem;
  font: inherit;
}
.post-box input:focus,
.post-box textarea:focus {
  outline: none;
  border-color: var(--gold);
}
.post-row {
  display: flex;
  gap: 0.5rem;
}
.post-row input {
  flex: 1;
}
.post-box button {
  background: linear-gradient(180deg, #d9b45c, #a8842f);
  border: 1px solid var(--gold);
  border-radius: 6px;
  color: #1a1508;
  font-weight: 700;
  padding: 0.5rem 1.4rem;
  cursor: pointer;
}
.post-box button:hover {
  filter: brightness(1.1);
}
.post-box button:disabled {
  opacity: 0.6;
  cursor: default;
}

/* ---------------- misc ---------------- */
.muted { color: var(--muted); }
.error { color: #ff6b6b; }
.ok { color: #7ec97e; }
footer {
  margin-top: 2.5rem;
  font-size: 0.78rem;
  color: var(--muted);
  text-align: center;
  opacity: 0.8;
}

/* ---------------- mobile ---------------- */
@media (max-width: 640px) {
  .banner-inner { padding: 1rem; }
  .titles h1 { font-size: 1.5rem; }
  .crest { font-size: 2rem; }

  /* table → stacked cards */
  .roster thead { display: none; }
  .roster, .roster tbody, .roster tr, .roster td { display: block; width: 100%; }
  .roster tr {
    border: 1px solid var(--edge);
    border-radius: 8px;
    margin-bottom: 0.6rem;
    padding: 0.4rem 0.6rem;
  }
  .roster td {
    border: none;
    padding: 0.25rem 0;
    display: flex;
    justify-content: space-between;
  }
  .roster td::before {
    content: attr(data-label);
    color: var(--muted);
    font-size: 0.75rem;
    text-transform: uppercase;
  }
  .realm { display: inline; margin-left: 0.4rem; }
}
</style>