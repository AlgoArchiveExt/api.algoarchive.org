import { Construct } from 'constructs';
import { healthFunction } from './functions/health';
import { commitFunction } from './functions/commit';
import { aws_apigatewayv2, aws_apigatewayv2_integrations, Stack, StackProps } from 'aws-cdk-lib';

export class APIStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
    
    const healthFunc = healthFunction(this);
    const healthIntegration = new aws_apigatewayv2_integrations.HttpLambdaIntegration('HealthIntegration', healthFunc);
    
    const api = new aws_apigatewayv2.HttpApi(this, 'AlgoArchive API', {
      apiName: 'AlgoArchiveAPI',
      description: 'API for AlgoArchive',
      createDefaultStage: true,
      defaultIntegration: healthIntegration,
    });


    const commitFunc = commitFunction(this);
    const commitIntegration = new aws_apigatewayv2_integrations.HttpLambdaIntegration('CommitIntegration', commitFunc);

    api.addRoutes({
      path: '/v1/solutions/commits',
      methods: [aws_apigatewayv2.HttpMethod.POST],
      integration: commitIntegration,
    });
  }
}
