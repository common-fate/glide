/**
 * Generated by orval v6.10.3 🍺
 * Do not edit manually.
 * Common Fate
 * Common Fate API
 * OpenAPI spec version: 1.0
 */
import type { TargetGroupDiagnostic } from './targetGroupDiagnostic';
import type { TargetGroupDeploymentActiveConfig } from './targetGroupDeploymentActiveConfig';
import type { TargetGroupDeploymentProvider } from './targetGroupDeploymentProvider';

export interface TargetGroupDeployment {
  id: string;
  functionArn: string;
  awsAccount: string;
  healthy: boolean;
  diagnostics: TargetGroupDiagnostic[];
  activeConfig: TargetGroupDeploymentActiveConfig;
  provider: TargetGroupDeploymentProvider;
}