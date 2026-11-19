-- MunkeHub v0.8.21 — private guild hub addon for Munkeklubben (WoW: Forever, RPPVP)
--
-- Two data paths (WoW addons CANNOT do HTTP):
--   1. OWN stats  -> SavedVariables (syncJSON)  -> helper exe uploads
--   2. GUILD relay -> CHAT_MSG_ADDON broadcast   -> Emil's helper uploads
--      (friends only need this addon — nothing to run)
local ADDON_NAME = ...

local PREFIX = "MUNKEHUB"
local RELAY_MAX = 50              -- max relayed members kept
local BROADCAST_INTERVAL = 300    -- seconds between guild broadcasts

local defaults = {
    apiURL = "",     -- legacy (unused — the helper handles the API address)
    syncJSON = "",   -- own stats payload (read by the helper)
    detailJSON = "", -- deep detail: spec/zone/gear (helper path ONLY — too big for relay)
    relayJSON = "",  -- relayed guild stats array (read by the helper)
}

local function Print(msg)
    print("|cff00ff96MunkeHub|r " .. tostring(msg))
end

-- ---------------------------------------------------------------------------
-- Stats collection — stored in SavedVariables, picked up by the helper exe
-- (WoW addons cannot do HTTP, so the helper bridges file -> API)
-- ---------------------------------------------------------------------------
-- WoW: Forever has first + last names. UnitNameUnmodified returns them as
-- two values when surnames are separate (see LibLatency/LibDurability).
-- `name` stays the row key (first name only) so UNIQUE(name, realm) never
-- changes — last_name is display-only.
local function LastName()
    if UnitNameUnmodified then
        local _, last = UnitNameUnmodified("player")
        if type(last) == "string" and last ~= "" then return last end
    end
    return ""
end

-- Blizzard character ID from the GUID (Player-6951-00EB3BEC -> 15416300).
-- The website builds the 3D model render URL from it.
local function BlizzardID()
    if not UnitGUID then return 0 end
    local guid = UnitGUID("player")
    if type(guid) ~= "string" then return 0 end
    local hex = guid:match("-([0-9A-Fa-f]+)$")
    if not hex then return 0 end
    return tonumber(hex, 16) or 0
end

