module.exports = {
  commonfate: {
    output: {
      clean: true,
      mode: "single",
      target: "./backend-client/registry/orval.ts",
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
      target: "../../../../openapi.yml",
    },
  },
};
