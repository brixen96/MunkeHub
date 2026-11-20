// Reproduce the file:/// console error in real Brave via raw CDP (no deps).
// Usage: node scripts/test-brave-console.mjs   (Brave must run with --remote-debugging-port=9333)
const list = await (await fetch("http://127.0.0.1:9333/json/list")).json();
const target = list.find((t) => t.type === "page" && t.url.includes("8081"))
  || list.find((t) => t.type === "page");
if (!target) throw new Error("no page target");

const errors = [];
const fileReqs = [];
const ws = new WebSocket(target.webSocketDebuggerUrl);
let id = 0;
const send = (method, params = {}) =>
  ws.send(JSON.stringify({ id: ++id, method, params }));

await new Promise((res) => { ws.onopen = res; });

ws.onmessage = (ev) => {
  const m = JSON.parse(ev.data);
  if (m.method === "Runtime.consoleAPICalled" && m.params.type === "error") {
    errors.push(m.params.args.map((a) => a.value ?? a.description ?? "").join(" "));
  } else if (m.method === "Runtime.exceptionThrown") {
    errors.push(m.params.exceptionDetails.text + " " +
      (m.params.exceptionDetails.exception?.description || ""));
  } else if (m.method === "Log.entryAdded" && m.params.entry.level === "error") {
    errors.push(m.params.entry.text + " " + (m.params.entry.url || ""));
  } else if (m.method === "Network.requestWillBeSent") {
    if (m.params.request.url.startsWith("file://")) fileReqs.push(m.params.request.url);
  } else if (m.method === "Network.responseReceived") {
    if (m.params.response.url.startsWith("file://")) fileReqs.push(m.params.response.url);
  }
};

send("Runtime.enable");
send("Log.enable");
send("Network.enable");
send("Page.enable");
send("Page.reload");
console.log("reloading...");

await new Promise((r) => setTimeout(r, 6000));

// expand first roster row (tooltips + 3D viewer)
send("Runtime.evaluate", {
  expression: `(() => { const r = document.querySelector('table tbody tr'); if (r) r.click(); return 'clicked'; })()`,
});
await new Promise((r) => setTimeout(r, 4000));

// hover gear links to trigger wowhead tooltips.js
send("Runtime.evaluate", {
  expression: `(() => {
    const links = [...document.querySelectorAll('a[href*="wowhead.com/forever/item"]')];
    links.slice(0, 5).forEach(a => {
      for (const t of ['mouseover', 'mouseenter']) {
        a.dispatchEvent(new MouseEvent(t, { bubbles: true, clientX: 100, clientY: 100 }));
      }
    });
    return links.length;
  })()`,
});
await new Promise((r) => setTimeout(r, 5000));

console.log("URL:", target.url);
console.log("file:// requests:", JSON.stringify([...new Set(fileReqs)], null, 2));
console.log("console errors:", JSON.stringify(errors, null, 2));
const fileErrs = errors.filter((e) => /file:|local resource/i.test(e));
console.log("FILE ERRORS:", fileErrs.length ? JSON.stringify(fileErrs, null, 2) : "none");

ws.close();