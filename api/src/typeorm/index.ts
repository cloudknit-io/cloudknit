import { Component } from './component.entity';
import { reconcileEntities } from './reconciliation';
import { User } from './User.entity';
import { Organization } from './Organization.entity';

export const entities = [
    User,
    Organization,
    Component, 
    ...reconcileEntities
];

export {
    User,
    Organization
}
