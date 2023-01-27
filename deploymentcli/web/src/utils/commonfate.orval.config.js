module.exports = {
  commonfate: {
    output: {
      clean: true,
      mode: "tags-split",
      target: "./common-fate-client/orval.ts",
      schemas: "./common-fate-client/types",
      client: "swr",
      mock: true,
      override: {
        mutator: {
          path: "./custom-instance.ts",
          name: "customInstanceCommonfate",
        },
      },
    },
    input: {
      target: "../../../../openapi.yml",
    },
  },
};
