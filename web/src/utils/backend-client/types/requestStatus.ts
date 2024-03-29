/**
 * Generated by orval v6.10.3 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */

/**
 * The status of an Access Request.

 */
export type RequestStatus = typeof RequestStatus[keyof typeof RequestStatus];


// eslint-disable-next-line @typescript-eslint/no-redeclare
export const RequestStatus = {
  APPROVED: 'APPROVED',
  PENDING: 'PENDING',
  CANCELLED: 'CANCELLED',
  DECLINED: 'DECLINED',
} as const;
