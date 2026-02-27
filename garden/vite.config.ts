import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";
import devtools from "solid-devtools/vite";
import { copyFileSync } from "node:fs";
import { resolve, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  plugins: [
    devtools(),
    solidPlugin(),
    {
      name: "copy-not-found",
      closeBundle() {
        const dist = resolve(__dirname, "dist");
        copyFileSync(resolve(dist, "index.html"), resolve(dist, "not_found.html"));
      },
    },
  ],
  server: {
    port: 5173,
  },
  build: {
    target: "esnext",
  },
});
