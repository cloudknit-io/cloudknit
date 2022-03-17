import { Company } from './company/Company';
import { costingEntities } from './costing';
import { User } from './entities/User';
import { reconcileEntities } from './reconciliation';
import { resourceEntities } from './resources';

export const entities = [Company, User, ...costingEntities, ...resourceEntities, ...reconcileEntities];
