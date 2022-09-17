# Onboard a New Team

1. Copy the SSH url to your `zl-config` repo
    ![SSH URL](../assets/images/team-onboard-clone-url.png)
1. **In your `zl-config` repo** create a `teams` directory
1. In the `teams` directory create the following yaml file. Name it `[[team-name]].yaml`:
    ```yaml
    apiVersion: stable.compuzest.com/v1
    kind: Team
    metadata:
      name: [[teamname]]
      namespace: [[companyname]]-config
    spec:
      teamName: [[teamname]]
      configRepo:
        # Paste the SSH git URL
        source: [[git@github.com:org/repo.git]]
        path: "."
    ```
1. Once you commit and push the change it will register the team repo with **{{ company_name }}** and watch for any updates
1. The `zl-config` repo should resemble:
    ```
    root
    |   README.md
    |___teams
    |   |   team-name.yaml
    ```
