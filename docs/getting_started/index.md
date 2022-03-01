# Getting Started

## Set AWS Credentials

Make sure you set [AWS Credentials](../settings/aws_credentials.md) that you want to use for provisioning your environment.

## Provision your first Environment

After you have setup AWS Credentials, you will want to provision a simple environment to see how it works end to end.

1. Clone your team repo locally on your machine
2. Go to root directory of the team repo and create a new folder called `hello-world`
3. Run below bash script using terminal and enter Company, Team and Environment Names when asked.

```
bash <(curl -s https://zlifecycle.github.io/docs/scripts/getting_started.sh)
```

4. Commit and push this file to your team repository
5. Go to Environments page on zLifecycle UI
   * After a few mins you should see the new environment with the team & env name you entered getting provisioned
   * It will start provisioning the `images` s3 bucket first and then `videos` s3 bucket
   * Once it starts provisioning, you can click on the `images` component and open the right panel
   * Right panel should show the terraform plan
6. Once the status changes to `Waiting For Approval` you will need to approve the changes by clicking on the `Approve` button below the terraform plan (as shown in the image below) to start provisoning the `images` s3 bucket (which is terraform apply)

![sample-right-panel](../assets/images/sample-right-panel.png "Sample Right Panel")

## Teardown your first Environment

After you have provisioned your first environment, let's go through the teardown exercise.

1. In the `hello-world.yaml` that you created in the Provision step above add a 'teardown' flag to `true` at the spec level
2. Commit and push changes to your team repository
    * The teardown will start at the bottom most leaf node
3. Approve the teardown plan when prompted
4. Monitor the progress on zLifecycle UI
