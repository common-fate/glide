import * as cdkAssets from "cdk-assets";
import * as path from "path";
import { AwsClient, ConsoleProgress } from "./client";

const p = path.join(__dirname, "cdk.out", `Granted.assets.json`);
console.log({ path: p });
const am = cdkAssets.AssetManifest.fromFile(p);
const pb = new cdkAssets.AssetPublishing(am, {
  aws: new AwsClient(),
  publishInParallel: true,
  progressListener: new ConsoleProgress(),
});
pb.publish().then((r) => console.log("success"));
