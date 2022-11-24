/**
 * Generated by orval v6.9.6 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import type { ProviderSetupStatus } from './providerSetupStatus';
import type { ProviderSetupStepOverview } from './providerSetupStepOverview';
import type { ProviderSetupConfigValues } from './providerSetupConfigValues';
import type { ProviderConfigValidation } from './accesshandler-openapi.yml/providerConfigValidation';

/**
 * A provider in the process of being set up through the guided setup workflow in Common Fate. These providers are **not** yet active.
 */
export interface ProviderSetup {
  /** A unique ID for the provider setup. This is a random KSUID to avoid potential conflicts with user-specified provider IDs in the `granted-deployment.yml` file. */
  id: string;
  /** The type of the Access Provider being set up. */
  type: string;
  /** The version of the provider. */
  version: string;
  /** The status of the setup process. */
  status: ProviderSetupStatus;
  /** An overview of the steps indicating whether they are complete. */
  steps: ProviderSetupStepOverview[];
  /** The current configuration values. */
  configValues: ProviderSetupConfigValues;
  configValidation: ProviderConfigValidation[];
}
