-- MunkeHub v0.1.0 — private guild hub addon for Munkeklubben (WoW: Forever, RPPVP)
local ADDON_NAME = ...

local defaults = {
    apiURL = "",   -- set via: /munkehub url http://host:port
    syncJSON = "", -- stats payload read by the MunkeHub helper (SavedVariables)
}

local function Print(msg)
    print("|cff00ff96MunkeHub|r " .. tostring(msg))
end

-- ---------------------------------------------------------------------------
-- Stats collection — stored in SavedVariables, picked up by the helper exe
-- (WoW addons cannot do HTTP, so the helper bridges file -> API)
-- ---------------------------------------------------------------------------
local function BuildSyncJSON()
    local name = UnitName("player") or ""
    local realm = GetRealmName() or ""
    local race = UnitRace("player") or ""
    local class = select(2, UnitClass("player")) or ""
    local level = UnitLevel("player") or 0
    local ilvl = 0
    if GetAverageItemLevel then
        local overall = GetAverageItemLevel()
        if type(overall) == "number" then ilvl = math.floor(overall * 10) / 10 end
    end
    -- Manual JSON: values are game-provided strings, escape quotes/backslashes
    local function esc(s)
        return (tostring(s):gsub('[\\"]', function(c) return "\\" .. c end))
    end
    return string.format(
        '{"name":"%s","realm":"%s","race":"%s","class":"%s","level":%d,"ilvl":%s}',
        esc(name), esc(realm), esc(race), esc(class), level, tostring(ilvl))
end

local function SyncStats()
    MunkeHubDB.syncJSON = BuildSyncJSON()
    Print("stats saved — the MunkeHub helper will upload them (or /reload to flush now).")
end

-- ---------------------------------------------------------------------------
-- Core frame: load SavedVariables, then we're ready
-- ---------------------------------------------------------------------------
local core = CreateFrame("Frame", "MunkeHubFrame", UIParent)
core:RegisterEvent("ADDON_LOADED")
core:RegisterEvent("PLAYER_ENTERING_WORLD")
core:SetScript("OnEvent", function(self, event, ...)
    if event == "PLAYER_ENTERING_WORLD" then
        -- silently refresh stats; WoW writes them to disk on /reload or logout
        if MunkeHubDB then MunkeHubDB.syncJSON = BuildSyncJSON() end
        return
    end
    if event ~= "ADDON_LOADED" or ... ~= ADDON_NAME then return end

    MunkeHubDB = MunkeHubDB or {}
    for k, v in pairs(defaults) do
        if MunkeHubDB[k] == nil then MunkeHubDB[k] = v end
    end

    self:UnregisterEvent("ADDON_LOADED")
    Print("loaded — type /munkehub for commands.")
end)

-- ---------------------------------------------------------------------------
-- Slash commands
-- ---------------------------------------------------------------------------
local function ShowHelp()
    Print("commands:")
    print("  |cffffffff/munkehub status|r            show current config")
    print("  |cffffffff/munkehub sync|r              save character stats for the helper to upload")
    print("  |cffffffff/munkehub url <address>|r     set API address (http://host:port)")
end

SLASH_MUNKEHUB1 = "/munkehub"
SLASH_MUNKEHUB2 = "/mh"
SlashCmdList["MUNKEHUB"] = function(msg)
    local cmd, rest = strmatch(msg or "", "^%s*(%S+)%s*(.-)%s*$")
    cmd = cmd and strlower(cmd)

    if cmd == "url" and rest ~= "" then
        MunkeHubDB.apiURL = rest
        Print("API set to " .. rest)
    elseif cmd == "sync" then
        SyncStats()
    elseif cmd == "status" then
        local url = MunkeHubDB.apiURL
        Print("API: " .. (url ~= "" and url or "|cffff4444(not set)|r"))
    else
        ShowHelp()
    end
end