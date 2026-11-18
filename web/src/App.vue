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
        <button
          class="download-btn"
          :class="{ active: tab === 'guide' }"
          @click="tab = 'guide'"
        >
          📖 Setup Guide
        </button>
      </div>
      <nav class="tabs">
        <button :class="{ active: tab === 'roster' }" @click="tab = 'roster'">
          Guild Roster
        </button>
        <button :class="{ active: tab === 'news' }" @click="tab = 'news'">
          Announcements
        </button>
        <button :class="{ active: tab === 'changelog' }" @click="tab = 'changelog'">
          Changelog
        </button>
      </nav>
    </header>

    <main>
      <!-- ================= ROSTER ================= -->
      <section v-show="tab === 'roster'">
        <!-- dashboard cards -->
        <div class="dash" v-if="stats">
          <div class="dash-card">
            <span class="dash-num">{{ stats.characters }}</span>
            <span class="dash-label">Characters</span>
          </div>
          <div class="dash-card">
            <span class="dash-num">{{ stats.avg_level }}</span>
            <span class="dash-label">Avg level</span>
          </div>
          <div class="dash-card">
            <span class="dash-num">{{ stats.avg_ilvl }}</span>
            <span class="dash-label">Avg iLvl</span>
          </div>
          <div class="dash-card">
            <span class="dash-num">{{ stats.classes.length }}</span>
            <span class="dash-label">Classes</span>
          </div>
        </div>
        <div class="classbar" v-if="stats && stats.classes.length">
          <span v-for="cc in stats.classes" :key="cc.class"
                class="class-chip" :class="classSlug(cc.class)">
            {{ cc.class }} ×{{ cc.count }}
          </span>
        </div>

        <div class="panel">
          <h2 class="panel-title">Characters <span class="hint">— click a row to expand</span></h2>
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
                <th>Updated</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="c in characters" :key="c.id">
                <tr class="clickable" @click="toggleRow(c.id)">
                  <td data-label="Character">
                    <span class="chevron" :class="{ open: openRow === c.id }">▸</span>
                    <span class="avatar" :class="classSlug(c.class)">
                      <img v-if="avatarSrc(c)" :src="avatarSrc(c)" alt=""
                           loading="lazy" @error="$event.target.style.display = 'none'">
                      <span v-else v-html="classGlyph(c.class)"></span>
                    </span>
                    <span class="char-name" :class="classSlug(c.class)">
                      {{ c.name }}<template v-if="c.last_name"> {{ c.last_name }}</template>
                    </span>
                    <span class="realm">{{ c.realm }}</span>
                  </td>
                  <td data-label="Race">{{ c.race }}</td>
                  <td data-label="Class" :class="classSlug(c.class)">
                    {{ c.class }}
                  </td>
                  <td data-label="Lvl" class="lvl">{{ c.level }}</td>
                  <td data-label="iLvl">{{ c.ilvl }}</td>
                  <td data-label="Updated" class="when">
                    {{ formatDate(c.updated_at) }}
                  </td>
                </tr>
                <tr v-if="openRow === c.id" class="detail-row">
                  <td colspan="6">
                    <div class="detail-panel">
                      <div class="detail-core">
                        <div class="detail-item">
                          <span class="detail-k">Spec</span>
                          <span class="detail-v">{{ c.spec || '—' }}</span>
                        </div>
                        <div class="detail-item">
                          <span class="detail-k">Zone</span>
                          <span class="detail-v">{{ c.zone || '—' }}</span>
                        </div>
                        <div class="detail-item">
                          <span class="detail-k">Played</span>
                          <span class="detail-v">{{ formatPlayed(c) }}</span>
                        </div>
                      </div>
                      <div v-if="parseGear(c).length" class="gear-grid">
                        <div v-for="g in parseGear(c)" :key="g.slot" class="gear-slot"
                             :class="'q' + g.q" :title="g.name">
                          <img v-if="g.fid || g.id" class="gear-icon"
                               :src="'/api/v1/icon/' + (g.fid || g.id)" alt=""
                               loading="lazy" @error="$event.target.style.display = 'none'">
                          <span class="gear-slot-name">{{ g.slot }}</span>
                          <span class="gear-item">{{ g.name || 'empty' }}</span>
                          <span class="gear-ilvl">{{ g.ilvl || '' }}</span>
                        </div>
                      </div>
                      <p v-else class="muted small">
                        No gear detail yet — /reload in game, then the helper uploads it.
                      </p>
                    </div>
                  </td>
                </tr>
              </template>
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

      <!-- ================= SETUP GUIDE ================= -->
      <section v-show="tab === 'guide'">
        <div class="panel">
          <h2 class="panel-title">Setup Guide</h2>

          <div class="guide-intro">
            <p>
              <strong>MunkeHub</strong> shows our guild on this website — every
              character and all guild news in one place. Getting your character
              on here takes about 2 minutes and you only do it once.
            </p>
          </div>

          <div class="explain-card">
            <h3>🧩 What is the Add-on?</h3>
            <p>
              A small file for World of Warcraft that reads your character
              (name, level, class, gear) and shares it with the guild.
              <strong>It does not change the game</strong> — no cheats, no
              gameplay changes. It just reads your character sheet.
            </p>
          </div>

          <div class="explain-card">
            <h3>🤖 What is the Helper?</h3>
            <p>
              A tiny program that uploads your stats to the website
              automatically. <strong>You do NOT need it</strong> — without it,
              your stats sync whenever Emil is online in-game. Run it only if
              you want your stats to update instantly.
            </p>
          </div>

          <h3 class="steps-title">Step by step</h3>
          <ol class="steps">
            <li>
              <strong>Download</strong> the package with the button below.
            </li>
            <li>
              <strong>Open the zip file</strong> → right-click it →
              <em>Extract All</em> → Extract.
            </li>
            <li>
              <strong>Copy the <code>MunkeHub</code> folder</strong> into your
              WoW AddOns folder:
              <code class="path">World of Warcraft\_forever_\Interface\AddOns</code>
            </li>
            <li>
              <strong>Start World of Warcraft.</strong> At the character select
              screen, click <em>AddOns</em> (bottom left) and make sure
              <strong>MunkeHub</strong> is ticked.
            </li>
            <li>
              <strong>Log in</strong> and type: <code>/reload</code>.
            </li>
            <li>
              <strong>Done!</strong> Your character appears on the
              <a href="#" @click.prevent="tab = 'roster'">Guild Roster</a>
              within a few minutes.
            </li>
          </ol>

          <div class="guide-download">
            <a class="download-btn big" href="/api/v1/download" download>
              ⬇ Download MunkeHub Setup
            </a>
            <p class="muted small">
              Contains: the add-on, the helper, and a README with these steps.
            </p>
          </div>

          <details class="post-box">
            <summary>Optional: instant sync with the Helper (for the impatient)</summary>
            <div class="helper-steps">
              <ol>
                <li>Ask Emil for your personal token.</li>
                <li>Double-click <code>munkehub-helper.exe</code> from the zip.</li>
                <li>Paste your token when asked — it's remembered from then on.</li>
                <li>Type <code>munkehub-helper.exe -install</code> once — it then
                  starts with Windows, silently in the background.</li>
              </ol>
              <p class="muted small">
                Your stats then appear on the website within ~10 seconds of
                logging in or typing <code>/reload</code> in game.
              </p>
              <p class="muted small">
                The helper updates itself automatically — you never need to
                download a new version.
              </p>
            </div>
          </details>
        </div>
      </section>

      <!-- ================= CHANGELOG ================= -->
      <section v-show="tab === 'changelog'">
        <div class="panel">
          <h2 class="panel-title">
            Changelog
            <span class="cl-current">v{{ webVersion }}</span>
          </h2>
          <p class="cl-intro muted">
            What changed in MunkeHub — newest first.
          </p>
          <div v-for="rel in changelog" :key="rel.version" class="cl-release">
            <div class="cl-head">
              <span class="cl-ver">v{{ rel.version }}</span>
              <span class="cl-date">{{ rel.date }}</span>
            </div>
            <ul class="cl-list">
              <li v-for="(e, i) in rel.entries" :key="i">{{ e }}</li>
            </ul>
          </div>
        </div>
      </section>

      </main>

    <footer>
      <span>
        MunkeHub v{{ webVersion }}
        <template v-if="apiVersion"> · API v{{ apiVersion }}</template>
        — private hall of Munkeklubben • WoW: Forever
      </span>
    </footer>
  </div>
