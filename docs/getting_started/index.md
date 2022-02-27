# Getting Started

## Provisioning your first Environment

After you have setup zLifecycle, you will want to provision a simple environment to see how it works end to end.

1. Download this [env.yaml](../examples/first-environment.yaml) file
2. Commit and push this file to your team repository
3. Go to Environments page on zLifecycle UI
   * After a few mins you should see the new environment `checkout-dev` getting provisioned
   * It will start provisioning the `images` s3 bucket first and then `videos` s3 bucket
   * Once it starts provisioning, you can click on the `images` component and open the right panel
   * Right panel should show the terraform plan
4. Once the status changes to `Waiting For Approval` you will need to approve the changes by clicking on the `Approve` button below the terraform plan (as shown in the image below) to start provisoning the `images` s3 bucket (which is terraform apply)

![sample-right-panel](../assets/images/sample-right-panel.png "Sample Right Panel")

## Teardown your first Environment

After you have provisioned your first environment, let's go through the exercise of tearing it down.

1. At the spec level, add/update the 'teardown' flag to true 
    * The teardown will start at the bottom most leaf node
    * Approve the teardown plan when prompted
    * Monitor the progress on zLifecycle UI

