// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";
import eslint from "vite-plugin-eslint";

export default defineConfig({
  plugins: [solidPlugin(), eslint()],
  server: {
    port: 8091,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8090",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
  build: {
    target: "esnext",
  },
  clearScreen: false,
});
