module.exports = {
  commonfate: {
    output: {
      clean: true,
      mode: "single",
      target: "./registry-client/orval.ts",
      client: "swr",
      mock: true,
      override: {
        mutator: {
          path: "./custom-instance.ts",
          name: "customInstanceRegistry",
        },
      },
    },
    input: {
      target:
        "https://raw.githubusercontent.com/common-fate/provider-registry-sdk-go/main/openapi.yml",
    },
  },
};
