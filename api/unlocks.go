package main

import "strings"

// --- Leveling unlocks (Tier 4 of the Quest Helper) -------------------------
//
// FACTS ONLY, scraped 2026-11-20 from the method.gg class guides listed in
// guides.md — their "Key Spells & Leveling Milestones" tables. No prose
// copied — just level, spell name and a short factual note.
//
// The Quests tab shows a "📜 Coming up" card: what the character unlocks in
// the next ~10 levels, so nobody misses a big power spike (pet at 10,
// Shadowform at 40, Windfury at 30...).

type unlock struct {
	Level int    `json:"level"`
	Spell string `json:"spell"`
	Note  string `json:"note,omitempty"`
}

// unlocksFor returns the character's upcoming unlocks: everything from
// level+1 to level+10 (empty for level 60).
func unlocksFor(class string, level int) []unlock {
	if level >= 60 {
		return []unlock{}
	}
	all, ok := classUnlocks[normalizeClass(class)]
	if !ok {
		return []unlock{}
	}
	out := []unlock{}
	for _, u := range all {
		if u.Level > level && u.Level <= level+10 {
			out = append(out, u)
		}
	}
	return out
}

// nextUnlock is the single closest unlock — fed to the AI so it can nudge
// ("2 levels from Shadowform").
func nextUnlock(class string, level int) *unlock {
	us := unlocksFor(class, level)
	if len(us) == 0 {
		return nil
	}
	return &us[0]
}

func normalizeClass(class string) string {
	switch {
	case class == "":
		return ""
	default:
		return strings.ToLower(strings.TrimSpace(class))
	}
}

