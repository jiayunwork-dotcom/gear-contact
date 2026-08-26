"use strict";
const $ = (id) => document.getElementById(id);
function showError(msg) { $("err-text").textContent = msg; $("err-panel").hidden = false; }
function hideError() { $("err-panel").hidden = true; }
function fill(rows) {
  const tb = $("out-table"); tb.innerHTML = "";
  for (const [k,v] of rows) { const tr = document.createElement("tr"); tr.innerHTML = "<td>"+k+"</td><td>"+v+"</td>"; tb.appendChild(tr); }
}
async function loadExample() {
  hideError();
  try {
    const resp = await fetch("/api/example"); const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "示例失败");
    $("body").value = JSON.stringify(data, null, 2);
    $("hint").textContent = "已加载 spur-20 算例。";
  } catch (e) { showError(String(e)); }
}
async function runMesh() {
  hideError();
  try {
    const resp = await fetch("/api/mesh", { method:"POST", headers:{"Content-Type":"application/json"}, body: $("body").value });
    const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "HTTP "+resp.status);
    $("out-panel").hidden = false;
    fill([
      ["中心距", data.center_distance.toPrecision(6)],
      ["法向齿距", data.base_pitch.toPrecision(6)],
      ["作用线长度", data.line_length.toPrecision(6)],
      ["重合度", data.contact_ratio.toPrecision(6)],
      ["顶隙", data.clearance.toPrecision(6)]
    ]);
  } catch (e) { showError(String(e)); }
}
$("btn-example").addEventListener("click", loadExample);
$("btn-run").addEventListener("click", runMesh);
loadExample();
