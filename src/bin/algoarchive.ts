#!/usr/bin/env node
import { App } from 'aws-cdk-lib';
import { APIStack } from '../lib/stacks/api/api-stack';

const app = new App();
new APIStack(app, 'AA-APIStack', {
  /* If you don't specify 'env', this stack will be environment-agnostic.
   * Account/Region-dependent features and context lookups will not work,
   * but a single synthesized template can be deployed anywhere. */

  /* Uncomment the next line to specialize this stack for the AWS Account
   * and Region that are implied by the current CLI configuration. */
  env: { account: process.env.CDK_DEFAULT_ACCOUNT, region: process.env.CDK_DEFAULT_REGION },
  description: 'API Stack for AlgoArchive. Defines the API Gateway and Lambda functions for its API.',
});