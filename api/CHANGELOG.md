# Changelog

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
