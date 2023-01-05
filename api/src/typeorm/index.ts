import { Component } from './component.entity';
import { User } from './User.entity';
import { Organization } from './Organization.entity';
import { ComponentReconcile } from './component-reconcile.entity';
import { EnvironmentReconcile } from './environment-reconcile.entity';
import { Environment } from './environment.entity';
import { EnvironmentSubscriber } from './environment.subscriber';
import { Team } from './team.entity';
import { ComponentSubscriber } from './component.subscriber';

export const entities = [
    User,
    Organization,
    Component, 
    ComponentReconcile,
    EnvironmentReconcile,
    Environment,
    Team
];

export const subscribers = [
    EnvironmentSubscriber,
    ComponentSubscriber
];

export {
    User,
    Organization,
    Component, 
    ComponentReconcile,
    EnvironmentReconcile,
    Environment,
    Team
}
