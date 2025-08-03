import * as cdk from 'aws-cdk-lib';
import { Template } from 'aws-cdk-lib/assertions';
import * as api from '../lib/stacks/api/api-stack';


test('SQS Queue Created', () => {
  const app = new cdk.App();

  const stack = new api.APIStack(app, 'MyTestStack');

  const template = Template.fromStack(stack);

  template.hasResourceProperties('AWS::SQS::Queue', {
    VisibilityTimeout: 300
  });
});
