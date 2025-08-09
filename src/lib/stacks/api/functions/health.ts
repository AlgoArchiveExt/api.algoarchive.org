import { aws_lambda, Duration } from 'aws-cdk-lib';
import { Code } from 'aws-cdk-lib/aws-lambda';
import { NodejsFunction } from 'aws-cdk-lib/aws-lambda-nodejs';
import { Construct } from 'constructs';


export const healthFunction = (stack: Construct) => {
  return new NodejsFunction(stack, 'Health',{
    description: 'Health check function for the API',
    code: Code.fromAsset('dist/lambda/health'),
    handler: 'index.handler',
    runtime: aws_lambda.Runtime.NODEJS_20_X,
    timeout: Duration.seconds(15),
    memorySize: 128,
  });
}