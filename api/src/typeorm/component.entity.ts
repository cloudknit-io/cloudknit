import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToOne,
  PrimaryGeneratedColumn,
  RelationId,
  UpdateDateColumn,
} from 'typeorm';
import { Organization } from './Organization.entity';
import { ComponentReconcile } from './component-reconcile.entity';
import { Environment } from './environment.entity';

@Entity({ name: 'components' })
@Index(['organization', 'environment', 'name'], { unique: true })
export class Component {
  @PrimaryGeneratedColumn()
  id: number;

  @ManyToOne(() => Environment, (environment) => environment.components)
  environment: Environment;

  @Column({
    name: 'component_name',
  })
  name: string;

  @Column({
    default: 'terraform',
  })
  type: string;

  @UpdateDateColumn({
    name: 'last_reconcile_datetime',
  })
  lastReconcileDatetime: string;

  @Column({
    name: 'is_deleted',
    default: null,
    type: 'boolean'
  })
  isDeleted: boolean;

  @OneToOne(() => ComponentReconcile, (compRecon) => compRecon.component, {
    eager: false
  })
  @JoinColumn({
    referencedColumnName: 'reconcileId',
    name: 'latest_comp_recon_id'
  })
  latestCompRecon: ComponentReconcile

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({
    referencedColumnName: 'id',
  })
  organization: Organization;

  @RelationId((comp: Component) => comp.environment)
  envId: number;

  @RelationId((comp: Component) => comp.organization)
  orgId: number;
}
