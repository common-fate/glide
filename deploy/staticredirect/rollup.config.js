import typescript from "@rollup/plugin-typescript";

/**
 * @type {import('rollup').RollupOptions}
 */
const config = {
  input: "src/index.ts",
  output: {
    file: "build/index.js",
    format: "cjs",
  },
  treeshake: false,
  plugins: [typescript()],
};

export default config;