-- Professions: primary 2 + cooking + fishing, comma-separated names.
-- GetProfessions() -> prof1, prof2, archaeology, fishing, cooking;
-- GetProfessionInfo(index) -> name, icon, skillLevel, ... (verified in
-- RSProfessionDB.lua / QuestieTrace — bare globals exist in Forever).
-- Feeds the guild profession list + the MunkeHub AI auction advisor.
local function BuildProfessions()
    if not GetProfessions then return "" end
    local p1, p2, _, fishing, cooking = GetProfessions()
    -- NOTE: build the index list manually — ipairs() stops at the first nil,
    -- so a table constructor with holes would drop later professions.
    local idxs = {}
    if p1 then idxs[#idxs + 1] = p1 end
    if p2 then idxs[#idxs + 1] = p2 end
    if cooking then idxs[#idxs + 1] = cooking end
    if fishing then idxs[#idxs + 1] = fishing end
    local out = {}
    for _, idx in ipairs(idxs) do
        if GetProfessionInfo then
            local name = GetProfessionInfo(idx)
            if type(name) == "string" and name ~= "" then
                out[#out + 1] = name
            end
        end
    end
    return table.concat(out, ",")
end

-- extras=true adds professions + gold (helper path always wants them).
-- The relay path passes extras=false when the message would exceed the
-- 240-byte addon limit — losing professions beats losing the whole relay.
local function BuildSyncJSON(extras)
    local name = UnitName("player") or ""
    local lastName = LastName()
    local realm = GetRealmName() or ""
    local bid = BlizzardID()
    local race = UnitRace("player") or ""
    local class = select(2, UnitClass("player")) or ""
    local level = UnitLevel("player") or 0
    local gold = GetMoney and GetMoney() or 0
    local ilvl = 0
    if GetAverageItemLevel then
        local overall = GetAverageItemLevel()
        if type(overall) == "number" then ilvl = math.floor(overall * 10) / 10 end
    end
    -- Manual JSON: values are game-provided strings, escape quotes/backslashes
    local function esc(s)
        return (tostring(s):gsub('[\\"]', function(c) return "\\" .. c end))
    end
    if extras then
        return string.format(
            '{"name":"%s","last_name":"%s","realm":"%s","race":"%s","class":"%s","level":%d,"ilvl":%s,"bid":%d,"professions":"%s","gold":%d}',
            esc(name), esc(lastName), esc(realm), esc(race), esc(class), level, tostring(ilvl), bid,
            esc(BuildProfessions()), gold)
    end
    return string.format(
        '{"name":"%s","last_name":"%s","realm":"%s","race":"%s","class":"%s","level":%d,"ilvl":%s,"bid":%d}',
        esc(name), esc(lastName), esc(realm), esc(race), esc(class), level, tostring(ilvl), bid)
end

-- Deep detail for the expandable roster rows. NEVER relayed (240-byte cap) —
-- only the helper uploads this via SavedVariables.
local GEAR_SLOTS = {
    [1] = "Head", [2] = "Neck", [3] = "Shoulder", [4] = "Shirt",
    [5] = "Chest", [6] = "Waist", [7] = "Legs", [8] = "Feet",
    [9] = "Wrist", [10] = "Hands", [11] = "Ring1", [12] = "Ring2",
    [13] = "Trinket1", [14] = "Trinket2", [15] = "Back",
    [16] = "MainHand", [17] = "OffHand", [18] = "Ranged", [19] = "Tabard",
}

-- WoW: Forever uses the modern API namespaces (C_Item, C_SpecializationInfo).
-- The bare globals (GetItemInfoInstant, GetInventoryItemLevel, GetPlayedTime,
-- GetInventoryItemQuality) DO NOT EXIST here — guarding on them silently
-- produced empty icons/ilvl/spec/played for every character.
local function GetSpecIndex()
    if C_SpecializationInfo and C_SpecializationInfo.GetSpecialization then
        return C_SpecializationInfo.GetSpecialization()
    end
    if GetSpecialization then return GetSpecialization() end
    return nil
end

local function GetSpecName(idx)
    if not idx then return "" end
    if C_SpecializationInfo and C_SpecializationInfo.GetSpecializationInfo then
        local _, name = C_SpecializationInfo.GetSpecializationInfo(idx)
        if name then return name end
    end
    if GetSpecializationInfo then
        local _, name = GetSpecializationInfo(idx)
        if name then return name end
    end
    return ""
end

-- Spec icon as a texture FILEID (4th return) — foreverdb serves it.
local function GetSpecIcon(idx)
    if not idx then return 0 end
    local fn
    if C_SpecializationInfo and C_SpecializationInfo.GetSpecializationInfo then
        fn = C_SpecializationInfo.GetSpecializationInfo
    elseif GetSpecializationInfo then
        fn = GetSpecializationInfo
    end
    if not fn then return 0 end
    local _, _, _, icon = fn(idx)
    return tonumber(icon) or 0
end

local function BuildDetailJSON()
    local function esc(s)
        return (tostring(s):gsub('[\\"]', function(c) return "\\" .. c end))
    end

    -- specialization (+ its icon fileID for the website avatar)
    local specIdx = GetSpecIndex()
    local spec = GetSpecName(specIdx)
    local specIcon = GetSpecIcon(specIdx)

    -- location
    local zone = GetZoneText() or ""

    -- played time (seconds) — async: RequestTimePlayed() fires TIME_PLAYED_MSG,
    -- the handler stores it in MunkeHubDB.played and rebuilds detailJSON.
    local played = tonumber(MunkeHubDB and MunkeHubDB.played) or 0

    -- gear: slot, name, ilvl, quality, icon
    local gear = {}
    for slot, slotName in pairs(GEAR_SLOTS) do
        local link = GetInventoryItemLink("player", slot)
        if link then
            local name = link:match("|h%[(.-)%]|h") or ""
            local itemID = tonumber(link:match("item:(%d+)")) or 0

            local ilvl = 0
            if C_Item and C_Item.GetDetailedItemLevelInfo then
                ilvl = C_Item.GetDetailedItemLevelInfo(link) or 0
            end

            local q = 0
            if C_Item and C_Item.GetItemQualityByID then
                q = C_Item.GetItemQualityByID(itemID) or 0
            end

            -- The game only exposes icon FILEIDS (numbers like 134592),
            -- never texture names. foreverdb.net serves icons keyed by that
            -- same fileID, so we send it; itemID stays as a fallback source.
            -- NOTE: select(5, ...) returns ALL remaining values, and
            -- tonumber(v, base) would take classID as the base -> nil.
            -- Truncate to one value first (local assignment or extra parens).
            local fid = 0
            if C_Item and C_Item.GetItemInfoInstant then
                local t = select(5, C_Item.GetItemInfoInstant(link))
                fid = tonumber(t) or 0
            elseif GetItemInfoInstant then
                local t = select(5, GetItemInfoInstant(link))
                fid = tonumber(t) or 0
            end
            gear[#gear + 1] = string.format(
                '{"slot":"%s","name":"%s","ilvl":%d,"q":%d,"id":%d,"fid":%d}',
                esc(slotName), esc(name), ilvl, q, itemID, fid)
        end
    end

    local name = UnitName("player") or ""
    local realm = GetRealmName() or ""
    return string.format(
        '{"name":"%s","realm":"%s","spec":"%s","spec_icon":%d,"zone":"%s","played":%d,"gear":[%s]}',
        esc(name), esc(realm), esc(spec), specIcon, esc(zone), played, table.concat(gear, ","))
end

local function SyncStats()
    MunkeHubDB.syncJSON = BuildSyncJSON(true)
    MunkeHubDB.detailJSON = BuildDetailJSON()
    if RequestTimePlayed then RequestTimePlayed() end -- fills played via TIME_PLAYED_MSG
    Print("stats saved — the helper uploads them (or /reload to flush now).")
end

-- ---------------------------------------------------------------------------
-- Guild relay: broadcast our stats, collect everyone else's
-- ---------------------------------------------------------------------------
local function BroadcastStats()
    if not C_ChatInfo or not IsInGuild() then return end
    -- full payload first (professions + gold), fall back to the plain one
    -- if the message would exceed the WoW addon limit
    local msg = BuildSyncJSON(true)
    if #msg + #PREFIX > 240 then msg = BuildSyncJSON(false) end
    if #msg + #PREFIX > 240 then return end -- WoW addon message limit
    C_ChatInfo.SendAddonMessage(PREFIX, msg, "GUILD")
end

local function RebuildRelayJSON()
    local map = MunkeHubDB.relayMap
    if not map then
        MunkeHubDB.relayJSON = ""
        return
    end
    local parts = {}
    for _, v in pairs(map) do
        parts[#parts + 1] = v
    end
    MunkeHubDB.relayJSON = (#parts > 0) and ("[" .. table.concat(parts, ",") .. "]") or ""
end

local function ReceiveRelay(prefix, message)
    if prefix ~= PREFIX then return end
    -- keyed by character name (owner grouping removed in v0.8.0)
    local key = message:match('"name":"([^"]*)"')
    if not key or key == "" then return end

    MunkeHubDB.relayMap = MunkeHubDB.relayMap or {}
    MunkeHubDB.relayMap[key] = message

    -- cap the map (drop one arbitrary entry when over)
    local count = 0
    for _ in pairs(MunkeHubDB.relayMap) do count = count + 1 end
    if count > RELAY_MAX then
        for k in pairs(MunkeHubDB.relayMap) do
            MunkeHubDB.relayMap[k] = nil
            break
        end
    end
    RebuildRelayJSON()
end

-- ---------------------------------------------------------------------------
-- Core frame: load SavedVariables, then we're ready
-- ---------------------------------------------------------------------------
local core = CreateFrame("Frame", "MunkeHubFrame", UIParent)
local tickerStarted = false
core:RegisterEvent("ADDON_LOADED")
core:RegisterEvent("PLAYER_ENTERING_WORLD")
core:RegisterEvent("CHAT_MSG_ADDON")
core:RegisterEvent("TIME_PLAYED_MSG")
core:SetScript("OnEvent", function(self, event, ...)
    if event == "CHAT_MSG_ADDON" then
        local prefix, message = ...
        ReceiveRelay(prefix, message)
        return
    end

    if event == "TIME_PLAYED_MSG" then
        -- async played time: (totalSeconds, levelSeconds) — store + rebuild
        local total = ...
        if MunkeHubDB and type(total) == "number" then
            MunkeHubDB.played = math.floor(total)
            MunkeHubDB.detailJSON = BuildDetailJSON()
        end
        return
    end

    if event == "PLAYER_ENTERING_WORLD" then
        -- silently refresh stats; WoW writes them to disk on /reload or logout
        if MunkeHubDB then
            MunkeHubDB.syncJSON = BuildSyncJSON(true)
            MunkeHubDB.detailJSON = BuildDetailJSON()
        end
        -- played time is async: ask for it, TIME_PLAYED_MSG fills it in
        if RequestTimePlayed then RequestTimePlayed() end
        -- announce ourselves to the guild shortly after loading, then periodically
        C_Timer.After(5, BroadcastStats)
        if not tickerStarted then
            tickerStarted = true
            C_Timer.NewTicker(BROADCAST_INTERVAL, BroadcastStats)
        end
        return
    end

    if event ~= "ADDON_LOADED" or ... ~= ADDON_NAME then return end

    MunkeHubDB = MunkeHubDB or {}
    for k, v in pairs(defaults) do
        if MunkeHubDB[k] == nil then MunkeHubDB[k] = v end
    end
    if C_ChatInfo then
        C_ChatInfo.RegisterAddonMessagePrefix(PREFIX)
    end

    self:UnregisterEvent("ADDON_LOADED")
    Print("loaded — type /munkehub for commands.")
end)

-- ---------------------------------------------------------------------------
-- Slash commands
-- ---------------------------------------------------------------------------
local function ShowHelp()
    Print("commands:")
    print("  |cffffffff/munkehub status|r            show sync + relay status")
    print("  |cffffffff/munkehub sync|r              save character stats now")
end

SLASH_MUNKEHUB1 = "/munkehub"
SLASH_MUNKEHUB2 = "/mh"
SlashCmdList["MUNKEHUB"] = function(msg)
    local cmd, rest = strmatch(msg or "", "^%s*(%S+)%s*(.-)%s*$")
    cmd = cmd and strlower(cmd)

    if cmd == "sync" then
        SyncStats()
    elseif cmd == "status" then
        local payload = MunkeHubDB and MunkeHubDB.syncJSON or ""
        local relay = MunkeHubDB and MunkeHubDB.relayJSON or ""
        if payload ~= "" then
            Print("own stats: |cffffffffsaved|r — helper uploads. /reload writes to disk.")
        else
            Print("own stats: |cffff4444not saved yet|r — type /mh sync")
        end
        local relayCount = 0
        for _ in gmatch(relay, '"owner"') do relayCount = relayCount + 1 end
        Print("relayed members: |cffffffff" .. relayCount .. "|r")
    else
        ShowHelp()
    end
end