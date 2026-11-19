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
        <button :class="{ active: tab === 'board' }" @click="tab = 'board'">
          Guild Board
        </button>
        <button :class="{ active: tab === 'ah' }" @click="tab = 'ah'">
          Auction House<span v-if="watchDealsCount" class="nav-badge">{{ watchDealsCount }}</span>
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
            <span class="dash-num">{{ fmtMoney(stats.total_gold) }}</span>
            <span class="dash-label">Guild gold</span>
          </div>
          <div class="dash-card">
            <span class="dash-num">{{ mountPct }}%</span>
            <span class="dash-label">Epic mount fund</span>
            <div class="gp-bar mount-bar"><div class="gp-fill" :style="{ width: mountPct + '%' }"></div></div>
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
                      <span class="online-dot" :class="{ on: isOnline(c.name) }"
                            :title="isOnline(c.name) ? 'In game now' : 'Offline'">●</span>{{ fullName(c) }}
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
        <div class="panel guides-panel">
          <h2 class="panel-title">Profession Guides <span class="hint">— WoW: Forever</span></h2>
          <div class="guides-layout">
            <!-- profession rail (sticky on wide screens) -->
            <nav class="guides-rail">
              <button v-for="g in guides" :key="g.id" class="rail-item"
                      :class="{ active: guideSel === g.id }"
                      @click="selectGuide(g.id)">
                <img class="rail-icon" :src="'/api/v1/icon/' + g.icon" alt=""
                     loading="lazy" @error="$event.target.style.display = 'none'">
                <span class="rail-name">{{ g.name }}</span>
                <span v-if="guidePct(g)" class="rail-pct">{{ guidePct(g) }}%</span>
              </button>
            </nav>

            <!-- guide content -->
            <div class="guides-content">
              <div v-for="g in visibleGuides" :key="g.id" class="guide">
                <div class="guide-head">
                  <img class="guide-icon" :src="'/api/v1/icon/' + g.icon" alt=""
                       loading="lazy" @error="$event.target.style.display = 'none'">
                  <div class="guide-headtext">
                    <h3 class="guide-title">{{ g.name }} 1-300</h3>
                    <p class="guide-blurb">{{ g.blurb }}</p>
                  </div>
                  <!-- progress tracker -->
                  <div v-if="guideTiers(g).length" class="guide-progress">
                    <div class="gp-bar"><div class="gp-fill" :style="{ width: guidePct(g) + '%' }"></div></div>
                    <span class="gp-label">{{ guidePct(g) }}% done</span>
                  </div>
                </div>

                <!-- guild crafters: who can help with this profession -->
                <div v-if="guideCrafters(g).length" class="guide-crafters">
                  <span class="gc-label">In the guild:</span>
                  <span v-for="(c, ci) in guideCrafters(g)" :key="c.id" class="gc-name">
                    <span v-if="ci">, </span>
                    <b :style="{ color: classColor(c) }">{{ fullName(c) }}</b>
                  </span>
                </div>

                <!-- your skill level (addon-synced) + remaining cost -->
                <div v-if="guideSkill" class="skill-banner">
                  <span class="sb-label">{{ guideSkill.name }} is at</span>
                  <b class="sb-level">{{ guideSkill.level }}/300</b>
                  <div class="gp-bar sb-bar">
                    <div class="gp-fill" :style="{ width: Math.min(100, Math.round(guideSkill.level / 3)) + '%' }"></div>
                  </div>
                  <span v-if="buyCart && guideSkill.level < 300" class="sb-cost">
                    to 300: ~{{ fmtMoney(buyCart.remaining) }}
                  </span>
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

                <!-- shopping list: live prices + tier totals + bag marks -->
                <details v-if="g.shopping && g.shopping.length" class="gsec" open>
                  <summary class="guide-h4">Shopping List</summary>
                  <div class="shop-grid">
                    <div v-for="s in g.shopping" :key="s.range" class="shop-card">
                      <div class="shop-cardhead">
                        <label class="tier-check" :title="'Mark ' + s.range + ' as done'">
                          <input type="checkbox" :checked="isTierDone(g, s.range)"
                                 @change="toggleTier(g, s.range)">
                          <span>{{ s.range }}</span>
                        </label>
                        <span v-if="tierTotal(s).sum" class="shop-total">
                          ~{{ fmtMoney(tierTotal(s).sum) }}
                          <span v-if="tierTotal(s).missing" class="shop-miss">+{{ tierTotal(s).missing }} unpriced</span>
                        </span>
                      </div>
                      <ul class="shop-list">
                        <li v-for="(it, i) in s.items" :key="i">
                          <b>{{ it.qty }}×</b>
                          <a v-if="it.id" :href="wowheadItem(it.id)" data-wowhead="" target="_blank" rel="noopener">{{ it.name }}</a>
                          <span v-else>{{ it.name }}</span>
                          <span v-if="it.id && priceOf(it.id)" class="shop-price">
                            {{ fmtMoney(priceOf(it.id).min) }}
                            <span v-if="dealPct(it.id)" class="ah-dealbadge">−{{ dealPct(it.id) }}%</span>
                          </span>
                          <span v-if="haveFull(it)" class="shop-have" title="Already in our bags">✓ have</span>
                          <button v-if="it.id" class="link-btn" :title="'Price details for ' + it.name"
                                  @click.stop="openItem(it.id, it.name)">📈</button>
                        </li>
                      </ul>
                    </div>
                  </div>
                </details>

                <!-- buy list cart: aggregated, minus bags, copyable -->
                <details v-if="buyCart && buyCart.rows.length" class="gsec" open>
                  <summary class="guide-h4">Buy List <span class="hint">— what to actually buy</span></summary>
                  <div class="cart-head">
                    <span>Total: <b class="ah-gold">{{ fmtMoney(buyCart.total) }}</b>
                      <span v-if="buyCart.missing" class="shop-miss">+{{ buyCart.missing }} unpriced</span>
                    </span>
                    <button class="chip" @click="copyCart">📋 Copy list</button>
                  </div>
                  <table class="guide-table">
                    <thead>
                      <tr><th>Item</th><th class="num">Need</th><th class="num">In bags</th><th class="num">Price</th><th class="num">Cost</th></tr>
                    </thead>
                    <tbody>
                      <tr v-for="r in buyCart.rows" :key="r.id">
                        <td>
                          <a :href="wowheadItem(r.id)" data-wowhead="" target="_blank" rel="noopener">{{ r.name }}</a>
                          <button class="link-btn" @click.stop="openItem(r.id, r.name)">📈</button>
                        </td>
                        <td class="num nowrap">{{ r.need > 0 ? r.need : '—' }}</td>
                        <td class="num nowrap" :class="{ 'shop-have-cell': r.have > 0 }">{{ r.have || '—' }}</td>
                        <td class="num nowrap">{{ r.priced ? fmtMoney(r.price || 0) : '—' }}</td>
                        <td class="num nowrap ah-gold">{{ r.need > 0 && r.priced ? fmtMoney(r.cost) : '—' }}</td>
                      </tr>
                    </tbody>
                  </table>
                </details>

                <!-- leveling steps (crafting professions) -->
                <details v-if="g.steps && g.steps.length" class="gsec" open>
                  <summary class="guide-h4">Leveling Route</summary>
                  <div v-for="st in g.steps" :key="st.range" class="step">
                    <div class="step-head">
                      <label class="tier-check" :title="'Mark ' + st.range + ' as done'">
                        <input type="checkbox" :checked="isTierDone(g, st.range)"
                               @change="toggleTier(g, st.range)">
                        <span>{{ st.range }}</span>
                      </label>
                    </div>
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
                              <span v-if="rg.id && priceOf(rg.id)" class="shop-price">
                                {{ fmtMoney(priceOf(rg.id).min) }}
                                <span v-if="dealPct(rg.id)" class="ah-dealbadge">−{{ dealPct(rg.id) }}%</span>
                              </span>
                              <span v-if="haveFull(rg)" class="shop-have" title="Already in our bags">✓</span>
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
                    <div class="step-head">
                      <label class="tier-check" :title="'Mark ' + st.range + ' as done'">
                        <input type="checkbox" :checked="isTierDone(g, st.range)"
                               @change="toggleTier(g, st.range)">
                        <span>{{ st.range }}</span>
                      </label>
                    </div>
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
                      <span v-if="priceOf(it.id)" class="shop-price">{{ fmtMoney(priceOf(it.id).min) }}</span>
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
          </div>
        </div>
      </section>

      <!-- ================= AUCTION HOUSE ================= -->
      <section v-show="tab === 'ah'">
        <div class="panel">
          <h2 class="panel-title">
            Auction House <span class="hint">— Horde · RPPVP · EU</span>
          </h2>

          <p v-if="ahError" class="error">Could not load auction data.</p>
          <p v-else-if="!ahSummary" class="muted">Loading auction data…</p>

          <template v-else>
            <!-- market overview -->
            <div class="dash ah-dash">
              <div class="dash-card">
                <span class="dash-num">{{ fmtNum(ahSummary.listing_count) }}</span>
                <span class="dash-label">Listings</span>
              </div>
              <div class="dash-card">
                <span class="dash-num">{{ fmtNum(ahSummary.items) }}</span>
                <span class="dash-label">Items</span>
              </div>
              <div class="dash-card">
                <span class="dash-num">{{ fmtNum(ahSummary.confirmed_24h) }}</span>
                <span class="dash-label">Sold, 24 h</span>
              </div>
              <div class="dash-card">
                <span class="dash-num ah-money">{{ fmtMoney(ahSummary.confirmed_value_24h) }}</span>
                <span class="dash-label">Sales value, 24 h</span>
              </div>
            </div>

            <!-- MunkeHub AI advisor -->
            <div class="ah-advisor">
              <div class="ah-adv-head">
                <h3>🐒 MunkeHub AI</h3>
                <button class="chip" :disabled="ahAdviceLoading"
                        @click="loadAdvice(true)">
                  {{ ahAdviceLoading ? "Thinking…" : "↻ Refresh advice" }}
                </button>
              </div>
              <p v-if="ahAdviceError" class="error">{{ ahAdviceError }}</p>
              <p v-else-if="!ahAdvice" class="muted">Asking the AI…</p>
              <template v-else>
                <p class="ah-adv-summary">{{ ahAdvice.summary }}</p>
                <div class="ah-adv-cols">
                  <div class="ah-adv-col">
                    <h4>💰 Sell</h4>
                    <ul><li v-for="(r, i) in ahAdvice.sell" :key="i">
                      <b>{{ r.item }}</b> — {{ r.why }}</li></ul>
                  </div>
                  <div class="ah-adv-col">
                    <h4>⛏ Farm</h4>
                    <ul><li v-for="(r, i) in ahAdvice.farm" :key="i">
                      <b>{{ r.item }}</b> — {{ r.why }}</li></ul>
                  </div>
                  <div class="ah-adv-col">
                    <h4>🚫 Avoid</h4>
                    <ul><li v-for="(r, i) in ahAdvice.avoid" :key="i">
                      <b>{{ r.item }}</b> — {{ r.why }}</li></ul>
                  </div>
                </div>
                <p class="ah-adv-mount">🐎 {{ ahAdvice.mount }}</p>
                <p v-if="guildProfs.length" class="ah-adv-profs">
                  Guild professions: {{ guildProfs.join(" · ") }}
                </p>
              </template>
            </div>

            <!-- section selector -->
            <div class="guide-chips">
              <button v-for="s in ahSections" :key="s.id" class="chip"
                      :class="{ active: ahTab === s.id }"
                      @click="selectAhSection(s.id)">{{ s.name }}</button>
            </div>

            <!-- ---------- Deals ---------- -->
            <div v-if="ahTab === 'deals'" class="guide">
              <details class="gsec" open>
                <summary class="guide-h4">
                  Deals right now — {{ ahDeals ? ahDeals.all : 0 }} items at least 20% under the 7-day price
                </summary>
                <p class="ah-note">
                  Buy low, sell high. Prices are per item at the cheapest listing.
                  <span v-if="ahDeals && ahDeals.captured_at">Scanned {{ fmtAgo(ahDeals.captured_at) }}.</span>
                </p>
                <table class="guide-table">
                  <thead>
                    <tr>
                      <th>Item</th><th class="num">Now</th><th class="num">Usual (7d)</th>
                      <th class="num">Under</th><th class="num">7-day trend</th><th class="num">Sold 24h</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="d in ahDealRows" :key="d.item_id">
                      <td class="ah-item">
                        <button class="star" :class="{ on: watchSet.has(d.item_id) }"
                                :title="watchSet.has(d.item_id) ? 'Remove from watchlist' : 'Add to watchlist'"
                                @click="toggleWatch(d.item_id)">★</button>
                        <img :src="d.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a :href="wowheadItem(d.item_id)" data-wowhead="" target="_blank" rel="noopener">{{ d.name }}</a>
                        <button class="link-btn" title="Price details" @click.stop="openItem(d.item_id, d.name)">📈</button>
                      </td>
                      <td class="num nowrap">{{ fmtMoney(d.minimum_ask_copper) }}
                        <span v-if="d.suspicious" class="ah-warn" :title="d.note">⚠️</span></td>
                      <td class="num nowrap muted">{{ fmtMoney(d.usual_copper) }}</td>
                      <td class="num ah-under">−{{ Math.round(d.discount * 100) }}%</td>
                      <td class="num ah-sparkcell">
                        <svg class="spark" viewBox="0 0 60 18" preserveAspectRatio="none" aria-hidden="true">
                          <polyline :points="sparkPath(ahHist[d.item_id])" />
                        </svg>
                      </td>
                      <td class="num">{{ fmtNum(d.possible_24h) }}</td>
                    </tr>
                  </tbody>
                </table>
              </details>
            </div>

            <!-- ---------- Best to sell ---------- -->
            <div v-if="ahTab === 'sell'" class="guide">
              <details class="gsec" open>
                <summary class="guide-h4">Best to sell right now</summary>
                <p class="ah-note">
                  What actually moved on our auction house — sorted by gold earned.
                  <span class="ah-days">
                    <button class="chip" :class="{ active: ahSalesDays === '1' }"
                            @click="loadSales('1')">24 hours</button>
                    <button class="chip" :class="{ active: ahSalesDays === '7' }"
                            @click="loadSales('7')">7 days</button>
                  </span>
                </p>
                <table class="guide-table">
                  <thead>
                    <tr>
                      <th>Item</th><th class="num">Sold</th><th class="num">Earned</th>
                      <th class="num">Avg price</th><th class="num">Cheapest now</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="s in ahSalesItems" :key="s.item_id">
                      <td class="ah-item">
                        <img :src="s.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a :href="wowheadItem(s.item_id)" data-wowhead="" target="_blank" rel="noopener">{{ s.name }}</a>
                      </td>
                      <td class="num">{{ fmtNum(s.confirmed) }}</td>
                      <td class="num nowrap ah-gold">{{ fmtMoney(s.confirmed_value) }}
                        <span v-if="s.suspicious" class="ah-warn" title="average price is skewed by a high listing">⚠️</span></td>
                      <td class="num nowrap muted">{{ fmtMoney(s.average_ask_copper) }}</td>
                      <td class="num nowrap">{{ fmtMoney(s.minimum_ask_copper) }}</td>
                    </tr>
                  </tbody>
                </table>
              </details>
            </div>

            <!-- ---------- Listings ---------- -->
            <div v-if="ahTab === 'listings'" class="guide">
              <details class="gsec" open>
                <summary class="guide-h4">
                  Listings — {{ ahListings && ahListings.query
                    ? fmtNum(ahListings.total) + ' matches'
                    : fmtNum(ahSummary.items) + ' items for sale' }}
                </summary>
                <p class="ah-note">
                  <input v-model="ahQuery" class="ah-search" type="search"
                         placeholder="Search items…" @input="searchListings">
                  <button class="chip" :class="{ active: ahOnlyOurs }"
                          :disabled="!ourSubclasses.size"
                          :title="ourSubclasses.size ? 'Show only what our professions farm' : 'No professions synced yet (addon /reload)'"
                          @click="ahOnlyOurs = !ahOnlyOurs">Only ours</button>
                  <span v-if="ahListings && ahListings.query">
                    {{ fmtNum(ahListings.total) }} matches for “{{ ahListings.query }}”.
                  </span>
                </p>
                <table class="guide-table">
                  <thead>
                    <tr>
                      <th>Item</th><th class="num">For sale</th><th class="num">Cheapest</th>
                      <th class="num">Avg price</th><th class="num">7-day trend</th><th class="num">Sold 24h</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="i in ahListingRows" :key="i.id">
                      <td class="ah-item">
                        <button class="star" :class="{ on: watchSet.has(i.id) }"
                                :title="watchSet.has(i.id) ? 'Remove from watchlist' : 'Add to watchlist'"
                                @click="toggleWatch(i.id)">★</button>
                        <img :src="i.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a :href="wowheadItem(i.id)" data-wowhead="" target="_blank" rel="noopener">{{ i.name }}</a>
                        <button class="link-btn" title="Price details" @click.stop="openItem(i.id, i.name)">📈</button>
                        <span class="ah-sub">{{ i.class_name }}<template v-if="i.subclass"> · {{ i.subclass }}</template></span>
                      </td>
                      <td class="num">{{ fmtNum(i.units) }} <span class="muted">({{ fmtNum(i.listings) }})</span></td>
                      <td class="num nowrap" :class="{ 'ah-below': i.below_vendor }">{{ fmtMoney(i.min_copper) }}</td>
                      <td class="num nowrap muted">{{ fmtMoney(i.avg_copper) }}
                        <span v-if="i.suspicious" class="ah-warn" title="average is skewed — someone listed far too high">⚠️</span></td>
                      <td class="num ah-sparkcell">
                        <svg class="spark" viewBox="0 0 60 18" preserveAspectRatio="none" aria-hidden="true">
                          <polyline :points="sparkPath(ahHist[i.id])" />
                        </svg>
                      </td>
                      <td class="num">{{ fmtNum(i.sold_24h) }}</td>
                    </tr>
                  </tbody>
                </table>
                <p v-if="ahListings && !ahListings.query"
                   class="ah-note muted">
                  Showing the hottest {{ ahListings.items.length }} of {{ fmtNum(ahListings.total) }} tracked items —
                  search above to find anything else.
                </p>
              </details>
            </div>

            <!-- ---------- Farming ---------- -->
            <div v-if="ahTab === 'farm'" class="guide">
              <details class="gsec" open>
                <summary class="guide-h4">
                  Farm this — ranked by daily market value
                </summary>
                <p class="ah-note">
                  Raw materials matching the guild's professions
                  <template v-if="ahFarmProfs.length">({{ ahFarmProfs.join(", ") }})</template>.
                  Market/day = cheapest price × sold in 24h.
                </p>
                <p v-if="!ahFarm.length" class="muted">Nothing farmable cached yet.</p>
                <table v-else class="guide-table">
                  <thead>
                    <tr>
                      <th>Item</th><th class="num">Cheapest</th><th class="num">Sold 24h</th>
                      <th class="num">Market/day</th><th>Where</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="f in ahFarm" :key="f.item_id">
                      <td class="ah-item">
                        <button class="star" :class="{ on: watchSet.has(f.item_id) }"
                                :title="watchSet.has(f.item_id) ? 'Remove from watchlist' : 'Add to watchlist'"
                                @click="toggleWatch(f.item_id)">★</button>
                        <img :src="f.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a :href="wowheadItem(f.item_id)" data-wowhead="" target="_blank" rel="noopener">{{ f.name }}</a>
                        <span class="ah-sub">{{ f.subclass }}</span>
                      </td>
                      <td class="num nowrap">{{ fmtMoney(f.min_copper) }}</td>
                      <td class="num">{{ fmtNum(f.sold_24h) }}</td>
                      <td class="num nowrap ah-gold">{{ fmtMoney(f.market_day) }}</td>
                      <td class="ah-where">{{ farmZones(f.subclass) }}</td>
                    </tr>
                  </tbody>
                </table>
              </details>
            </div>

            <!-- ---------- My bags ---------- -->
            <div v-if="ahTab === 'bags'" class="guide">
              <details class="gsec" open>
                <summary class="guide-h4">
                  My bags — gold sitting idle
                </summary>
                <p class="ah-note">
                  <select v-model="ahBagsName" class="ah-select" @change="loadBags">
                    <option v-for="c in characters" :key="c.id" :value="c.name">{{ fullName(c) }}</option>
                  </select>
                  <span v-if="ahBags">Total value: <b class="ah-gold">{{ fmtMoney(ahBags.total) }}</b>
                    · scanned {{ fmtAgo(ahBags.updated_at) }}</span>
                </p>
                <p v-if="ahBagsError" class="muted">{{ ahBagsError }}</p>
                <table v-else-if="ahBags" class="guide-table">
                  <thead>
                    <tr>
                      <th>Item</th><th class="num">Count</th>
                      <th class="num">Cheapest</th><th class="num">Worth</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="b in ahBags.items" :key="b.item_id">
                      <td class="ah-item">
                        <img v-if="b.icon" :src="b.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a v-if="b.tracked" :href="wowheadItem(b.item_id)" data-wowhead="" target="_blank" rel="noopener">{{ b.name }}</a>
                        <span v-else class="muted">{{ b.name }} <span class="ah-sub">(no price)</span></span>
                      </td>
                      <td class="num">{{ fmtNum(b.count) }}</td>
                      <td class="num nowrap">{{ b.tracked ? fmtMoney(b.min_copper) : "—" }}</td>
                      <td class="num nowrap" :class="b.tracked ? 'ah-gold' : 'muted'">
                        {{ b.tracked ? fmtMoney(b.value) : "—" }}
                      </td>
                    </tr>
                  </tbody>
                </table>
                <p v-else class="muted">Loading…</p>
              </details>
            </div>

            <!-- ---------- Money makers ---------- -->
            <div v-if="ahTab === 'money'" class="guide">
              <details class="gsec" open>
                <summary class="guide-h4">
                  Flipping — buy low now, resell at the usual price
                </summary>
                <p class="ah-note">
                  Ranked by expected daily profit (margin × sold in 24h), after
                  the auction house cut. ⚠️ rows lean on the 7-day baseline —
                  double-check before buying big.
                </p>
                <p v-if="!ahFlips.length" class="muted">No safe flips right now.</p>
                <table v-else class="guide-table">
                  <thead>
                    <tr>
                      <th>Item</th><th class="num">Buy now</th><th class="num">Resell at</th>
                      <th class="num">Profit each</th><th class="num">/day</th><th class="num">Sold 24h</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="f in ahFlips" :key="f.item_id">
                      <td class="ah-item">
                        <button class="star" :class="{ on: watchSet.has(f.item_id) }"
                                :title="watchSet.has(f.item_id) ? 'Remove from watchlist' : 'Add to watchlist'"
                                @click="toggleWatch(f.item_id)">★</button>
                        <img :src="f.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a :href="wowheadItem(f.item_id)" data-wowhead="" target="_blank" rel="noopener">{{ f.name }}</a>
                        <span class="ah-sub">−{{ Math.round(f.discount * 100) }}%</span>
                      </td>
                      <td class="num nowrap">{{ fmtMoney(f.buy_copper) }}</td>
                      <td class="num nowrap muted">{{ fmtMoney(f.sell_copper) }}</td>
                      <td class="num nowrap ah-gold">{{ fmtMoney(f.profit_each) }}</td>
                      <td class="num nowrap ah-gold">{{ fmtMoney(f.daily_profit) }}</td>
                      <td class="num">{{ fmtNum(f.sold_24h) }}</td>
                    </tr>
                  </tbody>
                </table>
              </details>

              <details class="gsec" open>
                <summary class="guide-h4">
                  Vendor snipes — guaranteed profit, no auctioning
                </summary>
                <p class="ah-note">
                  Cheapest listing is below what a vendor pays: buy → sell to a
                  vendor → instant gold.
                </p>
                <p v-if="!ahSnipes.length" class="muted">Nothing under vendor price right now.</p>
                <table v-else class="guide-table">
                  <thead>
                    <tr>
                      <th>Item</th><th class="num">Buy now</th><th class="num">Vendor pays</th>
                      <th class="num">Profit each</th><th class="num">All ({{ fmtNum(ahSnipesUnits) }} units)</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="s in ahSnipes" :key="s.item_id">
                      <td class="ah-item">
                        <img :src="s.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a :href="wowheadItem(s.item_id)" data-wowhead="" target="_blank" rel="noopener">{{ s.name }}</a>
                      </td>
                      <td class="num nowrap">{{ fmtMoney(s.buy_copper) }}</td>
                      <td class="num nowrap muted">{{ fmtMoney(s.vendor_copper) }}</td>
                      <td class="num nowrap ah-gold">{{ fmtMoney(s.profit_each) }}</td>
                      <td class="num nowrap ah-gold">{{ fmtMoney(s.total_profit) }}</td>
                    </tr>
                  </tbody>
                </table>
              </details>
            </div>

            <!-- ---------- Watchlist ---------- -->
            <div v-if="ahTab === 'watch'" class="guide">
              <details class="gsec" open>
                <summary class="guide-h4">
                  Watchlist — {{ ahWatch.length }} item{{ ahWatch.length === 1 ? "" : "s" }}
                </summary>
                <p class="ah-note">
                  Star items in Deals or Listings to track them here. The trend
                  shows the cheapest price over the last 7 days — a spike means
                  someone listed far too high.
                </p>
                <p v-if="!ahWatch.length" class="muted">
                  Nothing watched yet — click ★ on any item.
                </p>
                <table v-else class="guide-table">
                  <thead>
                    <tr>
                      <th>Item</th><th class="num">Cheapest now</th><th class="num">Avg</th>
                      <th class="num">7-day trend</th><th class="num">Sold 24h</th><th></th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="w in ahWatch" :key="w.item_id">
                      <td class="ah-item">
                        <img :src="w.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a :href="wowheadItem(w.item_id)" data-wowhead="" target="_blank" rel="noopener">{{ w.name }}</a>
                        <button class="link-btn" title="Price details" @click.stop="openItem(w.item_id, w.name)">📈</button>
                        <span v-if="w.in_deals" class="ah-dealbadge">−{{ Math.round(w.discount * 100) }}%</span>
                      </td>
                      <td class="num nowrap">{{ fmtMoney(w.min_copper) }}
                        <span v-if="w.suspicious" class="ah-warn" title="average is skewed — someone listed far too high">⚠️</span></td>
                      <td class="num nowrap muted">{{ fmtMoney(w.avg_copper) }}</td>
                      <td class="num ah-sparkcell">
                        <svg class="spark" viewBox="0 0 60 18" preserveAspectRatio="none" aria-hidden="true">
                          <polyline :points="sparkPath(w.points)" />
                        </svg>
                      </td>
                      <td class="num">{{ fmtNum(w.sold_24h) }}</td>
                      <td class="num">
                        <button class="star on" title="Remove from watchlist"
                                @click="toggleWatch(w.item_id)">★</button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </details>
            </div>

            <!-- ---------- Crafting ---------- -->
            <div v-if="ahTab === 'crafting'" class="guide">
              <div class="guide-chips">
                <button v-for="p in ahCraftChipProfs" :key="p.id" class="chip"
                        :class="{ active: ahCraftProf === String(p.id) }"
                        @click="loadCrafting(String(p.id))">{{ p.name }}</button>
              </div>
              <details class="gsec" open>
                <summary class="guide-h4">
                  Crafting profits — {{ ahCraft ? ahCraft.total : 0 }} recipes
                </summary>
                <p class="ah-note">
                  Profit = cheapest sell price minus reagent cost, per craft.
                  <span class="ah-days">
                    <button class="chip" :class="{ active: ahCraftSort === 'profit' }"
                            @click="setCraftSort('profit')">By profit</button>
                    <button class="chip" :class="{ active: ahCraftSort === 'roi' }"
                            @click="setCraftSort('roi')">By demand × profit</button>
                    <button class="chip" :class="{ active: ahOnlyOurs }"
                            :disabled="!guildProfs.length"
                            :title="guildProfs.length ? 'Show only our professions' : 'No professions synced yet (addon /reload)'"
                            @click="ahOnlyOurs = !ahOnlyOurs">Only ours</button>
                  </span>
                </p>
                <table class="guide-table">
                  <thead>
                    <tr>
                      <th>Recipe</th><th class="num">Reagents</th><th class="num">Sells for</th>
                      <th class="num">Profit</th><th class="num">Sold 24h</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="r in ahCraftRows" :key="r.id">
                      <td class="ah-item">
                        <img :src="r.icon" alt="" loading="lazy" @error="$event.target.style.display='none'">
                        <a :href="wowheadItem(r.id)" data-wowhead="" target="_blank" rel="noopener">{{ r.name }}</a>
                      </td>
                      <td class="num nowrap muted">{{ fmtMoney(r.cost_copper) }}</td>
                      <td class="num nowrap">{{ fmtMoney(r.minimum_ask_copper) }}</td>
                      <td class="num nowrap ah-gold">{{ fmtMoney(r.profit_copper) }}</td>
                      <td class="num">{{ fmtNum(r.possible_24h) }}</td>
                    </tr>
                  </tbody>
                </table>
              </details>
            </div>

            <p class="ah-credit">
              Data from
              <a href="https://blackwater.exchange/m/eu/15/124" target="_blank" rel="noopener">Blackwater Exchange</a>
              (scanned hourly, kept 29 days) · cached on MunkeHub
              <span v-if="ahStatus && ahStatus.last_poll"> · last refresh {{ fmtAgo(ahStatus.last_poll) }}</span>
            </p>
          </template>
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

      <!-- ================= GUILD BOARD ================= -->
      <section v-show="tab === 'board'">
        <div class="panel">
          <h2 class="panel-title">Guild Board <span class="hint">— orders & events</span></h2>

          <div class="board-top">
            <div class="guide-chips">
              <button class="chip" :class="{ active: boardTab === 'orders' }" @click="boardTab = 'orders'">
                Crafting Orders<span v-if="board && board.orders.filter(o => o.status !== 'done').length">
                  ({{ board.orders.filter(o => o.status !== 'done').length }})</span>
              </button>
              <button class="chip" :class="{ active: boardTab === 'events' }" @click="boardTab = 'events'">
                Events<span v-if="board && board.events.length"> ({{ board.events.length }})</span>
              </button>
            </div>
            <select v-model="rsvpName" class="ah-select" title="Your name (for RSVPs and claims)">
              <option value="">Your name…</option>
              <option v-for="c in characters" :key="c.id" :value="c.name">{{ fullName(c) }}</option>
            </select>
          </div>

          <p v-if="boardError" class="error">{{ boardError }}</p>
          <p v-else-if="!board" class="muted">Loading the board…</p>

          <!-- crafting orders -->
          <template v-else-if="boardTab === 'orders'">
            <details class="post-box">
              <summary>Post a crafting order</summary>
              <form @submit.prevent="addOrder" class="board-form">
                <input v-model="newOrder.item_name" type="text" placeholder="What do you need? (e.g. Thorium Bar)" required maxlength="80" />
                <input v-model="newOrder.qty" type="number" min="1" max="10000" class="qty-input" title="Quantity" />
                <input v-model="newOrder.wanted_by" type="text" placeholder="Needed by (optional, e.g. Sunday raid)" maxlength="60" />
                <input v-model="newOrder.note" type="text" placeholder="Note (optional)" maxlength="120" />
                <button type="submit">Post order</button>
              </form>
            </details>

            <p v-if="!board.orders.length" class="muted">No orders yet — be the first to ask.</p>
            <div v-for="o in board.orders" :key="o.id" class="board-card" :class="{ done: o.status === 'done' }">
              <div class="bc-main">
                <b>{{ o.qty }}× {{ o.item_name }}</b>
                <span v-if="o.wanted_by" class="bc-when">needed by {{ o.wanted_by }}</span>
                <span v-if="o.note" class="bc-note">{{ o.note }}</span>
              </div>
              <div class="bc-side">
                <span v-if="o.status === 'claimed'" class="bc-claim">🔨 {{ o.claimed_by }} is on it</span>
                <span v-else-if="o.status === 'done'" class="bc-done">✓ done</span>
                <button v-if="o.status === 'open'" class="chip" @click="setOrderStatus(o, 'claimed')">I'll craft it</button>
                <button v-if="o.status === 'claimed'" class="chip" @click="setOrderStatus(o, 'done')">Done</button>
                <button v-if="o.status !== 'done'" class="chip" @click="setOrderStatus(o, 'open')">Unclaim</button>
                <button class="star" title="Remove order" @click="removeOrder(o)">✕</button>
              </div>
            </div>
          </template>

          <!-- events -->
          <template v-else>
            <details class="post-box">
              <summary>Create an event</summary>
              <form @submit.prevent="addEvent" class="board-form">
                <input v-model="newEvent.title" type="text" placeholder="Title (e.g. Zul'Gurub raid)" required maxlength="80" />
                <input v-model="newEvent.when_text" type="text" placeholder="When (e.g. Saturday 20:00 server)" maxlength="80" />
                <input v-model="newEvent.note" type="text" placeholder="Note (optional)" maxlength="160" />
                <button type="submit">Create event</button>
              </form>
            </details>

            <p v-if="!board.events.length" class="muted">No events yet — plan the first one.</p>
            <div v-for="e in board.events" :key="e.id" class="board-card">
              <div class="bc-main">
                <b>{{ e.title }}</b>
                <span v-if="e.when_text" class="bc-when">🗓 {{ e.when_text }}</span>
                <span v-if="e.note" class="bc-note">{{ e.note }}</span>
              </div>
              <div class="bc-rsvp">
                <button class="rsvp-btn" :class="{ on: myRsvp(e) === 'yes' }" @click="rsvp(e, 'yes')">
                  ✓ Going {{ rsvpCount(e, 'yes') || '' }}
                </button>
                <button class="rsvp-btn maybe" :class="{ on: myRsvp(e) === 'maybe' }" @click="rsvp(e, 'maybe')">
                  ? Maybe {{ rsvpCount(e, 'maybe') || '' }}
                </button>
                <button class="rsvp-btn no" :class="{ on: myRsvp(e) === 'no' }" @click="rsvp(e, 'no')">
                  ✕ Can't {{ rsvpCount(e, 'no') || '' }}
                </button>
                <button class="star" title="Remove event" @click="removeEvent(e)">✕</button>
              </div>
            </div>
          </template>
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

      <!-- item detail popup: 7-day chart + price facts -->
      <div v-if="itemPopup" class="popup-backdrop" @click.self="closeItem">
        <div class="item-popup">
          <div class="popup-head">
            <h3>{{ itemPopup.name }}</h3>
            <button class="star" title="Close" @click="closeItem">✕</button>
          </div>
          <p v-if="!itemPopup.price" class="muted">No price data for this item.</p>
          <template v-else>
            <div class="popup-stats">
              <div class="popup-stat">
                <span class="dash-num">{{ fmtMoney(itemPopup.price.min) }}</span>
                <span class="dash-label">Cheapest now</span>
              </div>
              <div class="popup-stat">
                <span class="dash-num">{{ fmtMoney(itemPopup.price.avg7 || 0) }}</span>
                <span class="dash-label">7-day average</span>
              </div>
              <div class="popup-stat">
                <span class="dash-num">{{ fmtNum(itemPopup.price.units) }}</span>
                <span class="dash-label">Units for sale</span>
              </div>
              <div class="popup-stat">
                <span class="dash-num">{{ fmtNum(itemPopup.price.sold24h) }}</span>
                <span class="dash-label">Sold, 24 h</span>
              </div>
            </div>
            <svg v-if="itemPopup.points.length > 1" class="popup-spark" viewBox="0 0 300 60" preserveAspectRatio="none">
              <polyline :points="popupSpark(itemPopup.points)" />
            </svg>
            <p v-else class="muted small">Not enough history yet — the chart fills in as hourly snapshots accumulate.</p>
            <p class="muted small">
              <a :href="wowheadItem(itemPopup.id)" data-wowhead="" target="_blank" rel="noopener">View on Wowhead ↗</a>
            </p>
          </template>
        </div>
      </div>

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

