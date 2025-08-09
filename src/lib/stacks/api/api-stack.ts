import { Construct } from 'constructs';
import { healthFunction } from './functions/health';
import { commitFunction } from './functions/commit';
import { aws_apigatewayv2, aws_apigatewayv2_integrations, Stack, StackProps } from 'aws-cdk-lib';
import { CfnStage } from 'aws-cdk-lib/aws-apigatewayv2';

export class APIStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
    
    
    const api = new aws_apigatewayv2.HttpApi(this, 'AlgoArchive API', {
      apiName: 'AlgoArchiveAPI',
      description: 'API for AlgoArchive',
      createDefaultStage: false,
    });
    
    const healthFunc = healthFunction(this);
    const healthIntegration = new aws_apigatewayv2_integrations.HttpLambdaIntegration('HealthIntegration', healthFunc);

    api.addRoutes({
      path: '/health',
      methods: [aws_apigatewayv2.HttpMethod.GET],
      integration: healthIntegration,
    });

    const commitFunc = commitFunction(this);
    const commitIntegration = new aws_apigatewayv2_integrations.HttpLambdaIntegration('CommitIntegration', commitFunc);

    
    api.addRoutes({
      path: '/solutions/commits',
      methods: [aws_apigatewayv2.HttpMethod.POST],
      integration: commitIntegration,

    });

    new CfnStage(this, "DefaultStage", {
      apiId: api.apiId,
      stageName: "v1",
      description: "Default stage with throttling settings",
      autoDeploy: true,
      defaultRouteSettings: {
        throttlingBurstLimit: 20,
        throttlingRateLimit: 10
      }
    });
  }
}
