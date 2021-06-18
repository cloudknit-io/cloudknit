import { costingEntities } from './costing';
import { User } from './entities/User';
import { resourceEntities } from './resources';

export const entities = [User, ...costingEntities, ...resourceEntities];
