// MunkeHub profession guides — structured data, rendered with our armory
// styling (wowhead item links → tooltips work automatically).
// Facts only (skill ranges, item IDs, trainer coords) — written for the guild.

export const GUIDES = [
  {
    id: "cooking",
    name: "Cooking",
    icon: "inv_misc_food_01", // texture name → /api/v1/icon/
    blurb:
      "1-300 leveling route: what to train, what to buy, and what to cook at each tier. " +
      "Cooking pairs perfectly with Fishing — most reagents can be caught.",
    ranks: [
      { rank: "Apprentice", range: "1–75", skill: 0, level: 1 },
      { rank: "Journeyman", range: "75–150", skill: 50, level: 10 },
      { rank: "Expert", range: "150–225", skill: 125, level: 20 },
      { rank: "Artisan", range: "225–300", skill: 200, level: 35 },
    ],
    shopping: [
      {
        range: "1–50",
        items: [
          { qty: "~52", id: 6291, name: "Raw Brilliant Smallfish" },
          { qty: "1", id: 6325, name: "Recipe: Brilliant Smallfish" },
        ],
      },
      {
        range: "50–150",
        items: [
          { qty: "~55", id: 6289, name: "Raw Longjaw Mud Snapper" },
          { qty: "~55", id: 6308, name: "Raw Bristle Whisker Catfish" },
          { qty: "1", id: 6328, name: "Recipe: Longjaw Mud Snapper" },
          { qty: "1", id: 6330, name: "Recipe: Bristle Whisker Catfish" },
        ],
      },
      {
        range: "150–225",
        items: [
          { qty: "~27", id: 3685, name: "Raptor Egg" },
          { qty: "~27", id: 2692, name: "Hot Spices" },
          { qty: "1", id: 3682, name: "Recipe: Curiously Tasty Omelet" },
          { qty: "~55", id: 8365, name: "Raw Mithril Head Trout" },
          { qty: "1", id: 17062, name: "Recipe: Mithril Head Trout" },
          { qty: "1", id: 16072, name: "Expert Cookbook" },
        ],
      },
      {
        range: "225–300",
        items: [
          { qty: "12", id: 12207, name: "Giant Egg" },
          { qty: "10", id: 7974, name: "Zesty Clam Meat" },
          { qty: "20", id: 8932, name: "Alterac Swiss" },
          { qty: "~130", id: 4603, name: "Raw Spotted Yellowtail" },
          { qty: "1", id: 13939, name: "Recipe: Spotted Yellowtail" },
          { qty: "1", id: 13946, name: "Recipe: Poached Sunscale Salmon" },
        ],
      },
    ],
    steps: [
      {
        range: "1–50",
        note: "",
        rows: [
          {
            skill: "1-50",
            qty: "~52x",
            craft: { id: 6324, name: "Brilliant Smallfish" },
            reagents: [{ qty: "~52", id: 6291, name: "Raw Brilliant Smallfish" }],
          },
        ],
      },
      {
        range: "50–150",
        note: "Train Journeyman Cooking — requires level 10.",
        rows: [
          {
            skill: "50-100",
            qty: "~55x",
            craft: { id: 6329, name: "Longjaw Mud Snapper" },
            reagents: [{ qty: "~55", id: 6289, name: "Raw Longjaw Mud Snapper" }],
          },
          {
            skill: "100-150",
            qty: "~55x",
            craft: { id: 6331, name: "Bristle Whisker Catfish" },
            reagents: [{ qty: "~55", id: 6308, name: "Raw Bristle Whisker Catfish" }],
          },
        ],
      },
      {
        range: "150–225",
        note:
          "Buy the Expert Cookbook to raise your cap to 225 — requires level 20. " +
          "Alliance: Shandira in Ashenvale. Horde: Wulan in Desolace.",
        rows: [
          {
            skill: "150-175",
            qty: "~27x",
            craft: { id: 3683, name: "Curiously Tasty Omelet" },
            reagents: [
              { qty: "~27", id: 3685, name: "Raptor Egg" },
              { qty: "~27", id: 2692, name: "Hot Spices" },
            ],
          },
          {
            skill: "175-225",
            qty: "~55x",
            craft: { id: 17063, name: "Mithril Head Trout" },
            reagents: [{ qty: "~55", id: 8365, name: "Raw Mithril Head Trout" }],
          },
        ],
      },
      {
        range: "225–300",
        note:
          "Do the Clamlette Surprise quest in Gadgetzan (Dirge Quikcleave) to raise your " +
          "cap to 300 — requires level 35. Turn-in: 12 Giant Egg, 10 Zesty Clam Meat, " +
          "20 Alterac Swiss. Buy Recipe: Spotted Yellowtail from Gikkix while you're there.",
        rows: [
          {
            skill: "225-300",
            qty: "~130x",
            craft: { id: 13937, name: "Spotted Yellowtail" },
            reagents: [{ qty: "~130", id: 4603, name: "Raw Spotted Yellowtail" }],
          },
        ],
      },
    ],
    items: [
      {
        id: 5198,
        name: "Cookie's Stirring Rod",
        note: "+25% Cooking speed. Drops from Cookie in The Deadmines — great for big batches.",
      },
    ],
    trainers: {
      alliance: {
        basic: [
          { name: "Tomas", zone: "Elwynn Forest", coords: "44.4 66.0" },
          { name: "Zarrin", zone: "Teldrassil", coords: "57.1 61.3" },
          { name: "Gremlock Pilsnor", zone: "Dun Morogh", coords: "47.7 52.3" },
          { name: "Cook Ghilm", zone: "Dun Morogh", coords: "68.4 54.5" },
          { name: "Stephen Ryback", zone: "Stormwind", coords: "78.2 53.1" },
          { name: "Daryl Riknussun", zone: "Ironforge", coords: "60.1 36.4" },
          { name: "Alegorn", zone: "Darnassus", coords: "49.0 21.2" },
          { name: "Crystal Boughman", zone: "Redridge Mountains", coords: "17.7 43.5" },
        ],
        expert: [{ name: "Shandira", zone: "Ashenvale", coords: "49.5 67.1" }],
      },
      horde: {
        basic: [
          { name: "Cook Torka", zone: "Durotar", coords: "51.1 42.5" },
          { name: "Pyall Silentstride", zone: "Mulgore", coords: "44.9 61.7" },
          { name: "William Pickman", zone: "Tirisfal Glades", coords: "61.8 51.6" },
          { name: "Zamja", zone: "Orgrimmar", coords: "57.4 54.0" },
          { name: "Eunice Burch", zone: "Undercity", coords: "62.1 44.9" },
          { name: "Aska Mistrunner", zone: "Thunder Bluff", coords: "50.7 53.1" },
          { name: "Duhng", zone: "The Barrens", coords: "55.3 31.8" },
          { name: "Mudduk", zone: "Stranglethorn Vale", coords: "31.3 28.0" },
        ],
        expert: [{ name: "Wulan", zone: "Desolace", coords: "26.2 69.7" }],
      },
      artisan: [{ name: "Dirge Quikcleave", zone: "Tanaris (Gadgetzan)", coords: "52.6 28.1" }],
    },
  },
];