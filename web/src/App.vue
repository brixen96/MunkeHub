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
        <button :class="{ active: tab === 'guides' }" @click="tab = 'guides'">
          Guides
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
                      {{ fullName(c) }}
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
                      <!-- armory-style hero -->
                      <div class="armory-hero">
                        <div class="ah-left">
                          <img v-if="avatarSrc(c)" class="ah-avatar" :src="avatarSrc(c)" alt=""
                               loading="lazy" @error="$event.target.style.display = 'none'">
                          <div>
                            <p class="ah-meta">
                              Level {{ c.level }} · {{ c.race }} · {{ c.class }}
                            </p>
                            <h3 class="ah-name">{{ fullName(c) }}</h3>
                            <p class="ah-sub">
                              {{ c.realm }}
                              <span v-if="c.spec" class="ah-spec">{{ c.spec }}</span>
                            </p>
                          </div>
                        </div>
                        <div class="ah-stats">
                          <div class="ah-stat">
                            <b>{{ gearscore(c) }}</b>
                            <span>Gearscore</span>
                          </div>
                          <div class="ah-stat">
                            <b>{{ Math.round(c.ilvl) }}</b>
                            <span>Item level</span>
                          </div>
                        </div>
                      </div>
                      <div class="detail-core">
                        <div class="detail-item">
                          <span class="detail-k">Zone</span>
                          <span class="detail-v">{{ c.zone || '—' }}</span>
                        </div>
                        <div class="detail-item">
                          <span class="detail-k">Played</span>
                          <span class="detail-v">{{ formatPlayed(c) }}</span>
                        </div>
                      </div>
                      <!-- Character / Talents sub-tabs (armory style) -->
                      <div class="sheet-tabs">
                        <button :class="{ active: sheetTab[c.id] !== 'talents' }"
                                @click="setSheetTab(c, 'char')">Character</button>
                        <button :class="{ active: sheetTab[c.id] === 'talents' }"
                                @click="setSheetTab(c, 'talents')">Talents</button>
                      </div>
                      <template v-if="sheetTab[c.id] !== 'talents'">
                      <div v-if="parseGear(c).length" class="paperdoll">
                        <div class="pd-col">
                          <div v-for="g in pdSlots(c, PD_LEFT)" :key="g.slot"
                               class="gear-slot" :class="['q' + g.q, { empty: !g.id }]"
                               :title="g.id ? null : g.slot">
                            <a v-if="g.id" class="gear-link" :href="wowheadUrl(g)"
                               data-wowhead="" target="_blank" rel="noopener" @click.stop>
                              <img class="gear-icon" :src="'/api/v1/icon/' + (g.fid || g.id)" alt=""
                                   loading="lazy" @error="iconFallback($event, g)">
                              <span v-if="g.ilvl" class="gear-badge">{{ g.ilvl }}</span>
                            </a>
                          </div>
                        </div>
                        <div class="pd-center">
                          <div v-if="viewerActive(c)" :id="'mv-' + c.id" class="pd-mv"
                               :class="{ loading: viewerState[c.id] === 'loading' }"></div>
                          <img v-else-if="renderSrc(c)" class="pd-render" :src="renderSrc(c)" alt=""
                               loading="lazy" @error="markRenderFail(c)">
                          <template v-else>
                            <img v-if="avatarSrc(c)" class="pd-avatar" :src="avatarSrc(c)" alt=""
                                 loading="lazy" @error="$event.target.style.display = 'none'">
                            <span v-else class="pd-avatar pd-glyph" v-html="classGlyph(c.class)"></span>
                          </template>
                        </div>
                        <div class="pd-col">
                          <div v-for="g in pdSlots(c, PD_RIGHT)" :key="g.slot"
                               class="gear-slot" :class="['q' + g.q, { empty: !g.id }]"
                               :title="g.id ? null : g.slot">
                            <a v-if="g.id" class="gear-link" :href="wowheadUrl(g)"
                               data-wowhead="" target="_blank" rel="noopener" @click.stop>
                              <img class="gear-icon" :src="'/api/v1/icon/' + (g.fid || g.id)" alt=""
                                   loading="lazy" @error="iconFallback($event, g)">
                              <span v-if="g.ilvl" class="gear-badge">{{ g.ilvl }}</span>
                            </a>
                          </div>
                        </div>
                        <div class="pd-bottom">
                          <div v-for="g in pdSlots(c, PD_BOTTOM)" :key="g.slot"
                               class="gear-slot" :class="['q' + g.q, { empty: !g.id }]"
                               :title="g.id ? null : g.slot">
                            <a v-if="g.id" class="gear-link" :href="wowheadUrl(g)"
                               data-wowhead="" target="_blank" rel="noopener" @click.stop>
                              <img class="gear-icon" :src="'/api/v1/icon/' + (g.fid || g.id)" alt=""
                                   loading="lazy" @error="iconFallback($event, g)">
                              <span v-if="g.ilvl" class="gear-badge">{{ g.ilvl }}</span>
                            </a>
                          </div>
                        </div>
                      </div>
                      <p v-else class="muted small">
                        No gear detail yet — /reload in game, then the helper uploads it.
                      </p>
                      </template>
                      <!-- Talent trees (same look as armorywowforever.com) -->
                      <div v-else class="talents-pane">
                        <p v-if="talentState(c) === 'loading'" class="muted small">Loading talents…</p>
                        <p v-else-if="talentState(c) === 'error'" class="error small">Could not load talents.</p>
                        <p v-else-if="talentState(c) === 'empty'" class="muted small">
                          No talent data for this character yet.
                        </p>
                        <div v-else class="trees">
                          <div v-for="t in talentLayout(c)" :key="t.id" class="tree">
                            <h3>{{ t.name }} <span>{{ treePoints(c, t) }}</span></h3>
                            <div class="tal-grid" :style="gridStyle(t)">
                              <i v-for="(ln, i) in connectorLinks(t)" :key="'ln' + i"
                                 class="link" :style="ln"></i>
                              <a v-for="n in t.talents" :key="n.id" class="tal"
                                 :class="talentClass(c, t, n)" :style="talentPos(n)"
                                 :href="talentUrl(n)" data-wowhead=""
                                 target="_blank" rel="noopener" @click.stop>
                                <img :src="talentIcon(n)" alt="" loading="lazy"
                                     @error="$event.target.style.display = 'none'">
                                <span v-if="talentRank(c, t, n) > 0" class="rk">
                                  {{ talentRank(c, t, n) }}/{{ n.ranks.length }}
                                </span>
                              </a>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </section>

      <!-- ================= GUIDES ================= -->
      <section v-show="tab === 'guides'">
        <div class="panel">
          <h2 class="panel-title">Profession Guides <span class="hint">— WoW: Forever</span></h2>
          <!-- profession selector -->
          <div class="guide-chips">
            <button v-for="g in guides" :key="g.id" class="chip"
                    :class="{ active: guideSel === g.id }"
                    @click="guideSel = g.id">{{ g.name }}</button>
          </div>
          <div v-for="g in visibleGuides" :key="g.id" class="guide">
            <div class="guide-head">
              <img class="guide-icon" :src="'/api/v1/icon/' + g.icon" alt=""
                   loading="lazy" @error="$event.target.style.display = 'none'">
              <div>
                <h3 class="guide-title">{{ g.name }} 1-300</h3>
                <p class="guide-blurb">{{ g.blurb }}</p>
              </div>
            </div>

            <!-- ranks (only some guides have them) -->
            <details v-if="g.ranks && g.ranks.length" class="gsec" open>
              <summary class="guide-h4">Training & Ranks</summary>
              <table class="guide-table">
                <thead>
                  <tr><th>Rank</th><th>Skill range</th><th>Skill needed</th><th>Level</th></tr>
                </thead>
                <tbody>
                  <tr v-for="r in g.ranks" :key="r.rank">
                    <td>{{ r.rank }}</td><td>{{ r.range }}</td>
                    <td>{{ r.skill }}</td><td>{{ r.level }}</td>
                  </tr>
                </tbody>
              </table>
            </details>

            <!-- shopping list -->
            <details v-if="g.shopping && g.shopping.length" class="gsec" open>
              <summary class="guide-h4">Shopping List</summary>
              <div class="shop-grid">
                <div v-for="s in g.shopping" :key="s.range" class="shop-card">
                  <span class="shop-range">{{ s.range }}</span>
                  <ul class="shop-list">
                    <li v-for="(it, i) in s.items" :key="i">
                      <b>{{ it.qty }}×</b>
                      <a v-if="it.id" :href="wowheadItem(it.id)" data-wowhead="" target="_blank" rel="noopener">{{ it.name }}</a>
                      <span v-else>{{ it.name }}</span>
                    </li>
                  </ul>
                </div>
              </div>
            </details>

            <!-- leveling steps (crafting professions) -->
            <details v-if="g.steps && g.steps.length" class="gsec" open>
              <summary class="guide-h4">Leveling Route</summary>
              <div v-for="st in g.steps" :key="st.range" class="step">
                <div class="step-head">{{ st.range }}</div>
                <p v-if="st.note" class="step-note">{{ st.note }}</p>
                <table class="guide-table">
                  <thead>
                    <tr><th>Skill</th><th>Qty</th><th>Craft</th><th>Reagents</th></tr>
                  </thead>
                  <tbody>
                    <tr v-for="(row, ri) in st.rows" :key="ri">
                      <td class="nowrap">{{ row.skill }}</td>
                      <td class="nowrap">{{ row.qty }}</td>
                      <td>
                        <a v-if="row.craft.id" :href="wowheadLink(row.craft)" data-wowhead="" target="_blank" rel="noopener">{{ row.craft.name }}</a>
                        <span v-else>{{ row.craft.name }}</span>
                      </td>
                      <td>
                        <span v-for="(rg, i) in row.reagents" :key="i">
                          <span v-if="i"> · </span>
                          <b v-if="rg.qty">{{ rg.qty }}×</b>
                          <a v-if="rg.id" :href="wowheadItem(rg.id)" data-wowhead="" target="_blank" rel="noopener">{{ rg.name }}</a>
                          <span v-else>{{ rg.name }}</span>
                        </span>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </details>

            <!-- gathering route (mining/herbalism/skinning/fishing) -->
            <details v-if="g.gathering && g.gathering.length" class="gsec" open>
              <summary class="guide-h4">Leveling Route</summary>
              <div v-for="st in g.gathering" :key="st.range" class="step">
                <div class="step-head">{{ st.range }}</div>
                <p v-if="st.note" class="step-note">{{ st.note }}</p>
                <table v-if="st.zones && st.zones.length" class="guide-table">
                  <tbody>
                    <tr v-for="(z, zi) in st.zones" :key="zi">
                      <td v-for="(cell, ci) in z" :key="ci" :class="{ nowrap: ci === 0 }">{{ cell }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </details>

            <!-- notable items -->
            <details v-if="g.items && g.items.length" class="gsec" open>
              <summary class="guide-h4">Worth Knowing</summary>
              <ul class="guide-notes">
                <li v-for="it in g.items" :key="it.id">
                  <a :href="wowheadItem(it.id)" data-wowhead="" target="_blank" rel="noopener">{{ it.name }}</a>
                  — {{ it.note }}
                </li>
              </ul>
            </details>

            <!-- trainers: generic sections (new shape) -->
            <details v-if="g.trainers && g.trainers.sections" class="gsec" open>
              <summary class="guide-h4">Trainers</summary>
              <div v-for="(sec, si) in g.trainers.sections" :key="si" class="trainer-sec">
                <h5 v-if="sec.title" class="trainer-h5">{{ sec.title }}</h5>
                <p v-for="(n, ni) in sec.notes" :key="'n' + ni" class="step-note">{{ n }}</p>
                <ul v-if="sec.lines.length" class="trainer-lines">
                  <li v-for="(l, li) in sec.lines" :key="li">{{ l }}</li>
                </ul>
                <table v-if="sec.rows.length" class="guide-table">
                  <tbody>
                    <tr v-for="(r, ri) in sec.rows" :key="ri">
                      <td v-for="(cell, ci) in r" :key="ci" :class="{ nowrap: ci === 0 }">{{ cell }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </details>
            <!-- trainers: old shape (cooking) -->
            <details v-else class="gsec" open>
              <summary class="guide-h4">Trainers</summary>
              <div class="trainer-cols">
                <div>
                  <h5 class="trainer-h5">Alliance — Apprentice / Journeyman</h5>
                  <table class="guide-table">
                    <tbody>
                      <tr v-for="t in g.trainers.alliance.basic" :key="t.name">
                        <td>{{ t.name }}</td><td>{{ t.zone }}</td><td class="nowrap coords">{{ t.coords }}</td>
                      </tr>
                    </tbody>
                  </table>
                  <h5 class="trainer-h5">Alliance — Expert (book vendor)</h5>
                  <table class="guide-table">
                    <tbody>
                      <tr v-for="t in g.trainers.alliance.expert" :key="t.name">
                        <td>{{ t.name }}</td><td>{{ t.zone }}</td><td class="nowrap coords">{{ t.coords }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <div>
                  <h5 class="trainer-h5">Horde — Apprentice / Journeyman</h5>
                  <table class="guide-table">
                    <tbody>
                      <tr v-for="t in g.trainers.horde.basic" :key="t.name">
                        <td>{{ t.name }}</td><td>{{ t.zone }}</td><td class="nowrap coords">{{ t.coords }}</td>
                      </tr>
                    </tbody>
                  </table>
                  <h5 class="trainer-h5">Horde — Expert (book vendor)</h5>
                  <table class="guide-table">
                    <tbody>
                      <tr v-for="t in g.trainers.horde.expert" :key="t.name">
                        <td>{{ t.name }}</td><td>{{ t.zone }}</td><td class="nowrap coords">{{ t.coords }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
              <h5 class="trainer-h5">Artisan (quest — both factions)</h5>
              <table class="guide-table">
                <tbody>
                  <tr v-for="t in g.trainers.artisan" :key="t.name">
                    <td>{{ t.name }}</td><td>{{ t.zone }}</td><td class="nowrap coords">{{ t.coords }}</td>
                  </tr>
                </tbody>
              </table>
            </details>
          </div>
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
import { GUIDES } from "./guides";
import { GUIDES_EXTRA } from "./guides-extra";

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

// WoW character-sheet slot order (classic paper doll):
// left column top→bottom, bottom row right→left, right column top→bottom.
const PD_LEFT = ["Head", "Neck", "Shoulder", "Back", "Chest", "Shirt", "Tabard", "Wrist"];
const PD_RIGHT = ["Hands", "Waist", "Legs", "Feet", "Ring1", "Ring2", "Trinket1", "Trinket2"];
// displayed left→right = reverse of the requested right→left order
const PD_BOTTOM = ["Ranged", "OffHand", "MainHand"];

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

// 3D viewer state lives OUTSIDE Vue's reactivity: WebGL objects must never
// be proxied. VIEWER_SCRIPT memoizes the (large) viewer script load.
let VIEWER_SCRIPT = null;
const VIEWERS = {};

export default {
  name: "App",
  data() {
    return {
      PD_LEFT,
      PD_RIGHT,
      PD_BOTTOM,
      renderFailed: {},
      // char id -> 'loading' | 'ready' | 'static' (3D viewer state)
      viewerState: {},
      // char id -> 'char' | 'talents' (sub-tab inside the detail panel)
      sheetTab: {},
      // char id -> {layout, talents} | {error} | {loading} (talents tab data)
      talentData: {},
      tab: "roster",
      openRow: null,
      stats: null,
      guild: null,
      webVersion: WEB_VERSION,
      apiVersion: "",
      changelog: CHANGELOG,
      guides: [...GUIDES, ...GUIDES_EXTRA],
      // selected profession in the Guides tab (chip selector)
      guideSel: "cooking",
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
  computed: {
    // Guides tab: only the profession picked in the chip selector shows.
    visibleGuides() {
      return this.guides.filter((g) => g.id === this.guideSel);
    },
  },
  methods: {
    formatDate(iso) {
      if (!iso) return "";
      return new Date(iso).toLocaleString();
    },
    classSlug(cls) {
      return "c-" + (cls || "").toLowerCase().replace(/[^a-z]/g, "");
    },
    fullName(c) {
      return c.last_name ? c.name + " " + c.last_name : c.name;
    },
    // 3D character render from Blizzard's render service. Bucket = id % 256.
    renderSrc(c) {
      if (!c.bid || this.renderFailed[c.id]) return "";
      const bucket = c.bid % 256;
      return (
        "https://render.worldofwarcraft.com/forever-eu/character/" +
        bucket + "/" + c.bid + "-main-raw.png"
      );
    },
    markRenderFail(c) {
      this.renderFailed = { ...this.renderFailed, [c.id]: true };
    },
    // Armory gearscore: sum of round(ilvl * quality weight).
    // Weights verified against armorywowforever.com (302 for Brix).
    gearscore(c) {
      const w = { 1: 0.5, 2: 0.6, 3: 0.7, 4: 0.8 };
      return this.parseGear(c).reduce(
        (sum, g) => sum + Math.round((g.ilvl || 0) * (w[g.q] || 0)),
        0
      );
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
      if (this.openRow === id) {
        this.closeViewer(id);
        this.openRow = null;
        return;
      }
      if (this.openRow !== null) this.closeViewer(this.openRow);
      this.openRow = id;
      const c = (this.characters || []).find((x) => x.id === id);
      if (c) this.openViewer(c);
    },
    // --- 3D model viewer (wowhead ZamModelViewer — drag to rotate) ---------
    viewerActive(c) {
      const s = this.viewerState[c.id];
      return s === "loading" || s === "ready";
    },
    wowheadUrl(g) {
      return "https://www.wowhead.com/forever/item=" + g.id;
    },
    wowheadItem(id) {
      return "https://www.wowhead.com/forever/item=" + id;
    },
    // craft links: spells (recipes) vs items depending on link kind
    wowheadLink(c) {
      return "https://www.wowhead.com/forever/" + (c.k === "item" ? "item=" : "spell=") + c.id;
    },
    // --- talents tab -------------------------------------------------------
    setSheetTab(c, t) {
      const cur = this.sheetTab[c.id] || "char";
      if (cur === t) return;
      this.sheetTab = { ...this.sheetTab, [c.id]: t };
      if (t === "talents") {
        this.closeViewer(c.id); // hidden canvas must not keep rendering
        this.loadTalents(c);
      } else {
        this.openViewer(c); // re-init: v-if destroyed the container
      }
    },
    talentState(c) {
      const d = this.talentData[c.id];
      if (!d || d.loading) return "loading";
      if (d.error) return "error";
      if (d.empty || !d.layout || !d.talents || !Array.isArray(d.talents.trees)) {
        return "empty";
      }
      return "ok";
    },
    async loadTalents(c) {
      if (this.talentData[c.id]) return;
      this.talentData = { ...this.talentData, [c.id]: { loading: true } };
      try {
        const slug = c.last_name ? c.name + "-" + c.last_name : c.name;
        const sheet = await api.charModel(slug); // {model, talents, classId}
        // talents = array of builds [{active, trees}] — use the active one
        const builds = Array.isArray(sheet.talents) ? sheet.talents : sheet.talents ? [sheet.talents] : [];
        const build = builds.find((b) => b && b.active) || builds[0];
        if (!build || !Array.isArray(build.trees) || !sheet.classId) {
          this.talentData = { ...this.talentData, [c.id]: { empty: true } };
          return;
        }
        const layout = await api.talentLayout(sheet.classId);
        this.talentData = { ...this.talentData, [c.id]: { layout, talents: build } };
      } catch (err) {
        console.warn("MunkeHub talents:", err);
        this.talentData = { ...this.talentData, [c.id]: { error: true } };
      }
    },
    talentLayout(c) {
      const d = this.talentData[c.id];
      return d && d.layout ? d.layout : null;
    },
    charTalentTree(c, layoutTree) {
      const d = this.talentData[c.id];
      if (!d || !d.talents) return null;
      return (d.talents.trees || []).find((x) => x.id === layoutTree.id) || null;
    },
    treePoints(c, t) {
      const tree = this.charTalentTree(c, t);
      return tree ? tree.points : 0;
    },
    talentRank(c, t, n) {
      const tree = this.charTalentTree(c, t);
      if (!tree) return 0;
      const tal = (tree.talents || []).find((x) => x.id === n.id);
      return tal ? tal.rank : 0;
    },
    talentClass(c, t, n) {
      const r = this.talentRank(c, t, n);
      const max = n.ranks && n.ranks.length ? n.ranks.length : 1;
      if (r <= 0) return "off";
      return r >= max ? "max" : "part";
    },
    talentPos(n) {
      return { left: n.col * 60 + "px", top: n.row * 60 + "px" };
    },
    talentUrl(n) {
      const spell = n.ranks && n.ranks.length ? n.ranks[0] : 0;
      return "https://www.wowhead.com/forever/spell=" + spell;
    },
    // Layout gives a Blizzard icon URL; route through our icon cache by
    // texture name (zamimg, cached forever) instead of hotlinking.
    talentIcon(n) {
      if (!n.icon) return "";
      const base = n.icon.split("/").pop().split("?")[0];
      const name = base.replace(/\.(jpg|png|webp|gif)$/i, "");
      return /^[a-z0-9_]+$/i.test(name) ? "/api/v1/icon/" + name : n.icon;
    },
    // Grid box sized to the tree's actual rows/cols (armory uses a fixed
    // 4×7 grid; computing keeps odd trees centered).
    gridStyle(tree) {
      let mc = 0, mr = 0;
      tree.talents.forEach((n) => {
        if (n.col > mc) mc = n.col;
        if (n.row > mr) mr = n.row;
      });
      return { width: (mc + 1) * 60 + "px", height: (mr + 1) * 60 + "px" };
    },
    // Connector lines from prereq → talent (same column/row only). The line
    // spans the full gap, so consecutive edges merge into one long line —
    // exactly how the armory renders them.
    connectorLinks(tree) {
      const byId = {};
      tree.talents.forEach((n) => (byId[n.id] = n));
      const out = [];
      tree.talents.forEach((n) => {
        const p = n.prereq && byId[n.prereq.id];
        if (!p) return;
        if (p.col === n.col && n.row > p.row) {
          out.push(
            `left:${p.col * 60 + 29}px;top:${p.row * 60 + 30}px;` +
            `width:3px;height:${(n.row - p.row) * 60}px`
          );
        } else if (p.row === n.row && n.col > p.col) {
          out.push(
            `left:${p.col * 60 + 30}px;top:${p.row * 60 + 29}px;` +
            `width:${(n.col - p.col) * 60}px;height:3px`
          );
        }
      });
      return out;
    },
    // Icon chain: fileID (foreverdb) first — if that 404s (some items are
    // missing there), retry with the itemID (wowhead resolves the name).
    // Second failure = hide the img, keep the quality border.
    iconFallback(ev, g) {
      const img = ev.target;
      if (g.id && !img.dataset.triedId) {
        img.dataset.triedId = "1";
        img.src = "/api/v1/icon/" + g.id;
        return;
      }
      img.style.display = "none";
    },
    webglAvailable() {
      try {
        const cv = document.createElement("canvas");
        return !!(
          window.WebGLRenderingContext &&
          (cv.getContext("webgl") || cv.getContext("experimental-webgl"))
        );
      } catch (e) {
        return false;
      }
    },
    loadViewerScript() {
      if (!VIEWER_SCRIPT) {
        VIEWER_SCRIPT = new Promise((resolve, reject) => {
          const s = document.createElement("script");
          // wowhead's viewer, served through our asset proxy (their CDN
          // sends no CORS headers — but a <script> tag wouldn't care; the
          // proxy also covers every asset the viewer fetches via XHR).
          s.src = "/api/v1/mv/deployment/viewer/35ee72d/viewer.min.js";
          s.onload = () => resolve();
          s.onerror = () => {
            VIEWER_SCRIPT = null;
            reject(new Error("viewer script failed"));
          };
          document.head.appendChild(s);
        });
      }
      return VIEWER_SCRIPT;
    },
    async openViewer(c) {
      if (this.viewerState[c.id]) return;
      this.viewerState = { ...this.viewerState, [c.id]: "loading" };
      if (!this.webglAvailable()) {
        this.viewerState = { ...this.viewerState, [c.id]: "static" };
        return;
      }
      try {
        // The viewer expects Wowhead globals (tooltips.js defines WH but
        // not all of these) — same shim the armory uses.
        window.WH = window.WH || {};
        window.WH.debug = window.WH.debug || function () {};
        window.WH.defaultAnimation = window.WH.defaultAnimation || "Stand";
        window.WH.WebP = window.WH.WebP || {
          getImageExtension: function () { return ".webp"; },
        };
        await this.loadViewerScript();
        const slug = c.last_name ? c.name + "-" + c.last_name : c.name;
        const sheet = await api.charModel(slug); // {model, talents, classId}
        const model = sheet.model || {};
        if (this.openRow !== c.id || VIEWERS[c.id]) return;
        await this.$nextTick();
        const el = document.getElementById("mv-" + c.id);
        if (!el || this.openRow !== c.id || VIEWERS[c.id]) return;
        el.innerHTML = "";
        const viewer = new window.ZamModelViewer({
          type: 2, // Wowhead
          contentPath: "/api/v1/mv/",
          container: window.jQuery(el),
          // guard against a zero-size container (viewer throws on aspect 0)
          aspect: (el.clientWidth || 360) / (el.clientHeight || 400),
          hd: false,
          items: model.items || [],
          models: { id: model.modelId, type: 16 },
          charCustomization: { options: model.customizations || [] },
          background: "stage.png", // resolved as contentPath + background
          bgScale: [1, 1],
        });
        VIEWERS[c.id] = viewer;
        this.viewerState = { ...this.viewerState, [c.id]: "ready" };
      } catch (err) {
        console.warn("MunkeHub 3D viewer:", err);
        if (this.openRow === c.id) {
          // fall back to the static Blizzard render
          this.viewerState = { ...this.viewerState, [c.id]: "static" };
        }
      }
    },
    closeViewer(id) {
      const v = VIEWERS[id];
      if (v) {
        try {
          const gl = v.renderer && v.renderer.context;
          v.destroy();
          if (gl && gl.getExtension) {
            gl.getExtension("WEBGL_lose_context").loseContext();
          }
        } catch (e) {
          /* already gone */
        }
        delete VIEWERS[id];
      }
      if (this.viewerState[id]) {
        const s = { ...this.viewerState };
        delete s[id];
        this.viewerState = s;
      }
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
    // Paper-doll ordering: slots in `order`, unknown slots appended to the
    // right column so nothing ever disappears.
    pdSlots(c, order) {
      const gear = this.parseGear(c);
      const bySlot = {};
      gear.forEach((g) => (bySlot[g.slot] = g));
      const picked = order.map((s) => bySlot[s]).filter(Boolean);
      if (order === PD_RIGHT) {
        const known = new Set([...PD_LEFT, ...PD_RIGHT, ...PD_BOTTOM]);
        gear.forEach((g) => {
          if (!known.has(g.slot)) picked.push(g);
        });
      }
      return picked;
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
.avatar img {
  display: block;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  padding: 3px; /* keep the icon off the class-coloured ring */
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
/* ---------------- armory hero ---------------- */
.armory-hero {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 0.8rem;
  background: linear-gradient(180deg, #1b2430, #141a22);
  border: 1px solid var(--edge);
  border-radius: 8px;
  padding: 0.8rem 1rem;
  margin-bottom: 0.8rem;
}
.ah-left {
  display: flex;
  align-items: center;
  gap: 0.8rem;
  min-width: 0;
}
.ah-avatar {
  width: 56px;
  height: 56px;
  border-radius: 8px;
  border: 2px solid var(--gold);
  background: #131217;
  object-fit: cover;
  padding: 3px;
  flex: none;
}
.ah-meta {
  margin: 0;
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--muted);
}
.ah-name {
  margin: 0.1rem 0;
  font-family: Georgia, serif;
  font-size: 1.5rem;
  color: var(--gold-bright);
  line-height: 1.1;
}
.ah-sub {
  margin: 0;
  font-size: 0.8rem;
  color: var(--muted);
}
.ah-spec {
  display: inline-block;
  margin-left: 0.5rem;
  padding: 0.05rem 0.5rem;
  background: #23303f;
  border: 1px solid var(--edge);
  border-radius: 999px;
  color: var(--ink);
}
.ah-stats {
  display: flex;
  gap: 1.4rem;
}
.ah-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.ah-stat b {
  font-size: 2rem;
  line-height: 1;
  color: var(--gold-bright);
}
.ah-stat span {
  font-size: 0.65rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--muted);
}

/* ---------------- paper doll (armory style) ---------------- */
.paperdoll {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  gap: 0.7rem;
  align-items: start;
}
/* weapons row spans all columns, under the character render */
.pd-bottom {
  grid-column: 1 / -1;
  display: flex;
  justify-content: center;
  gap: 0.5rem;
}
.pd-col {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  min-width: 0;
}
.pd-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  min-width: 220px;
  min-height: 340px;
}
.pd-render {
  max-height: 420px;
  max-width: 100%;
  object-fit: contain;
  filter: drop-shadow(0 6px 14px rgba(0, 0, 0, 0.55));
}
.pd-avatar {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  border: 2px solid var(--gold);
  background: #131217;
  object-fit: cover;
  padding: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
}
.pd-avatar.pd-glyph {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--gold);
}
.pd-avatar.pd-glyph svg {
  width: 48px;
  height: 48px;
}

/* 3D model viewer (wowhead ZamModelViewer) — drag to rotate, scroll to zoom */
.pd-mv {
  width: 360px;
  height: 400px;
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  background: radial-gradient(ellipse at 50% 35%, #26252d 0%, #141318 78%);
}
.pd-mv.loading::after {
  content: "Summoning 3D model…";
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--muted);
  font-size: 0.85rem;
}
/* gear icon link: fills the slot, wowhead tooltip on hover */
.gear-link {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

/* --- Character / Talents sub-tabs --- */
.sheet-tabs {
  display: flex;
  gap: 0.5rem;
  margin: 0.9rem 0 0.8rem;
  border-bottom: 1px solid var(--edge);
  padding-bottom: 0.5rem;
}
.sheet-tabs button {
  background: #14131a;
  border: 1px solid var(--edge);
  color: var(--muted);
  padding: 0.35rem 1.1rem;
  border-radius: 6px 6px 0 0;
  cursor: pointer;
  font: inherit;
  font-size: 0.85rem;
}
.sheet-tabs button:hover {
  color: var(--ink);
}
.sheet-tabs button.active {
  color: var(--gold-bright);
  border-color: var(--gold);
  background: #1d1b22;
}

/* --- Talent trees (armory look: gold/green/brown gradients) --- */
.trees {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(265px, 1fr));
  gap: 1.4rem;
}
.tree h3 {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  margin: 0 0 12px;
  font-size: 1.05rem;
  font-weight: 700;
  color: #fff;
}
.tree h3 span {
  font-size: 1.35rem;
  font-weight: 800;
  color: #ffe2a6;
}
.tal-grid {
  position: relative;
  margin: 0 auto;
}
/* connector lines between prereq and talent */
.tal-grid .link {
  position: absolute;
  background: rgba(168, 188, 203, 0.25);
  border-radius: 1px;
}
.tal {
  position: absolute;
  display: block;
  box-sizing: border-box;
  width: 43px;
  height: 43px;
  padding: 2px;
  margin: 8.4px; /* centers the 43px box inside the 60px cell */
  border-radius: 4px;
  color: var(--ink);
}
.tal img {
  display: block;
  width: 39px;
  height: 39px;
  border-radius: 2px;
}
.tal.off {
  background: linear-gradient(145deg, rgb(138, 114, 80), rgb(58, 43, 22));
  box-shadow: 0 0 0 1px #0d0904, 0 3px 6px rgba(0, 0, 0, 0.45);
}
.tal.part {
  background: linear-gradient(145deg, rgb(125, 255, 106), rgb(28, 138, 16));
  box-shadow: 0 0 0 1px #000, 0 0 8px rgba(30, 255, 0, 0.35);
}
.tal.max {
  background: linear-gradient(145deg, rgb(247, 226, 173), rgb(176, 138, 74) 45%, rgb(111, 79, 39));
  box-shadow: 0 0 0 1px #000, 0 0 10px rgba(247, 226, 173, 0.45);
}
/* --- Guides section --- */
/* profession selector chips */
.guide-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
  margin-bottom: 1rem;
}
.chip {
  background: #14131a;
  border: 1px solid var(--edge);
  color: var(--muted);
  padding: 0.3rem 0.9rem;
  border-radius: 999px;
  cursor: pointer;
  font: inherit;
  font-size: 0.82rem;
  transition: all 0.12s ease;
}
.chip:hover {
  color: var(--ink);
  border-color: var(--gold);
}
.chip.active {
  color: #14131a;
  background: var(--gold);
  border-color: var(--gold-bright);
  font-weight: 700;
}
/* collapsible guide sections */
.gsec {
  border-top: 1px solid rgba(58, 56, 64, 0.6);
  padding-top: 0.15rem;
}
.gsec:first-of-type {
  border-top: 0;
}
summary.guide-h4 {
  cursor: pointer;
  list-style: none;
  user-select: none;
}
summary.guide-h4::-webkit-details-marker {
  display: none;
}
summary.guide-h4::after {
  content: " ▸";
  color: var(--gold);
  font-size: 0.8em;
}
details[open] > summary.guide-h4::after {
  content: " ▾";
}
.guide {
  border-top: 1px solid var(--edge);
  padding-top: 1.1rem;
  margin-top: 1.1rem;
}
.guide:first-of-type {
  border-top: 0;
  margin-top: 0.3rem;
  padding-top: 0;
}
.guide-head {
  display: flex;
  gap: 0.9rem;
  align-items: center;
  margin-bottom: 0.6rem;
}
.guide-icon {
  width: 44px;
  height: 44px;
  border: 2px solid var(--gold);
  border-radius: 6px;
  background: #0e0d11;
  padding: 2px;
}
.guide-title {
  margin: 0;
  font-size: 1.15rem;
  color: var(--gold-bright);
}
.guide-blurb {
  margin: 0.2rem 0 0;
  color: var(--muted);
  font-size: 0.85rem;
}
.guide-h4 {
  margin: 1.1rem 0 0.5rem;
  font-size: 0.95rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--gold);
}
.guide-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}
.guide-table th {
  text-align: left;
  color: var(--muted);
  font-weight: 600;
  font-size: 0.72rem;
  text-transform: uppercase;
  padding: 0.35rem 0.6rem;
  border-bottom: 1px solid var(--edge);
}
.guide-table td {
  padding: 0.4rem 0.6rem;
  border-bottom: 1px solid rgba(58, 56, 64, 0.5);
}
.guide-table a {
  color: var(--ink);
}
.guide-table a:hover {
  color: var(--gold-bright);
}
.nowrap {
  white-space: nowrap;
}
.coords {
  color: var(--muted);
  font-family: Consolas, monospace;
  font-size: 0.78rem;
}
/* shopping list cards */
.shop-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
  gap: 0.7rem;
}
.shop-card {
  background: #14131a;
  border: 1px solid var(--edge);
  border-radius: 8px;
  padding: 0.6rem 0.8rem;
}
.shop-range {
  display: inline-block;
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--gold-bright);
  background: rgba(200, 164, 77, 0.12);
  border: 1px solid rgba(200, 164, 77, 0.4);
  border-radius: 4px;
  padding: 0.1rem 0.5rem;
  margin-bottom: 0.4rem;
}
.shop-list {
  list-style: none;
  margin: 0;
  padding: 0;
  font-size: 0.82rem;
}
.shop-list li {
  padding: 0.15rem 0;
}
.shop-list b {
  color: var(--gold-bright);
  font-weight: 700;
}
.shop-list a {
  color: var(--ink);
}
.shop-list a:hover {
  color: var(--gold-bright);
}
/* leveling steps */
.step {
  margin-bottom: 0.9rem;
}
.step-head {
  display: inline-block;
  font-size: 0.8rem;
  font-weight: 700;
  color: #ffe2a6;
  background: rgba(247, 226, 173, 0.1);
  border: 1px solid rgba(247, 226, 173, 0.35);
  border-radius: 4px;
  padding: 0.15rem 0.7rem;
  margin-bottom: 0.4rem;
}
.step-note {
  margin: 0 0 0.5rem;
  font-size: 0.83rem;
  color: var(--muted);
}
.guide-notes {
  margin: 0;
  padding-left: 1.2rem;
  font-size: 0.85rem;
}
.guide-notes li {
  margin-bottom: 0.3rem;
}
.guide-notes a {
  color: var(--ink);
}
.guide-notes a:hover {
  color: var(--gold-bright);
}
/* trainer columns */
.trainer-cols {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
}
.trainer-h5 {
  margin: 0.7rem 0 0.35rem;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  color: var(--muted);
}
.trainer-cols .trainer-h5:first-child {
  margin-top: 0;
}
/* generic trainer sections (new guide shape) */
.trainer-sec {
  margin-bottom: 0.9rem;
}
.trainer-sec .trainer-h5 {
  margin-top: 0;
}
.trainer-lines {
  margin: 0.2rem 0 0.4rem;
  padding-left: 1.2rem;
  font-size: 0.85rem;
}
.trainer-lines li {
  margin-bottom: 0.2rem;
}

/* rank badge: bottom-right, overlapping the icon */
.tal .rk {
  position: absolute;
  left: 22px;
  top: 34px;
  right: -7px;
  bottom: -7px;
  font-size: 10px;
  font-weight: 700;
  line-height: 10px;
  padding: 2px 4px;
  background: #06141d;
  color: #ffe2a6;
  border: 1px solid #6f4f27;
  border-radius: 3px;
  white-space: nowrap;
}

/* square icon slots with quality border + ilvl badge (armory look) */
.gear-slot {
  position: relative;
  width: 46px;
  height: 46px;
  background: #0e0d11;
  border: 2px solid #6b7480;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex: none;
}
.gear-slot.empty {
  opacity: 0.35;
}
.gear-slot.q1 { border-color: #9d9d9d; }
.gear-slot.q2 { border-color: #1eff00; }
.gear-slot.q3 { border-color: #0070dd; }
.gear-slot.q4 { border-color: #a335ee; }
.gear-slot.q5 { border-color: #ff8000; }
.gear-icon {
  width: 38px;
  height: 38px;
  border-radius: 2px;
  image-rendering: auto;
}
.gear-badge {
  position: absolute;
  right: 1px;
  bottom: 0;
  font-size: 0.62rem;
  font-weight: 700;
  color: #fff;
  background: rgba(0, 0, 0, 0.75);
  padding: 0 3px;
  border-radius: 3px 0 0 0;
  line-height: 1.3;
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
  .avatar img { padding: 2px; }

  /* paper doll stacks: render first, then slot columns */
  .paperdoll { grid-template-columns: 1fr; }
  .pd-center { order: -1; min-width: 0; min-height: 0; }
  .pd-render { max-height: 260px; }
  .pd-mv { width: 100%; height: 260px; }
  .pd-avatar { width: 72px; height: 72px; }
  .pd-bottom { justify-content: flex-start; flex-wrap: wrap; }
  .ah-name { font-size: 1.2rem; }
  .ah-stat b { font-size: 1.5rem; }

  /* detail row spans full width, no fake label */
  .detail-row td { display: block; }
  .detail-row td::before { display: none; }
  .dash { grid-template-columns: repeat(2, 1fr); }
}
</style>