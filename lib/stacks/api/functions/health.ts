import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';

export const healthFunction = (stack: Construct) => new cdk.aws_lambda_nodejs.NodejsFunction(stack, 'AlgoArchiveHealth',{
  entry: 'lambdas/health/index.ts',
  handler: 'handler',
  runtime: cdk.aws_lambda.Runtime.NODEJS_22_X,
  timeout: cdk.Duration.seconds(15),
  memorySize: 128,
});