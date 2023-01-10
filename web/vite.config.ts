import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import checker from "vite-plugin-checker";
import rollupNodePolyFill from "rollup-plugin-node-polyfills";

export default defineConfig({
  plugins: [react(), checker({ typescript: true })],
  optimizeDeps: {
    esbuildOptions: {
      // Node.js global to browser globalThis
      define: {
        global: "globalThis", //<-- AWS SDK
      },
    },
  },
  build: {
    rollupOptions: {
      plugins: [
        // Enable rollup polyfills plugin
        // used during production bundling
        // @ts-ignore
        rollupNodePolyFill(),
      ],
    },
  },
  // https://github.com/aws/aws-sdk-js/issues/3673
  resolve: {
    alias: {
      "./runtimeConfig": "./runtimeConfig.browser",
    },
  },
  server: {
    host: "localhost",
    port: 3000,
  },
});
