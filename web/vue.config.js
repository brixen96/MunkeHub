const { defineConfig } = require("@vue/cli-service");

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8081,
    proxy: {
      // Dev only: forward API calls to the local API (or the live server)
      "/api": {
        target: process.env.MUNKEHUB_API || "http://localhost:4711",
        changeOrigin: true,
      },
      // Footer version check hits /health directly (nginx does this in prod)
      "/health": {
        target: process.env.MUNKEHUB_API || "http://localhost:4711",
        changeOrigin: true,
      },
    },
  },
});