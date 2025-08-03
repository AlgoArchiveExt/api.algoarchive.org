import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { healthFunction } from './functions/health';

export class APIStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    
    const healthFunc = healthFunction(this);
    const healthIntegration = new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration('HealthIntegration', healthFunc);
    
    const api = new cdk.aws_apigatewayv2.HttpApi(this, 'AlgoArchive API', {
      apiName: 'AlgoArchiveAPI',
      description: 'API for AlgoArchive',
      createDefaultStage: true,
      defaultIntegration: healthIntegration,
    });

  }
}
