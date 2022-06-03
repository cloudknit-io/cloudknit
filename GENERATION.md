# Operator ZL IL Generation

**Operator** generates [Argo Application](https://github.com/argoproj/argo-cd/blob/master/docs/operator-manual/application.yaml) yaml files which are then consumed by Argo and applied to the cluster. This doc describes how those files are generated.

## Initial Trigger

A git push to a customer's Team config repo triggers an [Argo CD Sync](https://argo-cd.readthedocs.io/en/stable/core_concepts/) (via Argo Watcher) which creates 3 k8s custom resources:


(i dont think this is correct. i think Company and Team are generated manually during onboarding)

- Company
- Team
- Environment

> NOTE: These resources are extracted from customer environment yamls

Each of these objects are created in their respective `[company]-config` k8s namespace which in turn trigger their respective [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) which live within [zLifecycle's Operator](https://github.com/CompuZest/zlifecycle-il-operator/).

## Operator Controllers

Once the Controllers managed by the Operator wake up they consume the incoming k8s objects, pull down the ZL IL repo, and start generating Argo Application yaml files based on the incoming k8s objects.

All Argo Application yamls will live in its own customer ZL IL repo that **zLifecycle** manages.

> NOTE: Customers do not have access to their IL repos

### Company

TODO

### Team

TODO

### Environment

TODO