// classUnlocks: method.gg "Key Spells & Leveling Milestones" facts per class.
var classUnlocks = map[string][]unlock{
	"priest": {
		{4, "Shadow Word: Pain", "primary DoT"},
		{5, "Lesser Magic Wand", "first wand — big early damage"},
		{6, "Power Word: Shield", "pre-shield before pulls"},
		{10, "Mind Blast", "replaces Smite on cooldown"},
		{12, "Inner Fire, Power Word: Fortitude rank 2", "survivability spike"},
		{13, "Greater Magic Wand", "~50% higher dps wand"},
		{14, "Psychic Scream", "fear — your escape button"},
		{18, "Shadow Word: Pain rank 3, Dispel Magic", ""},
		{26, "Lesser Mystic Wand, Shadow Word: Pain rank 4", ""},
		{30, "Inner Fire rank 3, Power Word: Shield rank 5", "large survivability boost"},
		{32, "Shadow Word: Death", "instant damage"},
		{40, "Shadowform + Mind Flay", "respec to Shadow — big spike"},
		{48, "Prayer of Fortitude", "group stamina buff"},
	},
	"warlock": {
		{1, "Immolate", "fire DoT from the trainer (10c)"},
		{4, "Corruption, Curse of Weakness", ""},
		{6, "Life Tap, Shadow Bolt rank 2", "carry food to eat back the tap"},
		{8, "Bane of Agony", "instant tag — great on busy mobs"},
		{10, "Drain Soul, Summon Voidwalker", "soul shards + your tank pet"},
		{14, "Drain Life", "core of drain tanking"},
		{18, "Create Soulstone, Searing Pain", ""},
		{20, "Demon Armor, Rain of Fire, Ritual of Summoning", "your only summon method in Forever"},
		{24, "Shadowburn", "spends a shard for burst"},
		{29, "Conflagrate", "Immolate → Conflag combo"},
		{30, "Siphon Life", "Affliction DoT that heals you"},
	},
	"mage": {
		{4, "Frostbolt, Conjure Water", "bread and butter + free drinks"},
		{6, "Conjure Food", ""},
		{8, "Polymorph, Arcane Missiles", "best CC in the game"},
		{10, "Frost Nova", "create distance"},
		{14, "Arcane Explosion", "first AoE"},
		{16, "Flamestrike", "AoE DoT"},
		{20, "Ice Lance, Blizzard, Blink, Teleports", "huge spike — biggest unlock level"},
		{24, "Counterspell", "your interrupt"},
		{25, "Ice Block", "first real emergency button"},
		{26, "Cone of Cold", ""},
		{30, "Ice Armor, Thunder Bluff teleport", ""},
		{34, "Cold Snap", "resets frost cooldowns"},
		{40, "Ice Barrier, Portals", "shield now also boosts frost damage"},
	},
	"rogue": {
		{4, "Backstab", "needs a dagger in main hand"},
		{6, "Gouge", "line up Backstab / lock casts"},
		{8, "Evasion", "main defensive"},
		{10, "Dual Wield, Slice and Dice", "priority train — huge damage boost"},
		{12, "Kick", "interrupt"},
		{14, "Expose Armor", ""},
		{16, "Feint", "threat drop"},
		{18, "Ambush", "big stealth opener"},
		{22, "Vanish", "drop threat completely"},
		{26, "Cheap Shot", "stun with Kick/Gouge"},
	},
	"warrior": {
		{4, "Rend, Charge", "Rend ignores armor (bleed)"},
		{8, "Hamstring", "your escape/slow — vital solo"},
		{10, "Defensive Stance quest", "rewards Sunder Armor, Taunt + a weapon"},
		{12, "Overpower", "punishes dodges"},
		{14, "Revenge, Demoralizing Shout", ""},
		{20, "Victory Rush, Retaliation", "VR = rage-free damage + self-heal"},
		{30, "Berserker Stance quest + Whirlwind Axe chain", "WW axe is a huge weapon spike"},
		{36, "Whirlwind", "hits up to 4 targets"},
	},
	"druid": {
		{10, "Bear Form", "class quest — major power spike, do it ASAP"},
		{14, "Cure Poison", "quest chain"},
		{16, "Aquatic Form", "quest chain"},
		{20, "Cat Form", "class quest — second major spike"},
		{30, "Travel Form", "40% speed outdoors, no mount needed"},
		{40, "Mount + Dire Bear Form", "upgrade of Bear Form"},
	},
	"hunter": {
		{4, "Serpent Sting", "ranged DoT to 60"},
		{6, "Arcane Shot, Hunter's Mark", "Mark buffs you AND your pet"},
		{8, "Concussive Shot", "slow for kiting pre-pet"},
		{10, "Aspect of the Hawk, Call Pet", "THE level — pet tanks and does most damage"},
		{12, "Mend Pet", "heal your pet"},
		{18, "Multi-Shot", "AoE cleave"},
		{20, "Aimed Shot, Aspect of the Cheetah, Freezing Trap", "big spike level"},
		{25, "Summon Hawk (talent)", "2 extra pets — 18s each"},
		{26, "Rapid Fire", "first real cooldown"},
		{28, "Frost Trap", "60% slow AoE"},
		{30, "Feign Death, Intimidation (talent)", "drop combat / pet stun"},
	},
	"paladin": {
		{4, "Judgement, Blessing of Might", "seals + judgement = your rotation"},
		{6, "Holy Strike, Divine Protection", "6s full immunity"},
		{8, "Hammer of Justice, Purify", "stun"},
		{10, "Seal of Fury, Righteous Fury, Blessing of Protection", ""},
		{12, "Redemption", "rez"},
		{14, "Blessing of Wisdom", "1h mana regen buff"},
		{16, "Retribution Aura", ""},
		{20, "Consecration, Blessing of Kings", "toolkit essentially complete at 20"},
		{44, "Hammer of Wrath", "execute on targets under 20%"},
	},
	"shaman": {
		{10, "Searing Totem", "single target damage boost"},
		{12, "Fire Nova", "main burst AoE"},
		{20, "Ghost Wolf", "+40% move speed"},
		{25, "Stormstrike (talent)", "main damage ability"},
		{26, "Magma Totem", "pulsing AoE"},
		{30, "Windfury Weapon", "huge damage increase — switch to 2H"},
		{32, "Windfury Totem", "party-wide windfury"},
		{40, "Rage of the Farseer", "personal bloodlust"},
	},
}