import { Construct } from "constructs";
import { App, TerraformStack } from "cdktf";
import * as google from '@cdktf/provider-google';

const project = 'ideal-pancake-380204';

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new google.provider.GoogleProvider(this, 'google', {
      project,
    });

    const workloadUser = new google.serviceAccount.ServiceAccount(this, 'workloadUser', {
      accountId: 'workloaduser',
    });

    new google.projectIamMember.ProjectIamMember(this, 'allowWorkloadIdentity' , {
      member: `serviceAccount:${workloadUser.email}`,
      project,
      role: 'roles/iam.workloadIdentityUser',
    });

    new google.projectIamMember.ProjectIamMember(this, 'allowPubSub' , {
      member: `serviceAccount:${workloadUser.email}`,
      project,
      role: 'roles/pubsub.editor',
    });

    new google.projectIamMember.ProjectIamMember(this, 'allowFireStore' , {
      member: `serviceAccount:${workloadUser.email}`,
      project,
      role: 'roles/datastore.user',
    });

    new google.pubsubTopic.PubsubTopic(this, 'exampleTopic', {
      name: 'example',
    });

  }
}

const app = new App();
new MyStack(app, "ideal-pancake");
app.synth();
