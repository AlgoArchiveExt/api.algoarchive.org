import { aws_lambda, aws_lambda_nodejs, Duration } from 'aws-cdk-lib';
import { Code } from 'aws-cdk-lib/aws-lambda';
import { NodejsFunction } from 'aws-cdk-lib/aws-lambda-nodejs';
import { Construct } from 'constructs';

export const commitFunction = (stack: Construct) => 
  new NodejsFunction(stack, 'Commit', {
    description: 'Function to handle committing LC solutions to GitHub',
    code: Code.fromAsset('dist/lambda/commit'),
    handler: 'index.handler',
    runtime: aws_lambda.Runtime.NODEJS_20_X,
    memorySize: 1024,
    timeout: Duration.seconds(30),
});

