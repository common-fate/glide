module.exports = {
  commonfate: {
    output: {
      clean: true,
      mode: "single",
      target: "./backend-client/local/orval.ts",
      client: "swr",
      mock: true,
      override: {
        mutator: {
          path: "./custom-instance.ts",
          name: "customInstanceLocal",
        },
      },
    },
    input: {
      target: "../../../openapi.yml",
    },
  },
};
