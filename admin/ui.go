package main

import "html/template"

var page = template.Must(template.New("admin").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>MunkeHub Admin</title>
<style>
  :root { color-scheme: dark; }
  * { box-sizing: border-box; }
  body {
    margin: 0; font-family: "Segoe UI", system-ui, sans-serif;
    background: radial-gradient(ellipse at top, #23212a 0%, #131217 70%) fixed;
    color: #dcd7cc;
  }
  .wrap { max-width: 720px; margin: 0 auto; padding: 2rem 1rem; }
  h1 { color: #ffd166; font-family: Georgia, serif; margin: 0 0 .25rem; }
  .sub { color: #8f8a7e; margin: 0 0 1.5rem; font-size: .9rem; }
  .panel {
    background: #1c1b1f; border: 1px solid #3a3840; border-radius: 10px;
    padding: 1.2rem 1.3rem; margin-bottom: 1.2rem;
    box-shadow: 0 6px 18px rgba(0,0,0,.4);
  }
  h2 {
    margin: 0 0 1rem; font-family: Georgia, serif; font-size: 1.05rem;
    color: #c8a44d; text-transform: uppercase; letter-spacing: .08em;
    border-bottom: 1px solid #3a3840; padding-bottom: .6rem;
  }
  label { display: block; font-size: .8rem; color: #8f8a7e; margin: .6rem 0 .2rem; text-transform: uppercase; }
  input {
    width: 100%; background: #131217; border: 1px solid #3a3840; border-radius: 6px;
    color: #dcd7cc; padding: .55rem .7rem; font: inherit;
  }
  input:focus { outline: none; border-color: #c8a44d; }
  button {
    background: linear-gradient(180deg, #d9b45c, #a8842f); border: 1px solid #c8a44d;
    border-radius: 6px; color: #1a1508; font-weight: 700; padding: .55rem 1.3rem;
    cursor: pointer; font: inherit; margin-top: .9rem;
  }
  button:hover { filter: brightness(1.1); }
  button.danger {
    background: linear-gradient(180deg, #b5483f, #7e2a24); border-color: #b5483f;
    color: #fff; padding: .3rem .8rem; font-size: .8rem; margin-top: 0;
  }
  table { width: 100%; border-collapse: collapse; font-size: .92rem; }
  th, td { text-align: left; padding: .5rem .6rem; border-bottom: 1px solid #3a3840; }
  th { color: #8f8a7e; font-size: .75rem; text-transform: uppercase; }
  .label { color: #ffd166; font-weight: 600; }
  .muted { color: #8f8a7e; font-size: .85rem; }
  .error { color: #ff6b6b; font-size: .85rem; margin-top: .5rem; }
  .ok { color: #7ec97e; font-size: .85rem; margin-top: .5rem; }
  .token {
    font-family: Consolas, monospace; background: #131217; border: 1px solid #3a3840;
    border-radius: 4px; padding: .1rem .4rem; color: #ffd166; font-size: .85rem;
  }
  .mint-result {
    background: #131217; border: 1px dashed #3a3840; border-radius: 6px;
    padding: .7rem .9rem; margin-top: .8rem; font-size: .88rem;
  }
</style>
</head>
<body>
<div class="wrap">
  <h1>🐒 MunkeHub Admin</h1>
  <p class="sub">Token management for Munkeklubben — local only, talks to {{.API}}</p>

  {{if not .Configured}}
  <div class="panel">
    <h2>Setup</h2>
    <label>API address (LAN)</label>
    <input id="cfg-api" value="{{.API}}" placeholder="http://brixmedia.local:4711">
    <label>Admin token (MUNKEHUB_ADMIN_TOKEN)</label>
    <input id="cfg-token" type="password" placeholder="paste your admin token">
    <button onclick="saveCfg()">Save & verify</button>
    <div id="cfg-msg"></div>
  </div>
  {{end}}

  <div class="panel">
    <h2>Mint a token</h2>
    <label>Member label (their roster name)</label>
    <input id="mint-label" placeholder="e.g. sven" onkeydown="if(event.key==='Enter')mint()">
    <button onclick="mint()">Mint token</button>
    <div id="mint-msg"></div>
    <div id="mint-result"></div>
  </div>

  <div class="panel">
    <h2>Minted tokens</h2>
    {{if .APIError}}<div class="error">{{.APIError}}</div>{{end}}
    <table>
      <thead><tr><th>Label</th><th>Characters</th><th>Created</th><th></th></tr></thead>
      <tbody id="tok-body"></tbody>
    </table>
    <p class="muted" id="tok-empty" style="display:none">No tokens minted yet.</p>
  </div>
</div>

<script>
const TOKENS = {{.Tokens}};

function render() {
  const body = document.getElementById('tok-body');
  const empty = document.getElementById('tok-empty');
  body.innerHTML = '';
  if (!TOKENS.length) { empty.style.display = 'block'; return; }
  empty.style.display = 'none';
  for (const t of TOKENS) {
    const tr = document.createElement('tr');
    tr.innerHTML = '<td class="label"></td><td></td><td class="muted"></td>' +
      '<td><button class="danger" onclick="revoke(' + t.id + ')">Revoke</button></td>';
    tr.children[0].textContent = t.label;
    tr.children[1].textContent = t.characters;
    tr.children[2].textContent = (t.created_at || '').slice(0, 10);
    body.appendChild(tr);
  }
}

async function saveCfg() {
  const msg = document.getElementById('cfg-msg');
  msg.className = '';
  const res = await fetch('/save', {
    method: 'POST', headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({
      api: document.getElementById('cfg-api').value,
      adminToken: document.getElementById('cfg-token').value,
    }),
  });
  if (res.ok) { location.reload(); }
  else { msg.className = 'error'; msg.textContent = await res.text(); }
}

async function mint() {
  const label = document.getElementById('mint-label').value.trim();
  const msg = document.getElementById('mint-msg');
  const out = document.getElementById('mint-result');
  msg.className = ''; out.innerHTML = '';
  if (!label) { msg.className = 'error'; msg.textContent = 'Enter a label.'; return; }
  const res = await fetch('/mint', {
    method: 'POST', headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({label}),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) { msg.className = 'error'; msg.textContent = data.error || ('HTTP ' + res.status); return; }
  out.className = 'mint-result';
  out.innerHTML = 'Token for <b class="label"></b> — send this to them (or use pack-guild):<br>' +
    '<span class="token" id="newtok"></span> ' +
    '<button class="danger" onclick="copyTok()">Copy</button>';
  out.querySelector('.label').textContent = data.label;
  document.getElementById('newtok').textContent = data.token;
  document.getElementById('mint-label').value = '';
  TOKENS.push({id: 0, label: data.label, characters: 0, created_at: new Date().toISOString()});
  render();
}

function copyTok() {
  navigator.clipboard.writeText(document.getElementById('newtok').textContent);
}

async function revoke(id) {
  if (!confirm('Revoke this token? Their helper will stop syncing.')) return;
  const res = await fetch('/revoke/' + id, {method: 'DELETE'});
  if (res.ok) location.reload();
}

render();
</script>
</body>
</html>`))
