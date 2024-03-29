/**
 * Generated by orval v6.10.3 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import type { CreateRequestWithSubRequest } from './createRequestWithSubRequest';
import type { RequestTiming } from './requestTiming';

/**
 * Detailed object for a Favorite. 
 */
export interface FavoriteDetail {
  id: string;
  name: string;
  with: CreateRequestWithSubRequest;
  reason?: string;
  timing: RequestTiming;
}