</template>

<script>
import { api, WEB_VERSION } from "./api";
import { CHANGELOG } from "./changelog";

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

// Class icon texture fileIDs (from the Forever client — served by
// foreverdb.net through our icon endpoint). Used when no spec icon exists.
const CLASS_ICON_FID = {
  warrior: 626008,
  paladin: 626003,
  hunter: 626000,
  rogue: 626005,
  priest: 626004,
  shaman: 626006,
  mage: 626001,
  warlock: 626007,
  druid: 625999,
};

export default {
  name: "App",
  data() {
    return {
      tab: "roster",
      openRow: null,
      stats: null,
      guild: null,
      webVersion: WEB_VERSION,
      apiVersion: "",
      changelog: CHANGELOG,
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
    try {
      this.stats = await api.stats();
    } catch {
      this.stats = null;
    }
    try {
      const h = await api.health();
      this.apiVersion = h.version || "";
    } catch {
      this.apiVersion = "?";
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
    // Avatar: spec icon when available, else the class icon (fileIDs served
    // by foreverdb.net through our icon endpoint), else the SVG glyph below.
    avatarSrc(c) {
      if (c.spec_icon) return "/api/v1/icon/" + c.spec_icon;
      const fid = CLASS_ICON_FID[(c.class || "").toLowerCase()];
      return fid ? "/api/v1/icon/" + fid : "";
    },
    // Original simple SVG glyphs per class (drawn for MunkeHub — no game assets)
    classGlyph(cls) {
      const g = {
        warrior: '<path d="M12 3l2 6 6 2-6 2-2 6-2-6-6-2 6-2z"/>', // blade star
        paladin: '<path d="M12 3l7 3v6c0 4-3 7-7 9-4-2-7-5-7-9V6z"/>', // shield
        hunter: '<path d="M4 20L18 6M18 6l-1 6M18 6l-6 1M18 6l1 6M18 6l-1-6" stroke-width="2" fill="none" stroke="currentColor"/>', // arrow
        rogue: '<path d="M12 3l3 9-3 9-3-9z M7 12h10" stroke-width="1.5" fill="currentColor"/>', // dagger
        priest: '<path d="M10 4h4v6h6v4h-6v6h-4v-6H4v-4h6z"/>', // cross
        deathknight: '<path d="M12 3v18M4.2 7.5l15.6 9M19.8 7.5l-15.6 9" stroke-width="2" stroke="currentColor" fill="none"/>', // rune star
        shaman: '<path d="M8 4h8v4H8zM7 10h10v4H7zM8 16h8v4H8z"/>', // totem
        mage: '<path d="M12 2l2.5 7.5L22 12l-7.5 2.5L12 22l-2.5-7.5L2 12l7.5-2.5z"/>', // arcane star
        warlock: '<path d="M12 3c4 5 6 7 6 11a6 6 0 11-12 0c0-4 2-6 6-11z"/>', // flame
        monk: '<circle cx="12" cy="12" r="7" fill="none" stroke="currentColor" stroke-width="3" stroke-dasharray="30 14"/>', // enso
        druid: '<path d="M12 3c5 3 7 7 6 12-4 1-8-1-9-5-1-3 1-6 3-7zM11 12c-1 3-3 5-6 6 1-5 3-7 6-6z"/>', // leaf
        demonhunter: '<path d="M4 6c3 0 5 2 6 6-2-1-4-1-6 1zM20 6c-3 0-5 2-6 6 2-1 4-1 6 1z"/>', // horns
        evoker: '<path d="M12 3l8 9-8 9-8-9z"/>', // scale
      };
      const slug = (cls || "").toLowerCase().replace(/[^a-z]/g, "");
      const body = g[slug] || g.warrior;
      return '<svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">' + body + '</svg>';
    },
    toggleRow(id) {
      this.openRow = this.openRow === id ? null : id;
    },
    parseGear(c) {
      if (!c.detail) return [];
      try {
        const d = JSON.parse(c.detail);
        return Array.isArray(d) ? d : [];
      } catch {
        return [];
      }
    },
    formatPlayed(c) {
      const h = Math.floor((c.played || 0) / 3600);
      return h > 0 ? h + "h" : "—";
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
  font-family: inherit;
  font-weight: 700;
  font-size: 0.9rem;
  padding: 0.55rem 1.1rem;
  text-decoration: none;
  white-space: nowrap;
  cursor: pointer;
  transition: filter 0.15s;
}
.download-btn:hover {
  filter: brightness(1.12);
}
.download-btn.active {
  box-shadow: 0 0 0 2px var(--gold-bright);
}
.download-btn.big {
  font-size: 1.1rem;
  padding: 0.8rem 1.8rem;
}

/* ---------------- setup guide ---------------- */
.guide-intro p {
  margin: 0 0 1rem;
  line-height: 1.6;
}
.explain-card {
  background: var(--stone-light);
  border: 1px solid var(--edge);
  border-radius: 8px;
  padding: 0.9rem 1.1rem;
  margin: 0.8rem 0;
}
.explain-card h3 {
  margin: 0 0 0.4rem;
  color: var(--gold-bright);
  font-size: 1rem;
}
.explain-card p {
  margin: 0;
  line-height: 1.55;
  color: var(--ink);
}
.steps-title {
  margin: 1.5rem 0 0.6rem;
  font-family: Georgia, serif;
  color: var(--gold);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-size: 1rem;
}
.steps {
  margin: 0;
  padding-left: 1.4rem;
}
.steps li {
  margin: 0.55rem 0;
  line-height: 1.55;
}
/* ---------------- changelog ---------------- */
.cl-current {
  float: right;
  font-size: 0.85rem;
  color: var(--gold-bright);
  letter-spacing: 0.04em;
}
.cl-intro {
  margin: -0.4rem 0 1rem;
  font-size: 0.92rem;
}
.cl-release {
  background: var(--stone-light);
  border: 1px solid var(--edge);
  border-radius: 8px;
  padding: 0.8rem 1.1rem;
  margin: 0.7rem 0;
}
.cl-head {
  display: flex;
  align-items: baseline;
  gap: 0.7rem;
  margin-bottom: 0.35rem;
}
.cl-ver {
  font-family: Georgia, serif;
  font-weight: bold;
  color: var(--gold-bright);
}
.cl-date {
  font-size: 0.82rem;
  color: var(--muted);
}
.cl-list {
  margin: 0;
  padding-left: 1.3rem;
}
.cl-list li {
  margin: 0.3rem 0;
  line-height: 1.5;
}

.steps code {
  background: #131217;
  border: 1px solid var(--edge);
  border-radius: 4px;
  padding: 0.1rem 0.4rem;
  color: var(--gold-bright);
  font-size: 0.88em;
}
.steps code.path {
  display: block;
  margin-top: 0.3rem;
  word-break: break-all;
}
.steps a {
  color: var(--gold-bright);
}
.guide-download {
  margin: 1.6rem 0 1rem;
  text-align: center;
}
.guide-download .small {
  margin-top: 0.5rem;
  font-size: 0.8rem;
}
.helper-steps {
  margin-top: 0.75rem;
}
.helper-steps ol {
  margin: 0 0 0.5rem;
  padding-left: 1.3rem;
}
.helper-steps li {
  margin: 0.3rem 0;
}
.helper-steps code {
  background: #131217;
  border: 1px solid var(--edge);
  border-radius: 4px;
  padding: 0.05rem 0.35rem;
  color: var(--gold-bright);
  font-size: 0.88em;
}
.small {
  font-size: 0.85rem;
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

/* ---------------- dashboard ---------------- */
.dash {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.8rem;
  margin-bottom: 0.8rem;
}
.dash-card {
  background: var(--stone);
  border: 1px solid var(--edge);
  border-top: 2px solid var(--gold);
  border-radius: 8px;
  padding: 0.9rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.35);
}
.dash-num {
  font-family: Georgia, serif;
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--gold-bright);
}
.dash-label {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--muted);
}
.classbar {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  margin-bottom: 1rem;
}
.class-chip {
  background: var(--stone);
  border: 1px solid var(--edge);
  border-radius: 999px;
  padding: 0.22rem 0.7rem;
  font-size: 0.78rem;
  font-weight: 600;
}
.hint {
  font-size: 0.75rem;
  color: var(--muted);
  text-transform: none;
  letter-spacing: 0;
  font-family: inherit;
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
.roster tr.clickable {
  cursor: pointer;
}
.chevron {
  display: inline-block;
  width: 1em;
  color: var(--gold);
  transition: transform 0.15s;
}
.avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: #131217;
  border: 2px solid currentColor;
  margin-right: 0.55rem;
  vertical-align: middle;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.45);
}
.avatar svg {
  display: block;
}
.chevron.open {
  transform: rotate(90deg);
}

/* ---------------- expandable detail ---------------- */
.detail-row td {
  padding: 0 !important;
  border-bottom: 1px solid var(--edge);
}
.detail-panel {
  background: #16151a;
  border-left: 3px solid var(--gold);
  padding: 0.9rem 1.1rem;
}
.detail-core {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 0.6rem;
  margin-bottom: 0.8rem;
}
.detail-item {
  display: flex;
  flex-direction: column;
}
.detail-k {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--muted);
}
.detail-v {
  color: var(--ink);
  font-weight: 600;
}
.gear-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
  gap: 0.45rem;
}
.gear-slot {
  background: var(--stone);
  border: 1px solid var(--edge);
  border-left: 3px solid #6b7480;
  border-radius: 5px;
  padding: 0.35rem 0.55rem;
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.gear-icon {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid var(--edge);
  image-rendering: auto;
  margin-bottom: 0.2rem;
}
.gear-slot.q4 { border-left-color: #ff8000; }  /* epic */
.gear-slot.q3 { border-left-color: #a335ee; }  /* rare */
.gear-slot.q2 { border-left-color: #1eff00; }  /* uncommon */
.gear-slot-name {
  font-size: 0.65rem;
  text-transform: uppercase;
  color: var(--muted);
}
.gear-item {
  font-size: 0.8rem;
  color: var(--ink);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.gear-ilvl {
  font-size: 0.7rem;
  color: var(--gold-bright);
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
  .avatar { width: 28px; height: 28px; }
  .avatar svg { width: 18px; height: 18px; }

  /* detail row spans full width, no fake label */
  .detail-row td { display: block; }
  .detail-row td::before { display: none; }
  .dash { grid-template-columns: repeat(2, 1fr); }
}
</style>