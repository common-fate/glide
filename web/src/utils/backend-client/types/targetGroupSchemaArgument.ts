/**
 * Generated by orval v6.7.1 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import type { TargetGroupSchemaArgumentResourceSchema } from './targetGroupSchemaArgumentResourceSchema';

/**
 * Define the metadata, data type and UI elements for the argument
 */
export interface TargetGroupSchemaArgument {
  id: string;
  title: string;
  description?: string;
  resourceSchema?: TargetGroupSchemaArgumentResourceSchema;
  resource?: string;
}