// Guides progress (tier checkboxes) — per-browser, no account needed.
function loadGuideProgress() {
  try {
    return JSON.parse(localStorage.getItem("munkehub_guide_progress")) || {};
  } catch {
    return {}
  }
}

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
      // selected profession in the Guides tab (rail selector)
      guideSel: "cooking",
      // Guides tab — live market data (loaded per selected guide)
      guidePrices: null, // item_id -> {min, avg7, units, sold24h}
      guideBags: null, // item_id -> {n, chars} (merged guild bag scans)
      guideProgress: loadGuideProgress(), // guideId -> [tier ranges done]
      // presence: name -> {online, last_seen} (helper heartbeats)
      presence: {},
      // Guild Board: crafting orders + events
      board: null, // {orders: [], events: []}
      boardTab: "orders",
      boardError: "",
      newOrder: { item_name: "", qty: 1, wanted_by: "", note: "" },
      newEvent: { title: "", when_text: "", note: "" },
      rsvpName: localStorage.getItem("munkehub_rsvp_name") || "",
      // item detail popup (7-day chart + price facts)
      itemPopup: null, // {id, name, price, points}
      // buy list cart (aggregated shopping list minus bags)
      cartOpen: false,
      // --- Auction House tab ---
      ahTab: "deals", // deals | sell | listings | crafting
      ahSummary: null,
      ahError: false,
      ahDeals: null,
      ahSales: null,
      ahSalesDays: "1",
      ahListings: null,
      ahQuery: "",
      ahCraft: null,
      ahCraftProf: "171",
      ahCraftProfs: [],
      ahStatus: null,
      // MunkeHub AI advisor
      ahAdvice: null,
      ahAdviceLoading: false,
      ahAdviceError: "",
      // price history + watchlist (⭐)
      ahWatch: [],
      ahHist: {}, // item_id -> [{t, min}] (7-day sparkline points)
      // money makers (flips + vendor snipes)
      ahFlips: [],
      ahSnipes: [],
      // farming table (guild professions × market)
      ahFarm: [],
      ahFarmProfs: [],
      // "only what WE can craft/farm" filter (Listings + Crafting chips)
      ahOnlyOurs: false,
      // my bags (addon → helper → API, valued against the market)
      ahBags: null,
      ahBagsName: "",
      ahBagsError: "",
      // crafting sort: "profit" (upstream default) | "roi" (profit × sold)
      ahCraftSort: "profit",
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
    // Guides: merged bag scans ("we already have this") + default guide prices.
    this.loadGuideBags();
    this.loadGuideData();
    // presence dots + Guild Board
    this.loadPresence();
    this.loadBoard();
    // Auction House: load the default section (deals) + overview + AI advice.
    this.loadAh();
    this.loadAdvice(false);
    this.loadWatch();
  },
  computed: {
    // Guides tab: only the profession picked in the chip selector shows.
    visibleGuides() {
      return this.guides.filter((g) => g.id === this.guideSel);
    },
    // Auction House sections (chip selector).
    ahSections() {
      return [
        { id: "deals", name: "Deals" },
        { id: "sell", name: "Best to sell" },
        { id: "listings", name: "Listings" },
        { id: "crafting", name: "Crafting" },
        { id: "farm", name: "Farming" },
        { id: "money", name: "Money makers" },
        { id: "bags", name: "My bags" },
        { id: "watch", name: `Watchlist${this.ahWatch.length ? " (" + this.ahWatch.length + ")" : ""}` },
      ];
    },
    // Watched items currently listed as deals → nav badge count.
    watchDealsCount() {
      return this.ahWatch.filter((w) => w.in_deals).length;
    },
    // Subclasses the guild can farm/consume (same mapping as the API).
    ourSubclasses() {
      const map = {
        herbalism: "Herb",
        mining: "Metal & Stone",
        skinning: "Leather",
        fishing: "Cooking",
        cooking: "Cooking",
        tailoring: "Cloth",
      };
      const set = new Set();
      this.guildProfs.forEach((p) => {
        const s = map[p.toLowerCase()];
        if (s) set.add(s);
      });
      return set;
    },
    // "Only what WE can craft/farm" filter for Listings.
    ahListingRows() {
      let rows = this.ahListings && this.ahListings.items ? this.ahListings.items : [];
      if (this.ahOnlyOurs && this.ourSubclasses.size) {
        rows = rows.filter((i) => this.ourSubclasses.has(i.subclass_name));
      }
      return rows;
    },
    // Set of watched item ids — for the ★ state in every table.
    watchSet() {
      return new Set(this.ahWatch.map((w) => w.item_id));
    },
    ahDealRows() {
      return (this.ahDeals && this.ahDeals.items ? this.ahDeals.items : []).slice(0, 50);
    },
    ahSalesItems() {
      return (this.ahSales && this.ahSales.items ? this.ahSales.items : []).slice(0, 50);
    },
    ahCraftRows() {
      return (this.ahCraft && this.ahCraft.recipes ? this.ahCraft.recipes : []).slice(0, 50);
    },
    // Profession chips: filtered to the guild's own professions when the
    // "only ours" toggle is on (recipes are per-profession, so the filter
    // lives on the selector, not the rows).
    ahCraftChipProfs() {
      if (!this.ahOnlyOurs || !this.guildProfs.length) return this.ahCraftProfs;
      const ours = new Set(this.guildProfs.map((p) => p.toLowerCase()));
      return this.ahCraftProfs.filter((p) => ours.has(String(p.name).toLowerCase()));
    },
    // total units available across all vendor snipes (for the "All" column)
    ahSnipesUnits() {
      return this.ahSnipes.reduce((n, s) => n + (s.units || 0), 0);
    },
    // Union of every roster character's professions (synced by the addon).
    guildProfs() {
      const set = new Set();
      (this.characters || []).forEach((c) => {
        (c.professions || "").split(",").forEach((p) => {
          p = p.trim();
          if (p) set.add(p);
        });
      });
      return [...set].sort();
    },
    // Epic mount fund progress (guild gold vs goal, from /stats).
    mountPct() {
      if (!this.stats || !this.stats.mount_goal) return 0;
      return Math.min(100, Math.round((this.stats.total_gold / this.stats.mount_goal) * 100));
    },
    // Names currently in game (helper heartbeats < 3 min old).
    onlineSet() {
      const s = new Set();
      Object.entries(this.presence || {}).forEach(([name, p]) => {
        if (p && p.online) s.add(name);
      });
      return s;
    },
    // Best skill level for the selected guide's profession (any character).
    guideSkill() {
      const g = this.guides.find((x) => x.id === this.guideSel);
      if (!g) return null;
      const want = g.name.toLowerCase();
      let best = null;
      (this.characters || []).forEach((c) => {
        (c.prof_levels || "").split(",").forEach((pair) => {
          const [pname, lvl] = pair.split(":");
          if (!pname || !lvl) return;
          if (pname.trim().toLowerCase() !== want) return;
          const n = parseInt(lvl, 10) || 0;
          if (!best || n > best.level) best = { name: c.name, level: n };
        });
      });
      return best;
    },
    // Buy list cart: every shopping tier aggregated, minus what's in bags.
    // `from` = only tiers starting at/after this skill (for "cost to 300").
    buyCart() {
      const g = this.guides.find((x) => x.id === this.guideSel);
      if (!g || !g.shopping) return null;
      const skill = this.guideSkill ? this.guideSkill.level : 0;
      const agg = new Map(); // id -> {name, qty, id}
      let remaining = 0;
      g.shopping.forEach((s) => {
        const start = parseInt(String(s.range).match(/\d+/)?.[0] || "0", 10);
        const fromHere = start >= skill;
        (s.items || []).forEach((it) => {
          if (!it.id) return;
          const cur = agg.get(it.id) || { id: it.id, name: it.name, qty: 0 };
          cur.qty += this.qtyNum(it.qty);
          agg.set(it.id, cur);
          if (fromHere) {
            const p = this.priceOf(it.id);
            const have = this.haveOf(it.id);
            const need = Math.max(0, this.qtyNum(it.qty) - (have ? have.n : 0));
            if (p) remaining += p.min * need;
          }
        });
      });
      const rows = [];
      let total = 0;
      let missing = 0;
      agg.forEach((it) => {
        const have = this.haveOf(it.id);
        const need = Math.max(0, it.qty - (have ? have.n : 0));
        const p = this.priceOf(it.id);
        const cost = p ? p.min * need : 0;
        if (!p) missing++;
        total += cost;
        rows.push({ ...it, have: have ? have.n : 0, need, cost, price: p ? p.min : 0, priced: !!p });
      });
      rows.sort((a, b) => b.cost - a.cost);
      return { rows, total, missing, remaining };
    },
  },
  methods: {
    formatDate(iso) {
      if (!iso) return "";
      return new Date(iso).toLocaleString();
    },
    // --- Auction House -------------------------------------------------
    async loadAh() {
      try {
        this.ahSummary = await api.ahSummary();
        this.ahDeals = await api.ahDeals();
        this.ahStatus = await api.ahStatus().catch(() => null);
        this.ahError = false;
        // sparklines for the deals table (one batched request)
        this.loadHistory(this.ahDealRows.map((d) => d.item_id));
      } catch {
        this.ahError = true;
      }
    },
    selectAhSection(id) {
      this.ahTab = id;
      if (id === "sell" && !this.ahSales) this.loadSales("1");
      if (id === "listings" && !this.ahListings) this.searchListings();
      if (id === "crafting" && !this.ahCraft) this.loadCrafting(this.ahCraftProf);
      if (id === "money" && !this.ahFlips.length && !this.ahSnipes.length) this.loadMoney();
      if (id === "farm" && !this.ahFarm.length) this.loadFarm();
      if (id === "bags" && !this.ahBags) {
        this.ahBagsName = this.ahBagsName || (this.characters[0] && this.characters[0].name) || "";
        this.loadBags();
      }
      if (id === "watch") this.loadWatch();
    },
    // Bag contents valued against the market (needs addon + helper running).
    async loadBags() {
      this.ahBagsError = "";
      this.ahBags = null;
      if (!this.ahBagsName) {
        this.ahBagsError = "No character selected.";
        return;
      }
      try {
        this.ahBags = await api.ahBags(this.ahBagsName);
      } catch (e) {
        this.ahBags = null;
        this.ahBagsError = String(e.message || e).includes("404")
          ? "No bag data yet — /reload in-game, then wait for the helper to sync."
          : "Could not load bags — try again.";
      }
    },
    // Farming table: guild professions × market (server ranks by market/day).
    async loadFarm() {
      try {
        const res = await api.ahFarming();
        this.ahFarm = res.items || [];
        this.ahFarmProfs = res.professions || [];
      } catch {
        this.ahFarm = [];
      }
    },
    // Zone hints per subclass — from the guides' leveling routes (frontend).
    farmZones(subclass) {
      // cloth drops from humanoids — no gathering route exists
      if (subclass === "Cloth") return "Humanoid mobs (dungeons, quests)";
      const guideId = {
        Herb: "herbalism",
        "Metal & Stone": "mining",
        Leather: "skinning",
        Cooking: "fishing",
      }[subclass];
      const g = this.guides.find((x) => x.id === guideId);
      if (!g || !g.gathering || !g.gathering.length) return "—";
      // last tier = endgame zones (where the valuable stuff drops)
      const last = g.gathering[g.gathering.length - 1];
      const zones = (last.zones || []).map((z) => (Array.isArray(z) ? z[0] : z)).filter(Boolean);
      return zones.slice(0, 3).join(", ") || "—";
    },
    // Money makers: flips + vendor snipes (two parallel calls, no cache needed).
    async loadMoney() {
      try {
        const [flips, snipes] = await Promise.all([api.ahFlips(), api.ahVendorSnipes()]);
        this.ahFlips = flips.items || [];
        this.ahSnipes = snipes.items || [];
      } catch {
        this.ahFlips = [];
        this.ahSnipes = [];
      }
    },
    // --- price history + watchlist ------------------------------------
    // One batched request for up to 60 sparklines (server downsamples to 40).
    async loadHistory(ids) {
      const missing = ids.filter((id) => !this.ahHist[id]);
      if (!missing.length) return;
      try {
        const res = await api.ahHistory(missing.slice(0, 60));
        this.ahHist = { ...this.ahHist, ...res.points };
      } catch {
        /* sparklines are decoration — no error state needed */
      }
    },
    async loadWatch() {
      try {
        const res = await api.ahWatchlist();
        this.ahWatch = res.items || [];
      } catch {
        this.ahWatch = [];
      }
    },
    async toggleWatch(id) {
      const watched = this.watchSet.has(id);
      try {
        if (watched) {
          await api.ahWatchRemove(id);
          this.ahWatch = this.ahWatch.filter((w) => w.item_id !== id);
        } else {
          await api.ahWatchAdd(id);
          await this.loadWatch();
        }
      } catch {
        /* friends-only site — ignore, next load resyncs */
      }
    },
    // Larger sparkline for the item popup (300×60 viewBox).
    popupSpark(points) {
      if (!points || points.length < 2) return "";
      const mins = points.map((p) => p.min);
      const lo = Math.min(...mins);
      const hi = Math.max(...mins);
      const span = hi - lo || 1;
      return points
        .map((p, i) => {
          const x = (i / (points.length - 1)) * 300;
          const y = 56 - ((p.min - lo) / span) * 52;
          return `${x.toFixed(1)},${y.toFixed(1)}`;
        })
        .join(" ");
    },
    // Sparkline polyline points from [{t, min}] — flat line when no history.
    sparkPath(points) {
      if (!points || points.length < 2) return "0,9 60,9";
      const mins = points.map((p) => p.min);
      const lo = Math.min(...mins);
      const hi = Math.max(...mins);
      const span = hi - lo || 1;
      return points
        .map((p, i) => {
          const x = (i / (points.length - 1)) * 60;
          const y = 17 - ((p.min - lo) / span) * 16;
          return `${x.toFixed(1)},${y.toFixed(1)}`;
        })
        .join(" ");
    },
    async loadSales(days) {
      this.ahSalesDays = days;
      try {
        this.ahSales = await api.ahSales(days);
      } catch {
        this.ahSales = null;
      }
    },
    // MunkeHub AI advisor — cached 6h server-side; refresh=1 forces a new call.
    async loadAdvice(refresh) {
      this.ahAdviceLoading = true;
      this.ahAdviceError = "";
      try {
        this.ahAdvice = await api.ahAdvice(refresh);
      } catch (e) {
        this.ahAdviceError =
          String(e.message || e).includes("503")
            ? "AI advisor not configured yet (server is missing its key)."
            : "AI advisor unavailable right now — try Refresh.";
      } finally {
        this.ahAdviceLoading = false;
      }
    },
    async searchListings() {
      try {
        this.ahListings = await api.ahListings(this.ahQuery, 1, 50);
        this.loadHistory(this.ahListingRows.map((i) => i.id));
      } catch {
        this.ahListings = null;
      }
    },
    async loadCrafting(prof) {
      this.ahCraftProf = prof;
      try {
        this.ahCraft = await api.ahCrafting(prof, this.ahCraftSort);
        if (!this.ahCraftProfs.length && this.ahCraft.professions) {
          this.ahCraftProfs = this.ahCraft.professions;
        }
      } catch {
        this.ahCraft = null;
      }
    },
    // Toggle crafting sort: raw profit vs profit × sold/24h (real demand).
    setCraftSort(s) {
      if (this.ahCraftSort === s) return;
      this.ahCraftSort = s;
      this.loadCrafting(this.ahCraftProf);
    },
    // Copper → "1g 23s 45c" (gold hidden below 1g, like the game).
    fmtMoney(copper) {
      const c = Math.round(Number(copper) || 0);
      if (c < 100) return c + "c";
      const g = Math.floor(c / 10000);
      const s = Math.floor((c % 10000) / 100);
      const r = c % 100;
      if (!g) return `${s}s ${r}c`;
      return `${g}g ${s}s ${r}c`;
    },
    fmtNum(n) {
      const v = Number(n) || 0;
      return v.toLocaleString("en-GB");
    },
    // "2026-11-19T13:42:12Z" → "12 minutes ago".
    fmtAgo(iso) {
      const t = new Date(iso).getTime();
      if (!t) return "";
      const mins = Math.max(0, Math.round((Date.now() - t) / 60000));
      if (mins < 1) return "just now";
      if (mins < 60) return `${mins} min ago`;
      const h = Math.floor(mins / 60);
      if (h < 24) return `${h}h ago`;
      return `${Math.floor(h / 24)}d ago`;
    },
    // --- Guides tab: live prices, bags, progress, crafters ----------------
    selectGuide(id) {
      this.guideSel = id;
      this.loadGuideData();
    },
    // Prices for every item the visible guide references (one batched call).
    async loadGuideData() {
      const g = this.guides.find((x) => x.id === this.guideSel);
      if (!g) return;
      const ids = new Set();
      (g.shopping || []).forEach((s) =>
        (s.items || []).forEach((it) => it.id && ids.add(it.id))
      );
      (g.steps || []).forEach((st) =>
        (st.rows || []).forEach((row) =>
          (row.reagents || []).forEach((rg) => rg.id && ids.add(rg.id))
        )
      );
      (g.items || []).forEach((it) => it.id && ids.add(it.id));
      if (!ids.size) return;
      try {
        const res = await api.ahPrices([...ids].slice(0, 200));
        this.guidePrices = { ...(this.guidePrices || {}), ...(res.prices || {}) };
      } catch {
        /* prices are decoration — the guide works without them */
      }
    },
    // Merged bag scans from every character — "we already have 60 of these".
    async loadGuideBags() {
      try {
        const res = await api.ahBagsAll();
        const map = {};
        (res.items || []).forEach((b) => {
          map[b.item_id] = { n: b.count, chars: b.chars };
        });
        this.guideBags = map;
      } catch {
        this.guideBags = null;
      }
    },
    priceOf(id) {
      return (this.guidePrices && this.guidePrices[id]) || null;
    },
    haveOf(id) {
      return (this.guideBags && this.guideBags[id]) || null;
    },
    // "~52" / "12" / "408-459" → number (leading qty; ranges take the low end)
    qtyNum(q) {
      const m = String(q == null ? "" : q).match(/\d+/);
      return m ? parseInt(m[0], 10) : 0;
    },
    // Enough in our bags to cover this line?
    haveFull(it) {
      const need = this.qtyNum(it.qty);
      const h = it.id ? this.haveOf(it.id) : null;
      return !!(h && need > 0 && h.n >= need);
    },
    // Cost of one shopping-list tier at today's cheapest prices.
    tierTotal(s) {
      let sum = 0;
      let missing = 0;
      (s.items || []).forEach((it) => {
        const p = it.id ? this.priceOf(it.id) : null;
        if (p) sum += p.min * this.qtyNum(it.qty);
        else missing++;
      });
      return { sum, missing };
    },
    // −% badge when today's cheapest is 20%+ under the 7-day average.
    dealPct(id) {
      const p = this.priceOf(id);
      if (!p || !p.avg7 || !(p.min > 0) || p.min > 0.8 * p.avg7) return 0;
      return Math.round((1 - p.min / p.avg7) * 100);
    },
    // Tier checkboxes (progress tracker) — shopping ranges, else route ranges.
    guideTiers(g) {
      const src =
        (g.shopping && g.shopping.length && g.shopping) ||
        (g.steps && g.steps.length && g.steps) ||
        (g.gathering && g.gathering.length && g.gathering) ||
        [];
      return src.map((x) => x.range);
    },
    guidePct(g) {
      const tiers = this.guideTiers(g);
      if (!tiers.length) return 0;
      const done = (this.guideProgress[g.id] || []).filter((t) =>
        tiers.includes(t)
      ).length;
      return Math.round((done / tiers.length) * 100);
    },
    isTierDone(g, t) {
      return (this.guideProgress[g.id] || []).includes(t);
    },
    toggleTier(g, t) {
      const cur = new Set(this.guideProgress[g.id] || []);
      if (cur.has(t)) cur.delete(t);
      else cur.add(t);
      this.guideProgress = { ...this.guideProgress, [g.id]: [...cur] };
      try {
        localStorage.setItem(
          "munkehub_guide_progress",
          JSON.stringify(this.guideProgress)
        );
      } catch {
        /* private mode — progress just won't persist */
      }
    },
    // Guild members with this profession (addon-synced professions column).
    guideCrafters(g) {
      const want = (g.name || "").toLowerCase();
      return (this.characters || []).filter((c) =>
        (c.professions || "")
          .split(",")
          .some((p) => p.trim().toLowerCase() === want)
      );
    },
    classColor(c) {
      return CLASS_COLORS[(c.class || "").toLowerCase()] || "var(--ink)";
    },
    // --- presence ----------------------------------------------------------
    async loadPresence() {
      try {
        const res = await api.presence();
        this.presence = res.presence || {};
      } catch {
        this.presence = {};
      }
    },
    isOnline(name) {
      return this.onlineSet.has(name);
    },
    // --- Guild Board -------------------------------------------------------
    async loadBoard() {
      try {
        this.board = await api.board();
        this.boardError = "";
      } catch {
        this.boardError = "Could not load the board.";
      }
    },
    async addOrder() {
      if (!this.newOrder.item_name.trim()) return;
      try {
        await api.boardOrderAdd(
          { ...this.newOrder, qty: parseInt(this.newOrder.qty, 10) || 1 },
          this.postToken
        );
        this.newOrder = { item_name: "", qty: 1, wanted_by: "", note: "" };
        await this.loadBoard();
      } catch {
        this.boardError = "Posting failed — check your member token.";
      }
    },
    async setOrderStatus(o, status) {
      try {
        await api.boardOrderStatus(o.id, { name: this.rsvpName || "Someone", status }, this.postToken);
        await this.loadBoard();
      } catch {
        this.boardError = "Update failed — check your member token.";
      }
    },
    async removeOrder(o) {
      try {
        await api.boardOrderDelete(o.id, this.postToken);
        await this.loadBoard();
      } catch {
        /* friends-only site — ignore */
      }
    },
    async addEvent() {
      if (!this.newEvent.title.trim()) return;
      try {
        await api.boardEventAdd(this.newEvent, this.postToken);
        this.newEvent = { title: "", when_text: "", note: "" };
        await this.loadBoard();
      } catch {
        this.boardError = "Posting failed — check your member token.";
      }
    },
    async removeEvent(e) {
      try {
        await api.boardEventDelete(e.id, this.postToken);
        await this.loadBoard();
      } catch {
        /* ignore */
      }
    },
    async rsvp(e, status) {
      if (!this.rsvpName.trim()) {
        this.boardError = "Pick your name first (top of the board).";
        return;
      }
      localStorage.setItem("munkehub_rsvp_name", this.rsvpName.trim());
      try {
        await api.boardRsvp(e.id, { name: this.rsvpName.trim(), status }, this.postToken);
        await this.loadBoard();
      } catch {
        this.boardError = "RSVP failed — check your member token.";
      }
    },
    rsvpCount(e, status) {
      return (e.rsvps || []).filter((r) => r.status === status).length;
    },
    myRsvp(e) {
      const me = (e.rsvps || []).find((r) => r.name === this.rsvpName.trim());
      return me ? me.status : "";
    },
    // --- item detail popup -------------------------------------------------
    async openItem(id, name) {
      this.itemPopup = { id, name, price: null, points: [] };
      try {
        const [pr, hi] = await Promise.all([api.ahPrices([id]), api.ahHistory([id])]);
        this.itemPopup = {
          id,
          name,
          price: (pr.prices && pr.prices[id]) || null,
          points: (hi.points && hi.points[id]) || [],
        };
      } catch {
        /* keep the popup open with what we have */
      }
    },
    closeItem() {
      this.itemPopup = null;
    },
    // --- buy list cart -----------------------------------------------------
    cartText() {
      const c = this.buyCart;
      if (!c) return "";
      const lines = c.rows
        .filter((r) => r.need > 0)
        .map((r) => `${r.need}x ${r.name}`);
      return lines.join("\n") + `\nTotal: ${this.fmtMoney(c.total)}`;
    },
    async copyCart() {
      try {
        await navigator.clipboard.writeText(this.cartText());
      } catch {
        /* clipboard blocked — the text is visible for manual copy */
      }
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
  max-width: 1360px;
  margin: 0 auto;
  padding: 1.5rem 1.5rem 3rem;
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
/* text-heavy tabs stay readable on wide screens */
.guide-intro,
.steps,
.cl-list,
.notice p,
.explain-card,
.helper-steps {
  max-width: 760px;
}
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
/* deal-alert badge on the Auction House tab (watched items on sale) */
.nav-badge {
  display: inline-block;
  margin-left: 0.4rem;
  background: #b3402f;
  color: #fff;
  border-radius: 999px;
  font-size: 0.68rem;
  font-weight: 700;
  min-width: 1.15rem;
  padding: 0.05rem 0.3rem;
  text-align: center;
}
.ah-where {
  color: var(--muted);
  font-size: 0.78rem;
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
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  gap: 0.8rem;
  margin-bottom: 0.8rem;
}
/* online dot on roster names */
.online-dot {
  color: #555;
  font-size: 0.7rem;
  margin-right: 0.35rem;
  vertical-align: 1px;
}
.online-dot.on {
  color: #7fd18a;
  text-shadow: 0 0 6px rgba(127, 209, 138, 0.8);
}
.mount-bar {
  width: 100%;
  margin-top: 0.35rem;
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
/* rail (left) + content (right) — rail sticks while reading on wide screens */
.guides-layout {
  display: grid;
  grid-template-columns: 210px minmax(0, 1fr);
  gap: 1.4rem;
  align-items: start;
}
.guides-rail {
  position: sticky;
  top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
.rail-item {
  display: flex;
  align-items: center;
  gap: 0.55rem;
  width: 100%;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 8px;
  color: var(--muted);
  font: inherit;
  font-size: 0.88rem;
  padding: 0.4rem 0.6rem;
  cursor: pointer;
  text-align: left;
  transition: all 0.12s ease;
}
.rail-item:hover {
  color: var(--ink);
  border-color: var(--edge);
}
.rail-item.active {
  color: var(--gold-bright);
  background: rgba(200, 164, 77, 0.1);
  border-color: var(--gold);
  font-weight: 700;
}
.rail-icon {
  width: 22px;
  height: 22px;
  border: 1px solid var(--edge);
  border-radius: 4px;
  background: #0e0d11;
  padding: 1px;
  flex: none;
}
.rail-item.active .rail-icon {
  border-color: var(--gold);
}
.rail-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.rail-pct {
  font-size: 0.68rem;
  font-weight: 700;
  color: var(--gold-bright);
  background: rgba(200, 164, 77, 0.14);
  border-radius: 999px;
  padding: 0.05rem 0.4rem;
}
.guide-head {
  align-items: flex-start;
}
.guide-headtext {
  flex: 1;
  min-width: 0;
}
/* progress tracker (tier checkboxes) */
.guide-progress {
  margin-left: auto;
  flex: none;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.25rem;
  min-width: 130px;
}
.gp-bar {
  width: 130px;
  height: 8px;
  background: #14131a;
  border: 1px solid var(--edge);
  border-radius: 999px;
  overflow: hidden;
}
.gp-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--gold), var(--gold-bright));
  border-radius: 999px;
  transition: width 0.25s ease;
}
.gp-label {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--muted);
}
/* tier checkbox inside range badges */
.tier-check {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  cursor: pointer;
  user-select: none;
}
.tier-check input {
  accent-color: var(--gold);
  width: 14px;
  height: 14px;
  margin: 0;
  cursor: pointer;
}
.tier-check input:checked + span {
  text-decoration: line-through;
  opacity: 0.55;
}
/* inside .step-head the badge is the head itself — no nested badge */
.step-head .tier-check span {
  background: none;
  border: none;
  padding: 0;
  font-size: inherit;
  color: inherit;
}
/* guild crafters line */
.guide-crafters {
  display: flex;
  flex-wrap: wrap;
  gap: 0.35rem;
  align-items: baseline;
  background: rgba(200, 164, 77, 0.06);
  border: 1px solid rgba(200, 164, 77, 0.3);
  border-radius: 8px;
  padding: 0.5rem 0.8rem;
  margin-bottom: 0.8rem;
  font-size: 0.85rem;
}
.gc-label {
  color: var(--gold);
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-right: 0.15rem;
}
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

