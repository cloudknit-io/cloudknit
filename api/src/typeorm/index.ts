import { costingEntities } from './costing';
import { User } from './entities/User';
import { reconcileEntities } from './reconciliation';
import { resourceEntities } from './resources';

export const entities = [User, ...costingEntities, ...resourceEntities, ...reconcileEntities];
