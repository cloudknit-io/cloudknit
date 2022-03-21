# Teardown

You can teardown an entire environment using one of the options below:

## Keep Environment Config in repo

This option is helpful if you want to easily re-provision the environment at a later time.

1. Update/Add `teardown` flag with value `true` 
2. Commit & Push your changes to the repo
3. The environment components will start destroying one by one (monitor progress on the zLifecycle UI)

## Remove Environment Config from repo

1. Delete the Environment config (yaml, tfvars, tf files etc)
2. Commit & Push changes to the repo
3. The environment components will start destroying one by one (monitor progress on the zLifecycle UI) 
