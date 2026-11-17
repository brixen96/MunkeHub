-- MunkeHub v0.1.0 — private guild hub addon for Munkeklubben (WoW: Forever, RPPVP)
local ADDON_NAME = ...

local defaults = {
    apiURL = "", -- set via: /munkehub url http://host:port
}

local function Print(msg)
    print("|cff00ff96MunkeHub|r " .. tostring(msg))
end

-- ---------------------------------------------------------------------------
-- Core frame: load SavedVariables, then we're ready
-- ---------------------------------------------------------------------------
local core = CreateFrame("Frame", "MunkeHubFrame", UIParent)
core:RegisterEvent("ADDON_LOADED")
core:SetScript("OnEvent", function(self, event, ...)
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
    print("  |cffffffff/munkehub url <address>|r     set API address (http://host:port)")
    print("  |cffffffff/munkehub sync|r              fetch guild data from MunkeHub")
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
        if MunkeHubDB.apiURL == "" then
            Print("no API set — run /munkehub url <address> first")
        else
            Print("sync not implemented yet — API is coming")
        end
    elseif cmd == "status" then
        local url = MunkeHubDB.apiURL
        Print("API: " .. (url ~= "" and url or "|cffff4444(not set)|r"))
    else
        ShowHelp()
    end
end