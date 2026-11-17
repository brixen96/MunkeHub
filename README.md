# MunkeHub

Private hub for the **Munkeklubben** guild — World of Warcraft: Forever (RPPVP server).
Only for Emil and friends. Not public, not commercial.

## Parts

| Part | Folder | Stack |
|------|--------|-------|
| API | `api/` | Go, Gin, SQLite, SSE |
| Website | `web/` | Vue 3, Docker, self-hosted + port-forwarded |
| WoW addon | `addon/MunkeHub/` | Lua, Interface `16001` |

## Status

- [x] Repo scaffold
- [x] Addon skeleton (loads in game, slash commands)
- [ ] API (health endpoint first)
- [ ] Website
- [ ] Addon ↔ API sync

## Addon development

Source of truth is `addon/MunkeHub/`. Sync to the game:

```powershell
./scripts/sync-addon.ps1
```

Then `/reload` in game (or restart the client).
