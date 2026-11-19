import { createApp } from "vue";
import jQuery from "jquery";
import App from "./App.vue";

// The wowhead 3D viewer (ZamModelViewer) expects global jQuery — it
// registers its binary ajax transport at script-eval time using bare `$`.
// Same setup as armorywowforever.com (tooltips.js loads first in
// index.html, then jQuery, exactly like their page order).
window.jQuery = window.$ = jQuery;

createApp(App).mount("#app");