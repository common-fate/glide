/**
 * Generated by orval v6.7.1 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import type { TargetGroupSchema } from './targetGroupSchema';
import type { TargetGroupFrom } from './targetGroupFrom';

export interface TargetGroup {
  id: string;
  schema: TargetGroupSchema;
  from: TargetGroupFrom;
  icon: string;
  createdAt?: string;
  updatedAt?: string;
}
