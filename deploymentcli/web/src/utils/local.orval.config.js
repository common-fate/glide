module.exports = {
  commonfate: {
    output: {
      clean: true,
      mode: "tags-split",
      target: "./local-client/orval.ts",
      schemas: "./local-client/types",
      client: "swr",
      override: {
        mutator: {
          path: "./custom-instance.ts",
          name: "customInstanceLocal",
        },
      },
    },
    input: {
      target: "../../../../deploymentcli.openapi.yml",
    },
  },
};
