/**
 * Generated by orval v6.10.3 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import type { ProviderV2Status } from './providerV2Status';

/**
 * ProviderV2
 */
export interface ProviderV2 {
  name: string;
  team: string;
  version: string;
  stackId: string;
  status: ProviderV2Status;
  type: string;
  id: string;
}
