# MunkeHub — What we can fetch from the game

Overview of everything the addon can read from WoW: Forever and push to the
website. Written 2026-11-17, based on the standard WoW addon API — some APIs
may differ in Forever (custom Classic fork), so verify in-game before
building on them (`/dump FunctionName`).

**Hard constraints (always):**
- ❌ **No HTTP from the addon** — everything goes through SavedVariables → helper → API
- ⏱️ Data only exists while the game runs; SavedVariables hit disk on `/reload` or logout
- 📡 Guild relay messages are small (~240 bytes each) and throttled — big data must go through the helper, not `CHAT_MSG_ADDON`
- 👁️ Only what the addon can see on *your* client (your character, guild roster, group members)

---

## ✅ Already used (v0.5+)

| Data | API | Used for |
|------|-----|----------|
| Name, realm, race, class, level | `UnitName`, `GetRealmName`, `UnitRace`, `UnitClass`, `UnitLevel` | Roster |
| Item level | `GetAverageItemLevel` | Roster |
| Owner label | `/mh label` (SavedVariables) | Grouping characters per member |

---

## 🎮 Character core

| Data | API (standard WoW) | Website idea |
|------|--------------------|--------------|
| Specialization / spec | `GetSpecialization`, `GetSpecializationInfo` | Show spec in roster (e.g. "Brix — Shadow") |
| Talents | `C_SpecializationInfo.GetConfigID`, talent tree APIs | Talent build pages, compare builds |
| Played time | `GetPlayedTime` | "Veteran" badges, activity stats |
| Rest/XP state | `GetRestState`, XP to next level | Leveling progress bars for alts |
| Titles | `GetNumTitles`, `GetTitleName` | RP flair on profiles |
| Current location | `GetZoneText`, `GetSubZoneText` | "Where is everyone now" map |
| Coordinates | `GetPlayerMapPosition` | Live map (needs frequent sync) |

## ⚔️ Gear & equipment

| Data | API | Website idea |
|------|-----|--------------|
| Equipped items (slot, name, ilvl, quality) | `GetInventoryItemLink`, `GetInventoryItemID` | Full gear page per character, slot-by-slot |
| Item stats on gear | `GetItemStats` | Gear analysis, best-in-slot checks |
| Equipment manager sets | `C_EquipmentSet` | Show saved sets (PvP/raid) |
| Transmog appearance | `C_TransmogCollection` | Cosmetic showcase (big data — helper only) |
| Bags / bank contents | `GetContainerItemInfo` | Guild bank overview, alt gear lookup |

## 📈 Progression

| Data | API | Website idea |
|------|-----|--------------|
| Raid boss kills | `GetAchievementInfo` (raid achievements) / `GetLootSpecialization` | Progress tracker: "Molten Core 8/10" |
| Mythic+ rating & best runs | `C_ChallengeMode` | M+ ladder, weekly bests |
| PvP rating (2v2/3v3/RBG) | `GetPersonalRatedInfo` | PvP ladder |
| Honor level | `GetHonorLevel` | PvP showcase |
| Quest log | `GetNumQuestLogEntries` | Alt questing progress |
| Achievements (points, recent) | `GetTotalAchievementPoints`, `GetAchievementInfo` | Achievement showcase, "rare" badges |

## 🏦 Economy

| Data | API | Website idea |
|------|-----|--------------|
| Gold on character | `GetMoney` | Wealth ranking (⚠️ privacy — maybe officer-only) |
| Currencies (badges, tokens) | `GetCurrencyInfo` | "What do I still need" tracker |
| Professions + skill level | `GetProfessions`, `GetProfessionInfo` | Crafting directory: "who can make X?" |
| Auction house (own listings) | `C_AuctionHouse` | Guild crafting orders (limited) |

## 🏰 Guild data (read anytime — no group needed)

| Data | API | Website idea |
|------|-----|--------------|
| Guild roster: name, class, rank, note, officer note, zone, **online status** | `GetGuildRosterInfo` | **Live roster with online dots** — big win, works even when members don't have the addon |
| Guild rank permissions | `GetGuildRankInfo` | Rank display |
| Guild news/events feed | `GetGuildEventInfo` | Auto news: "X reached level 60" |
| Guild bank tabs | `GetGuildBankInfo` (if permitted) | Bank overview |
| Guild MOTD | `GetGuildRosterMOTD` | Show on website header |
| Guild bank money | `GetGuildBankMoney` | Guild funds display |

> 💡 **Big insight:** the guild roster API gives us **everyone's online status,
> level, class and rank** from *your* client alone — no addon needed on their
> side. That's a live member list for free (synced while you're online).

## 👥 Social / group

| Data | API | Website idea |
|------|-----|--------------|
| Group/raid members + roles | `GetNumGroupMembers`, `UnitGroupRolesAssigned` | "Who's in the raid right now" |
| Friends list | `GetNumFriends`, `GetFriendInfo` | Alt tracking (if friends are alts) |
| Chat messages (guild/officer) | `CHAT_MSG_GUILD` etc. | ⚠️ RP server — chat logs are sensitive; only with clear consent, maybe just keyword highlights ("raid tonight") |

## 📡 Communication (addon messages — already used for relay)

| Channel | Use |
|---------|-----|
| `CHAT_MSG_ADDON` guild broadcast | ✅ Stats relay (v0.5) |
| Addon message channels | Future: live sign-ups ("I'm in for raid"), polls, loot council votes — all in-game, no website needed |
| Addon message limits | ~240 bytes/message, throttled — send small events, not bulk data |

---

## 💡 My ideas beyond raw data (features)

1. **Live online roster** — guild roster API → green dots on the website (works with just you online)
2. **Event calendar + sign-ups** — raid nights, who's coming; sign-ups via addon message or website
3. **Progress tracker** — boss kills per raid tier, guild bests
4. **Crafting directory** — "who has Alchemy 300?" from professions data
5. **Alt tracker** — one member, many characters (already supported by `owner` label)
6. **Loot council / DKP** — loot history + points (needs officer workflow design)
7. **RP profiles** — RPPVP server: character bios, RP stories, titles (fits the server culture)
8. **Attendance log** — who showed up to raids (parse guild events or manual officer marks)
9. **Weekly digest** — auto-announcement: levels gained, new members, raid results
10. **In-game sign-up slash command** — `/mh signup friday` → broadcasts to guild → website shows the list

## ⛔ Not possible (don't plan on these)

| Idea | Why not |
|------|---------|
| Addon does HTTP directly | WoW blocks all network calls from Lua |
| Data while game is closed | Addon only runs inside the client |
| Other players' private data | Only what the API exposes to your client |
| Real-time push to website | Chain is: game → SavedVariables (on reload/logout) → helper → API |
| Reading *their* SavedVariables | Each player's file is on their own PC — only their helper can read it |

---

## Suggested next steps (by value/effort)

| Feature | Value | Effort | Note |
|---------|-------|--------|------|
| Live online roster (guild API) | 🔥 high | low | Works with only Emil online |
| Spec + gear in roster | 🔥 high | low | One more field in syncJSON |
| Event sign-ups | high | medium | Addon command + API table + website section |
| Progress tracker | medium | medium | Needs achievement parsing |
| Crafting directory | medium | low | Professions field in syncJSON |
| RP profiles | high (RP server!) | medium | Text fields, officer-editable |