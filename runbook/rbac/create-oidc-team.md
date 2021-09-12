# Create OIDC team

## Overview

zLifecycle supports SSO using ArgoCD's Dex server

## When to use this runbook

When you want to add a new OIDC team to sign in using SSO with RBAC

## Prerequisites

1. Admin access to your GitHub Organisation profile

## Initial Steps Overview

- [Create a GitHub team](#create-a-github-team)
- [Create a zLifecycle team](#create-a-zlifecycle-team)

## Detailed Steps

### Create a GitHub team
1. Go to your Organization's GitHub profile and select Teams setting
2. Create a team and assign people to it

### Create a zLifecycle team
1. Create a Team Custom Resource and name it EXACTLY like the GitHub team
2. Commit and push it to your config repo
