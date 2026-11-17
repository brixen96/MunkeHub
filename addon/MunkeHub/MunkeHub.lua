-- MunkeHub v0.5.0 — private guild hub addon for Munkeklubben (WoW: Forever, RPPVP)
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
    label = "",      -- roster owner label, set via: /mh label sven
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
    local owner = (MunkeHubDB and MunkeHubDB.label) or ""
    if owner == "" then owner = name end
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
        '{"owner":"%s","name":"%s","realm":"%s","race":"%s","class":"%s","level":%d,"ilvl":%s}',
        esc(owner), esc(name), esc(realm), esc(race), esc(class), level, tostring(ilvl))
end

-- Deep detail for the expandable roster rows. NEVER relayed (240-byte cap) —
-- only the helper uploads this via SavedVariables.
local GEAR_SLOTS = {
    [1] = "Head", [2] = "Neck", [3] = "Shoulder", [5] = "Chest",
    [6] = "Waist", [7] = "Legs", [8] = "Feet", [9] = "Wrist",
    [10] = "Hands", [15] = "Back", [16] = "MainHand", [17] = "OffHand",
}

local function BuildDetailJSON()
    local function esc(s)
        return (tostring(s):gsub('[\\"]', function(c) return "\\" .. c end))
    end

    -- specialization
    local spec = ""
    if GetSpecialization then
        local idx = GetSpecialization()
        if idx then
            local _, name = GetSpecializationInfo(idx)
            spec = name or ""
        end
    end

    -- location
    local zone = GetZoneText() or ""

    -- played time (seconds)
    local played = 0
    if GetPlayedTime then played = GetPlayedTime() or 0 end

    -- gear: slot, name, ilvl, quality
    local gear = {}
    for slot, slotName in pairs(GEAR_SLOTS) do
        local link = GetInventoryItemLink("player", slot)
        if link then
            local name = link:match("|h%[(.-)%]|h") or ""
            local ilvl = 0
            if GetInventoryItemLevel then
                ilvl = GetInventoryItemLevel("player", slot) or 0
            end
            local q = 0
            if GetInventoryItemQuality then
                q = GetInventoryItemQuality("player", slot) or 0
            end
            gear[#gear + 1] = string.format(
                '{"slot":"%s","name":"%s","ilvl":%d,"q":%d}',
                esc(slotName), esc(name), ilvl, q)
        end
    end

    return string.format(
        '{"spec":"%s","zone":"%s","played":%d,"gear":[%s]}',
        esc(spec), esc(zone), played, table.concat(gear, ","))
end

local function SyncStats()
    MunkeHubDB.syncJSON = BuildSyncJSON()
    MunkeHubDB.detailJSON = BuildDetailJSON()
    Print("stats saved — the helper uploads them (or /reload to flush now).")
end

-- ---------------------------------------------------------------------------
-- Guild relay: broadcast our stats, collect everyone else's
-- ---------------------------------------------------------------------------
local function BroadcastStats()
    if not C_ChatInfo or not IsInGuild() then return end
    local msg = BuildSyncJSON()
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
    local owner = message:match('"owner":"([^"]*)"')
    if not owner or owner == "" then return end

    MunkeHubDB.relayMap = MunkeHubDB.relayMap or {}
    MunkeHubDB.relayMap[owner] = message

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
core:SetScript("OnEvent", function(self, event, ...)
    if event == "CHAT_MSG_ADDON" then
        local prefix, message = ...
        ReceiveRelay(prefix, message)
        return
    end

    if event == "PLAYER_ENTERING_WORLD" then
        -- silently refresh stats; WoW writes them to disk on /reload or logout
        if MunkeHubDB then
            MunkeHubDB.syncJSON = BuildSyncJSON()
            MunkeHubDB.detailJSON = BuildDetailJSON()
        end
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
    print("  |cffffffff/munkehub label <name>|r      set your roster label (e.g. /mh label sven)")
end

SLASH_MUNKEHUB1 = "/munkehub"
SLASH_MUNKEHUB2 = "/mh"
SlashCmdList["MUNKEHUB"] = function(msg)
    local cmd, rest = strmatch(msg or "", "^%s*(%S+)%s*(.-)%s*$")
    cmd = cmd and strlower(cmd)

    if cmd == "label" and rest ~= "" then
        MunkeHubDB.label = rest
        MunkeHubDB.syncJSON = BuildSyncJSON()
        Print("label set to |cffffffff" .. rest .. "|r — all your characters group under it.")
        Print("/reload or logout to save it to disk.")
    elseif cmd == "sync" then
        SyncStats()
    elseif cmd == "status" then
        local label = MunkeHubDB and MunkeHubDB.label or ""
        local payload = MunkeHubDB and MunkeHubDB.syncJSON or ""
        local relay = MunkeHubDB and MunkeHubDB.relayJSON or ""
        Print("label: " .. (label ~= "" and ("|cffffffff" .. label .. "|r")
            or "|cffff4444not set|r — /mh label <name>"))
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