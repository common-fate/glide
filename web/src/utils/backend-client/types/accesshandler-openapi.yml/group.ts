/**
 * Generated by orval v6.9.6 🍺
 * Do not edit manually.
 * Approvals
 * Granted Approvals API
 * OpenAPI spec version: 1.0
 */
import type { GroupOption } from '../accesshandler-openapi.yml/groupOption';

export interface Group {
  id: string;
  title: string;
  description?: string;
  options: GroupOption[];
}