import { costingEntities } from './costing';
import { reconcileEntities } from './reconciliation';
import { resourceEntities } from './resources';
import { User } from './User.entity';
import { Organization } from './Organization.entity';

export const entities = [
    User,
    Organization,
    ...costingEntities, 
    ...resourceEntities,
    ...reconcileEntities
];

export {
    User,
    Organization
}
