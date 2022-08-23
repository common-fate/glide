module.exports = {
  env: {
    browser: true,
    es2020: true,
    node: true,
  },
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:prettier/recommended",
  ],
  parser: "@typescript-eslint/parser",
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 11,
    sourceType: "module",
    project: "./tsconfig.json",
  },
  plugins: ["@typescript-eslint", "jest"],
  rules: {
    "@typescript-eslint/await-thenable": "error",
    "@typescript-eslint/no-floating-promises": "error",
    /**
     * I've disabled this due to issues with
     * PrismJS in PolicyDiff and its module resolution.
     * Very hard to work around with type safety
     *
     * Github permalink:
     * https://github.com/Exponent-Labs/iamzero-enterprise/blob/fae2999ce6b883be89b345964bb2e31de2327fe7/web/src/pages/findingDetails/PolicyDiff.tsx#L12-L17
     */
    "@typescript-eslint/ban-ts-comment": "off",
    "@typescript-eslint/explicit-module-boundary-types": "off",
    "@typescript-eslint/no-explicit-any": "off",
  },
};