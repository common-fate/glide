/**
 * Generated by orval v6.7.1 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import type { Diagnostic } from './diagnostic';

/**
 * Handler represents a deployment of a provider. 
Handlers can be linked to target groups via routes
 */
export interface TGHandler {
  id: string;
  runtime: string;
  functionArn: string;
  awsAccount: string;
  awsRegion: string;
  healthy: boolean;
  diagnostics: Diagnostic[];
}
