import { StackOutputs } from "../lib/helpers/outputs";
import { writeFileSync } from "fs";

// testOutputs will have a type error when new fields are added to stack outputs
// a new entry should be added here as this is used in a test to ensure consistency with the corresponding go type
// in pkg/deploy/output.go
function generateDefaultOutputs(outputs: StackOutputs): StackOutputs {
  const defaultOutputs: StackOutputs = {} as StackOutputs;

  Object.keys(outputs).forEach((key) => {
    defaultOutputs[key as keyof StackOutputs] = "abcdefg";
  });

  return defaultOutputs;
}
const defaultOutputs = generateDefaultOutputs({} as StackOutputs);
// Write the json object to ./testOutputs.json so that it can be parsed by a go test in pkg/deploy.output_test.go
writeFileSync("./testOutputs.json", JSON.stringify(defaultOutputs));
