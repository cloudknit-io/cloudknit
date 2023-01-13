# Changelog

## 2022-12-28 Move ArgoCD fields to API

### Created

* teams `orgs/:orgId/teams/:teamId` routing
* environments `orgs/:orgId/teams/:teamId/environments/:environmentId` routing

### Modified

* `costing` and `reconciliation` were moved to `/orgs/:orgId/teams/:teamId`

### Removed

* `costing/info/environment` endpoint
* `GET reconciliation/components`- use `component` controller instead
* `POST reconciliation/component/save` - all components are now saved via `POST /orgs/:orgId/team/:teamId/reconciliation/spec`
* `POST reconciliation/environment/save` - all environments are now saved via `POST /orgs/:orgId/team/:teamId/reconciliation/spec`
* `POST costing/savecomponent` - all components are now saved via `POST /orgs/:orgId/team/:teamId/reconciliation/spec`
* `GET costing/team/:name` - use `GET teams/:teamId?withCost=true`
* `GET costing/environment/:name` - use `GET teams/:teamId/environments/:environmentId?withCost=true`

## Multi-tenancy Changes

### Created

- All routes now start with one of the following:
    - `/v1/orgs`
    - `/v1/orgs/{orgId}`
    - `/v1/users`
    - `/v1/users/{userId}`
- added `created` to Organization table
- added `updated` to Users table

### Modified

- `company` table is now `organization`
- `/organization/oath` routes now live at `/v1/orgs/{orgId}/secrets/oauth`
- `/organization/github` routes now live at `/v1/orgs/{orgId}/secrets/github`
- `/auth/users/:orgid` is now `/v1/orgs/{orgId}/auth/users`
- all other routes are now prepended with `/v1/orgs/{orgId}/`
- `timeStamp` is now `updated` on Organization table
- `timeStamp` is now `created` on Users table
- `/costing/v1/api` is not just `/costing`
- `/reconciliation/v1/api` is not just `/reconciliation`
- secrets now operate on parameter paths that include the organization id - `/{orgId}/param/path`

### Removed
