// GENERATED from scraped guide facts (node scripts/gen-guides.js) — do not edit by hand.
// Facts only: skill ranges, quantities, item IDs, trainer names/zones.
export const GUIDES_EXTRA = [
  {
    "id": "alchemy",
    "name": "Alchemy",
    "icon": "trade_alchemy",
    "blurb": "Cheapest 1-300 route with full shopping lists. Pairs with Herbalism — grow your own herbs and save gold.",
    "shopping": [
      {
        "range": "1-150",
        "items": [
          {
            "qty": "83",
            "id": 2447,
            "name": "Peacebloom"
          },
          {
            "qty": "54",
            "id": 765,
            "name": "Silverleaf"
          },
          {
            "qty": "50",
            "id": 785,
            "name": "Mageroyal"
          },
          {
            "qty": "56",
            "id": 2449,
            "name": "Earthroot"
          },
          {
            "qty": "60",
            "id": 2450,
            "name": "Briarthorn"
          },
          {
            "qty": "10",
            "id": 3820,
            "name": "Stranglekelp"
          },
          {
            "qty": "40",
            "id": 6359,
            "name": "Firefin Snapper"
          }
        ]
      },
      {
        "range": "150-225",
        "items": [
          {
            "qty": "72",
            "id": 3820,
            "name": "Stranglekelp"
          },
          {
            "qty": "51",
            "id": 3356,
            "name": "Kingsblood"
          },
          {
            "qty": "36",
            "id": 3357,
            "name": "Liferoot"
          }
        ]
      },
      {
        "range": "225-300",
        "items": [
          {
            "qty": "26",
            "id": 3357,
            "name": "Liferoot"
          },
          {
            "qty": "26",
            "id": 3820,
            "name": "Stranglekelp"
          },
          {
            "qty": "25",
            "id": 3358,
            "name": "Khadgar's Whisker"
          },
          {
            "qty": "47",
            "id": 8838,
            "name": "Sungrass"
          },
          {
            "qty": "25",
            "id": 8839,
            "name": "Blindweed"
          },
          {
            "qty": "47",
            "id": 13422,
            "name": "Stonescale Eel"
          }
        ]
      }
    ],
    "steps": [
      {
        "range": "1-55",
        "note": "",
        "rows": [
          {
            "skill": "1-10",
            "qty": "9x",
            "craft": {
              "name": "Elixir of Minor Force",
              "id": 1245250,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "9",
                "name": "Silverleaf",
                "id": null
              },
              {
                "qty": "9",
                "name": "Empty Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "10-55",
            "qty": "45x",
            "craft": {
              "name": "Elixir of Minor Defense",
              "id": 7183,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "45",
                "name": "Silverleaf",
                "id": null
              },
              {
                "qty": "45",
                "name": "Peacebloom",
                "id": null
              },
              {
                "qty": "45",
                "name": "Empty Vial",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "55-150",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "rows": [
          {
            "skill": "55-65",
            "qty": "10x",
            "craft": {
              "name": "Minor Mana Potion",
              "id": 2331,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "10",
                "name": "Mageroyal",
                "id": null
              },
              {
                "qty": "10",
                "name": "Peacebloom",
                "id": null
              },
              {
                "qty": "10",
                "name": "Empty Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "65-90",
            "qty": "28x",
            "craft": {
              "name": "Elixir of Minor Fortitude",
              "id": 2334,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "56",
                "name": "Earthroot",
                "id": null
              },
              {
                "qty": "28",
                "name": "Peacebloom",
                "id": null
              },
              {
                "qty": "28",
                "name": "Empty Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "90-120",
            "qty": "30x",
            "craft": {
              "name": "Elixir of Wisdom",
              "id": 3171,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "30",
                "name": "Mageroyal",
                "id": null
              },
              {
                "qty": "60",
                "name": "Briarthorn",
                "id": null
              },
              {
                "qty": "30",
                "name": "Empty Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "120-130",
            "qty": "10x",
            "craft": {
              "name": "Lesser Mana Potion",
              "id": 3173,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "10",
                "name": "Mageroyal",
                "id": null
              },
              {
                "qty": "10",
                "name": "Stranglekelp",
                "id": null
              },
              {
                "qty": "10",
                "name": "Empty Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "130-150",
            "qty": "20x",
            "craft": {
              "name": "Fire Oil",
              "id": 7837,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "40",
                "name": "Firefin Snapper",
                "id": null
              },
              {
                "qty": "20",
                "name": "Empty Vial",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "150-225",
        "note": "Train Expert - requires skill 125 and level 20.",
        "rows": [
          {
            "skill": "150-165",
            "qty": "15x",
            "craft": {
              "name": "Elixir of Fire Power",
              "id": 7845,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "30",
                "name": "Fire Oil",
                "id": null
              },
              {
                "qty": "15",
                "name": "Kingsblood",
                "id": null
              },
              {
                "qty": "15",
                "name": "Leaded Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "165-195",
            "qty": "36x",
            "craft": {
              "name": "Mana Potion",
              "id": 3452,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "36",
                "name": "Stranglekelp",
                "id": null
              },
              {
                "qty": "36",
                "name": "Kingsblood",
                "id": null
              },
              {
                "qty": "36",
                "name": "Leaded Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "195-225",
            "qty": "36x",
            "craft": {
              "name": "Nature Protection Potion",
              "id": 7259,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "36",
                "name": "Liferoot",
                "id": null
              },
              {
                "qty": "36",
                "name": "Stranglekelp",
                "id": null
              },
              {
                "qty": "36",
                "name": "Leaded Vial",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "225-300",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "rows": [
          {
            "skill": "225-235",
            "qty": "26x",
            "craft": {
              "name": "Nature Protection Potion",
              "id": 7259,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "26",
                "name": "Liferoot",
                "id": null
              },
              {
                "qty": "26",
                "name": "Stranglekelp",
                "id": null
              },
              {
                "qty": "26",
                "name": "Leaded Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "235-250",
            "qty": "15x",
            "craft": {
              "name": "Elixir of Greater Intellect",
              "id": 11465,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "15",
                "name": "Blindweed",
                "id": null
              },
              {
                "qty": "15",
                "name": "Khadgar's Whisker",
                "id": null
              },
              {
                "qty": "15",
                "name": "Crystal Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "250-260",
            "qty": "47x",
            "craft": {
              "name": "Stonescale Oil",
              "id": 17551,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "47",
                "name": "Stonescale Eel",
                "id": null
              },
              {
                "qty": "47",
                "name": "Leaded Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "260-265",
            "qty": "10x",
            "craft": {
              "name": "Elixir of Greater Intellect",
              "id": 11465,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "10",
                "name": "Blindweed",
                "id": null
              },
              {
                "qty": "10",
                "name": "Khadgar's Whisker",
                "id": null
              },
              {
                "qty": "10",
                "name": "Crystal Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "265-300",
            "qty": "47x",
            "craft": {
              "name": "Elixir of Greater Defense",
              "id": 17554,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "94",
                "name": "Stonescale Oil",
                "id": null
              },
              {
                "qty": "47",
                "name": "Sungrass",
                "id": null
              },
              {
                "qty": "47",
                "name": "Crystal Vial",
                "id": null
              }
            ]
          }
        ]
      }
    ],
    "gathering": [],
    "trainers": {
      "sections": [
        {
          "title": "Capital Cities",
          "lines": [
            "Orgrimmar: Whuut",
            "Undercity: Doctor Martin Felben",
            "Thunder Bluff: Kray"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Miao'zan in Sen'jin Village"
            ],
            [
              "Tauren",
              "Mulgore",
              "There is no Alchemy trainer in Mulgore, go to Kray on the middle rise in Thunder Bluff"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Carolai Anise in the small house at the north-west end of Brill"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Iaadaria Bitterwind or Valiena Swiftgale in Valanaar, Nyassa Swiftdraught in Shen'dar Village"
            ]
          ],
          "notes": []
        },
        {
          "title": "Capital Cities",
          "lines": [
            "Stormwind City: Tel'Athir",
            "Ironforge: Vosur Brakthel",
            "Darnassus: Milla Fairancora"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Human",
              "Elwynn Forest",
              "Alchemist Mallory in a cottage in the hills between Stormwind and Northshire"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "There is no Alchemy trainer in Dun Morogh, go to Vosur Brakthel in Tinker Town in Ironforge"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "Cyndra Kindwhisper in a small house east of the inn in Dolanaar"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Iaadaria Bitterwind or Valiena Swiftgale in Valanaar, Nyassa Swiftdraught in Shen'dar Village"
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "blacksmithing",
    "name": "Blacksmithing",
    "icon": "trade_blacksmithing",
    "blurb": "1-300 route covering armor and weapon smithing, plus specializations.",
    "shopping": [
      {
        "range": "1-150",
        "items": [
          {
            "qty": "90",
            "id": 2835,
            "name": "Rough Stone"
          },
          {
            "qty": "180",
            "id": 2840,
            "name": "Copper Bar"
          },
          {
            "qty": "5",
            "id": 2842,
            "name": "Silver Bar"
          },
          {
            "qty": "72",
            "id": 2836,
            "name": "Coarse Stone"
          },
          {
            "qty": "290",
            "id": 2841,
            "name": "Bronze Bar"
          },
          {
            "qty": "60",
            "id": 2838,
            "name": "Heavy Stone"
          },
          {
            "qty": "58",
            "id": 2319,
            "name": "Medium Leather"
          }
        ]
      },
      {
        "range": "150-225",
        "items": [
          {
            "qty": "180",
            "id": 3575,
            "name": "Iron Bar"
          },
          {
            "qty": "30",
            "id": 2605,
            "name": "Green Dye (vendor)"
          },
          {
            "qty": "32",
            "id": 7912,
            "name": "Solid Stone"
          },
          {
            "qty": "100",
            "id": 3859,
            "name": "Steel Bar (100 Iron Bar and 100 Coal if you have Mining 165)"
          },
          {
            "qty": "60",
            "id": 2838,
            "name": "Heavy Stone"
          },
          {
            "qty": "436",
            "id": 3860,
            "name": "Mithril Bar"
          },
          {
            "qty": "24",
            "id": 4234,
            "name": "Heavy Leather"
          },
          {
            "qty": "12",
            "id": 3864,
            "name": "Citrine"
          }
        ]
      },
      {
        "range": "225-300",
        "items": [
          {
            "qty": "5",
            "id": 12365,
            "name": "Dense Stone"
          },
          {
            "qty": "224",
            "id": 3860,
            "name": "Mithril Bar"
          },
          {
            "qty": "64",
            "id": 4304,
            "name": "Thick Leather"
          },
          {
            "qty": "1910",
            "id": 12359,
            "name": "Thorium Bar"
          },
          {
            "qty": "30",
            "id": 8170,
            "name": "Rugged Leather"
          },
          {
            "qty": "67",
            "id": 7910,
            "name": "Star Ruby"
          },
          {
            "qty": "25",
            "id": 7909,
            "name": "Aquamarine"
          }
        ]
      }
    ],
    "steps": [
      {
        "range": "1-75",
        "note": "",
        "rows": [
          {
            "skill": "1-25",
            "qty": "30x",
            "craft": {
              "name": "Rough Sharpening Stone",
              "id": 2660,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "30",
                "name": "Rough Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "25-50",
            "qty": "30x",
            "craft": {
              "name": "Rough Grinding Stone",
              "id": 3320,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "60",
                "name": "Rough Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "50-75",
            "qty": "36x",
            "craft": {
              "name": "Coarse Grinding Stone",
              "id": 3326,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "72",
                "name": "Coarse Stone",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "75-150",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "rows": [
          {
            "skill": "75-80",
            "qty": "5x",
            "craft": {
              "name": "Silver Rod",
              "id": 7818,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "5",
                "name": "Silver Bar",
                "id": null
              },
              {
                "qty": "10",
                "name": "Rough Grinding Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "80-95",
            "qty": "18x",
            "craft": {
              "name": "Runed Copper Belt",
              "id": 2666,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "180",
                "name": "Copper Bar",
                "id": null
              }
            ]
          },
          {
            "skill": "95-100",
            "qty": "5x",
            "craft": {
              "name": "Patterned Bronze Bracers",
              "id": 2672,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "25",
                "name": "Bronze Bar",
                "id": null
              },
              {
                "qty": "10",
                "name": "Coarse Grinding Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "100-112",
            "qty": "20x",
            "craft": {
              "name": "Heavy Grinding Stone",
              "id": 3337,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "60",
                "name": "Heavy Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "112-125",
            "qty": "13x",
            "craft": {
              "name": "Patterned Bronze Bracers",
              "id": 2672,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "65",
                "name": "Bronze Bar",
                "id": null
              },
              {
                "qty": "26",
                "name": "Coarse Grinding Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "125-140",
            "qty": "24x",
            "craft": {
              "name": "Bronze Shortsword",
              "id": 2742,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "120",
                "name": "Bronze Bar",
                "id": null
              },
              {
                "qty": "96",
                "name": "Weak Flux",
                "id": null
              },
              {
                "qty": "48",
                "name": "Medium Leather",
                "id": null
              }
            ]
          },
          {
            "skill": "140-150",
            "qty": "10x",
            "craft": {
              "name": "Bronze Warhammer",
              "id": 9985,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "80",
                "name": "Bronze Bar",
                "id": null
              },
              {
                "qty": "10",
                "name": "Strong Flux",
                "id": null
              },
              {
                "qty": "10",
                "name": "Medium Leather",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "150-215",
        "note": "Train Expert - requires skill 125 and level 20.",
        "rows": [
          {
            "skill": "150-175",
            "qty": "30x",
            "craft": {
              "name": "Green Iron Bracers",
              "id": 3501,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "180",
                "name": "Iron Bar",
                "id": null
              },
              {
                "qty": "30",
                "name": "Green Dye",
                "id": null
              }
            ]
          },
          {
            "skill": "175-180",
            "qty": "8x",
            "craft": {
              "name": "Solid Grinding Stone",
              "id": 9920,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "32",
                "name": "Solid Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "180-195",
            "qty": "20x",
            "craft": {
              "name": "Golden Scale Bracers",
              "id": 7223,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "100",
                "name": "Steel Bar",
                "id": null
              },
              {
                "qty": "40",
                "name": "Heavy Grinding Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "195-210",
            "qty": "17x",
            "craft": {
              "name": "Mithril Plate Pants",
              "id": 9931,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "204",
                "name": "Mithril Bar",
                "id": null
              }
            ]
          },
          {
            "skill": "210-215",
            "qty": "6x",
            "craft": {
              "name": "Heavy Mithril Axe",
              "id": 9993,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "72",
                "name": "Mithril Bar",
                "id": null
              },
              {
                "qty": "12",
                "name": "Citrine",
                "id": null
              },
              {
                "qty": "6",
                "name": "Solid Grinding Stone",
                "id": null
              },
              {
                "qty": "24",
                "name": "Heavy Leather",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "215-300",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "rows": [
          {
            "skill": "245-270",
            "qty": "27x",
            "craft": {
              "name": "Imperial Plate Bracers",
              "id": 16649,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "540",
                "name": "Thorium Bar",
                "id": null
              },
              {
                "qty": "27",
                "name": "Star Ruby",
                "id": null
              }
            ]
          },
          {
            "skill": "270-290",
            "qty": "20x",
            "craft": {
              "name": "Imperial Plate Boots",
              "id": 16657,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "680",
                "name": "Thorium Bar",
                "id": null
              },
              {
                "qty": "20",
                "name": "Star Ruby",
                "id": null
              },
              {
                "qty": "20",
                "name": "Aquamarine",
                "id": null
              }
            ]
          },
          {
            "skill": "290-300",
            "qty": "10x",
            "craft": {
              "name": "Imperial Plate Leggings",
              "id": 16730,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "440",
                "name": "Thorium Bar",
                "id": null
              },
              {
                "qty": "20",
                "name": "Star Ruby",
                "id": null
              }
            ]
          }
        ]
      }
    ],
    "gathering": [],
    "trainers": {
      "sections": [
        {
          "title": "Capital Cities",
          "lines": [
            "Orgrimmar: Ug'thok",
            "Undercity: Basil Frye",
            "Thunder Bluff: Thrag Stonehoof"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Dwukk at the forge in Razor Hill"
            ],
            [
              "Tauren",
              "Mulgore",
              "There is no Blacksmithing trainer in Mulgore, go to Thrag Stonehoof on the lower rise in Thunder Bluff"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Angus Hammerhand at Bandarion Keep. Nobody near him sells a Blacksmith Hammer, buy one before you go from Walter Mason in Deathknell or Abigail Shiel in Brill."
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Daeann Steelwind in Valanaar, Aedi Thriceforged in Shen'dar Village"
            ]
          ],
          "notes": []
        },
        {
          "title": "Capital Cities",
          "lines": [
            "Stormwind City: Dane Lindgren",
            "Ironforge: Groum Stonebeard"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Human",
              "Elwynn Forest",
              "Smith Argus in Goldshire, across from the inn"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "Tognus Flintfire in Kharanos, across from the Thunderbrew Distillery"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "There is no Blacksmithing trainer in Teldrassil or Darnassus, the closest one is Delfrum Flintbeard in Auberdine, Darkshore"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Daeann Steelwind in Valanaar, Aedi Thriceforged in Shen'dar Village"
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "enchanting",
    "name": "Enchanting",
    "icon": "trade_enchanting",
    "blurb": "1-300 route with a disenchanting loop — level Enchanting while leveling other professions.",
    "shopping": [
      {
        "range": "1-150",
        "items": [
          {
            "qty": "19",
            "id": 247786,
            "name": "Mote of Magic (vendor)"
          },
          {
            "qty": "1",
            "id": 6217,
            "name": "Copper Rod (vendor)"
          },
          {
            "qty": "265",
            "id": 10940,
            "name": "Strange Dust"
          },
          {
            "qty": "1",
            "id": 6338,
            "name": "Silver Rod"
          },
          {
            "qty": "1",
            "id": 1210,
            "name": "Shadowgem"
          },
          {
            "qty": "4",
            "id": 10938,
            "name": "Lesser Magic Essence"
          },
          {
            "qty": "3",
            "id": 10940,
            "name": "Strange Dust"
          }
        ]
      },
      {
        "range": "150-225",
        "items": [
          {
            "qty": "64",
            "id": 11083,
            "name": "Soul Dust"
          },
          {
            "qty": "1",
            "id": 11128,
            "name": "Golden Rod"
          },
          {
            "qty": "1",
            "id": 5500,
            "name": "Iridescent Pearl"
          },
          {
            "qty": "1",
            "id": 11082,
            "name": "Greater Astral Essence"
          },
          {
            "qty": "32",
            "id": 17034,
            "name": "Maple Seed (vendor)"
          },
          {
            "qty": "41",
            "id": 3372,
            "name": "Leaded Vial (vendor)"
          },
          {
            "qty": "31",
            "id": 11137,
            "name": "Vision Dust"
          },
          {
            "qty": "1",
            "id": 11144,
            "name": "Truesilver Rod"
          },
          {
            "qty": "1",
            "id": 7971,
            "name": "Black Pearl"
          },
          {
            "qty": "1",
            "id": 11135,
            "name": "Greater Mystic Essence"
          },
          {
            "qty": "50",
            "id": 17035,
            "name": "Stranglethorn Seed (vendor)"
          }
        ]
      },
      {
        "range": "225-300",
        "items": [
          {
            "qty": "37",
            "id": 11137,
            "name": "Vision Dust"
          },
          {
            "qty": "260",
            "id": 11176,
            "name": "Dream Dust"
          },
          {
            "qty": "20",
            "id": 8831,
            "name": "Purple Lotus"
          },
          {
            "qty": "20",
            "id": 8925,
            "name": "Crystal Vial (vendor)"
          },
          {
            "qty": "1",
            "id": 16206,
            "name": "Arcanite Rod"
          },
          {
            "qty": "1",
            "id": 13926,
            "name": "Golden Pearl"
          },
          {
            "qty": "10",
            "id": 16204,
            "name": "Illusion Dust"
          },
          {
            "qty": "4",
            "id": 16203,
            "name": "Greater Eternal Essence"
          },
          {
            "qty": "4",
            "id": 14343,
            "name": "Small Brilliant Shard"
          },
          {
            "qty": "2",
            "id": 14344,
            "name": "Large Brilliant Shard"
          }
        ]
      }
    ],
    "steps": [
      {
        "range": "1-70",
        "note": "",
        "rows": [
          {
            "skill": "1-2",
            "qty": "1x",
            "craft": {
              "name": "Runed Copper Rod",
              "id": 7421,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "1",
                "name": "Copper Rod",
                "id": null
              },
              {
                "qty": "1",
                "name": "Mote of Magic",
                "id": null
              }
            ]
          },
          {
            "skill": "2-20",
            "qty": "18x",
            "craft": {
              "name": "Enchant Bracer",
              "id": 7418,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "",
                "name": "Inferior Stamina - 18 Mote of Magic",
                "id": null
              }
            ]
          },
          {
            "skill": "20-70",
            "qty": "",
            "craft": {
              "name": "Make one of these:",
              "id": null,
              "k": "spell"
            },
            "reagents": []
          }
        ]
      },
      {
        "range": "150-225",
        "note": "Train Expert - requires skill 125 and level 20.",
        "rows": [
          {
            "skill": "150-151",
            "qty": "1x",
            "craft": {
              "name": "Runed Golden Rod",
              "id": 13628,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "1",
                "name": "Golden Rod",
                "id": null
              },
              {
                "qty": "1",
                "name": "Iridescent Pearl",
                "id": null
              },
              {
                "qty": "1",
                "name": "Greater Astral Essence",
                "id": null
              },
              {
                "qty": "1",
                "name": "Soul Dust",
                "id": null
              }
            ]
          },
          {
            "skill": "151-165",
            "qty": "16x",
            "craft": {
              "name": "Minor Mana Oil",
              "id": 25125,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "16",
                "name": "Soul Dust",
                "id": null
              },
              {
                "qty": "32",
                "name": "Maple Seed",
                "id": null
              },
              {
                "qty": "16",
                "name": "Leaded Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "165-200",
            "qty": "47x",
            "craft": {
              "name": "Enchant Boots",
              "id": 13637,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "",
                "name": "Lesser Agility - 47 Soul Dust",
                "id": null
              }
            ]
          },
          {
            "skill": "200-201",
            "qty": "1x",
            "craft": {
              "name": "Runed Truesilver Rod",
              "id": 13702,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "1",
                "name": "Truesilver Rod",
                "id": null
              },
              {
                "qty": "1",
                "name": "Black Pearl",
                "id": null
              },
              {
                "qty": "1",
                "name": "Greater Mystic Essence",
                "id": null
              },
              {
                "qty": "1",
                "name": "Vision Dust",
                "id": null
              }
            ]
          },
          {
            "skill": "201-220",
            "qty": "25x",
            "craft": {
              "name": "Lesser Wizard Oil",
              "id": 25126,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "25",
                "name": "Vision Dust",
                "id": null
              },
              {
                "qty": "50",
                "name": "Stranglethorn Seed",
                "id": null
              },
              {
                "qty": "25",
                "name": "Leaded Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "220-225",
            "qty": "5x",
            "craft": {
              "name": "Enchant Cloak",
              "id": 13746,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "",
                "name": "Greater Defense - 5 Vision Dust",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "225-300",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "rows": [
          {
            "skill": "225-250",
            "qty": "37x",
            "craft": {
              "name": "Enchant Gloves",
              "id": 13815,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "",
                "name": "Agility - 37 Vision Dust",
                "id": null
              }
            ]
          },
          {
            "skill": "250-265",
            "qty": "20x",
            "craft": {
              "name": "Lesser Mana Oil",
              "id": 25127,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "20",
                "name": "Dream Dust",
                "id": null
              },
              {
                "qty": "20",
                "name": "Purple Lotus",
                "id": null
              },
              {
                "qty": "20",
                "name": "Crystal Vial",
                "id": null
              }
            ]
          },
          {
            "skill": "265-299",
            "qty": "40x",
            "craft": {
              "name": "Enchant Shield",
              "id": 20017,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "",
                "name": "Greater Stamina - 240 Dream Dust",
                "id": null
              }
            ]
          },
          {
            "skill": "299-300",
            "qty": "1x",
            "craft": {
              "name": "Runed Arcanite Rod",
              "id": 20051,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "1",
                "name": "Arcanite Rod",
                "id": null
              },
              {
                "qty": "1",
                "name": "Golden Pearl",
                "id": null
              },
              {
                "qty": "10",
                "name": "Illusion Dust",
                "id": null
              },
              {
                "qty": "4",
                "name": "Greater Eternal Essence",
                "id": null
              },
              {
                "qty": "4",
                "name": "Small Brilliant Shard",
                "id": null
              },
              {
                "qty": "2",
                "name": "Large Brilliant Shard",
                "id": null
              }
            ]
          }
        ]
      }
    ],
    "gathering": [],
    "trainers": {
      "sections": [
        {
          "title": "Capital Cities",
          "lines": [
            "Orgrimmar: Jhag",
            "Undercity: Malcomb Wynn",
            "Thunder Bluff: Mot Dawnstrider"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Pa'zula in Sen'jin Village"
            ],
            [
              "Tauren",
              "Mulgore",
              "There is no Enchanting trainer in Mulgore, go to Mot Dawnstrider on the middle rise in Thunder Bluff"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Vance Undergloom upstairs in the inn in Brill"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Eaysaa Brightgust in Valanaar, Nasalanna Windsinger in Shen'dar Village"
            ]
          ],
          "notes": []
        },
        {
          "title": "Capital Cities",
          "lines": [
            "Stormwind City: Betty Quin",
            "Ironforge: Thonys Pillarstone",
            "Darnassus: Lalina Summermoon"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Human",
              "Elwynn Forest",
              "There is no Enchanting trainer in Elwynn Forest, go to Betty Quin in the Mage Quarter in Stormwind City"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "There is no Enchanting trainer in Dun Morogh, go to Thonys Pillarstone at the Great Forge in Ironforge"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "Alanna Raveneye at the Oracle Glade, in the north-west of Teldrassil"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Eaysaa Brightgust in Valanaar, Nasalanna Windsinger in Shen'dar Village"
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "engineering",
    "name": "Engineering",
    "icon": "trade_engineering",
    "blurb": "1-300 route with the cheapest parts, plus specializations.",
    "shopping": [
      {
        "range": "1-150",
        "items": [
          {
            "qty": "66",
            "id": 2835,
            "name": "Rough Stone"
          },
          {
            "qty": "72",
            "id": 2840,
            "name": "Copper Bar"
          },
          {
            "qty": "53",
            "id": 2589,
            "name": "Linen Cloth"
          },
          {
            "qty": "60",
            "id": 2836,
            "name": "Coarse Stone"
          },
          {
            "qty": "5",
            "id": 2842,
            "name": "Silver Bar"
          },
          {
            "qty": "80",
            "id": 2841,
            "name": "Bronze Bar"
          },
          {
            "qty": "10",
            "id": 1206,
            "name": "Moss Agate"
          },
          {
            "qty": "30",
            "id": 2838,
            "name": "Heavy Stone"
          },
          {
            "qty": "15",
            "id": 2592,
            "name": "Wool Cloth"
          }
        ]
      },
      {
        "range": "150-225",
        "items": [
          {
            "qty": "30",
            "id": 2841,
            "name": "Bronze Bar"
          },
          {
            "qty": "45",
            "id": 2592,
            "name": "Wool Cloth"
          },
          {
            "qty": "15",
            "id": 2319,
            "name": "Medium Leather"
          },
          {
            "qty": "4",
            "id": 3859,
            "name": "Steel Bar"
          },
          {
            "qty": "120",
            "id": 7912,
            "name": "Solid Stone"
          },
          {
            "qty": "74",
            "id": 3860,
            "name": "Mithril Bar"
          },
          {
            "qty": "20",
            "id": 4338,
            "name": "Mageweave Cloth"
          }
        ]
      },
      {
        "range": "225-300",
        "items": [
          {
            "qty": "84",
            "id": 3860,
            "name": "Mithril Bar"
          },
          {
            "qty": "88",
            "id": 12365,
            "name": "Dense Stone"
          },
          {
            "qty": "225",
            "id": 12359,
            "name": "Thorium Bar"
          },
          {
            "qty": "35",
            "id": 14047,
            "name": "Runecloth"
          }
        ]
      }
    ],
    "steps": [
      {
        "range": "1-50",
        "note": "",
        "rows": [
          {
            "skill": "1-30",
            "qty": "66x",
            "craft": {
              "name": "Rough Blasting Powder",
              "id": 3918,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "66",
                "name": "Rough Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "30-50",
            "qty": "33x",
            "craft": {
              "name": "Handful of Copper Bolts",
              "id": 3922,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "33",
                "name": "Copper Bar",
                "id": null
              }
            ]
          },
          {
            "skill": "50-51",
            "qty": "1x",
            "craft": {
              "name": "Arclight Spanner",
              "id": 7430,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "6",
                "name": "Copper Bar",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "50-150",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "rows": [
          {
            "skill": "51-75",
            "qty": "33x",
            "craft": {
              "name": "Rough Copper Bomb",
              "id": 3923,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "33",
                "name": "Copper Bar",
                "id": null
              },
              {
                "qty": "33",
                "name": "Handful of Copper Bolts",
                "id": null
              },
              {
                "qty": "66",
                "name": "Rough Blasting Powder",
                "id": null
              },
              {
                "qty": "33",
                "name": "Linen Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "75-90",
            "qty": "60x",
            "craft": {
              "name": "Coarse Blasting Powder",
              "id": 3929,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "60",
                "name": "Coarse Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "90-100",
            "qty": "20x",
            "craft": {
              "name": "Coarse Dynamite",
              "id": 3931,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "60",
                "name": "Coarse Blasting Powder",
                "id": null
              },
              {
                "qty": "20",
                "name": "Linen Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "100-105",
            "qty": "5x",
            "craft": {
              "name": "Silver Contact",
              "id": 3973,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "5",
                "name": "Silver Bar",
                "id": null
              }
            ]
          },
          {
            "skill": "105-125",
            "qty": "25x",
            "craft": {
              "name": "Bronze Tube",
              "id": 3938,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "50",
                "name": "Bronze Bar",
                "id": null
              },
              {
                "qty": "25",
                "name": "Weak Flux",
                "id": null
              }
            ]
          },
          {
            "skill": "125-135",
            "qty": "10x",
            "craft": {
              "name": "Standard Scope",
              "id": 3978,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "10",
                "name": "Bronze Tube",
                "id": null
              },
              {
                "qty": "10",
                "name": "Moss Agate",
                "id": null
              }
            ]
          },
          {
            "skill": "135-150",
            "qty": "",
            "craft": {
              "name": "Make both of these:",
              "id": null,
              "k": "spell"
            },
            "reagents": []
          }
        ]
      },
      {
        "range": "150-200",
        "note": "Train Expert - requires skill 125 and level 20.",
        "rows": [
          {
            "skill": "150-160",
            "qty": "15x",
            "craft": {
              "name": "Bronze Framework",
              "id": 3953,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "30",
                "name": "Bronze Bar",
                "id": null
              },
              {
                "qty": "15",
                "name": "Medium Leather",
                "id": null
              },
              {
                "qty": "15",
                "name": "Wool Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "160-175",
            "qty": "15x",
            "craft": {
              "name": "Explosive Sheep",
              "id": 3955,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "30",
                "name": "Heavy Blasting Powder",
                "id": null
              },
              {
                "qty": "15",
                "name": "Whirring Bronze Gizmo",
                "id": null
              },
              {
                "qty": "15",
                "name": "Bronze Framework",
                "id": null
              },
              {
                "qty": "30",
                "name": "Wool Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "175-194",
            "qty": "60x",
            "craft": {
              "name": "Solid Blasting Powder",
              "id": 12585,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "120",
                "name": "Solid Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "194-195",
            "qty": "1x",
            "craft": {
              "name": "Gyromatic Micro-Adjustor",
              "id": 12590,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "4",
                "name": "Steel Bar",
                "id": null
              }
            ]
          },
          {
            "skill": "195-200",
            "qty": "6x",
            "craft": {
              "name": "Mithril Tube",
              "id": 12589,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "18",
                "name": "Mithril Bar",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "200-300",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "rows": [
          {
            "skill": "200-215",
            "qty": "20x",
            "craft": {
              "name": "Unstable Trigger",
              "id": 12591,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "20",
                "name": "Mithril Bar",
                "id": null
              },
              {
                "qty": "20",
                "name": "Mageweave Cloth",
                "id": null
              },
              {
                "qty": "20",
                "name": "Solid Blasting Powder",
                "id": null
              }
            ]
          },
          {
            "skill": "215-238",
            "qty": "40x",
            "craft": {
              "name": "Mithril Casing",
              "id": 12599,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "120",
                "name": "Mithril Bar",
                "id": null
              }
            ]
          },
          {
            "skill": "238-250",
            "qty": "20x",
            "craft": {
              "name": "Hi-Explosive Bomb",
              "id": 12619,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "40",
                "name": "Mithril Casing",
                "id": null
              },
              {
                "qty": "20",
                "name": "Unstable Trigger",
                "id": null
              },
              {
                "qty": "40",
                "name": "Solid Blasting Powder",
                "id": null
              }
            ]
          },
          {
            "skill": "250-260",
            "qty": "44x",
            "craft": {
              "name": "Dense Blasting Powder",
              "id": 19788,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "88",
                "name": "Dense Stone",
                "id": null
              }
            ]
          },
          {
            "skill": "260-285",
            "qty": "35x",
            "craft": {
              "name": "Thorium Widget",
              "id": 19791,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "105",
                "name": "Thorium Bar",
                "id": null
              },
              {
                "qty": "35",
                "name": "Runecloth",
                "id": null
              }
            ]
          },
          {
            "skill": "285-300",
            "qty": "20x",
            "craft": {
              "name": "Thorium Tube",
              "id": 19795,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "120",
                "name": "Thorium Bar",
                "id": null
              }
            ]
          }
        ]
      }
    ],
    "gathering": [],
    "trainers": {
      "sections": [
        {
          "title": "Capital Cities",
          "lines": [
            "Orgrimmar: Thund",
            "Undercity: Graham Van Talen"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Mukdrak in Razor Hill, near the forge next to the road to Orgrimmar"
            ],
            [
              "Tauren",
              "Mulgore",
              "Twizwick Sprocketgrind on the plateau above the Venture Co. Mine, hard to reach at low level. There is no anvil near him and nobody there sells the Blacksmith Hammer, so it's probably best to go to Thund in Orgrimmar instead."
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "There is no Engineering trainer in Tirisfal Glades, go to Graham Van Talen in the Rogues' Quarter in Undercity"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Strange Hermit in Shadowgale Forest"
            ]
          ],
          "notes": []
        },
        {
          "title": "Capital Cities",
          "lines": [
            "Stormwind City: Sprite Jumpsprocket",
            "Ironforge: Jemma Quikswitch"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Human",
              "Elwynn Forest",
              "There is no Engineering trainer in Elwynn Forest, go to Sprite Jumpsprocket in the Dwarven District in Stormwind City"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "Bronk Guzzlegear at Steelgrill's Depot"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "There is no Engineering trainer in Teldrassil or Darnassus, the closest one is Jenna Lemkenilli in Auberdine in Darkshore"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Strange Hermit in Shadowgale Forest"
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "leatherworking",
    "name": "Leatherworking",
    "icon": "trade_leatherworking",
    "blurb": "1-300 route for light/medium/heavy leathers, plus specializations. Pairs with Skinning.",
    "shopping": [
      {
        "range": "1-150",
        "items": [
          {
            "qty": "98-105",
            "id": 2934,
            "name": "Ruined Leather Scraps"
          },
          {
            "qty": "348-414",
            "id": 2318,
            "name": "Light Leather"
          },
          {
            "qty": "27",
            "id": 783,
            "name": "Light Hide"
          },
          {
            "qty": "120-139",
            "id": 2319,
            "name": "Medium Leather"
          },
          {
            "qty": "10",
            "id": 4232,
            "name": "Medium Hide"
          },
          {
            "qty": "30",
            "id": 4235,
            "name": "Heavy Hide"
          }
        ]
      },
      {
        "range": "150-225",
        "items": [
          {
            "qty": "109-126",
            "id": 2319,
            "name": "Medium Leather"
          },
          {
            "qty": "153",
            "id": 4234,
            "name": "Heavy Leather"
          },
          {
            "qty": "363-390",
            "id": 4304,
            "name": "Thick Leather"
          }
        ]
      },
      {
        "range": "225-300",
        "items": [
          {
            "qty": "208-224",
            "id": 4304,
            "name": "Thick Leather"
          },
          {
            "qty": "916-1030",
            "id": 8170,
            "name": "Rugged Leather"
          },
          {
            "qty": "464-530",
            "id": 14047,
            "name": "Runecloth"
          }
        ]
      }
    ],
    "steps": [
      {
        "range": "1-55",
        "note": "",
        "rows": [
          {
            "skill": "1-30",
            "qty": "35x",
            "craft": {
              "name": "Light Leather",
              "id": 2881,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "105",
                "name": "Ruined Leather Scraps",
                "id": null
              }
            ]
          },
          {
            "skill": "30-55",
            "qty": "27x",
            "craft": {
              "name": "Cured Light Hide",
              "id": 3816,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "27",
                "name": "Light Hide",
                "id": null
              },
              {
                "qty": "27",
                "name": "Salt",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "55-130",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "rows": [
          {
            "skill": "55-75",
            "qty": "21x",
            "craft": {
              "name": "Embossed Leather Gloves",
              "id": 3756,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "63",
                "name": "Light Leather",
                "id": null
              },
              {
                "qty": "42",
                "name": "Coarse Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "75-80",
            "qty": "8x",
            "craft": {
              "name": "Medium Leather",
              "id": 20648,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "32",
                "name": "Light Leather",
                "id": null
              }
            ]
          },
          {
            "skill": "80-90",
            "qty": "10x",
            "craft": {
              "name": "Cured Medium Hide",
              "id": 3817,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "10",
                "name": "Medium Hide",
                "id": null
              },
              {
                "qty": "10",
                "name": "Salt",
                "id": null
              }
            ]
          },
          {
            "skill": "90-100",
            "qty": "19x",
            "craft": {
              "name": "Fine Leather Belt",
              "id": 3763,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "114",
                "name": "Light Leather",
                "id": null
              },
              {
                "qty": "38",
                "name": "Coarse Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "100-115",
            "qty": "24x",
            "craft": {
              "name": "Light Leather Pants",
              "id": 9068,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "240",
                "name": "Light Leather",
                "id": null
              },
              {
                "qty": "24",
                "name": "Cured Light Hide",
                "id": null
              },
              {
                "qty": "24",
                "name": "Fine Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "115-125",
            "qty": "10x",
            "craft": {
              "name": "Dark Leather Belt",
              "id": 3766,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "10",
                "name": "Fine Leather Belt",
                "id": null
              },
              {
                "qty": "10",
                "name": "Cured Medium Hide",
                "id": null
              },
              {
                "qty": "20",
                "name": "Fine Thread",
                "id": null
              },
              {
                "qty": "10",
                "name": "Gray Dye",
                "id": null
              }
            ]
          },
          {
            "skill": "125-130",
            "qty": "7x",
            "craft": {
              "name": "Heavy Leather",
              "id": 20649,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "35",
                "name": "Medium Leather",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "130-205",
        "note": "Train Expert - requires skill 125 and level 20.",
        "rows": [
          {
            "skill": "130-143",
            "qty": "30x",
            "craft": {
              "name": "Cured Heavy Hide",
              "id": 3818,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "30",
                "name": "Heavy Hide",
                "id": null
              },
              {
                "qty": "90",
                "name": "Salt",
                "id": null
              }
            ]
          },
          {
            "skill": "143-155",
            "qty": "17x",
            "craft": {
              "name": "Hillman's Leather Gloves",
              "id": 3764,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "238",
                "name": "Medium Leather",
                "id": null
              },
              {
                "qty": "68",
                "name": "Fine Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "155-165",
            "qty": "10x",
            "craft": {
              "name": "Barbaric Shoulders",
              "id": 7151,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "80",
                "name": "Heavy Leather",
                "id": null
              },
              {
                "qty": "10",
                "name": "Cured Heavy Hide",
                "id": null
              },
              {
                "qty": "20",
                "name": "Fine Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "165-185",
            "qty": "20x",
            "craft": {
              "name": "Guardian Gloves",
              "id": 7156,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "80",
                "name": "Heavy Leather",
                "id": null
              },
              {
                "qty": "20",
                "name": "Cured Heavy Hide",
                "id": null
              },
              {
                "qty": "20",
                "name": "Silken Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "185-205",
            "qty": "22x",
            "craft": {
              "name": "Nightscape Headband",
              "id": 10507,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "110",
                "name": "Thick Leather",
                "id": null
              },
              {
                "qty": "44",
                "name": "Silken Thread",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "205-300",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "rows": [
          {
            "skill": "205-235",
            "qty": "36x",
            "craft": {
              "name": "Nightscape Pants",
              "id": 10548,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "504",
                "name": "Thick Leather",
                "id": null
              },
              {
                "qty": "144",
                "name": "Silken Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "235-265",
            "qty": "36x",
            "craft": {
              "name": "Wicked Leather Gauntlets",
              "id": 19049,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "288",
                "name": "Rugged Leather",
                "id": null
              },
              {
                "qty": "36",
                "name": "Black Dye",
                "id": null
              },
              {
                "qty": "36",
                "name": "Rune Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "265-300",
            "qty": "53x",
            "craft": {
              "name": "Runic Leather Headband",
              "id": 19082,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "742",
                "name": "Rugged Leather",
                "id": null
              },
              {
                "qty": "530",
                "name": "Runecloth",
                "id": null
              },
              {
                "qty": "53",
                "name": "Rune Thread",
                "id": null
              }
            ]
          }
        ]
      }
    ],
    "gathering": [],
    "trainers": {
      "sections": [
        {
          "title": "Capital Cities",
          "lines": [
            "Orgrimmar: Kamari",
            "Undercity: Dan Golthas",
            "Thunder Bluff: Mak"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "There is no Leatherworking trainer in Durotar, go to Kamari in the Drag in Orgrimmar"
            ],
            [
              "Tauren",
              "Mulgore",
              "Chaw Stronghide in Bloodhoof Village"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Shelene Rhobart at the road junction southeast of Brill"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Emerii Tallgust in Valanaar, Indari Sunseam in Shen'dar Village"
            ]
          ],
          "notes": []
        },
        {
          "title": "Capital Cities",
          "lines": [
            "Stormwind City: Randal Worth",
            "Ironforge: Gretta Finespindle",
            "Darnassus: Darianna"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Human",
              "Elwynn Forest",
              "Adele Fielder in the small house by Crystal Lake, northeast of Goldshire"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "There is no Leatherworking trainer in Dun Morogh, go to Gretta Finespindle at the Great Forge in Ironforge"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "Nadyia Maneweaver in the small house near the road to Darnassus"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Emerii Tallgust in Valanaar, Indari Sunseam in Shen'dar Village"
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "tailoring",
    "name": "Tailoring",
    "icon": "trade_tailoring",
    "blurb": "1-300 route — cheap cloth crafts, no gathering needed.",
    "shopping": [
      {
        "range": "1-150",
        "items": [
          {
            "qty": "60",
            "id": 2589,
            "name": "Linen Cloth"
          },
          {
            "qty": "408-459",
            "id": 2592,
            "name": "Wool Cloth"
          },
          {
            "qty": "29",
            "id": 3383,
            "name": "Elixir of Wisdom"
          },
          {
            "qty": "205-216",
            "id": 4306,
            "name": "Silk Cloth"
          },
          {
            "qty": "94-115",
            "id": 4338,
            "name": "Mageweave Cloth"
          }
        ]
      },
      {
        "range": "150-225",
        "items": [
          {
            "qty": "492-520",
            "id": 4306,
            "name": "Silk Cloth"
          },
          {
            "qty": "507-625",
            "id": 4338,
            "name": "Mageweave Cloth"
          },
          {
            "qty": "104-115",
            "id": 14047,
            "name": "Runecloth"
          }
        ]
      },
      {
        "range": "225-280",
        "items": [
          {
            "qty": "1292-1440",
            "id": 14047,
            "name": "Runecloth"
          },
          {
            "qty": "44",
            "id": 8170,
            "name": "Rugged Leather"
          },
          {
            "qty": "39",
            "id": 14227,
            "name": "Ironweb Spider Silk"
          }
        ]
      }
    ],
    "steps": [
      {
        "range": "1-55",
        "note": "",
        "rows": [
          {
            "skill": "1-30",
            "qty": "30x",
            "craft": {
              "name": "Bolt of Linen Cloth",
              "id": 2963,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "60",
                "name": "Linen Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "30-55",
            "qty": "28x",
            "craft": {
              "name": "Linen Belt",
              "id": 8776,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "28",
                "name": "Bolt of Linen Cloth",
                "id": null
              },
              {
                "qty": "28",
                "name": "Coarse Thread",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "55-140",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "rows": [
          {
            "skill": "55-80",
            "qty": "",
            "craft": {
              "name": "Learn Journeyman Tailoring before you make these. You can reach past 75 with them, but you won't get any skill points above 75 until you learn it, so don't go AFK while crafting before you have it.",
              "id": null,
              "k": "spell"
            },
            "reagents": []
          },
          {
            "skill": "80-95",
            "qty": "22x",
            "craft": {
              "name": "Woolen Cape",
              "id": 2402,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "22",
                "name": "Bolt of Woolen Cloth",
                "id": null
              },
              {
                "qty": "22",
                "name": "Fine Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "95-100",
            "qty": "5x",
            "craft": {
              "name": "Double-stitched Woolen Shoulders",
              "id": 3848,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "15",
                "name": "Bolt of Woolen Cloth",
                "id": null
              },
              {
                "qty": "10",
                "name": "Fine Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "100-120",
            "qty": "54x",
            "craft": {
              "name": "Bolt of Silk Cloth",
              "id": 3839,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "216",
                "name": "Silk Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "120-140",
            "qty": "29x",
            "craft": {
              "name": "Gloves of Meditation",
              "id": 3852,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "116",
                "name": "Bolt of Woolen Cloth",
                "id": null
              },
              {
                "qty": "87",
                "name": "Fine Thread",
                "id": null
              },
              {
                "qty": "29",
                "name": "Elixir of Wisdom",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "140-200",
        "note": "Train Expert - requires skill 125 and level 20.",
        "rows": [
          {
            "skill": "140-150",
            "qty": "23x",
            "craft": {
              "name": "Bolt of Mageweave",
              "id": 3865,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "115",
                "name": "Mageweave Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "150-170",
            "qty": "20x",
            "craft": {
              "name": "Crimson Silk Vest",
              "id": 8791,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "80",
                "name": "Bolt of Silk Cloth",
                "id": null
              },
              {
                "qty": "40",
                "name": "Red Dye",
                "id": null
              },
              {
                "qty": "40",
                "name": "Fine Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "170-190",
            "qty": "26x",
            "craft": {
              "name": "Crimson Silk Pantaloons",
              "id": 8799,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "104",
                "name": "Bolt of Silk Cloth",
                "id": null
              },
              {
                "qty": "52",
                "name": "Red Dye",
                "id": null
              },
              {
                "qty": "52",
                "name": "Silken Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "190-200",
            "qty": "19x",
            "craft": {
              "name": "Black Mageweave Leggings",
              "id": 12049,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "38",
                "name": "Bolt of Mageweave",
                "id": null
              },
              {
                "qty": "57",
                "name": "Silken Thread",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "200-300",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "rows": [
          {
            "skill": "200-205",
            "qty": "5x",
            "craft": {
              "name": "Black Mageweave Shoulders",
              "id": 12074,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "15",
                "name": "Bolt of Mageweave",
                "id": null
              },
              {
                "qty": "10",
                "name": "Heavy Silken Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "205-215",
            "qty": "23x",
            "craft": {
              "name": "Bolt of Runecloth",
              "id": 18401,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "115",
                "name": "Runecloth",
                "id": null
              }
            ]
          },
          {
            "skill": "215-225",
            "qty": "19x",
            "craft": {
              "name": "Red Mageweave Bag",
              "id": 12079,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "95",
                "name": "Bolt of Mageweave",
                "id": null
              },
              {
                "qty": "38",
                "name": "Red Dye",
                "id": null
              },
              {
                "qty": "38",
                "name": "Heavy Silken Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "225-240",
            "qty": "24x",
            "craft": {
              "name": "Runecloth Belt",
              "id": 18402,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "72",
                "name": "Bolt of Runecloth",
                "id": null
              },
              {
                "qty": "24",
                "name": "Rune Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "240-250",
            "qty": "11x",
            "craft": {
              "name": "Runecloth Gloves",
              "id": 18417,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "44",
                "name": "Bolt of Runecloth",
                "id": null
              },
              {
                "qty": "44",
                "name": "Rugged Leather",
                "id": null
              },
              {
                "qty": "11",
                "name": "Rune Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "250-280",
            "qty": "39x",
            "craft": {
              "name": "Runecloth Cuffs",
              "id": 1257487,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "195",
                "name": "Bolt of Runecloth",
                "id": null
              },
              {
                "qty": "39",
                "name": "Ironweb Spider Silk",
                "id": null
              },
              {
                "qty": "39",
                "name": "Rune Thread",
                "id": null
              }
            ]
          },
          {
            "skill": "280-300",
            "qty": "",
            "craft": {
              "name": "Make the Runecloth Shoulders if you can get the pattern. These are all the recipes I found so far that still give points after 280. I left out the ones that aren't worth it for leveling, and there might be new ones that nobody has found yet.",
              "id": 18449,
              "k": "spell"
            },
            "reagents": []
          }
        ]
      }
    ],
    "gathering": [],
    "trainers": {
      "sections": [
        {
          "title": "Capital Cities",
          "lines": [
            "Orgrimmar: Snang",
            "Undercity: Victor Ward",
            "Thunder Bluff: Vhan"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "There is no Tailoring trainer in Durotar, go to Snang in the Drag in Orgrimmar"
            ],
            [
              "Tauren",
              "Mulgore",
              "There is no Tailoring trainer in Mulgore, go to Vhan on the middle rise in Thunder Bluff"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Bowen Brisboise in the small house at Cold Hearth Manor, southwest of Brill"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Othesia Evengale in Valanaar, Taleen Shimmerthread in Shen'dar Village"
            ]
          ],
          "notes": []
        },
        {
          "title": "Capital Cities",
          "lines": [
            "Stormwind City: Lawrence Schneider",
            "Ironforge: Uthrar Threx",
            "Darnassus: Trianna"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Human",
              "Elwynn Forest",
              "Eldrin in a house at the Eastvale Logging Camp"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "There is no Tailoring trainer in Dun Morogh, go to Uthrar Threx at the Great Forge in Ironforge"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "There is no Tailoring trainer in Teldrassil, go to Trianna at the Craftsmen's Terrace in Darnassus"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Othesia Evengale in Valanaar, Taleen Shimmerthread in Shen'dar Village"
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "mining",
    "name": "Mining",
    "icon": "trade_mining",
    "blurb": "Zone-by-zone mining route 1-300 with smelting shortcuts and skill bonus items.",
    "shopping": [],
    "steps": [],
    "gathering": [
      {
        "range": "1-65",
        "note": "",
        "zones": [
          [
            "Durotar"
          ],
          [
            "Mulgore"
          ],
          [
            "Tirisfal Glades"
          ],
          [
            "Zephras Isle"
          ],
          [
            "Elwynn Forest"
          ],
          [
            "Dun Morogh"
          ],
          [
            "Teldrassil"
          ]
        ]
      },
      {
        "range": "65-125",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "zones": [
          [
            "Wetlands"
          ],
          [
            "Silverpine Forest"
          ],
          [
            "Westfall"
          ],
          [
            "Hillsbrad Foothills"
          ],
          [
            "The Barrens"
          ],
          [
            "Ashenvale"
          ],
          [
            "Stonetalon Mountains"
          ]
        ]
      },
      {
        "range": "125-175",
        "note": "Train Expert - requires skill 125 and level 20.",
        "zones": [
          [
            "Arathi Highlands"
          ],
          [
            "Thousand Needles"
          ],
          [
            "Desolace"
          ]
        ]
      },
      {
        "range": "175-245",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "zones": [
          [
            "The Hinterlands"
          ],
          [
            "Tanaris"
          ]
        ]
      },
      {
        "range": "245-275",
        "note": "",
        "zones": [
          [
            "Un'Goro Crater"
          ],
          [
            "Blasted Lands"
          ],
          [
            "Felwood"
          ]
        ]
      },
      {
        "range": "275-300",
        "note": "",
        "zones": [
          [
            "Un'Goro Crater"
          ],
          [
            "Eastern Plaguelands"
          ],
          [
            "Winterspring"
          ],
          [
            "Burning Steppes"
          ]
        ]
      }
    ],
    "trainers": {
      "sections": [
        {
          "title": "",
          "lines": [
            "Orgrimmar: Makaru",
            "Undercity: Brom Killian",
            "Thunder Bluff: Brek Stonehoof",
            "Norzsh in the Valley of Trials (only up to 75)",
            "Krunn in Razor Hill",
            "Vartha Rockmane in Camp Narache (only up to 75)",
            "Brek Stonehoof in Thunder Bluff",
            "Walter Mason in Deathknell (only up to 75)",
            "Brom Killian in Undercity",
            "Stormwind City: Gelman Stonehand",
            "Ironforge: Geofram Bouldertoe",
            "Darnassus: no trainer, the closest one is Kurdram Stonehammer in Auberdine, Darkshore",
            "Kelsey Fargo (only up to 75)",
            "Gelman Stonehand in Stormwind City",
            "Sally Swiftwrench (only up to 75)",
            "Yarr Hammerstone at Steelgrill's Depot",
            "Dank Drizzlecut at Gol'Bolar Quarry",
            "Fanorran Stilloak at Aldrassil in Shadowglen (only up to 75)",
            "Kurdram Stonehammer in Auberdine, Darkshore"
          ],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Norzsh in the Valley of Trials (only up to 75)\nKrunn in Razor Hill"
            ],
            [
              "Tauren",
              "Mulgore",
              "Vartha Rockmane in Camp Narache (only up to 75)\n\tBrek Stonehoof in Thunder Bluff"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Walter Mason in Deathknell (only up to 75)\n\tBrom Killian in Undercity"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "There is no trainer in Thendal Village, where you start, but you don't need one. Harvesting Windstones is one of the first quests you get in Thendal Village. Pick Mining for Dummies as your reward, it teaches you Mining.\n            The trainers on the island are Halavuul Cragwind in Valanaar and Messana Crestwind in Shen'dar Village."
            ],
            [
              "Human",
              "Elwynn Forest",
              "Kelsey Fargo (only up to 75)\nGelman Stonehand in Stormwind City"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "Sally Swiftwrench (only up to 75)\n\tYarr Hammerstone at Steelgrill's Depot\n\tDank Drizzlecut at Gol'Bolar Quarry"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "Fanorran Stilloak at Aldrassil in Shadowglen (only up to 75)\n\tKurdram Stonehammer in Auberdine, Darkshore"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "There is no trainer in Thendal Village, where you start, but you don't need one. Harvesting Windstones is one of the first quests you get in Thendal Village. Pick Mining for Dummies as your reward, it teaches you Mining.\n            The trainers on the island are Halavuul Cragwind in Valanaar and Messana Crestwind in Shen'dar Village."
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "herbalism",
    "name": "Herbalism",
    "icon": "trade_herbalism",
    "blurb": "Zone-by-zone herb route 1-300 with skill bonus items.",
    "shopping": [],
    "steps": [],
    "gathering": [
      {
        "range": "1-70",
        "note": "",
        "zones": [
          [
            "Durotar"
          ],
          [
            "Mulgore"
          ],
          [
            "Tirisfal Glades"
          ],
          [
            "Zephras Isle"
          ],
          [
            "Elwynn Forest"
          ],
          [
            "Dun Morogh"
          ],
          [
            "Teldrassil"
          ]
        ]
      },
      {
        "range": "70-115",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "zones": [
          [
            "Silverpine Forest"
          ],
          [
            "Loch Modan"
          ],
          [
            "Darkshore"
          ]
        ]
      },
      {
        "range": "115-170",
        "note": "Train Expert - requires skill 125 and level 20.",
        "zones": [
          [
            "Hillsbrad Foothills"
          ],
          [
            "Stonetalon Mountains"
          ]
        ]
      },
      {
        "range": "170-205",
        "note": "",
        "zones": [
          [
            "Stranglethorn Vale"
          ],
          [
            "Arathi Highlands"
          ]
        ]
      },
      {
        "range": "205-230",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "zones": [
          [
            "Tanaris"
          ],
          [
            "Searing Gorge"
          ]
        ]
      }
    ],
    "trainers": {
      "sections": [
        {
          "title": "",
          "lines": [
            "Orgrimmar: Jandi",
            "Undercity: Martha Alliestar",
            "Thunder Bluff: Komin Winterhoof",
            "Zor'la in the Valley of Trials (only up to 75)",
            "Mishiki in Sen'jin Village",
            "Garan Sunstrider in Camp Narache (only up to 75)",
            "Komin Winterhoof in Thunder Bluff",
            "Florence Nightshade in Deathknell (only up to 75)",
            "Faruza in Brill",
            "Stormwind City: Tannysa",
            "Ironforge: Reyna Stonebranch",
            "Darnassus: Firodren Mooncaller",
            "Nancy Songflower (only up to 75)",
            "Herbalist Pomeroy",
            "Emrys Flintbeard (only up to 75)",
            "Reyna Stonebranch in Ironforge",
            "Eleyna Duskbreeze at Aldrassil in Shadowglen (only up to 75)",
            "Malorne Bladeleaf in Dolanaar"
          ],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Zor'la in the Valley of Trials (only up to 75)\nMishiki in Sen'jin Village"
            ],
            [
              "Tauren",
              "Mulgore",
              "Garan Sunstrider in Camp Narache (only up to 75)\n\tKomin Winterhoof in Thunder Bluff"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Florence Nightshade in Deathknell (only up to 75)\n\tFaruza in Brill"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "There is no trainer in Thendal Village, where you start, but you don't need one. Harvesting Windstones is one of the first quests you get in Thendal Village. Pick Wild Harvest as your reward, it teaches you Herbalism.\n\t\t\t\tThe trainers on the island are Syriel Nightrain in Valanaar and Halassa Fernbreeze in Shen'dar Village."
            ],
            [
              "Human",
              "Elwynn Forest",
              "Nancy Songflower (only up to 75)\nHerbalist Pomeroy"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "Emrys Flintbeard (only up to 75)\n\tReyna Stonebranch in Ironforge"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "Eleyna Duskbreeze at Aldrassil in Shadowglen (only up to 75)\n\tMalorne Bladeleaf in Dolanaar"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "There is no trainer in Thendal Village, where you start, but you don't need one. Harvesting Windstones is one of the first quests you get in Thendal Village. Pick Wild Harvest as your reward, it teaches you Herbalism.\n\t\t\t\tThe trainers on the island are Syriel Nightrain in Valanaar and Halassa Fernbreeze in Shen'dar Village."
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "skinning",
    "name": "Skinning",
    "icon": "trade_skinning",
    "blurb": "Zone-by-zone skinning route 1-300 with skill requirements per rank.",
    "shopping": [],
    "steps": [],
    "gathering": [
      {
        "range": "1-75",
        "note": "Buy a Fishing Pole and Shiny Baubles from a Fishing Supply vendor.",
        "zones": [
          [
            "Durotar"
          ],
          [
            "Mulgore"
          ],
          [
            "Tirisfal Glades"
          ],
          [
            "Zephras Isle"
          ],
          [
            "Elwynn Forest"
          ],
          [
            "Dun Morogh"
          ],
          [
            "Teldrassil"
          ]
        ]
      },
      {
        "range": "75-150",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "zones": [
          [
            "Darnassus"
          ],
          [
            "Stormwind City"
          ],
          [
            "Orgrimmar"
          ],
          [
            "Undercity"
          ],
          [
            "Thunder Bluff"
          ]
        ]
      },
      {
        "range": "150-225",
        "note": "Expert Fishing: buy the book from Old Man Heming in Booty Bay (requires level 20 and Fishing 125). Train Expert - requires skill 125 and level 20.",
        "zones": [
          [
            "Hillsbrad Foothills"
          ],
          [
            "Duskwood"
          ]
        ]
      },
      {
        "range": "225-280",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "zones": []
      }
    ],
    "trainers": {
      "sections": [
        {
          "title": "",
          "lines": [
            "Orgrimmar: Thuwd",
            "Undercity: Killian Hagey",
            "Thunder Bluff: Mooranta",
            "Kagil in the Valley of Trials (only up to 75)",
            "Thuwd in Orgrimmar",
            "Nawka Wildsong in Camp Narache (only up to 75)",
            "Yonn Deepcut in Bloodhoof Village",
            "Margaret Weaver in Deathknell (only up to 75)",
            "Rand Rhobart in Tirisfal Glades",
            "Stormwind City: Maris Granger",
            "Ironforge: Balthus Stoneflayer",
            "Darnassus: Eladriel",
            "Riley Peltskinner in Northshire Valley (only up to 75)",
            "Helene Peltskinner in Goldshire",
            "Brighid Stormflayer (only up to 75)",
            "Balthus Stoneflayer in Ironforge",
            "Terunne Bearshaper at Aldrassil in Shadowglen (only up to 75)",
            "Radnaal Maneweaver in Teldrassil"
          ],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Kagil in the Valley of Trials (only up to 75)\nThuwd in Orgrimmar"
            ],
            [
              "Tauren",
              "Mulgore",
              "Nawka Wildsong in Camp Narache (only up to 75)\n\tYonn Deepcut in Bloodhoof Village"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Margaret Weaver in Deathknell (only up to 75)\n\tRand Rhobart in Tirisfal Glades"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "There is no trainer in Thendal Village, where you start, but you don't need one. Harvesting Windstones is one of the first quests you get in Thendal Village. Pick Pelt Collecting for Beginners as your reward, it teaches you Skinning.\n\t\t\t\tThe trainers on the island are Ishlee Breezewhisper in Valanaar and Mendalass Tattermend in Shen'dar Village."
            ],
            [
              "Human",
              "Elwynn Forest",
              "Riley Peltskinner in Northshire Valley (only up to 75)\nHelene Peltskinner in Goldshire"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "Brighid Stormflayer (only up to 75)\n\tBalthus Stoneflayer in Ironforge"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "Terunne Bearshaper at Aldrassil in Shadowglen (only up to 75)\n\tRadnaal Maneweaver in Teldrassil"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "There is no trainer in Thendal Village, where you start, but you don't need one. Harvesting Windstones is one of the first quests you get in Thendal Village. Pick Pelt Collecting for Beginners as your reward, it teaches you Skinning.\n\t\t\t\tThe trainers on the island are Ishlee Breezewhisper in Valanaar and Mendalass Tattermend in Shen'dar Village."
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "first-aid",
    "name": "First Aid",
    "icon": "trade_alchemy",
    "blurb": "1-300 route — bandages from linen to runecloth, cheapest cloth first.",
    "shopping": [
      {
        "range": "1-225",
        "items": [
          {
            "qty": "148",
            "id": 2589,
            "name": "Linen Cloth"
          },
          {
            "qty": "32",
            "id": 2589,
            "name": "Linen Cloth"
          },
          {
            "qty": "12",
            "id": 2592,
            "name": "Wool Cloth"
          },
          {
            "qty": "151",
            "id": 2592,
            "name": "Wool Cloth"
          },
          {
            "qty": "141",
            "id": 4306,
            "name": "Silk Cloth"
          },
          {
            "qty": "19",
            "id": 4338,
            "name": "Mageweave Cloth"
          }
        ]
      },
      {
        "range": "225-300",
        "items": [
          {
            "qty": "89",
            "id": 4338,
            "name": "Mageweave Cloth"
          },
          {
            "qty": "20",
            "id": 14047,
            "name": "Runecloth"
          },
          {
            "qty": "30",
            "id": 13464,
            "name": "Golden Sansam"
          }
        ]
      }
    ],
    "steps": [
      {
        "range": "75-150",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "rows": [
          {
            "skill": "75-80",
            "qty": "16x",
            "craft": {
              "name": "Heavy Linen Bandage",
              "id": 3276,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "32",
                "name": "Linen Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "80-90",
            "qty": "12x",
            "craft": {
              "name": "Wool Bandage",
              "id": 3277,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "12",
                "name": "Wool Cloth",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "150-225",
        "note": "Train Expert - requires skill 125 and level 20.",
        "rows": [
          {
            "skill": "150-180",
            "qty": "47x",
            "craft": {
              "name": "Silk Bandage",
              "id": 7928,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "47",
                "name": "Silk Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "180-210",
            "qty": "47x",
            "craft": {
              "name": "Heavy Silk Bandage",
              "id": 7929,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "94",
                "name": "Silk Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "210-225",
            "qty": "19x",
            "craft": {
              "name": "Mageweave Bandage",
              "id": 10840,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "19",
                "name": "Mageweave Cloth",
                "id": null
              }
            ]
          }
        ]
      },
      {
        "range": "225-300",
        "note": "Train Artisan - requires skill 200 and level 35.",
        "rows": [
          {
            "skill": "225-240",
            "qty": "29x",
            "craft": {
              "name": "Mageweave Bandage",
              "id": 10840,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "29",
                "name": "Mageweave Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "240-260",
            "qty": "30x",
            "craft": {
              "name": "Heavy Mageweave Bandage",
              "id": 10841,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "60",
                "name": "Mageweave Cloth",
                "id": null
              }
            ]
          },
          {
            "skill": "260-275",
            "qty": "20x",
            "craft": {
              "name": "Runecloth Bandage",
              "id": 18629,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "20",
                "name": "Runecloth",
                "id": null
              }
            ]
          },
          {
            "skill": "275-300",
            "qty": "30x",
            "craft": {
              "name": "Major Healing Potion",
              "id": 1244436,
              "k": "spell"
            },
            "reagents": [
              {
                "qty": "30",
                "name": "Golden Sansam",
                "id": null
              },
              {
                "qty": "90",
                "name": "Soothing Spices",
                "id": null
              },
              {
                "qty": "30",
                "name": "Crystal Vial",
                "id": null
              }
            ]
          }
        ]
      }
    ],
    "gathering": [],
    "trainers": {
      "sections": [
        {
          "title": "Capital Cities",
          "lines": [
            "Orgrimmar: Arnok",
            "Undercity: Mary Edras",
            "Thunder Bluff: Pand Stonebinder"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Rawrk in the barracks in Razor Hill"
            ],
            [
              "Tauren",
              "Mulgore",
              "Vira Younghoof inside the inn in Bloodhoof Village"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Nurse Neela inside the inn in Brill"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Melasa Fairmend in Valanaar, Naleeia Tattermend in Shen'dar Village"
            ]
          ],
          "notes": []
        },
        {
          "title": "Capital Cities",
          "lines": [
            "Stormwind City: Shaina Fuller",
            "Ironforge: Nissa Firestone",
            "Darnassus: Dannelor"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Human",
              "Elwynn Forest",
              "Michelle Belle upstairs in the inn in Goldshire"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "Thamner Pol in the Thunderbrew Distillery in Kharanos"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "Byancie in Dolanaar, in the building north of the inn"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Melasa Fairmend in Valanaar, Naleeia Tattermend in Shen'dar Village"
            ]
          ],
          "notes": []
        }
      ]
    }
  },
  {
    "id": "fishing",
    "name": "Fishing",
    "icon": "trade_fishing",
    "blurb": "Zone-by-zone fishing route 1-300 — pair with Cooking for free food.",
    "shopping": [],
    "steps": [],
    "gathering": [
      {
        "range": "1-75",
        "note": "Buy a Fishing Pole and Shiny Baubles from a Fishing Supply vendor.",
        "zones": [
          [
            "Durotar"
          ],
          [
            "Mulgore"
          ],
          [
            "Tirisfal Glades"
          ],
          [
            "Zephras Isle"
          ],
          [
            "Elwynn Forest"
          ],
          [
            "Dun Morogh"
          ],
          [
            "Teldrassil"
          ]
        ]
      },
      {
        "range": "75-150",
        "note": "Train Journeyman - requires skill 50 and level 10.",
        "zones": [
          [
            "Darnassus"
          ],
          [
            "Stormwind City"
          ],
          [
            "Orgrimmar"
          ],
          [
            "Undercity"
          ],
          [
            "Thunder Bluff"
          ]
        ]
      },
      {
        "range": "150-225",
        "note": "Expert Fishing: buy the book from Old Man Heming in Booty Bay (requires level 20 and Fishing 125). Train Expert - requires skill 125 and level 20.",
        "zones": [
          [
            "Dustwallow Marsh"
          ],
          [
            "Alterac Mountains"
          ],
          [
            "Arathi Highlands"
          ],
          [
            "Desolace"
          ],
          [
            "Stranglethorn Vale"
          ],
          [
            "Swamp of Sorrows"
          ],
          [
            "Thousand Needles"
          ]
        ]
      },
      {
        "range": "225-300",
        "note": "Artisan Fishing is learned via quest. Train Artisan - requires skill 200 and level 35.",
        "zones": [
          [
            "Felwood"
          ],
          [
            "Feralas"
          ],
          [
            "The Hinterlands"
          ],
          [
            "Tanaris"
          ],
          [
            "Un'Goro Crater"
          ],
          [
            "Western Plaguelands"
          ]
        ]
      }
    ],
    "trainers": {
      "sections": [
        {
          "title": "Capital Cities",
          "lines": [
            "Orgrimmar: Lumak",
            "Undercity: Armand Cromwell",
            "Thunder Bluff: Kah Mistrunner"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Orc, Troll",
              "Durotar",
              "Lau'Tiki on the coast at Sen'jin Village"
            ],
            [
              "Tauren",
              "Mulgore",
              "Uthan Stillwater at Stonebull Lake, next to Bloodhoof Village"
            ],
            [
              "Undead",
              "Tirisfal Glades",
              "Clyde Kellen at Brightwater Lake, east of Brill"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Baelann Swiftcurrent in Valanaar, Fenn Fairweather in Shen'dar Village"
            ]
          ],
          "notes": []
        },
        {
          "title": "Capital Cities",
          "lines": [
            "Stormwind City: Arnold Leland",
            "Ironforge: Grimnur Stonebrand",
            "Darnassus: Astaia"
          ],
          "rows": [],
          "notes": []
        },
        {
          "title": "Starting Zones",
          "lines": [],
          "rows": [
            [
              "Human",
              "Elwynn Forest",
              "Lee Brown at Crystal Lake, east of Goldshire"
            ],
            [
              "Dwarf, Gnome",
              "Dun Morogh",
              "Paxton Ganter at Iceflow Lake, on the way to Gnomeregan"
            ],
            [
              "Night Elf",
              "Teldrassil",
              "Androl Oakhand in Rut'theran Village, on the shore next to the dock"
            ],
            [
              "Skyborne",
              "Zephras Isle",
              "Baelann Swiftcurrent in Valanaar, Fenn Fairweather in Shen'dar Village"
            ]
          ],
          "notes": []
        }
      ]
    }
  }
];
