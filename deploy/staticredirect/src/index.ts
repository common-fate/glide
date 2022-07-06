import { CloudFrontFunctionsEvent } from "aws-lambda";
import { fixURI } from "./uri";

function handler(event: CloudFrontFunctionsEvent) {
  var request = event.request;
  request.uri = fixURI(request.uri);
  return request;
}