/* ---------------- auction house ---------------- */
.ah-dash {
  grid-template-columns: repeat(4, 1fr);
}
/* "1870g 61s 2c" is long — keep it on one line in the card */
.ah-money {
  font-size: 1.15rem;
  white-space: nowrap;
}
.guide-table .num {
  text-align: right;
}
.ah-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.ah-item img {
  width: 22px;
  height: 22px;
  border-radius: 3px;
  flex: none;
}
.ah-sub {
  color: var(--muted);
  font-size: 0.72rem;
  margin-left: 0.35rem;
}
.ah-under {
  color: #7fd18a;
  font-weight: 700;
}
.ah-gold {
  color: var(--gold-bright);
  font-weight: 700;
}
.ah-below {
  color: #e07a7a;
}
.ah-note {
  font-size: 0.8rem;
  color: var(--muted);
  margin: 0.4rem 0 0.7rem;
}
.ah-days {
  display: inline-flex;
  gap: 0.4rem;
  margin-left: 0.6rem;
}
.ah-days .chip {
  padding: 0.15rem 0.6rem;
  font-size: 0.72rem;
}
.ah-search {
  background: #14131a;
  border: 1px solid var(--edge);
  border-radius: 6px;
  color: var(--ink);
  font: inherit;
  font-size: 0.85rem;
  padding: 0.35rem 0.7rem;
  width: 240px;
  max-width: 100%;
}
.ah-search:focus {
  outline: none;
  border-color: var(--gold);
}
.ah-select {
  background: #14131a;
  border: 1px solid var(--edge);
  border-radius: 6px;
  color: var(--ink);
  font: inherit;
  font-size: 0.85rem;
  margin-right: 0.6rem;
  padding: 0.3rem 0.5rem;
}
.ah-select:focus {
  outline: none;
  border-color: var(--gold);
}
.ah-credit {
  margin-top: 1rem;
  font-size: 0.72rem;
  color: var(--muted);
  text-align: center;
}
/* ---------------- Guild Board ---------------- */
.board-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.8rem;
  flex-wrap: wrap;
  margin-bottom: 0.8rem;
}
.board-form {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  padding: 0.6rem 0;
}
.board-form input {
  background: #14131a;
  border: 1px solid var(--edge);
  border-radius: 6px;
  color: var(--ink);
  font: inherit;
  font-size: 0.85rem;
  padding: 0.35rem 0.7rem;
  flex: 1;
  min-width: 140px;
}
.board-form input:focus {
  outline: none;
  border-color: var(--gold);
}
.board-form .qty-input {
  flex: 0 0 70px;
  min-width: 70px;
}
.board-form button {
  background: linear-gradient(180deg, #d9b45c, #a8842f);
  border: 1px solid var(--gold);
  border-radius: 6px;
  color: #1a1508;
  font: inherit;
  font-weight: 700;
  font-size: 0.82rem;
  padding: 0.35rem 0.9rem;
  cursor: pointer;
}
.board-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.8rem;
  flex-wrap: wrap;
  background: var(--stone-light);
  border: 1px solid var(--edge);
  border-radius: 8px;
  padding: 0.65rem 0.9rem;
  margin-bottom: 0.55rem;
}
.board-card.done {
  opacity: 0.55;
}
.board-card.done .bc-main b {
  text-decoration: line-through;
}
.bc-main {
  display: flex;
  align-items: baseline;
  gap: 0.6rem;
  flex-wrap: wrap;
  min-width: 0;
}
.bc-when {
  color: var(--gold-bright);
  font-size: 0.8rem;
}
.bc-note {
  color: var(--muted);
  font-size: 0.8rem;
}
.bc-side,
.bc-rsvp {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  flex-wrap: wrap;
}
.bc-claim {
  color: var(--gold-bright);
  font-size: 0.82rem;
}
.bc-done {
  color: #7fd18a;
  font-weight: 700;
  font-size: 0.82rem;
}
.rsvp-btn {
  background: #14131a;
  border: 1px solid var(--edge);
  border-radius: 6px;
  color: var(--muted);
  font: inherit;
  font-size: 0.78rem;
  padding: 0.25rem 0.6rem;
  cursor: pointer;
  transition: all 0.12s;
}
.rsvp-btn:hover {
  border-color: var(--gold);
  color: var(--ink);
}
.rsvp-btn.on {
  background: rgba(200, 164, 77, 0.15);
  border-color: var(--gold);
  color: var(--gold-bright);
  font-weight: 700;
}
.rsvp-btn.maybe.on {
  background: rgba(247, 226, 173, 0.12);
  border-color: rgba(247, 226, 173, 0.5);
  color: #ffe2a6;
}
.rsvp-btn.no.on {
  background: rgba(224, 122, 122, 0.12);
  border-color: rgba(224, 122, 122, 0.5);
  color: #e07a7a;
}
/* ---------------- item popup ---------------- */
.popup-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
  padding: 1rem;
}
.item-popup {
  background: var(--stone);
  border: 1px solid var(--gold);
  border-radius: 10px;
  padding: 1.1rem 1.3rem;
  width: 100%;
  max-width: 520px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.6);
}
.popup-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.6rem;
  margin-bottom: 0.7rem;
}
.popup-head h3 {
  margin: 0;
  color: var(--gold-bright);
  font-family: Georgia, serif;
}
.popup-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.5rem;
  margin-bottom: 0.8rem;
}
.popup-stat {
  background: var(--stone-light);
  border: 1px solid var(--edge);
  border-top: 2px solid var(--gold);
  border-radius: 6px;
  padding: 0.5rem 0.6rem;
  display: flex;
  flex-direction: column;
  gap: 0.1rem;
}
.popup-stat .dash-num {
  font-size: 1rem;
}
.popup-spark {
  width: 100%;
  height: 60px;
  margin-bottom: 0.6rem;
}
.popup-spark polyline {
  fill: none;
  stroke: var(--gold);
  stroke-width: 2;
  stroke-linejoin: round;
  stroke-linecap: round;
  vector-effect: non-scaling-stroke;
}
/* MunkeHub AI advisor card */
.ah-advisor {
  background: rgba(200, 164, 77, 0.06);
  border: 1px solid rgba(200, 164, 77, 0.35);
  border-radius: 8px;
  padding: 0.9rem 1rem;
  margin-bottom: 1rem;
}
.ah-adv-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.6rem;
}
.ah-adv-head h3 {
  margin: 0;
  font-family: Georgia, serif;
  font-size: 1rem;
  color: var(--gold-bright);
}
.ah-adv-summary {
  margin: 0.5rem 0 0.7rem;
  font-size: 0.88rem;
}
.ah-adv-cols {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.9rem;
}
.ah-adv-col h4 {
  margin: 0 0 0.35rem;
  font-size: 0.78rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--gold);
}
.ah-adv-col ul {
  margin: 0;
  padding-left: 1.1rem;
  font-size: 0.82rem;
}
.ah-adv-col li {
  margin-bottom: 0.3rem;
}
.ah-adv-mount {
  margin: 0.8rem 0 0;
  font-size: 0.85rem;
  color: var(--gold-bright);
}
.ah-adv-profs {
  margin: 0.35rem 0 0;
  font-size: 0.75rem;
  color: var(--muted);
}
.ah-warn {
  cursor: help;
  margin-left: 0.25rem;
}
/* watchlist star */
.star {
  background: none;
  border: none;
  color: var(--muted);
  cursor: pointer;
  font-size: 0.95rem;
  line-height: 1;
  padding: 0 0.15rem 0 0;
  flex: none;
}
.star:hover {
  color: var(--gold);
}
.star.on {
  color: var(--gold-bright);
}
/* 7-day sparkline */
.ah-sparkcell {
  width: 72px;
}
.spark {
  display: block;
  width: 60px;
  height: 18px;
  margin-left: auto;
}
.spark polyline {
  fill: none;
  stroke: var(--gold);
  stroke-width: 1.5;
  stroke-linejoin: round;
  stroke-linecap: round;
  vector-effect: non-scaling-stroke;
}
.ah-dealbadge {
  background: rgba(127, 209, 138, 0.15);
  border: 1px solid rgba(127, 209, 138, 0.5);
  border-radius: 4px;
  color: #7fd18a;
  font-size: 0.7rem;
  font-weight: 700;
  margin-left: 0.4rem;
  padding: 0.05rem 0.3rem;
}
.ah-credit a {
  color: var(--muted);
  text-decoration: underline;
}
.ah-credit a:hover {
  color: var(--gold-bright);
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
.shop-range,
.tier-check span {
  display: inline-block;
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--gold-bright);
  background: rgba(200, 164, 77, 0.12);
  border: 1px solid rgba(200, 164, 77, 0.4);
  border-radius: 4px;
  padding: 0.1rem 0.5rem;
}
.shop-card .shop-range {
  margin-bottom: 0.4rem;
}
/* tier header: checkbox range + tier total */
.shop-cardhead {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 0.4rem;
}
.shop-cardhead .tier-check span {
  margin-bottom: 0;
}
.shop-total {
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--gold-bright);
  white-space: nowrap;
}
.shop-miss {
  color: var(--muted);
  font-weight: 400;
  font-size: 0.7rem;
}
/* live price next to a reagent */
.shop-price {
  font-size: 0.72rem;
  color: var(--muted);
  margin-left: 0.4rem;
  white-space: nowrap;
}
/* skill banner: your level + remaining cost */
.skill-banner {
  display: flex;
  align-items: center;
  gap: 0.7rem;
  flex-wrap: wrap;
  background: rgba(127, 209, 138, 0.07);
  border: 1px solid rgba(127, 209, 138, 0.35);
  border-radius: 8px;
  padding: 0.55rem 0.9rem;
  margin-bottom: 0.8rem;
  font-size: 0.88rem;
}
.sb-label {
  color: var(--muted);
}
.sb-level {
  color: #7fd18a;
  font-family: Georgia, serif;
  font-size: 1.1rem;
}
.sb-bar {
  width: 140px;
}
.sb-cost {
  margin-left: auto;
  color: var(--gold-bright);
  font-weight: 700;
}
/* buy list cart */
.cart-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.6rem;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}
.shop-have-cell {
  color: #7fd18a;
  font-weight: 700;
}
.link-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.78rem;
  padding: 0 0.15rem;
  opacity: 0.6;
  vertical-align: middle;
}
.link-btn:hover {
  opacity: 1;
}
/* "✓ have" — already in our bags */
.shop-have {
  font-size: 0.68rem;
  font-weight: 700;
  color: #7fd18a;
  background: rgba(127, 209, 138, 0.12);
  border: 1px solid rgba(127, 209, 138, 0.4);
  border-radius: 4px;
  padding: 0 0.35rem;
  margin-left: 0.35rem;
  white-space: nowrap;
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

  /* guides: rail becomes a horizontal scroll strip above the content */
  .guides-layout { grid-template-columns: 1fr; gap: 0.9rem; }
  .guides-rail {
    position: static;
    flex-direction: row;
    overflow-x: auto;
    padding-bottom: 0.3rem;
    gap: 0.35rem;
  }
  .rail-item { width: auto; flex: none; }
  .guide-progress { display: none; }
}
</style>