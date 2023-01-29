import { Function } from "aws-cdk-lib/aws-lambda";
import { AccessHandler } from "../constructs/access-handler";
import { AppBackend } from "../constructs/app-backend";

export interface CFService {
  id: string;
  /** human-readable label for the service */
  label: string;
  /** What happens if the service is unavailable? */
  failureImpact: string;
  /** A short description for what the service does */
  description: string;
  function: Function;
}

/**
 * returns the registry of all services involved in a
 * Common Fate deployment.
 *
 * Used between both dev and prod stacks so we don't forget
 * to register a service in prod.
 */
export const getServices = ({
  appBackend,
  accessHandler,
}: {
  appBackend: AppBackend;
  accessHandler: AccessHandler;
}): CFService[] => {
  return [...appBackend.getServices(), accessHandler.getService()];
};
