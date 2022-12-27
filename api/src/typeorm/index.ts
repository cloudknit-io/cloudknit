import { Component } from './component.entity';
import { User } from './User.entity';
import { Organization } from './Organization.entity';
import { ComponentReconcile } from './component-reconcile.entity';
import { EnvironmentReconcile } from './environment-reconcile.entity';
import { Environment } from './environment.entity';

export const entities = [
    User,
    Organization,
    Component, 
    ComponentReconcile,
    EnvironmentReconcile,
    Environment,
];

export {
    User,
    Organization
}
