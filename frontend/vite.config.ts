import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { resolve } from "node:path";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  base: "/",
  root: "src",
  publicDir: resolve(__dirname, "public"),
  plugins: [react(), tailwindcss()],
  build: {
    outDir: resolve(__dirname, "dist"),
    emptyOutDir: true,
    rollupOptions: {
      input: {
        "": resolve(__dirname, "src/index.html"),
      },
      output: {
        assetFileNames: "assets/bundle.js",
      },
    },
  },
});
