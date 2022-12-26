import { CostResource } from 'src/costing/dtos/Resource.dto';
import { Column, Entity, JoinColumn, ManyToOne, UpdateDateColumn } from 'typeorm';
import { Organization } from './Organization.entity';
import { Environment } from './reconciliation/environment.entity';

@Entity({ name: 'components' })
export class Component {
  // TODO : Get rid of this.
  @Column({
    primary: true,
    name: 'id',
  })
  id: string;

  @Column({
    name: 'team_name',
  })
  teamName: string;

  @ManyToOne(() => Environment, (environment) => environment.components, {
    eager: true
  })
  environment: Environment;

  @Column({
    name: 'component_name',
  })
  componentName: string;

  @Column({
    name: 'status'
  })
  status: string;

  @Column({
    name: 'estimated_cost',
    type: 'decimal',
    precision: 10,
    scale: 3,
  })
  estimatedCost: number = 0;

  @UpdateDateColumn({
    name: 'last_reconcile_datetime'
  })
  lastReconcileDatetime: string;

  @Column({
    default: -1
  })
  duration: number;

  @Column({
    default: false,
    type: 'boolean'
  })
  isDestroyed?: boolean

  @Column({
    name: 'cost_resources',
    default: null,
    type: 'json'
  })
  costResources: CostResource[];

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization
}
