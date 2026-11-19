// One-off: inspect the gathering guides' shape (zones for the farming table).
import { GUIDES_EXTRA } from "../web/src/guides-extra.js";

const pick = GUIDES_EXTRA.filter((x) => x.gathering).map((x) => ({
  id: x.id,
  name: x.name,
  keys: Object.keys(x),
  gkeys: x.gathering[0] ? Object.keys(x.gathering[0]) : null,
  tiers: x.gathering.length,
  zones: x.gathering[0] ? (x.gathering[0].zones || []).slice(0, 4) : null,
  notes: (x.notes || []).slice(0, 2),
}));
console.log(JSON.stringify(pick, null, 1).slice(0, 3000));