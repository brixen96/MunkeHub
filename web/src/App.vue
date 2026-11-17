<template>
  <div class="page">
    <header>
      <h1>🐒 MunkeHub</h1>
      <p v-if="guild">
        <strong>{{ guild.name }}</strong> — {{ guild.game }} · {{ guild.server }}
      </p>
      <p v-else-if="guildError" class="error">Could not reach the MunkeHub API.</p>
      <p v-else>Loading guild info…</p>
    </header>

    <main>
      <h2>Guild Roster</h2>
      <p v-if="charError" class="error">Could not load the roster.</p>
      <p v-else-if="!characters">Loading…</p>
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
            <td class="name">{{ c.name }} <span class="realm">{{ c.realm }}</span></td>
            <td>{{ c.race }}</td>
            <td :class="'class-' + (c.class || '').toLowerCase()">{{ c.class }}</td>
            <td>{{ c.level }}</td>
            <td>{{ c.ilvl }}</td>
            <td>{{ c.owner }}</td>
            <td>{{ formatDate(c.updated_at) }}</td>
          </tr>
        </tbody>
      </table>

      <h2>Announcements</h2>
      <p v-if="annError" class="error">Could not load announcements.</p>
      <p v-else-if="!announcements">Loading…</p>
      <p v-else-if="announcements.length === 0" class="muted">
        No announcements yet.
      </p>
      <article v-for="a in announcements" :key="a.id" class="card">
        <h3>{{ a.title }}</h3>
        <p>{{ a.body }}</p>
        <time>{{ formatDate(a.created_at) }}</time>
      </article>
    </main>

    <footer>
      <span>MunkeHub — private hub for Munkeklubben</span>
    </footer>
  </div>
</template>

<script>
import { api } from "./api";

export default {
  name: "App",
  data() {
    return {
      guild: null,
      guildError: false,
      announcements: null,
      annError: false,
      characters: null,
      charError: false,
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
  },
};
</script>

<style>
:root {
  color-scheme: dark;
}
* {
  box-sizing: border-box;
}
body {
  margin: 0;
  font-family: system-ui, -apple-system, "Segoe UI", sans-serif;
  background: #14171c;
  color: #e6e6e6;
}
.page {
  max-width: 720px;
  margin: 0 auto;
  padding: 2rem 1rem;
}
header h1 {
  margin: 0 0 0.25rem;
}
header p {
  margin: 0.25rem 0;
  color: #b9c2cc;
}
main h2 {
  margin-top: 2rem;
  border-bottom: 1px solid #2a2f36;
  padding-bottom: 0.5rem;
}
.card {
  background: #1c2128;
  border: 1px solid #2a2f36;
  border-radius: 8px;
  padding: 1rem;
  margin: 0.75rem 0;
}
.card h3 {
  margin: 0 0 0.5rem;
  color: #ffd166;
}
.card time {
  font-size: 0.8rem;
  color: #8a939e;
}
.muted {
  color: #8a939e;
}
.roster {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}
.roster th,
.roster td {
  text-align: left;
  padding: 0.5rem 0.6rem;
  border-bottom: 1px solid #2a2f36;
}
.roster th {
  color: #8a939e;
  font-weight: 600;
  font-size: 0.8rem;
  text-transform: uppercase;
}
.roster .name {
  color: #ffd166;
  font-weight: 600;
}
.roster .realm {
  color: #8a939e;
  font-weight: 400;
  font-size: 0.8rem;
}
.error {
  color: #ff6b6b;
}
footer {
  margin-top: 3rem;
  font-size: 0.8rem;
  color: #6b7480;
  text-align: center;
}
</style>