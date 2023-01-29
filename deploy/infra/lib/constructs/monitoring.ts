import { Duration } from "aws-cdk-lib";
import * as cloudwatch from "aws-cdk-lib/aws-cloudwatch";
import { Construct } from "constructs";
import { CFService } from "../helpers/service";

interface Props {
  deploymentSuffix: string;
  services: CFService[];
}

export class Monitoring extends Construct {
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    const health = new cloudwatch.Dashboard(
      this,
      "CommonFateDeploymentHealth",
      { dashboardName: "CommonFateHealth" + props.deploymentSuffix }
    );

    // add a row for each service
    props.services.forEach((s) => {
      health.addWidgets(
        new cloudwatch.TextWidget({
          height: 2,
          markdown: `# ${s.label}`,
        })
      );

      health.addWidgets(
        new cloudwatch.TextWidget({
          height: 6,
          markdown: `${s.description}

**Failure Impact:** ${s.failureImpact}
`,
        }),
        new cloudwatch.GraphWidget({
          height: 6,
          title: `${s.label} Errors`,
          left: [
            new cloudwatch.Metric(
              s.function.metricErrors({ period: Duration.hours(24) })
            ),
          ],
        }),
        new cloudwatch.GraphWidget({
          height: 6,
          title: `${s.label} Execution Duration`,
          left: [
            new cloudwatch.Metric(
              s.function.metricDuration({ period: Duration.hours(24) })
            ),
          ],
        }),
        new cloudwatch.GraphWidget({
          height: 6,
          title: `${s.label} Invocations`,
          left: [
            new cloudwatch.Metric(
              s.function.metricInvocations({ period: Duration.hours(24) })
            ),
          ],
        })
      );
    });
  }
}
