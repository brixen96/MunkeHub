package main

import "strings"

// --- Leveling talent hints (Tier 3 of the Quest Helper) --------------------
//
// FACTS ONLY, scraped 2026-11-20 from the method.gg class guides listed in
// guides.md (9 classes). No prose copied — just spec names, talent names,
// level milestones and point priorities, same rule as leveling.go.
//
// The quest advisor gets the bracket matching the character's level so it can
// say "spend your next points into X" alongside quest advice.

type talentBracket struct {
	Min   int    `json:"min"`
	Max   int    `json:"max"`
	Hint  string `json:"hint"`
}

type classTalents struct {
	Spec     string           `json:"spec"`
	Brackets []talentBracket  `json:"brackets"`
}

// talentGuide: method.gg facts per class (keyed lowercase).
func talentGuide(class string) (classTalents, bool) {
	g, ok := talentGuides[strings.ToLower(strings.TrimSpace(class))]
	return g, ok
}

// talentHintFor returns the bracket matching the character's level (the
// bracket whose range contains level, else the last one below it).
func talentHintFor(class string, level int) (classTalents, *talentBracket) {
	g, ok := talentGuide(class)
	if !ok {
		return classTalents{}, nil
	}
	var pick *talentBracket
	for i := range g.Brackets {
		b := &g.Brackets[i]
		if level >= b.Min && level <= b.Max {
			return g, b
		}
		if b.Max <= level {
			pick = b // fallback: latest bracket already passed
		}
	}
	if pick == nil && len(g.Brackets) > 0 {
		pick = &g.Brackets[0]
	}
	return g, pick
}

var talentGuides = map[string]classTalents{
	"priest": {
		Spec: "Holy/ wand build 1-39, respec Shadow (Shadowform) at 40",
		Brackets: []talentBracket{
			{Min: 10, Max: 39, Hint: "Holy tree for mana efficiency + tankiness; keep upgrading your wand — damage falls off in the late 30s."},
			{Min: 40, Max: 60, Hint: "Respec to Shadow at 40 for Shadowform + Mind Flay; still fine to heal dungeons in shadow spec."},
		},
	},
	"warlock": {
		Spec: "Affliction (drain tanking), Demonology from 46",
		Brackets: []talentBracket{
			{Min: 10, Max: 29, Hint: "Affliction: 4 Suppression → 5 Improved Corruption → 1 Improved Life Tap → Pandemic."},
			{Min: 30, Max: 45, Hint: "Fill Fel Concentration, 2 Nightfall, 2 Malevolence → 3 Soul Siphon → 5 Shadow Mastery."},
			{Min: 46, Max: 60, Hint: "Swap to Demonology: 5 Unholy Power → Demonic Energies → Fel Vitality; finish in Affliction (Grim Reach, last Suppression)."},
		},
	},
	"mage": {
		Spec: "Frost, Arcane from 43",
		Brackets: []talentBracket{
			{Min: 10, Max: 25, Hint: "Frost: Improved Frostbolt → Precision → Frostbite → Ice Lance; 3/5 Ice Shards only, save the rest. Ice Block at 25."},
			{Min: 26, Max: 42, Hint: "Shatter + Fingers of Frost, Improved Cone of Cold, Winter's Chill, Cold Snap, then Ice Barrier."},
			{Min: 43, Max: 60, Hint: "Arcane: Improved Channeling → Arcane Concentration → Missile Barrage → Arcane Subtlety; Arcane Meditation last."},
		},
	},
	"rogue": {
		Spec: "Combat for solo, Subtlety if grouped (Assassination with a good dagger)",
		Brackets: []talentBracket{
			{Min: 10, Max: 29, Hint: "Combat for solo play; Assassination works if you have a dagger (Backstab/Mutilate)."},
			{Min: 30, Max: 60, Hint: "Swap to Subtlety if you looted a strong dagger by 30 (and can afford the respec); otherwise stay Combat."},
		},
	},
	"warrior": {
		Spec: "Arms (2H, with a group) or Protection (solo grind)",
		Brackets: []talentBracket{
			{Min: 10, Max: 19, Hint: "Arms: Improved Rend → Tactical Mastery → Deep Wounds. Protection alternative: Shield Specialization → Improved Thunder Clap → Improved Bloodrage."},
			{Min: 20, Max: 39, Hint: "Arms: Bloodthrill + Sweeping Strikes + Weaponmaster. Protection: Master of Defense → Improved Revenge → Vanguard (stay in Defensive Stance for rage)."},
			{Min: 40, Max: 60, Hint: "Arms capstone Mortal Strike, then back to Fury for Unbridled Wrath + Enrage. Protection: Bastion + Focused Rage."},
		},
	},
	"druid": {
		Spec: "Feral (damage build) for solo; Balance only past 30",
		Brackets: []talentBracket{
			{Min: 10, Max: 29, Hint: "Feral: Ferocity first → Feline Swiftness (great once Cat Form unlocks at 20) → deeper for Leader of the Pack + Predatory Strikes."},
			{Min: 30, Max: 60, Hint: "Feral capstones Primal Bite / Shifting Power; dip into Balance (5 Genesis → Nature's Majesty/Reach) and Resto (Heart of the Wild)."},
		},
	},
	"hunter": {
		Spec: "Beast Mastery, Marksmanship from 41",
		Brackets: []talentBracket{
			{Min: 10, Max: 24, Hint: "Beast Mastery: Deadly Aspects → Focused Fire → Improved Revive Pet → 1 Endurance Training (unlocks next row) → Pathfinding at 20."},
			{Min: 25, Max: 40, Hint: "Summon Hawk at 25, then Bestial Swiftness → Intimidation → Bestial Wrath capstone; leave Unleashed Fury/Frenzy at 4/5."},
			{Min: 41, Max: 60, Hint: "31 points in Beast Mastery, last 20 into Marksmanship: Lethal Attacks → Careful Aim → Rapid Killing → Efficiency → Trueshot Aura."},
		},
	},
	"paladin": {
		Spec: "Retribution first (mana + Pursuit of Justice), Holy/Protection later",
		Brackets: []talentBracket{
			{Min: 10, Max: 20, Hint: "Retribution: Benediction → Sanctified Judgement → Improved Judgement, racing to Pursuit of Justice (15% move speed, by ~21)."},
			{Min: 21, Max: 39, Hint: "Seal of Command → Conviction → Sacred Arbiter → Two-Handed Spec → Vengeance; or dip Holy (Divine Strength → Improved Seals → Holy Conduit)."},
			{Min: 40, Max: 60, Hint: "Capstone Unsealed Wrath (needs 1 Unyielding Faith); finish in Holy (Holy Conduit) and Protection (Precision). Dual spec unlocks at 40 (50g)."},
		},
	},
	"shaman": {
		Spec: "Enhancement",
		Brackets: []talentBracket{
			{Min: 10, Max: 24, Hint: "Enhancement tree toward Stormstrike (unlocks at 25)."},
			{Min: 25, Max: 60, Hint: "Stormstrike at 25, then 2/2 Improved Ghost Wolf right after (instant Ghost Wolf)."},
		},
	},
}