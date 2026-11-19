package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// --- Leveling route (Tier 2 of the Quest Helper) ---------------------------
//
// FACTS ONLY, scraped 2026-11-19 from the sources in guides.md:
//   - conquestcapped.com/guides/wow-forever/wow-forever-leveling/
//     (zone route tables: 1-30 per faction, 30-60 shared, live-tested)
//   - expcarry.com/wow-forever-leveling-guide-horde
//     (Horde route by level bracket + race-specific starts)
// No prose copied — just level brackets, zone names, dungeon names.
//
// The quest advisor gets the brackets matching the character's level so it
// can answer "my log is empty — where do I go?" with real route data.

type levelingBracket struct {
	Min      int      `json:"min"`
	Max      int      `json:"max"`
	Zones    []string `json:"zones"`
	Dungeons []string `json:"dungeons,omitempty"`
	Note     string   `json:"note,omitempty"`
}

// hordeBranch: race-specific route for levels 1-30 (expcarry).
func hordeBranch(race string) []levelingBracket {
	switch strings.ToLower(race) {
	case "undead":
		// Forsaken stay on the Eastern Kingdoms — no early Kalimdor trip.
		return []levelingBracket{
			{Min: 1, Max: 10, Zones: []string{"Tirisfal Glades"}},
			{Min: 10, Max: 15, Zones: []string{"Tirisfal Glades", "Silverpine Forest"}},
			{Min: 15, Max: 20, Zones: []string{"Silverpine Forest"}, Dungeons: []string{"Ruins of Lordaeron 15-20"}, Note: "Best new early Horde dungeon — right above the Undercity."},
			{Min: 20, Max: 30, Zones: []string{"Hillsbrad Foothills", "Arathi Highlands"}, Dungeons: []string{"City of Dalaran 28-33"}},
		}
	case "tauren":
		return []levelingBracket{
			{Min: 1, Max: 10, Zones: []string{"Mulgore"}},
			{Min: 10, Max: 20, Zones: []string{"The Barrens"}, Dungeons: []string{"Wailing Caverns 17-24"}},
			{Min: 20, Max: 30, Zones: []string{"Stonetalon Mountains", "Thousand Needles", "The Barrens"}},
		}
	default: // Orc, Troll, Skyborne and anything unknown
		return []levelingBracket{
			{Min: 1, Max: 10, Zones: []string{"Durotar"}},
			{Min: 10, Max: 15, Zones: []string{"The Barrens"}, Dungeons: []string{"Ragefire Chasm 13-18"}},
			{Min: 15, Max: 25, Zones: []string{"The Barrens", "Stonetalon Mountains", "Ashenvale"}, Dungeons: []string{"Wailing Caverns 17-24"}},
			{Min: 25, Max: 30, Zones: []string{"Thousand Needles", "Ashenvale", "Hillsbrad Foothills"}},
		}
	}
}

// sharedRoute: both factions from 30 onward (conquestcapped, live-tested).
// Brackets overlap on purpose — pick whichever zones you're near.
func sharedRoute() []levelingBracket {
	return []levelingBracket{
		{Min: 30, Max: 40, Zones: []string{"Arathi Highlands", "Desolace", "Stranglethorn Vale", "Alterac Mountains", "Thousand Needles"},
			Dungeons: []string{"Scarlet Monastery", "Razorfen Kraul", "Gnomeregan", "City of Dalaran 28-33", "The Drowned City 35-40"}},
		{Min: 35, Max: 45, Zones: []string{"Riverglades (36-44)", "Shen'dralas (35-45)", "Dustwallow Marsh", "Badlands", "Swamp of Sorrows"},
			Dungeons: []string{"Krol'dok Stronghold 40-45", "Razorfen Downs", "Uldaman"},
			Note: "Riverglades and Shen'dralas are new Forever zones — fresh quest pools."},
		{Min: 40, Max: 50, Zones: []string{"Tanaris", "Feralas", "The Hinterlands", "Searing Gorge"},
			Dungeons: []string{"Zul'Farrak", "Maraudon", "Alcaz Prison 48-53"}},
		{Min: 48, Max: 55, Zones: []string{"Un'Goro Crater", "Felwood", "Azshara", "Blasted Lands"},
			Dungeons: []string{"The Temple of Atal'Hakkar", "Blackrock Depths"}},
		{Min: 50, Max: 60, Zones: []string{"Burning Steppes", "Western Plaguelands", "Eastern Plaguelands", "Winterspring", "Silithus"},
			Dungeons: []string{"Blackrock Spire", "Blackmaw Hold 55-60", "Shaper's Terrace 58-60"}},
		{Min: 54, Max: 60, Zones: []string{"Mount Hyjal"},
			Dungeons: []string{"Dire Maul", "Scholomance", "Stratholme"},
			Note: "Last stop before 60 — most quests need level 54+."},
	}
}

// routeFor returns every bracket overlapping the character's level:
// the race branch (below 30) plus the shared route.
func routeFor(level int, race string) []levelingBracket {
	out := []levelingBracket{}
	if level < 30 {
		for _, b := range hordeBranch(race) {
			if b.Max >= level {
				out = append(out, b)
			}
		}
	}
	for _, b := range sharedRoute() {
		if b.Max >= level && b.Min <= level+5 { // show current + next bracket
			out = append(out, b)
		}
	}
	return out
}

// levelingHandler serves GET /api/v1/leveling?level=37&race=Undead.
func levelingHandler(c *gin.Context) {
	level := atoiDefault(c.Query("level"), 0)
	race := c.Query("race")
	if level <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "level required"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"level":    level,
		"race":     race,
		"brackets": routeFor(level, race),
	})
}