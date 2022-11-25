module.exports = {
  commonfate: {
    output: {
      clean: true,
      mode: "tags-split",
      target: "./backend-client/orval.ts",
      schemas: "./backend-client/types",
      client: "swr",
      mock: true,
      override: {
        mutator: {
          path: "./custom-instance.ts",
          name: "customInstance",
        },
      },
    },
    input: {
      target: "../../../openapi.yml",
    },
  },
};
